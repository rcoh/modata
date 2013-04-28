package modata

import (
    "goweb"
    "fmt"
    "sync"
)

type ReplicationServer struct {
    mu sync.Mutex
    port string
}

func (rs *ReplicationServer) Replicate(node int) {
    fmt.Println("Replicating!")
}

func (rs *ReplicationServer) HelloController (c *goweb.Context) {
    fmt.Fprintf(c.ResponseWriter, "Hello there %s\n", c.PathParams["name"])
}

type KeyOwner struct {
    Key string
    Node string
}

// Requests to know the node
func (rs *ReplicationServer) WhoHasNode (c *goweb.Context) {
    key := c.PathParams["key"]
    c.RespondWithData(KeyOwner{key, "none"})
}

func (rs *ReplicationServer) AddNode (c *goweb.Context) {
    node := c.PathParams["node"]
    fmt.Println(node)
    c.RespondWithOK()
}

func StartReplicationServer(port string) *ReplicationServer{
    rs := new(ReplicationServer)

    rs.port = port

    go func() {
        goweb.ConfigureDefaultFormatters()
        goweb.MapFunc("/hello/{name}", func(c * goweb.Context) {
            rs.HelloController(c)
        })
        goweb.MapFunc("/whohas/{key}", func(c *goweb.Context) {
            rs.WhoHasNode(c)
        })
        goweb.ListenAndServe(port)
    }();

    return rs
}

