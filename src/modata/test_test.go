package modata

import "testing"
import "fmt"
import "time"


func TestSetup(t *testing.T) {
    fmt.Println("Test: Basic setup of clients and servers")
}

func TestBlock(t *testing.T) {
    fmt.Println("Test: Local Block service is correct")
    bs := StartBlockServer("localhost:1234")

    time.Sleep(2 * time.Second)
    fmt.Println(bs)


    status, data, nodeid := JsonGet("http://localhost:1234/find-value/foo")
    if (status != NOTFOUND) {
        fmt.Println(data)
        t.Errorf("Nonexistent key exists\n")
    }

    status, data, nodeid = JsonPostUrl("http://localhost:1234/store?key=foo&data=bar")
    if (status != OK) {
        t.Errorf("Could not post\n")
    }

    status, data, nodeid = JsonGet("http://localhost:1234/find-value/foo")
    if (status != OK && data != "bar") {
        t.Errorf("Incorrect data exists\n")
    }
    fmt.Println(nodeid)

    // a := NodeID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    b := NodeID{64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
    c := NodeID{64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
    d := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 1}

    ct := Contact{b, "", 0}
    bs.routingTable.Update(ct)
    ct.Port = 1
    bs.routingTable.Update(ct)

    ct.ID = c
    bs.routingTable.Update(ct)

    ct.ID = d
    bs.routingTable.Update(ct)

    status, data, nodeid = JsonGet("http://localhost:1234/find-node/400000000000000000000000000000000000000000")
    fmt.Printf("Find node: %v\n", data)

    fmt.Println("... Pass")
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

    c := NodeID{64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

    d := NodeID{128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 1}
    
    rt := RoutingTable{}
    rt.k = 20
    rt.me = a

    ct := Contact{b, "", 0}
    rt.Update(ct)
    ct.Port = 1
    rt.Update(ct)

    ct.ID = c
    rt.Update(ct)

    ct.ID = d
    rt.Update(ct)

    fmt.Printf("%v\n", rt.buckets[1].Front().Value);

    fmt.Printf("%v\n", rt.FindClosest(k, 4))

    fmt.Println("... Pass")
}
