package modata

import (
    "math/rand"
    "strconv"
    "encoding/hex"
    "crypto/sha1"
    "time"
    "encoding/json"
    "fmt"
)

// REST convenience to marshall all the things to json
func RespondWithStatus(status string, data interface{}) string {
    response := make(map[string]interface{})
    response["status"] = status
    response["data"] = data
    edata, _ := json.Marshal(response)
    fmt.Println(string(edata))
    return string(edata)
}

func RespondWithData(data interface{}) string {
    return RespondWithStatus(OK, data);
}

func RespondOk() string {
    return RespondWithStatus(OK, nil)
}

func RespondNotFound() string {
    return RespondWithStatus(NOTFOUND, nil)
}

func RespondNotOk() string {
    return RespondWithStatus(NOTOK, nil)
}


// String -> []Byte and visaversa
func MakeHex(barray []byte) string {
    return hex.EncodeToString(barray)
}

func MakeByteArray(str string) []byte {
    data, _ := hex.DecodeString(str)
    return data
}

// Creates a dictionary with a single key and a single value
func KeyValue(key string, value string) map[string]string {
    response := make(map[string]string)
    response[key] = value
    return response
}


// GUID generation
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

// SHA1s a string
func Hash(value string) []byte {
    // Generate hash of the provided string
	hasher := sha1.New()

    // Add some randomness
	hasher.Write([]byte(value))

    // Hash the result.
	return hasher.Sum(nil)
}

func HashByte(value []byte) []byte {
    hasher:= sha1.New()
    hasher.Write(value)
    return hasher.Sum(nil)
}

func MakeKey(source []byte) Key {
    var pre []byte
    pre = source
    if (len(source) != IDLength) {
        fmt.Println("Fixing source byte array")
        pre = HashByte(source)
    }
    fmt.Println(pre)
    k := Key{}
    for i := range k {
        k[i] = pre[i]
    }
    return k
}
