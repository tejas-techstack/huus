package engine

import (
  "errors"
  "fmt"
  "strings"
  // "log"
)

const blockSize = 4096

type valueOffset int

// struct to hold kv pair
type KV struct {
  key int
  valOff valueOffset
}

// struct for each node of the b+tree
// TODO: instead of nextLeaf, make the last element of <Node>.children as the nextLeaf
type Node struct{
  kvStore []KV
  children []*Node
  isLeaf bool
  nextLeaf *Node 
}

// B+ tree struct
type BPtree struct {
  root *Node
  minNode int
}


/* EXTERNAL FUNCTIONS */
// get a value of the key, error if key does not exist
// func (tree *BPtree) Get(key int) (valueOffset, error){} 

// insert a key and value return error if it fails
func (tree *BPtree) Insert(key int, value []byte) (error) {
  if tree.root == nil {
    return errors.New("Empty tree")
  }
  // if value greater than blockSize return error
  if len(value) > blockSize {
    return errors.New("Value greater than blockSize")
  }

  root := tree.root

  // preemptively split root if it is full (makes it such that root is never full)
  if len(root.kvStore) == 2*tree.minNode - 1{
    // split root
    // make newNode, this will be the new root.
    // add current root as its child and then split it.
    newNode := &Node{
      kvStore : []KV{},
      children : []*Node{root},
      isLeaf : false,
      nextLeaf : nil,
    }

    err := tree.splitChild(newNode, 0)
    if err != nil{
      return err
    }
    tree.root = newNode
  }

  // insert the key into the store
  // replace _ with valOff here
  _, err := tree.insertNonFull(tree.root, key)
  if err != nil {
    return err
  }
  /*
  else {
    if err := writeVal(valOff, value); err != nil{
      // BUG: if write fails, the inserted key needs to be deleted or write needs to happen again.
      return err
    }
  }
  */

  return nil
}


// delete a key and value return error if key does not exist or if deletion fails
func (tree BPtree) Delete(key int) (error) {
  deletekey.Delete(key)
}


/* HELPER FUNCTIONS */

// function to initiate a new Tree

func CreateNewTree(minNode int) (*BPtree) {
  return &BPtree{
    root : &Node{
            kvStore : []KV{},
            children : []*Node{},
            isLeaf : true,
            nextLeaf : nil,
          },
    minNode : minNode,
  }
}

// function to write value into value store
// func writeVal(valOff valueOffset, value []byte) (error) {}

// function to read value from value store
// func readVal (valOff valueOffset) ([]byte, error) {}


// function to split a node in a b plus tree
func (tree *BPtree) splitChild(parent *Node, index int) (error) {
  child := parent.children[index]
  minNode := tree.minNode

  newNode := &Node{
    kvStore : []KV{},
    children : []*Node{},
    isLeaf : child.isLeaf,
    nextLeaf : nil,
  }

  if child.isLeaf {
    newNode.kvStore = append(newNode.kvStore, child.kvStore[minNode-1:]...)
    child.kvStore = child.kvStore[:minNode-1]
    
    newNode.nextLeaf = child.nextLeaf
    child.nextLeaf = newNode
    
    parent.kvStore = append(parent.kvStore, KV{})
    copy(parent.kvStore[index+1:], parent.kvStore[index:])
    parent.kvStore[index] = newNode.kvStore[0]
  } else {
    midKey := child.kvStore[minNode-1]
    newNode.kvStore = child.kvStore[minNode:]
    newNode.children = child.children[minNode:]
    child.kvStore = child.kvStore[:minNode-1]
    child.children = child.children[:minNode]
    
    parent.kvStore = append(parent.kvStore, KV{})
    copy(parent.kvStore[index+1:], parent.kvStore[index:])
    parent.kvStore[index] = midKey
  }
  
  parent.children = append(parent.children, nil)
  copy(parent.children[index+2:], parent.children[index+1:])
  parent.children[index+1] = newNode

  return nil
}


// recursive function to insert into a node
func (tree *BPtree) insertNonFull(node *Node, key int) (valueOffset, error) {
  // find node to insert into
  // if current node is full, split it before moving onto the child
  // if it is leaf insert into the leaf
	if node.isLeaf {
		i := 0
		// Find the position for the key in the leaf node
		for i < len(node.kvStore) && key > node.kvStore[i].key {
			i++
		}

		// If the key already exists, return an error
		if i < len(node.kvStore) && key == node.kvStore[i].key {
			return -1, errors.New("Key already exists")
		}

    kvPair := KV{
      key : key,
      valOff : -1,
    }

		// Insert the key into the correct position in the node
		node.kvStore = append(node.kvStore[:i], append([]KV{kvPair}, node.kvStore[i:]...)...)
		return -1, nil
	}

	// For internal nodes, find the correct child to insert the key
	i := 0
	for i < len(node.kvStore) && key > node.kvStore[i].key {
		i++
	}

	// If the child is full, split it
	if len(node.children[i].kvStore) == 2*tree.minNode-1 {
		tree.splitChild(node, i)
		// After splitting, if the key is greater than the promoted key, move to the next child
		if key > node.kvStore[i].key {
			i++
		}
	}

	// Recursively call the function to fill a child
  _, err := tree.insertNonFull(node.children[i], key)
	return -1, err
}

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
