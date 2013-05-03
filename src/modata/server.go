package modata

import (
    "fmt"
    "sync"
    "web"
    "strconv"
    "strings"
    "container/heap"
    "encoding/json"
)

type BlockServer struct {
    mu sync.Mutex
    contact Contact
    id NodeID
    server *web.Server
    data map[Key]string
    routingTable *RoutingTable

    MyContact Contact
}

//
// Store a key,value pair locally
//
func (bs *BlockServer) Store (c *web.Context) string {
    key, exists := c.Params["key"]
    file, _ := c.Params["data"]
    // Should verify the key is the hash of the data

    if exists {
        hashedKey := MakeKey(Hash(key))
        bs.data[hashedKey] = file
        return RespondOk()
    } else {
        return RespondWithStatus("FAIL", "NO KEY")
    }
    return RespondNotOk()
}

//
// Do a distributed store
//
func (bs *BlockServer) IterativeStore (c *web.Context) string {
    key, exists := c.Params["key"]

    if exists {
        hashedKey := Hash(key)
        replicated := make([]string, 0)
        nodelist := bs.IterativeFindNode(c, MakeHex(hashedKey))
        result := make(map[string] interface{})
        err := json.Unmarshal([]byte(nodelist), &result)
        if (err == nil) {
            contactList := MakeContactList(result["data"].([]interface{}))
            for _, contact := range contactList {
                status, _, dcontact := JsonPostUrl(contact.ToHttpAddress() +
                                                   "/store?key=" + c.Params["key"] +
                                                   "&data=" + c.Params["data"],
                                                   bs.contact)
                if (status == OK) {
                    bs.UpdateContact(dcontact)
                    replicated = append(replicated,
                                        dcontact.Addr + ":" +
                                        strconv.Itoa(dcontact.Port))
                }
            }
            fmt.Println(contactList)
        }
        return RespondWithData(map[string] interface{} {
            "replication-nodes" : replicated,
            "replication-count" : len(replicated),
        })
    } else {
        return RespondWithStatus("FAIL", "NO KEY")
    }
    return RespondOk()
}

// Locally find a value
func (bs *BlockServer) FindValue(c *web.Context, key string) string {
    hashedKey := MakeKey(Hash(key))
    value, ok := bs.data[hashedKey]
    if ok {
        return RespondWithData(value)
    }
    return RespondNotFound()
}

