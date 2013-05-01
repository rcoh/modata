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

// Requests to know the node
func (bs *BlockServer) WhoHasNode (c *web.Context, node string) {
    c.WriteString("Don't know who owns " + node)
}

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

func (bs *BlockServer) FindValue(c *web.Context, key string) string {
    hashedKey := MakeKey(Hash(key))
    value, ok := bs.data[hashedKey]
    if ok {
        return RespondWithData(value)
    }
    return RespondNotFound()
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

func (bs *BlockServer) FindNode(c *web.Context, node string) string {
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
        // Identifier for this node
        bs.server.Post("/store", func (c *web.Context) string {
            c.ContentType("json")
            return bs.Store(c)
        })
        bs.server.Get("/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            return bs.FindValue(c, key)
        })
        bs.server.Get("/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            return bs.FindNode(c, node)
        })

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


