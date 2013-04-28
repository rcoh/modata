package modata

import (
    "goweb"
    "fmt"
    "sync"
)

type ReplicationServer struct {
    mu sync.Mutex
}

func (rs *ReplicationServer) Replicate(node int) {
    fmt.Println("Replicating!")
}

func (rs *ReplicationServer) HelloHandler(c *goweb.Context) {
    fmt.Fprintf(c.ResponseWriter, "Hello there %s\n", c.PathParams["name"])
}

func StartReplicationServer(port string) {
    rs := new(ReplicationServer)
    rs.Replicate(4)

    //goweb.MapFunc("/hello/{name}/", rs.HelloHandler)
    goweb.ListenAndServe(port)
}
