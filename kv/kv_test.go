package kv

import (
  "testing"
)
/*
  "path"
  "fmt"
  "os"
*/

// function to create test tree in TempDir
func createTestTree() (*BPTree) {
  // set the config.
  return nil
}

// Open(path string) error
func TestOpen(t *testing.T) {
  /*
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  } else {
    t.Log(tree)
  }
*/
return 
}

// t.Get(key []byte) ([]byte, error)

// t.Put(key, value []byte) (error) 

// t.insertIntoNode(cur *node, key []byte, pointer) error 

// t.findLeaf(key []byte) (*node, error) {}

// t.findLeafToInsert(key []byte) (*node, error) {}

// t.findChildIndex(cur *node, key []byte) (int, error){}

// t.findChild(parent *node, key []byte) (*node, error) {}

// t.splitRoot() error {}

// t.splitNode(cur *node, parent *node) error {}

// t.insertKeyAt(cur, index, key[]byte) error {}

// t.insertValueAt(cur, index, value []byte) error {}

// t. insertNodeAt(cur, index, nodeId uint32) error {}

// compare(x, y []byte) int {}

// calcMinOrder(order uint16) {}
