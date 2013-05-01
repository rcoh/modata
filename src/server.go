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
func (rs *BlockServer) WhoHasNode (c *web.Context, node string) {
    c.WriteString("Don't know who owns " + node)
}

func (rs *BlockServer) Store (c *web.Context) string {
    key, exists := c.Params["key"]
    file, _ := c.Params["data"]
    // Should verify the key is the hash of the data
    fmt.Println("Verify: ", VerifyKV(key, file))
    fmt.Printf("Hashing %v to %v\n", key, Hash(key))


    if exists {
        hashedKey := MakeKey(Hash(key))
        rs.data[hashedKey] = file
        return RespondWithData(KeyValue(key, file))
    } else {
        return RespondWithStatus("FAIL", "NO KEY")
    }
    return RespondNotOk()
}

func (rs *BlockServer) FindValue(c *web.Context, key string) string {
    hashedKey := MakeKey(Hash(key))
    value, ok := rs.data[hashedKey]
    if ok {
        return RespondWithData(value)
    }
    return RespondNotFound()
}

func VerifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

func (rs *BlockServer) FindNode(c *web.Context, node string) string {
    return RespondOk()
}

func StartBlockServer(name string) *BlockServer{
    rs := new(BlockServer)

    // Node identifier for chord
    rs.identifier = MakeHex(MakeGUID())
    rs.name = name
    rs.data = make(map[Key]string)
    rs.server = web.NewServer()

    go func() {
        // Identifier for this node
        rs.server.Post("/store", func (c *web.Context) string {
            c.ContentType("json")
            return rs.Store(c)
        })
        rs.server.Get("/find-value/(.*)", func (c *web.Context, key string) string {
            c.ContentType("json")
            return rs.FindValue(c, key)
        })
        rs.server.Get("/find-node/(.*)", func (c *web.Context, node string) string {
            c.ContentType("json")
            return rs.FindNode(c, node)
        })

        rs.server.Get("/ping", func (c *web.Context) string {
            c.ContentType("json")
            response := fmt.Sprintf("Ping from %v acknowledged by %v\n",
                                     c.Request.RemoteAddr, rs.name)
            return RespondWithData(response)
        })

        fmt.Printf("Listening on %v\n", name)
        rs.server.Run(name)
    }();

    return rs
}


