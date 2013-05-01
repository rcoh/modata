package modata

import "bytes"
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
    time.Sleep(100 * time.Second)
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

func TestUtilXOR(t *testing.T) {
    fmt.Println("Testing XOR")

    x := XOR([]byte("123"), []byte("123"))
    xg := []byte{0,0,0}
    y := XOR([]byte("012"), []byte("112"))
    yg := []byte{1,0,0}
    if (bytes.Compare(x, xg) != 0) {
        fmt.Println("Bad Xor on equal items")
    }
    if (bytes.Compare(y, yg) != 0) {
        fmt.Println("Bad Xor on unequal items")
    }

}
