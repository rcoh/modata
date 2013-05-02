package main

import (
    "flag"
    "modata"
    "fmt"
    "time"
    "strings"
)

func main() {
    isBlock := flag.Bool("block", true, "Run the Block server")
    isReplication := flag.Bool("replication", false, "Run the replication server")

    blockName := flag.String("block-server",
                             "localhost:1234",
                             "Name of block server")
    replicationName := flag.String("replication-server",
                                   "localhost:1337",
                                   "Name of replication server")
    bootstrap := flag.String("bootstrap",
                             "", "Comma-separated list of addr:port bootstrap servers")

    flag.Parse()

    fmt.Printf("Block: %v at %v\n", *isBlock, *blockName)
    fmt.Printf("Replication: %v at %v\n", *isReplication, *replicationName)

    // Block Server
    var bs *modata.BlockServer
    if (*isBlock) {
        fmt.Println("Starting block server!")
        bs = modata.StartBlockServer(*blockName)
    }

    // Replication Server
    var rs *modata.ReplicationServer
    if (*isReplication) {
        fmt.Println("Starting replication server")
        rs = modata.StartReplicationServer(*replicationName)
    }

    fmt.Println(bs)
    fmt.Println(rs)

    // Bootstrap if requested
    if *bootstrap != "" {
        servers := strings.Split(*bootstrap, ",")
        for _, s := range servers {
                status, _, _ := modata.JsonGet("http://" + s + "/ping", bs.MyContact)
                fmt.Printf("Ping status from %v: %v\n", s, status)
        }
    }

    // Spin loop foreverz
    for *isBlock || *isReplication {
        // Yield to other threads, usually the server threads
        time.Sleep(100 * time.Second)
    }
}
