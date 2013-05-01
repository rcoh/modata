package main

import (
    "flag"
    "modata"
    "fmt"
)

func main() {
    isBlock := flag.Bool("block", true, "Run the Block server")
    isReplication := flag.Bool("replocation", false, "Run the replication server")

    serverName := flag.String("server", "localhost:1234", "Name of the server")

    flag.Parse()

    fmt.Printf("Block: %v\n", *isBlock)
    fmt.Printf("Replication: %v\n", *isReplication)
    fmt.Printf("Name: %v\n", *serverName)

    var bs *modata.BlockServer
    if (*isBlock) {
        fmt.Println("Starting block server!")
        bs = modata.StartBlockServer(*serverName)
    }
    fmt.Println(bs)

    // Spin loop foreverz
    for {
        fmt.Println(bs.server)
    }
}
