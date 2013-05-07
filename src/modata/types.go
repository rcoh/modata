package modata

import "container/heap"

const Alpha = 3

// Bytes are 8 bits, so a size 20 byte array means a 160 bit nodeid
const IDLength = 20

// NodeIDs are 160 bit identifiers
type NodeID [IDLength]byte

// Keys are also 160 bit, so that they can be compared with Keys
type Key NodeID

// Response codes
const (
    OK = "OK"
    NOTFOUND = "NOTFOUND"
    NOTOK = "NOTOK"
    ERROR = "ERROR"
)

type KeyInfo struct {
    ClosestNodes ContactList
}

// Contact List, implements sort interface

type ContactList []Contact

func (l ContactList) Len() int {
	return len(l)
}

func (l ContactList) Less(i, j int) bool {
	return l[i].ID.LessThan(&l[j].ID)
}

func (l ContactList) Swap(i, j int) {
    if (i >= 0 && j >= 0) {
        l[i], l[j] = l[j], l[i]
    }
}


// ContactDistance List, implements heap interface

type ContactDistanceList []ContactDistance

func (l ContactDistanceList) Len() int {
	return len(l)
}

func (l ContactDistanceList) Less(i, j int) bool {
    return l[i].Distance.LessThan(&l[j].Distance)
}

func (l ContactDistanceList) Swap(i, j int) {
    if (i >= 0 && j >= 0) {
        l[i], l[j] = l[j], l[i]
    }
}

func (l *ContactDistanceList) Peek() interface{} {
    element := heap.Pop(l)
    heap.Push(l, element)
    return element
}

func (l *ContactDistanceList) Pop() interface{} {
    a := *l
    n := len(a)
    if (n == 0) { return ContactDistance{} }
    item := a[n-1]
    *l = a[0:n-1]
    return item
}

func (l *ContactDistanceList) Push(x interface{}) {
    *l = append(*l, x.(ContactDistance))
}