//
// Do a distributed lookup
//
func (bs *BlockServer) IterativeFindValue (c *web.Context, key string) (value string) {
    hashedKey := Hash(key)
    nodelist := bs.IterativeFindNode(c, MakeHex(hashedKey))
    result := make(map[string] interface{})
    err := json.Unmarshal([]byte(nodelist), &result)
    if (err == nil) {
        contactList := MakeContactList(result["data"].([]interface{}))
        for _, contact := range contactList {
            status, ddata, dcontact := JsonGet(contact.ToHttpAddress() +
            "/find-value/"+key,
            bs.contact)
            if (status == OK) {
                bs.UpdateContact(dcontact)
                value = ddata.(string)
                break
            }
        }
    }
    return RespondWithData(value)
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

//
// Locally find a node
//
func (bs *BlockServer) FindNode(c *web.Context, node string) string {
    nodeID := MakeNodeID(node)
    if (nodeID == NodeID{}) {
        return RespondWithStatus("NOTOK", "Invalid node identifier")
    }
    results := bs.routingTable.FindClosest(MakeNodeID(node), bs.routingTable.k)
    return RespondWithData(results)
}

//
// Do a distributed find node
//
func (bs *BlockServer) IterativeFindNode (c *web.Context, node string) string {
    nodeID := MakeNodeID(node)

    if (nodeID == NodeID{}) {
        return RespondWithStatus("NOTOK", "Invalid node identifier")
    }

    // Keep track of contacts we need to find
    contacts := bs.routingTable.FindClosest(nodeID, Alpha)
    heap.Init(&contacts)

    // Keeps nodes that we are finding
    results := make(ContactDistanceList, 0, bs.routingTable.k)
    heap.Init(&results)

    // Have we contacted this node before
    contacted := make(map[ContactDistance]bool)

    done := false

    for !done {
        closest := results.Peek()

        for i := 0; i < Alpha; i++ {
            if (len(contacts) == 0) {
                // No more nodes to check, we're done
                done = true
                break
            }

            contact := heap.Pop(&contacts).(ContactDistance)
            if (contact == ContactDistance{} || contacted[contact]) {
                // Don't contact this node, try another one
                i -= 1
                continue
            }

            fmt.Printf("Making request to %v\n", contact)
            status, data, _ := JsonGet(contact.Contact.ToHttpAddress() +
                                              "/find-node/" + 
                                              node,
                                              bs.contact)
            contacted[contact] = true
            if (status == OK) {
                heap.Push(&results, contact)
                bs.UpdateContact(contact.Contact)
                contactList := MakeContactDistanceList(data.([]interface{}))
                for _, dcontact := range contactList {
                    heap.Push(&contacts, dcontact)
                }
            } else {
                fmt.Println("All of the problems")
                // Handle removing bad nodes
            }
        }
        if closest == results.Peek() {
            done = true
            fmt.Println("Didn't find anything closer, done")
        }
    }

    fmt.Printf("Results: %v\n", results)

    result := make(ContactList, bs.routingTable.k)
    contact := ContactDistance{}
    last := 0
    for index := range(result) {
        for (contact == ContactDistance{} && len(results) > 0) {
            contact = heap.Pop(&results).(ContactDistance)
        }
        if (contact != ContactDistance{}) {
            result[index] = contact.Contact
            contact = ContactDistance{}
        } else {
            last = index
            break
        }
    }

    return RespondWithData(result[:last])
}

func (bs *BlockServer) Contact() Contact {
    return bs.contact
}

func StartBlockServer(name string) *BlockServer{
    bs := new(BlockServer)

    // Node id for kademlia
    id := MakeNode(MakeGUID())
    addr := strings.Split(name, ":")[0]
    port, _ := strconv.Atoi(strings.Split(name, ":")[1])

    bs.contact = Contact{id, addr, port}
    bs.data = make(map[Key]string)
    bs.server = web.NewServer()
    bs.routingTable = NewRoutingTable(20, MakeByteSlice(id))

    go func() {
        // Primitive store, stores the key,value in the local data
        bs.server.Post("/store", func (c *web.Context) string {
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.Store(c)
        })

        // Kademlia based store, does smart stuff
        bs.server.Post("/distributed/store", func (c *web.Context) string {
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.IterativeStore(c)
        })

        // Primitive find-value
        bs.server.Get("/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.FindValue(c, key)
        })


        // Kademlia based find-value, does smart stuff
        bs.server.Get("/distributed/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.IterativeFindValue(c, key)
        })


        // Primitive find-node
        bs.server.Get("/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.FindNode(c, node)
        })

        // Kademlia based find-node, does smart stuff
        bs.server.Get("/distributed/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return bs.IterativeFindNode(c, node)
        })

        // Ping endpoint, sees if the server is alive
        bs.server.Get("/ping", func (c *web.Context) string {
            bs.updateContactFromHeader(c)
            bs.prepareResponse(c)
            return RespondOk()
        })

        fmt.Printf("Listening on %v\n", name)
        bs.server.Run(name)
    }();

    return bs
}

func (bs *BlockServer) prepareResponse(c *web.Context) {
    c.ContentType("json")
    c.SetHeader("Modata-NodeID", HexNodeID(bs.contact.ID), true)
    c.SetHeader("Modata-Address", bs.contact.Addr, true)
    c.SetHeader("Modata-Port", fmt.Sprintf("%v", bs.contact.Port), true)
}

func (bs *BlockServer) updateContactFromHeader(c *web.Context) {
    srcContact, err := DecodeModataHeaders(c.Request.Header)
    if (err == nil) {
        // We have a valid contact, update it
        bs.UpdateContact(srcContact)
    } else {
        fmt.Printf("Could not decode modata headers: %v\n", err)
    }
}

func (bs *BlockServer) UpdateContact(contact Contact) {
    emptyContact := Contact{}
    if (contact != emptyContact) {
        // We propbably have a valid contact, update it
        fmt.Printf("Updating: %v\n", contact)
        bs.routingTable.Update(contact)
    }
}


