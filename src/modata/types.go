package modata

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

