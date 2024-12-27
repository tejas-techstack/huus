package main

import (
  "errors"
)

type Node struct {
  keys []int
  children []*Node
  isLeaf bool
} 

type Btree struct {
  root *Node
  minNode int
}


// set env variable for minNode

func createNewTree(minDegree int) *Btree {
  return &Btree{
    root : &Node{
      keys : make([]int, 0),
      children : make([]*Node, 0),
      isLeaf : true,
    },
    minNode: minDegree,
  }
}
func (tree *Btree) SplitChild(parent *Node, index int) {
    child := parent.children[index]

    // Create new node that will store right half of the child's keys
    newNode := &Node{
        keys:     make([]int, 0),
        children: make([]*Node, 0),
        isLeaf:   child.isLeaf,
    }

    midKey := child.keys[tree.minNode-1]

    newNode.keys = append(newNode.keys, child.keys[tree.minNode:]...)
    if !child.isLeaf {
        newNode.children = append(newNode.children, child.children[tree.minNode:]...)
    }

    child.keys = child.keys[:tree.minNode-1]
    if !child.isLeaf {
        child.children = child.children[:tree.minNode]
    }

    parent.keys = append(parent.keys, 0)
    copy(parent.keys[index+1:], parent.keys[index:])
    parent.keys[index] = midKey

    parent.children = append(parent.children, nil)
    copy(parent.children[index+2:], parent.children[index+1:])
    parent.children[index+1] = newNode
}

func (tree *Btree) Insert(key int) (error) {
  // if the root is full pre-emptively split the root
  // this makes it such that the root will never be full
  // TODO: need to run tests on this to check if it is faster
  // than splitting the root only when the promotion of split happens to root

  root := tree.root

  if len(root.keys) == tree.minNode*2 - 1 {
    newRoot := &Node{
      keys : []int{},
      children: []*Node{root},
      isLeaf : false,
    }
    tree.root = newRoot
    tree.SplitChild(newRoot, 0)
    return tree.InsertNonFull(newRoot, key)
  }

  // once the root has been split
  // insert the key by searching
  return tree.InsertNonFull(root, key)
}

func (tree *Btree) InsertNonFull(node *Node, key int) (error){
  // find the node to insert into.
  
  i := len(node.keys) - 1
  if node.isLeaf {
    for i >= 0 && key < node.keys[i]{
      if key == node.keys[i]{
        return errors.New("Key already exists")
      }
      i--
    }

    if i >= 0 && node.keys[i] == key{
      return errors.New("Key already exists")
    }

    node.keys = append(node.keys[:i+1], append([]int{key}, node.keys[i+1:]...)...)
    return nil
  }

  for i >= 0 && key < node.keys[i]{
    if key == node.keys[i]{
      return errors.New("Key already exists")
    }
    i--
  }

  if i >= 0 && node.keys[i] == key{
    return errors.New("key already exists")
  }

  childIndex := i+1

  if len(node.children[childIndex].keys) == tree.minNode * 2-1{
    tree.SplitChild(node, childIndex)
    if node.keys[childIndex] < key{
      (childIndex)++
    }
  }

  // recursively call function to fill a child
  return tree.InsertNonFull(node.children[childIndex], key)
}


