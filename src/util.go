package modata

import (
    "math/rand"
    "strconv"
    "encoding/hex"
    "crypto/sha1"
    "time"
)

func MakeGUID() string {
	// UID generation, basically simplified UUID rfc spec
	hasher := sha1.New()

    // Add some randomness
	hasher.Write([]byte(strconv.Itoa(rand.Int())))

    // Add the time in nanoseconds
	hasher.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))

    // Hash the result. This should result in something that is entirely unique
    // unless we get a hash collision or two random numbers are generated the
    // same at the same exact time
	return hex.EncodeToString(hasher.Sum(nil))
}

func Hash(value string) string {
    // Generate hash of the provided string
	hasher := sha1.New()

    // Add some randomness
	hasher.Write([]byte(value))

    // Hash the result.
	return hex.EncodeToString(hasher.Sum(nil))
}


