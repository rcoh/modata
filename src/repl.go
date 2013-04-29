package modata

import (
    "fmt"
    "sync"
    "web"
)

type ReplicationServer struct {
    mu sync.Mutex
    name string
    identifier string
    server *web.Server
}

// Requests to know the node
func (rs *ReplicationServer) WhoHasNode (c *web.Context, node string) {
    c.WriteString("Don't know who owns " + node)
}

func (rs *ReplicationServer) Identifier (c *web.Context) {
    c.WriteString(rs.identifier)
}

func StartReplicationServer(name string) *ReplicationServer{
    rs := new(ReplicationServer)

    // Node identifier for chord
    rs.identifier = MakeGUID()
    rs.name = name

    rs.server = web.NewServer()

    go func() {
        // Identifier for this node
        rs.server.Get("/id", func (c *web.Context) {
            rs.Identifier(c)
        })
        rs.server.Get("/whohas/(.*)", func (c *web.Context, node string) {
            rs.WhoHasNode(c, node)
        })

        rs.server.Run(name)
    }();

    return rs
}

