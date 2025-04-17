package engine

import (
  "fmt"
)

func printTree(t *BPTree) error {
  // print root.
  if t.metadata == nil {
    fmt.Println("Tree is empty.")
    return nil
  }
  root, err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return fmt.Errorf("Error loading node.")
  }

  if len(root.key) == 0 {
    fmt.Println("Tree is empty.")
    return nil
  }

  printLevels(t, root , 0)

  return nil
}

func printSpaces(level int ){
  for i := 0; i < (4*level); i++ {
    fmt.Printf(" ")
  }
}

func printNode(level int, cur *node) {

  fmt.Printf("Level %v: { %v keys : %v ", level, cur.id, cur.key)

  if cur.isLeaf {
    fmt.Printf("Values : [ ")
    for _, v := range cur.pointers {
      fmt.Printf("%v ", v.asValue())
    }
  } else {
    fmt.Printf("Node ids : [ ")
    for _, v := range cur.pointers {
      fmt.Printf("%v ", v.asNodeId())
    }
  }

  fmt.Printf("] %v %v }\n", cur.isLeaf, cur.sibling)
}

func printLevels(t *BPTree, cur *node, level int) {

  if cur == nil {
    return
  }

  if cur.isLeaf {
    printSpaces(level)
    printNode(level, cur)
    return
  }

  printSpaces(level)
  printNode(level, cur)
  for _, v := range cur.pointers {
    child, _ := t.storage.loadNode(v.asNodeId())
    printLevels(t, child, level+1)
  }
}

