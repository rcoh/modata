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
    time.Sleep(100 * time.Second)
}

func xTestReplication(t *testing.T) {
    fmt.Println("Test: Initialization of replication service is correct")
    rs := StartReplicationServer("localhost:8080")
    rs2 := StartReplicationServer("localhost:8081")
    fmt.Println(rs)
    fmt.Println(rs2)
    time.Sleep(100 * time.Second)
    fmt.Println("... Pass")
}

func TestDistance(t *testing.T) {
    fmt.Println("Test: Distance function between Nodes")
    a := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
    b := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

    c := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

    if Distance(&a, &b) != c {
        t.Errorf("XOR: %v\n", Distance(&a, &b))
    }
    fmt.Println("... Pass")
}

func TestBucketing(t *testing.T) {
    fmt.Println("Test: Bucketing")
    a := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    b := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    rt := RoutingTable{}
    rt.me = a
    fmt.Printf("Bucket %d\n", rt.BucketForNode(b))
    fmt.Println("... Pass")
}

func ExampleUpdate() {
    a := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    b := NodeID{192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    rt := RoutingTable{}
    rt.k = 20
    rt.me = a

    c1 := Contact{b, "", 0}
    rt.Update(c1)

    // Make sure 
    fmt.Printf("Update 1: Size of bucket 1: %d\n", rt.buckets[1].Len())

    c2 := Contact{b, "", 0}
    rt.Update(c2)

    fmt.Printf("Update 2: Size of bucket 1: %d\n", rt.buckets[1].Len())

    // Output: Update 1: Size of bucket 1: 1
    // Update 2: Size of bucket 1: 1
}

func TestShortlist(t *testing.T) {
    fmt.Println("Test: Shortlist")

    a := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

    k := NodeID{64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

    b := NodeID{64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
    
    rt := RoutingTable{}
    rt.k = 20
    rt.me = a

    c := Contact{b, "", 0}
    rt.Update(c)
    c.port = 1
    rt.Update(c)

    fmt.Printf("%v\n", rt.buckets[1].Front().Value);

    fmt.Printf("%v\n", rt.SelectShortlist(k, 2))

    fmt.Println("... Pass")
}
