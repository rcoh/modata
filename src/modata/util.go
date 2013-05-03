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
    "strings"
)

// REST convenience to marshall all the things to json
func RespondWithStatus(status string, data interface{}) string {
    response := make(map[string]interface{})
    response["status"] = status
    response["data"] = data
    edata, _ := json.Marshal(response)
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

func DecodeModataHeaders (header http.Header) (Contact, error) {
    node := header.Get("Modata-NodeID")
    nodeID := MakeNodeID(node)
    addr := header.Get("Modata-Address")
    port, err := strconv.Atoi(header.Get("Modata-Port"))

    contact := Contact {nodeID, addr, port}
    if (node != "" && addr != "" && err == nil) {
        return contact, nil
    }
    return contact, fmt.Errorf("Invalid Modata headers")
}

// REST convenience to unmarshall all the things from json
func JsonGet(uri string, self Contact) (string, interface{}, Contact) {
    dstContact := Contact{}

    // Make the http request
    client := &http.Client{}
    req, _ := http.NewRequest("GET", uri, nil)

    nullID := NodeID{}
    if self.ID != nullID {
        req.Header.Set("Modata-NodeID", HexNodeID(self.ID))
        req.Header.Set("Modata-Address", self.Addr)
        req.Header.Set("Modata-Port", fmt.Sprintf("%v", self.Port))
    }
    resp, err := client.Do(req)
    if (err != nil) { return ERROR, err, dstContact}

    // Decode the headers into a contact
    dstContact, err = DecodeModataHeaders(resp.Header)
    if (err != nil) {
        dstContact = Contact{}
    }

    // Decode the body of the response into a []byte
    response := make(map[string]interface{})
    raw, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if (err != nil) { return ERROR, err, dstContact}

    // Decode the []byte into the standard mapping
    err = json.Unmarshal(raw, &response)
    if (err != nil) { return ERROR, err, dstContact}

    // Return the values TODO: check for existence of correct keys
    return response["status"].(string), response["data"], dstContact
}

func JsonPost(uri string, data map[string]string, self Contact) (string, interface{}, Contact) {
    dstContact := Contact{}

    // Make the http post request
    values := make(url.Values)
    for k,v := range data {
        values.Set(k, v)
    }

    client := &http.Client{}
    req, err := http.NewRequest("POST", uri, strings.NewReader(values.Encode()))
    if (err != nil) { return ERROR, err, dstContact}

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    nullID := NodeID{}
    if self.ID != nullID {
        req.Header.Set("Modata-NodeID", HexNodeID(self.ID))
        req.Header.Set("Modata-Address", self.Addr)
        req.Header.Set("Modata-Port", fmt.Sprintf("%v", self.Port))
    }

    resp, err := client.Do(req)
    if (err != nil) { return ERROR, err, dstContact}

    // Decode the headers into a contact
    dstContact, err = DecodeModataHeaders(resp.Header)
    if (err != nil) {
        dstContact = Contact{}
    }

    // Decode the body of the response into a []byte
    response := make(map[string]interface{})
    raw, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if (err != nil) { return ERROR, err, dstContact }

    // Decode the []byte into the standard mapping
    err = json.Unmarshal(raw, &response)
    if (err != nil) { return ERROR, err, dstContact }

    // Return the values TODO: check for existence of correct keys
    return response["status"].(string), response["data"], dstContact
}

func JsonPostUrl(uri string, self Contact) (string, interface{}, Contact) {
    blank := make(map[string]string)
    return JsonPost(uri, blank, self)
}

// Done with REST convenience methods

// String -> []Byte and visaversa
func MakeHex(barray []byte) string {
    return hex.EncodeToString(barray)
}

func HexNodeID(id NodeID) string {
    return MakeHex(MakeByteSlice(id))
}

func MakeByteSlice(id NodeID) []byte {
    var a [IDLength]byte = id
    return a[:]
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

func MakeKey(source []byte) Key {
    var pre []byte
    pre = source
    if (len(source) != IDLength) {
        fmt.Println("Fixing source byte array")
        pre = HashByte(source)
    }
    k := Key{}
    for i := range k {
        k[i] = pre[i]
    }
    return k
}

func MakeNode(source []byte) NodeID {
    var pre []byte
    pre = source
    k := NodeID{}
    if (len(source) != IDLength) {
        return k
    }
    for i := range k {
        k[i] = pre[i]
    }
    return k
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

func Pretty(v interface{}) string {
    jsonData,_ := json.MarshalIndent(v, "", "\t")
    return string(jsonData)
}
