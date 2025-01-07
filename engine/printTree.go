package engine

import (
  "fmt"
  "strings"
)

// helper function to print the tree.
func (tree *BPtree) Print() {
   if tree.root == nil {
       fmt.Println("Empty tree")
       return
   }
   printNode(tree.root, 0)
}

func printNode(node *Node, level int) {
   indent := strings.Repeat("  ", level)
   
   fmt.Printf("%sNode(leaf:%v): ", indent, node.isLeaf)
   for _, kv := range node.kvStore {
       fmt.Printf("%d ", kv.key)
   }
   fmt.Println()

   if !node.isLeaf {
       for _, child := range node.children {
           printNode(child, level+1)
       }
   }
}
