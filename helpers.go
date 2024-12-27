package main

import (
  "fmt"
  "strings"
)

/*
func (tree *Btree) PrintTree(node *Node, level int) {

    if node == nil {
        return
    }

    // Create indentation based on level
    indent := strings.Repeat("    ", level)
    
    // Print current level and node keys
    fmt.Printf("Level %d: %s", level, indent)
    
    // Print keys with brackets and separators
    fmt.Print("[")
    for i, key := range node.keys {
        if i > 0 {
            fmt.Print("|")
        }
        fmt.Printf("%d", key)
    }
    fmt.Println("]")

    // Print all children with increased indentation
    for _, child := range node.children {
        if child != nil {
            tree.PrintTree(child, level+1)
        }
    }
}
*/
func (tree *Btree) printTreeRecursive(node *Node, prefix string, isLast bool, stringBuilder *strings.Builder) {
    if node == nil {
        return
    }

    nodeStr := "["
    for i, key := range node.keys {
        if i > 0 {
            nodeStr += "|"
        }
        nodeStr += fmt.Sprintf("%d", key)
    }
    nodeStr += "]"

    if isLast {
        stringBuilder.WriteString(fmt.Sprintf("%s└── %s\n", prefix, nodeStr))
    } else {
        stringBuilder.WriteString(fmt.Sprintf("%s├── %s\n", prefix, nodeStr))
    }

    childPrefix := prefix
    if isLast {
        childPrefix += "    "
    } else {
        childPrefix += "│   "
    }

    for i, child := range node.children {
        isLastChild := (i == len(node.children)-1)
        tree.printTreeRecursive(child, childPrefix, isLastChild, stringBuilder)
    }
}

func (tree *Btree) getHeight(node *Node) int {
    if node == nil {
        return 0
    }
    maxHeight := 0
    for _, child := range node.children {
        height := tree.getHeight(child)
        if height > maxHeight {
            maxHeight = height
        }
    }
    return maxHeight + 1
}

func (tree *Btree) PrintTree(node *Node, level int) {
    var sb strings.Builder
    fmt.Fprintln(&sb, "\nB-tree Structure:")
    fmt.Fprintf(&sb, "Height: %d\n", tree.getHeight(node)-1)
    tree.printTreeRecursive(node, "", true, &sb)
    fmt.Print(sb.String())
}
