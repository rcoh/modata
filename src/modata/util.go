package modata

import (
    "math/rand"
    "strconv"
    "encoding/hex"
    "crypto/sha1"
    "time"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
)

// REST convenience to marshall all the things to json
func RespondWithStatus(status string, data interface{}, node NodeID) string {
    response := make(map[string]interface{})
    response["status"] = status
    response["data"] = data
    response["node"] = node
    edata, _ := json.Marshal(response)
    fmt.Println(string(edata))
    return string(edata)
}

func RespondWithData(data interface{}, node NodeID) string {
    return RespondWithStatus(OK, data, node);
}

func RespondOk(node NodeID) string {
    return RespondWithStatus(OK, nil, node)
}

func RespondNotFound(node NodeID) string {
    return RespondWithStatus(NOTFOUND, nil, node)
}

func RespondNotOk(node NodeID) string {
    return RespondWithStatus(NOTOK, nil, node)
}


// REST convenience to unmarshall all the things from json
func JsonGet(uri string) (string, interface{}, interface{}) {
    // Make the http request
    resp, err := http.Get(uri)
    if (err != nil) { return ERROR, err }

    // Decode the body of the response into a []byte
    response := make(map[string]interface{})
    raw, err := ioutil.ReadAll(resp.Body)
    if (err != nil) { return ERROR, err }

    // Decode the []byte into the standard mapping
    err = json.Unmarshal(raw, &response)
    if (err != nil) { return ERROR, err }

    // Return the values TODO: check for existence of correct keys
    return response["status"].(string), response["data"], response["node"]
}

func JsonPost(uri string, data map[string]string) (string, interface{}, interface{}) {
    // Make the http post request
    values := make(url.Values)
    for k,v := range data {
        values.Set(k, v)
    }
    resp, err := http.PostForm(uri, values)
    if (err != nil) { return ERROR, err }

    // Decode the body of the response into a []byte
    response := make(map[string]interface{})
    raw, err := ioutil.ReadAll(resp.Body)
    if (err != nil) { return ERROR, err }

    // Decode the []byte into the standard mapping
    err = json.Unmarshal(raw, &response)
    if (err != nil) { return ERROR, err }

    // Return the values TODO: check for existence of correct keys
    return response["status"].(string), response["data"], response["node"]
}

func JsonPostUrl(uri string) (string, interface{}, interface{}) {
    blank := make(map[string]string)
    return JsonPost(uri, blank)
}

// Done with REST convenience methods

// String -> []Byte and visaversa
func MakeHex(barray []byte) string {
    return hex.EncodeToString(barray)
}

func MakeByteArray(str string) []byte {
    data, _ := hex.DecodeString(str)
    return data
}

func MakeNodeID(str string) (id NodeID) {
    s := MakeByteArray(str)
    if len(s) != IDLength {
        return id
    }
    for i, v := range s {
        id[i] = v
    }
    return id
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
