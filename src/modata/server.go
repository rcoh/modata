package modata

import (
    "fmt"
    "sync"
    "web"
)

type BlockServer struct {
    mu sync.Mutex
    name string
    identifier string
    server *web.Server
    data map[Key]string
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
    hashedKey := MakeKey(Hash(key))
    fmt.Println(hashedKey)
    return RespondOk()
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

//
// Locally find a node
//
func (bs *BlockServer) FindNode(c *web.Context, node string) string {
    return RespondOk()
}

//
// Do a distributed find node
//
func (bs *BlockServer) IterativeFindNode (c *web.Context, node string) string {
    return RespondOk()
}

func StartBlockServer(name string) *BlockServer{
    bs := new(BlockServer)

    // Node identifier for chord
    bs.identifier = MakeHex(MakeGUID())
    bs.name = name
    bs.data = make(map[Key]string)
    bs.server = web.NewServer()

    go func() {
        // Primitive store, stores the key,value in the local data
        bs.server.Post("/store", func (c *web.Context) string {
            c.ContentType("json")
            return bs.Store(c)
        })

        // Kademlia based store, does smart stuff
        bs.server.Post("/distributed/store", func (c *web.Context) string {
            c.ContentType("json")
            return bs.IterativeStore(c)
        })

        // Primitive find-value
        bs.server.Get("/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            return bs.FindValue(c, key)
        })


        // Kademlia based find-value, does smart stuff
        bs.server.Get("/distributed/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            return bs.IterativeFindValue(c, key)
        })


        // Primitive find-node
        bs.server.Get("/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            return bs.FindNode(c, node)
        })

        // Kademlia based find-node, does smart stuff
        bs.server.Get("/distributed/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            return bs.IterativeFindNode(c, node)
        })

        // Ping endpoint, sees if the server is alive
        bs.server.Get("/ping", func (c *web.Context) string {
            c.ContentType("json")
            response := fmt.Sprintf("Ping from %v acknowledged by %v\n",
                                     c.Request.RemoteAddr, bs.name)
            return RespondWithData(response)
        })

        fmt.Printf("Listening on %v\n", name)
        bs.server.Run(name)
    }();

    return bs
}


