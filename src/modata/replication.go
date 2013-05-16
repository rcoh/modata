package modata

import (
  "fmt"
  "web"
  "sort"
  "time"
  "os"
  "log"
  "io/ioutil"
)

const REPL_INTERVAL = 10 * time.Second

type ReplicationServer struct {
  MyContact Contact

  server *web.Server
  blockServer *BlockServer
}

func BuildKeyMap(contacts ContactList) map[string][]string {
  keyMapping := make(map[string][]string)
  for _, contact := range contacts {
    status, data, _ := JsonGet(contact.ToHttpAddress() + "/keys", Contact{})
    if (status == OK) {
      for _, key := range data.([]interface{}) {
        skey := key.(string)
        lst, ok := keyMapping[skey]
        if ok {
          keyMapping[skey] = append(lst, contact.ToHttpAddress())
        } else {
          keyMapping[skey] = []string{contact.ToHttpAddress()}
        }
      }
    }
  }
  return keyMapping
}

func (rs *ReplicationServer) IntelligentReplication () {
  contacts := rs.blockServer.routingTable.AllContacts()
  if len(contacts) < rs.blockServer.routingTable.k {
    // Not enough nodes to replicate up to k
    return
  }
  contacts = append(contacts, rs.blockServer.contact)
  keyMapping := BuildKeyMap(contacts)

  // Stores repl_count -> keys 
  replicationMap := make(map[int][]string)
  for key, value := range keyMapping {
    replicationCount := len(value)
    lst, ok := replicationMap[replicationCount]
    if ok {
      replicationMap[replicationCount] = append(lst, key)
    } else {
      replicationMap[replicationCount] = []string{key}
    }
  }

  replicationCounts := make([]int,0)
  for key, _ := range replicationMap {
    replicationCounts = append(replicationCounts, key)
  }
  sort.Ints(replicationCounts)

  max := 3
  for _, value := range replicationCounts {
    if value >= rs.blockServer.routingTable.k {
      // Don't replicate fully replicated nodes
      break
    }
    keysToReplicate := replicationMap[value]
    for _, key := range keysToReplicate {
      fmt.Printf("Re-replicating %v with replication %v\n", key, value)
      server := keyMapping[key][0]
      JsonPostUrl(server + "/replicate?key=" + key, Contact{})
    }
    max -= 1
    if (max <= 0) {
      break
    }
  }


}

func (rs *ReplicationServer) KeyMap(c *web.Context) string{
  contacts := rs.blockServer.routingTable.AllContacts()
  contacts = append(contacts, rs.blockServer.contact)
  keyMapping := BuildKeyMap(contacts)
  finalData := make(map[string]interface{})
  finalData["stats"] = map[string]int{"numKeys": len(keyMapping), "contacts": len(contacts),
                        "averageKeysPerContact": len(keyMapping) /
                        len(contacts)}
  finalData["keys"] = make(map[string]interface{})
  for key, value := range keyMapping {
    sort.Strings(value)
    finalData["keys"].(map[string]interface{})[key] = make(map[string]interface{})
    finalData["keys"].(map[string]interface{})[key].(map[string]interface{})["replicationCount"] = len(value)
    finalData["keys"].(map[string]interface{})[key].(map[string]interface{})["nodes"] = value
  }
  return RespondWithData(finalData)
}

func StartReplicationServer(name string, bs *BlockServer) *ReplicationServer{
  rs := new(ReplicationServer)

  rs.server = web.NewServer()
  rs.blockServer = bs

  logFile, err := os.OpenFile(name + ".log", os.O_CREATE | os.O_APPEND, 0644)
  if err == nil {
    rs.server.SetLogger(log.New(logFile, "", log.Ldate | log.Ltime))
    fmt.Printf("Set logger to %v\n", logFile)
    defer logFile.Close()
  }

  go func() {
    rs.server.Get("/keymap", func (c *web.Context) string {
      c.ContentType("json")
      c.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
      return rs.KeyMap(c)
    })

    rs.server.Get("/", func(c *web.Context) string {
      c.ContentType("text/html")
      b, _ := ioutil.ReadFile("../viz/viz.html")
      return string(b)
    })

    fmt.Printf("Listening on %v\n", name)
    rs.server.Run(name)
  }();

  // create a thread to do intelligent replication every now and again
  go func() {
    time.Sleep(10 * time.Second)
    for {
      rs.IntelligentReplication()
      time.Sleep(REPL_INTERVAL)
    }
  }()

  return rs
}

