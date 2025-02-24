package kv

import (
  "testing"
  "path"
  "fmt"
  "os"
)

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
    err = tree.Put(encodeUint64(i), encodeUint64(i))
    if err != nil {
      t.Fatalf("Error inserting key : %s", err)
    }
  }

  printTree(tree)
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

  // example of key not existing:
  val, exists, err := tree.Get(encodeUint64(101))
  if err != nil {
    t.Fatalf("Error getting value : %s", err)
  }
  if !exists {
    t.Log("Key does not exist")
  } else {
    t.Log("key exists:", val)
  }

  // example of key exisiting:
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

  printTree(tree)
}

func TestDelete(t *testing.T){
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


  for i := 1; i < 14; i++{
    key := encodeUint64(i)
    val := encodeUint64(i)
    err = tree.Put(key, val)
    if err != nil {
      t.Fatalf("Could not insert key : %s", err)
    }
  }

  fmt.Println("Tree before Deletion : ")
  printTree(tree)

  _, err = tree.Delete(encodeUint64(13))
  _, err = tree.Delete(encodeUint64(12))
  if err != nil {
    t.Fatalf("Error deleting key : %s", err)
  }

  fmt.Println("Tree after Deletion : ")
  printTree(tree)
}
