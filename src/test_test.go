package modata

import "testing"
import "fmt"
import "time"

func TestSetup(t *testing.T) {
    fmt.Println("Test: Basic setup of clients and servers")

}

func TestReplication(t *testing.T) {
    fmt.Println("Test: Initialization of replication service is correct")
    rs := StartReplicationServer(":8080")
    fmt.Println(rs)
    time.Sleep(10 * time.Second)
    fmt.Printf("... Pass \n")
}
