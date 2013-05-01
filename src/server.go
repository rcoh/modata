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
}

// Requests to know the node
func (rs *BlockServer) WhoHasNode (c *web.Context, node string) {
    c.WriteString("Don't know who owns " + node)
}

func (rs *BlockServer) upload (c *web.Context) {
    key, exists := c.Params["key"]
    file, _ := c.Params["data"]
    if exists {
        c.WriteString(key)
    } else {
        c.WriteString("NO KEY")
    }
    fmt.Println("VERIFY: ", verifyKV(key, file))
}

func (rs *BlockServer) get(c *web.Context, key string) {
    c.WriteString("Looks like you wanted: " + key)
}

func verifyKV(key string, value string) bool {
    return MakeHex(Hash(value)) == key
}

func StartBlockServer(name string) *BlockServer{
    rs := new(BlockServer)

    // Node identifier for chord
    rs.identifier = MakeHex(MakeGUID())
    rs.name = name

    rs.server = web.NewServer()

    go func() {
        // Identifier for this node
        rs.server.Post("/upload", func (c *web.Context) {
            rs.upload(c)
        })
        rs.server.Get("/get/(.*)", func (c *web.Context, node string) {
            rs.get(c, node)
        })

        fmt.Printf("Listening on %v\n", name)
        rs.server.Run(name)
    }();

    return rs
}


