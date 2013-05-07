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

        fmt.Printf("Listening on %v\n", name)
        rs.server.Run(name)
    }();

    return rs
}

