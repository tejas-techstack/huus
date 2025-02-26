package kv

import (
  "fmt"
)

func printTree(t *BPTree) error {
  // print root.
  root, err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return fmt.Errorf("Error loading node.")
  }

  if root.key == nil {
    return fmt.Errorf("Tree is empty.")
  }

  printLevels(t, root , 0)

  return nil
}

func printSpaces(level int ){
  for i := 0; i < (4*level); i++ {
    fmt.Printf(" ")
  }
}

func printLevels(t *BPTree, cur *node, level int) {

  if cur == nil {
    return
  }

  if cur.isLeaf {
    printSpaces(level)
    fmt.Printf("Level %d : %v\n",level, cur)
    return
  }

  printSpaces(level)
  fmt.Printf("Level %d : %v\n", level, cur)
  for _, v := range cur.pointers {
    child, _ := t.storage.loadNode(v.asNodeId())
    printLevels(t, child, level+1)
  }
}

