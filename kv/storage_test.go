package kv

import (
  "testing"
  "os"
  "path"
  "fmt"
)

func TestNewStorage(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")
  defer func () {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s : %s", dbDir, err))
    }
  }()

  s, err := newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating storage : %s", err)
  }

  t.Log(s.pageSize,s.freePages,s.lastPageId)

  // reload to check if it loads from pre existing file correctly.
  s, err = newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating storage : %s", err)
  }
  t.Log(s.pageSize,s.freePages,s.lastPageId)

}

func TestNewNode(t *testing.T) {

  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s: %s", dbDir, err))
    }
  }()

  s, err := newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating newStorage : %s", err)
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

  s, err := newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating newStorage : %s", err)
  }

  newNodeId, err := s.newNode()
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  }

  node, err := s.loadNode(newNodeId)
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  } else {
    t.Log("Node data :", node)
  }
}

func TestUpdateNode(t *testing.T) {

  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")
  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  s, err := newStorage(path.Join(dbDir, "test.db"), 4096)
  if err != nil {
    t.Fatalf("Error creating newStorage : %s", err)
  }

  newNodeId, err := s.newNode()
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  } else {
    t.Log(newNodeId)
  }

  // t.Log("id thats gonna load with raw data:", newNodeId)
  node, err := s.loadNode(newNodeId)
  if err != nil {
    t.Fatalf("Error creating newNode : %s", err)
  }

  t.Log("Node before updating : ", node)

  for i := 0; i<100; i++ {
    node.key = append(node.key, encodeUint32(uint32(i)))
  }

  err = s.updateNode(node)
  if err != nil {
    t.Fatalf("Error updating node : %s", err)
  }

  nodeAfterChange, err := s.loadNode(node.id)
  if err != nil {
    t.Fatalf("Error updating node")
  }

  t.Log("Node after updating with more keys: ", nodeAfterChange)

  node.key = nil

  err = s.updateNode(node)
  if err != nil {
    t.Fatalf("Error updaing node : %s", err)
  }

  nodeAfterChange, err = s.loadNode(node.id)
  if err != nil {
    t.Fatalf("Error updaing node : %s", err)
  }

  t.Log("Node after updaing with 0 keys: ", nodeAfterChange) 
}
