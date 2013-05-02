package modata

import (
    "fmt"
    "sync"
    "web"
    "strconv"
    "strings"
)

type BlockServer struct {
    mu sync.Mutex
    name string
    addr string
    port int
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
    fmt.Println("Verify: ", VerifyKV(key, file))
    fmt.Printf("Hashing %v to %v\n", key, Hash(key))


    if exists {
        hashedKey := MakeKey(Hash(key))
        bs.data[hashedKey] = file
        return RespondWithData(KeyValue(key, file), bs.id)
    } else {
        return RespondWithStatus("FAIL", "NO KEY", bs.id)
    }
    return RespondNotOk(bs.id)
}

//
// Do a distributed store
//
func (bs *BlockServer) IterativeStore (c *web.Context) string {
    return RespondOk(bs.id)
}

// Locally find a value
func (bs *BlockServer) FindValue(c *web.Context, key string) string {
    hashedKey := MakeKey(Hash(key))
    value, ok := bs.data[hashedKey]
    if ok {
        return RespondWithData(value, bs.id)
    }
    return RespondNotFound(bs.id)
}

//
// Do a distributed lookup
//
func (bs *BlockServer) IterativeFindValue (c *web.Context, key string) string {
    hashedKey := MakeKey(Hash(key))
    fmt.Println(hashedKey)
    return RespondOk(bs.id)
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

//
// Locally find a node
//
func (bs *BlockServer) FindNode(c *web.Context, node string) string {
    results := bs.routingTable.FindClosest(MakeNodeID(node), bs.routingTable.k)
    return RespondWithStatus(OK, results, bs.id)
}

//
// Do a distributed find node
//
func (bs *BlockServer) IterativeFindNode (c *web.Context, node string) string {
    return RespondOk(bs.id)
}

func StartBlockServer(name string) *BlockServer{
    bs := new(BlockServer)

    // Node id for kademlia
    id := MakeGUID()

    bs.id = MakeNode(id)
    bs.name = name
    comps := strings.Split(name, ":")
    bs.addr = comps[0]
    bs.port, _ = strconv.Atoi(comps[1])
    bs.MyContact = Contact{bs.id, bs.addr, bs.port}
    bs.data = make(map[Key]string)
    bs.server = web.NewServer()
    bs.routingTable = NewRoutingTable(20, id)

    go func() {
        // Primitive store, stores the key,value in the local data
        bs.server.Post("/store", func (c *web.Context) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.Store(c)
        })

        // Kademlia based store, does smart stuff
        bs.server.Post("/distributed/store", func (c *web.Context) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.IterativeStore(c)
        })

        // Primitive find-value
        bs.server.Get("/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.FindValue(c, key)
        })


        // Kademlia based find-value, does smart stuff
        bs.server.Get("/distributed/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.IterativeFindValue(c, key)
        })


        // Primitive find-node
        bs.server.Get("/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.FindNode(c, node)
        })

        // Kademlia based find-node, does smart stuff
        bs.server.Get("/distributed/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            bs.updateContact(c)
            return bs.IterativeFindNode(c, node)
        })

        // Ping endpoint, sees if the server is alive
        bs.server.Get("/ping", func (c *web.Context) string {
            c.ContentType("json")
            bs.updateContact(c)
            response := fmt.Sprintf("Ping from %v acknowledged by %v\n",
                                     c.Request.RemoteAddr, bs.name)

            return RespondWithData(response, bs.id)
        })

        fmt.Printf("Listening on %v\n", name)
        bs.server.Run(name)
    }();

    return bs
}

func (bs *BlockServer) updateContact(c *web.Context) {
    h := c.Request.Header
    contact := Contact{}

    if h.Get("Modata-Address") == "" {
        // No valid contact data, ignore
        return
    }

    contact.ID = MakeNodeID(h.Get("Modata-NodeID"))
    contact.Addr = h.Get("Modata-Address")
    contact.Port, _ = strconv.Atoi(h.Get("Modata-Port"))

    fmt.Printf("Got contact from header: %v\n", contact)

    bs.routingTable.Update(contact)
}


