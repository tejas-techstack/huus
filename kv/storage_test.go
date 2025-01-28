package kv

import (
  "testing"
  "os"
  "path"
  "fmt"
)

func TestNewNode(t *testing.T) {

  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s: %s", dbDir, err))
    }
  }()

  /*
  s, err := newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating newStorage : %s", err)
  } else {
    t.Log(s)
  }
  */

  fo, _ := os.OpenFile(path.Join(dbDir, "test.db"), os.O_RDWR|os.O_CREATE, 0600)
  s := &storage{
    fo : fo,
    pageSize : 4096,
    lastPageId : 1001,
    metadata : &storageMetadata{
      pageSize : 4096,
      custom : nil,
    },
  }

  newNodeId, err := s.newNode()
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  } else {
    t.Log(newNodeId)
  }


}

func TestLoadNode(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  fo, _ := os.OpenFile(path.Join(dbDir, "test.db"), os.O_RDWR|os.O_CREATE, 0600)
  s := &storage{
    fo : fo,
    pageSize : 4096,
    lastPageId : 1001,
    metadata : &storageMetadata{
      pageSize : 4096,
      custom : nil,
    },
  }

  newNodeId, err := s.newNode()
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  }

  /*
  node_raw, err := s.loadNodeRaw(newNodeId)
  if err != nil {
    t.Fatalf("Error loading raw node data : %s", err)
  } else {
    t.Log(len(node_raw))
  }
  */

  node, err := s.loadNode(newNodeId)
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  } else {
    t.Log("Node data :", node)
  }



}
