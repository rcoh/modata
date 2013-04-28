package modata

import "testing"
import "fmt"

func TestSetup(t *testing.T) {
    fmt.Println("Test: Basic setup of clients and servers")

}

func TestReplication(t *testing.T) {
    fmt.Println("Test: Initialization of replication service is correct")
    StartReplicationServer(":8080")
}
