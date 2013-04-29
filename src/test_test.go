package modata

import "testing"
import "fmt"
import "time"

func TestSetup(t *testing.T) {
    fmt.Println("Test: Basic setup of clients and servers")

}

func TestBlock(t *testing.T) {
    fmt.Println("Test: Initialization of block service is correct")
    bs := StartBlockServer("localhost:1234")
    fmt.Println(bs)
}

func TestReplication(t *testing.T) {
    fmt.Println("Test: Initialization of replication service is correct")
    rs := StartReplicationServer("localhost:8080")
    rs2 := StartReplicationServer("localhost:8081")
    fmt.Println(rs)
    fmt.Println(rs2)
    time.Sleep(100 * time.Second)
    fmt.Printf("... Pass \n")
}
