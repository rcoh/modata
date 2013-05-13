package modata

import (
  "fmt"
  "web"
  "sort"
)

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
    } else {
      fmt.Printf("Failed to get keys stored on %v\n", contact)
      fmt.Printf("Attempted to hit %v\n", contact.ToHttpAddress() + "/keys")
      fmt.Println(data)
      fmt.Println(status)
    }
  }
  return keyMapping
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

  go func() {
    // Identifier for this node

    rs.server.Get("/contacts", func (c *web.Context) string {
      c.ContentType("json")
      return RespondWithData(bs.routingTable.AllContacts())
    })

    rs.server.Get("/keymap", func (c *web.Context) string {
      c.ContentType("json")
      return rs.KeyMap(c)
    })

    fmt.Printf("Listening on %v\n", name)
    rs.server.Run(name)
  }();

  // create a thread to call tick() periodically
  /*
  go func() {
  }()
*/

  return rs
}

