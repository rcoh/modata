package modata

import (
    "math/rand"
    "strconv"
    "encoding/hex"
    "crypto/sha1"
    "time"
)

func MakeHex(barray []byte) string {
    return hex.EncodeToString(barray)
}

func MakeGUID() []byte {
	// UID generation, basically simplified UUID rfc spec
	hasher := sha1.New()

    // Add some randomness
	hasher.Write([]byte(strconv.Itoa(rand.Int())))

    // Add the time in nanoseconds
	hasher.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))

    // Hash the result. This should result in something that is entirely unique
    // unless we get a hash collision or two random numbers are generated the
    // same at the same exact time
	return hasher.Sum(nil)
}

func Hash(value string) []byte {
    // Generate hash of the provided string
	hasher := sha1.New()

    // Add some randomness
	hasher.Write([]byte(value))

    // Hash the result.
	return hasher.Sum(nil)
}

func XOR(left []byte, right []byte) []byte {
    barray := make([]byte, len(left))
    for i := range barray {
        barray[i] = left[i] ^ right[i]
    }
    return barray
}


