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
    file, _ := c.Params["file"]
    if exists {
        c.WriteString(key)
    } else {
        c.WriteString("NO KEY")
    }
}

func StartBlockServer(name string) *BlockServer{
    rs := new(BlockServer)

    // Node identifier for chord
    rs.identifier = MakeGUID()
    rs.name = name

    rs.server = web.NewServer()

    go func() {
        // Identifier for this node
        rs.server.Post("/upload", func (c *web.Context) {
            rs.upload(c)
        })
        rs.server.Get("/whohas/(.*)", func (c *web.Context, node string) {
            rs.WhoHasNode(c, node)
        })

        fmt.Printf("Listening on %v\n", name)
        rs.server.Run(name)
    }();

    return rs
}


