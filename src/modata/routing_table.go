package modata

import (
	"fmt"
	"container/list"
	"sort"
    "strconv"
)

type Contact struct {
	ID NodeID
	Addr string
	Port int
}

type ContactDistance struct {
    Contact Contact
    Distance NodeID
}

func (c *Contact) ToHttpAddress() string {
    return "http://" + c.Addr + ":"+ strconv.Itoa(c.Port)
}

type Bucket struct {
	list.List;
}

type RoutingTable struct {
	k int
	me NodeID
	buckets [IDLength*8]Bucket
}

func NewRoutingTable(k int, id []byte) (rt *RoutingTable) {
	if len(id) != IDLength {
		return nil
	}

	rt = &RoutingTable{}
	rt.k = k
	for i, v := range id {
		rt.me[i] = v
	}

	return rt
}

func Distance(a *NodeID, b *NodeID) (d NodeID) {
	for i, _ := range a {
		d[i] = a[i] ^ b[i]
	}
	return d
}

func Compare(a *NodeID, b *NodeID) int {
	for i := 0; i < IDLength; i++ {
		if a[i] > b[i] {
			return 1
		} else if a[i] < b[i] {
			return -1
		}
	}
	return 0
}

func (a *NodeID) GreaterThan(b *NodeID) bool {
	return Compare(a, b) == 1
}

func (a *NodeID) LessThan(b *NodeID) bool {
	return Compare(a, b) == -1
}

func (a *NodeID) Equals(b *NodeID) bool {
	return Compare(a, b) == 0
}


func MakeContactList(m []interface{}) ContactList {
	cl := make(ContactList, len(m))

	for i, v := range m {
		cl[i] = MakeContact(v)
	}
	return cl
}

func MakeContactDistanceList(m []interface{}) ContactDistanceList {
    cdl := make(ContactDistanceList, len(m))

    for i, v := range m {
        cdl[i] = MakeContactDistance(v)
    }
    return cdl
}

func MakeContact(v interface{}) Contact {
	m := v.(map[string]interface{})
	id := m["ID"].([]interface{})
	nid := [IDLength]byte{}
	for i, v := range id {
		nid[i] = byte(v.(float64))
	}

	return Contact{nid, m["Addr"].(string), int(m["Port"].(float64))}
}

func MakeContactDistance(v interface{}) ContactDistance {
	m := v.(map[string]interface{})

    id := m["Distance"].([]interface{})
	nid := [IDLength]byte{}
	for i, v := range id {
		nid[i] = byte(v.(float64))
	}

	return ContactDistance{MakeContact(m["Contact"]), nid}
}

func (rt *RoutingTable) Init() {
	for _, b := range rt.buckets {
		b.Init()
	}
}

func (rt *RoutingTable) BucketForNode(n NodeID) (i int) {
	for i := 0; i < (IDLength * 8); i++ {
		// Are the ith bits set?
		a := (rt.me[i / 8] & (128 >> (uint) (i % 8))) > 0
		b := (n[i / 8] & (128 >> (uint) (i % 8))) > 0

		if a != b {
			return i
		}
	}
	// Bucket for self?
	return -1
}

func (rt *RoutingTable) Update(c Contact) {
	if rt.me == c.ID {
		// Ignore contacts with self
		return
	}

	fmt.Printf("Update %v\n", c.ID)

	// See if we already have this contact
	bucket := &rt.buckets[rt.BucketForNode(c.ID)]
	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value == c {
			// New most recently contacted
			bucket.MoveToBack(e)
			return
		}
	}

	// Do we have space for a new contact?
	if bucket.Len() < rt.k {
		bucket.PushBack(c)
	} else {
		// Ping the least recently contacted
		// FIXME: PING
		if true {	// No response from ping
			bucket.Remove(bucket.Front())
			bucket.PushBack(c)
		} // Else ignore the new contact
	}
}

func (rt *RoutingTable) FindClosest(k NodeID, alpha int) (ContactDistanceList) {
    shortlist := make(ContactDistanceList, 0, alpha)

	idx := rt.BucketForNode(k)

	// Grab entire buckets until we have >= alpha nodes, then sort and return the alpha closest
	for len(shortlist) < alpha && idx >= 0 {
		bucket := &rt.buckets[idx]

		for e := bucket.Front(); e != nil; e = e.Next() {
            value := e.Value.(Contact)
            distance := Distance(&(value.ID), &k)
			shortlist = append(shortlist, ContactDistance{value, distance})
		}

		idx -= 1
	}

    // Search upward too if we don't have enough
    idx = rt.BucketForNode(k) + 1
    for len(shortlist) < alpha && idx < IDLength * 8 {
        bucket := &rt.buckets[idx]

        for e := bucket.Front(); e != nil; e = e.Next() {
            value := e.Value.(Contact)
            distance := Distance(&(value.ID), &k)
            shortlist = append(shortlist, ContactDistance{value, distance})
        }
        idx += 1
    }

	sort.Sort(shortlist)

	if alpha > len(shortlist) {
		return shortlist
	}

    return shortlist[:alpha]
}
