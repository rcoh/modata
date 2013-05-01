package modata

import (
	// "fmt"
	"container/list"
	"sort"
)

type Contact struct {
	id NodeID
	addr string
	port int
}

type Bucket struct {
	list.List;
}

type RoutingTable struct {
	k int
	me NodeID
	buckets [IDLength*8]Bucket
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

// Sort support for contact lists

type ContactList []Contact

func (l ContactList) Len() int {
	return len(l)
}

func (l ContactList) Less(i, j int) bool {
	return l[i].id.LessThan(&l[j].id)
}

func (l ContactList) Swap(i, j int) {
	t := l[i]
	l[i] = l[j]
	l[j] = t
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
	if rt.me == c.id {
		// Ignore contacts with self
		return
	}

	// See if we already have this contact
	bucket := &rt.buckets[rt.BucketForNode(c.id)]
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

func (rt *RoutingTable) FindClosest(k NodeID, alpha int) (shortlist ContactList) {
	shortlist = make(ContactList, 0, alpha)

	idx := rt.BucketForNode(k)

	// Grab entire buckets until we have >= alpha nodes, then sort and return the alpha closest
	for len(shortlist) < alpha && idx >= 0 {
		bucket := &rt.buckets[idx]

		for e := bucket.Front(); e != nil; e = e.Next() {
			shortlist = append(shortlist, e.Value.(Contact))
		}

		idx -= 1
	}

	sort.Sort(shortlist)

	return shortlist[:alpha]
}
