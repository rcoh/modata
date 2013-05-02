package modata

import (
    "fmt"
    "sync"
    "web"
    "strings"
    "strconv"
)

type BlockServer struct {
    mu sync.Mutex
    contact Contact
    id NodeID
    server *web.Server
    data map[Key]string
    routingTable *RoutingTable
}

//
// Store a key,value pair locally
//
func (bs *BlockServer) Store (c *web.Context) string {
    key, exists := c.Params["key"]
    file, _ := c.Params["data"]
    // Should verify the key is the hash of the data
    fmt.Println("Verify: ", VerifyKV(key, file))
    fmt.Printf("Hashing %v to %v\n", key, Hash(key))


    if exists {
        hashedKey := MakeKey(Hash(key))
        bs.data[hashedKey] = file
        return RespondWithData(KeyValue(key, file))
    } else {
        return RespondWithStatus("FAIL", "NO KEY")
    }
    return RespondNotOk()
}

//
// Do a distributed store
//
func (bs *BlockServer) IterativeStore (c *web.Context) string {
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
func (bs *BlockServer) IterativeFindValue (c *web.Context, key string) string {

    return RespondOk()
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

//
// Locally find a node
//
func (bs *BlockServer) FindNode(c *web.Context, node string) string {
    results := bs.routingTable.FindClosest(MakeNodeID(node), bs.routingTable.k)
    return RespondWithData(results)
}

//
// Do a distributed find node
//
func (bs *BlockServer) IterativeFindNode (c *web.Context, node string) string {
    nodeID := MakeNodeID(node)
    contacts := bs.routingTable.FindClosest(nodeID, Alpha)
    closest := Closest(nodeID, contacts)

    for _, contact := range contacts {
        status, data, nodeid := JsonGet(contact.ToHttpAddress(), bs.contact)
        fmt.Printf("%v, %v, %v\n", status, data, nodeid)
    }

    fmt.Println(contacts)
    fmt.Println(closest)

    return RespondOk()
}

func Closest (node NodeID, list ContactList) Contact {
    if (len(list) == 0) { return Contact{} }
    closest := list[0]
    for _, contact := range list {
        x := Distance(&(contact.ID), &node)
        y := Distance(&(closest.ID), &node)
        if x.LessThan(&y) {
            closest = contact
        }
    }
    return closest
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


