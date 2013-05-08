package modata

import (
    "fmt"
    "web"
)

type ReplicationServer struct {
    MyContact Contact

    server *web.Server
    blockServer *BlockServer
}

func StartReplicationServer(name string, bs *BlockServer) *ReplicationServer{
    rs := new(ReplicationServer)

    rs.server = web.NewServer()

    go func() {
        // Identifier for this node

        rs.server.Get("/stats/", func (c *web.Context) string {
            return fmt.Sprintf("%v\n", bs.data)
        })

        fmt.Printf("Listening on %v\n", name)
        rs.server.Run(name)
    }();

    return rs
}

