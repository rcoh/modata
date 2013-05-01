package modata

import (
	// "fmt"
	"container/list"
)

type Contact struct {
	other NodeID
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
	// See if we already have this contact
	bucket := &rt.buckets[rt.BucketForNode(c.other)]
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

func (rt *RoutingTable) SelectShortlist(k NodeID, alpha int) (shortlist []Contact) {
	shortlist = make([]Contact, 0, alpha)

	idx := rt.BucketForNode(k)

	for len(shortlist) < alpha && idx >= 0 {
		bucket := &rt.buckets[idx]

		for e := bucket.Front(); e != nil && len(shortlist) < alpha; e = e.Next() {
			shortlist = append(shortlist, e.Value.(Contact))
		}

		idx -= 1
	}

	return shortlist
}
