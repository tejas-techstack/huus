package kv

import (
  "testing"
  "path"
  "fmt"
  "os"
)

/*
// Open(path string) error
func TestOpen(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  t.Log(tree.order, tree.storage, tree.metadata, tree.minKeyNum)

}
*/

/*
func TestGet(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  key := []byte{1}
  val, exists, err := tree.Get(key)
  if err != nil {
    t.Fatalf("Error getting value : %s", err)
  }

  if exists {
    t.Log(val)
  } else {
    t.Log("Value does not exist.")
  }

}
*/

/*
func TestInitRoot(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  key := []byte{15}
  val := []byte{18}
  err = tree.initializeRoot(key, val)
  if err != nil {
    t.Fatalf("Error initializing root : %s", err)
  }


  // t.Log(tree.storage.metadata)
  root, _ := tree.storage.loadNode(tree.metadata.rootId)
  t.Log(root)
}
*/

// t.Put(key, value []byte) (error) 
func TestPut(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  for i := 1; i < 5; i++{
    key := []byte{byte(i)}
    val := []byte{byte(i)}
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Error inserting key : %s", err)
    }
  }

  root, _ := tree.storage.loadNode(tree.metadata.rootId)
  
  for i := 0; i < 4; i++{
    key := root.key[i]
    val := root.pointers[i].asValue()
    t.Logf("%d %d", key, val)
  }
}

func TestPutAndGet(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 100, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  for i := 1; i < 100; i++{
    key := encodeUint64(i)
    val := encodeUint64(i)
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Error inserting key : %s", err)
    }
  }

  val, exists, err := tree.Get(encodeUint64(101))
  if err != nil {
    t.Fatalf("Error getting value : %s", err)
  }
  if !exists {
    t.Log("Key does not exist")
  } else {
    t.Log("key exists:", val)
  }

  val, exists, err = tree.Get(encodeUint64(99))
  if err != nil {
    t.Fatalf("Error getting value : %s", err)
  }
  if !exists {
    t.Log("Key does not exist")
  } else {
    t.Log("key exists:", val)
  }
}

func TestSplitRoot(t *testing.T) {
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 5, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }

  for i := 1; i < 100; i++{
    key := encodeUint64(i)
    val := encodeUint64(i)
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Could not insert key : %s", err)
    }
  }
}

/*
func TestChaining(t *testing.T){
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 1000, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }


  for i := 1; i < 10000; i++{
    key := encodeUint64(i)
    val := encodeUint64(i)
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Could not insert key : %s", err)
    }
  }  
}
*/

func TestDelete(t *testing.T){
  dbDir, _ := os.MkdirTemp(os.TempDir(), "example")

  defer func() {
    if err := os.RemoveAll(dbDir); err != nil {
      panic(fmt.Errorf("failed to remove %s:%s", dbDir, err))
    } 
  }()

  tree, err := Open(path.Join(dbDir, "example.db"), 20, 4096)
  if err != nil {
    t.Fatalf("Error opening tree : %s", err)
  }


  for i := 1; i < 25; i++{
    key := encodeUint64(i)
    val := encodeUint64(i)
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Could not insert key : %s", err)
    }
  }

  node, err := tree.storage.loadNode(tree.metadata.rootId)
  node, err = tree.storage.loadNode(node.pointers[0].asNodeId())
  t.Log("Node before deletion:", node.key)

  _, err = tree.Delete(encodeUint64(9))
  if err != nil {
    t.Fatalf("Error deleting key : %s", err)
  }

  node, err = tree.storage.loadNode(tree.metadata.rootId)
  node, err = tree.storage.loadNode(node.pointers[0].asNodeId())
  t.Log("Node after deletion:", node.key)
}
