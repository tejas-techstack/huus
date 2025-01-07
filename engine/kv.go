package engine

import (
  "errors"
  "fmt"
)

const (
  blockSize = 4096
  tombstone = "tombstone"
)

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
// returns nil, -1, -1 and an error if error occurs
func (tree *BPtree) Get(key int) (*Node, int, valueOffset, error){
  if tree.root == nil {
    return nil,-1, -1, errors.New("Tree is empty")
  }

  temp := tree.root
  for !temp.isLeaf {
    i := len(temp.kvStore) - 1
    for i >= 0 && key < temp.kvStore[i].key{
      i--
    }
    i++
    temp = temp.children[i]
  }

  if len(temp.kvStore) == 0{
    return nil, -1, -1, errors.New("Empty node (deletion needs to be implemented properly)")
  }

  if key < temp.kvStore[0].key || key > temp.kvStore[len(temp.kvStore)-1].key {
    return nil,-1, -1, errors.New("Key does not exist.")
  } else {
    // HACK: can be improved using binary search
    for i:=0 ;i < len(temp.kvStore); i++{
      if key == temp.kvStore[i].key{
        return temp, temp.kvStore[i].key, temp.kvStore[i].valOff, nil
      }
    }

    return nil,-1, -1, errors.New("Key does not exist.")
  }
  

  return nil, -1 , -1, errors.New("Unknown error")
} 

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
  valOff, err := tree.insertNonFull(tree.root, key)
  if err != nil {
    return err
  } else {
    err = WriteVal(value, valOff)
    if err != nil{
      fmt.Println("Error writing value :", err)
      err = tree.Delete(key)
      return err
    }
  }

  return nil
}


// delete a key and value return error if key does not exist or if deletion fails
func (tree BPtree) Delete(key int) (error) {
  // as of now, if deletion occurs
  // we only delete the key in the leaf nodes, (ignore the properties of the tree)
  // we donot delete any of the parents and ignore the minimum node requirements in the leaf nodes.
  // TODO: delete duplicates and merge nodes.
  node,index, valOff, err := tree.Get(key)
  if err != nil {
    return err
  }

  node.kvStore = append(node.kvStore[:index], node.kvStore[index+1:]...)
  ts := []byte(tombstone)
  WriteVal(ts, valOff)
  return nil
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


// function to read value from value store
// func readVal (valOff valueOffset) ([]byte, error) {}


// function to split a node in a b plus tree
func (tree *BPtree) splitChild(parent *Node, index int) error {
  child := parent.children[index]
  minNode := tree.minNode

  newNode := &Node{
    kvStore: []KV{},
    children: []*Node{},
    isLeaf: child.isLeaf,
    nextLeaf: nil,
  }

  if child.isLeaf {
    splitPoint := minNode - 1
    newNode.kvStore = append(newNode.kvStore, child.kvStore[splitPoint:]...)
    child.kvStore = child.kvStore[:splitPoint]
    
    newNode.nextLeaf = child.nextLeaf
    child.nextLeaf = newNode
    
    parent.kvStore = append(parent.kvStore, KV{})
    copy(parent.kvStore[index+1:], parent.kvStore[index:])
    parent.kvStore[index] = newNode.kvStore[0]
  } else {
    splitPoint := minNode - 1
    midKey := child.kvStore[splitPoint]
    newNode.kvStore = append(newNode.kvStore, child.kvStore[splitPoint+1:]...)
    newNode.children = append(newNode.children, child.children[splitPoint+1:]...)
    child.kvStore = child.kvStore[:splitPoint]
    child.children = child.children[:splitPoint+1]
    
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

		i := len(node.kvStore) - 1
		// Find the position for the key in the leaf node
		for i >= 0 && key < node.kvStore[i].key {
			i--
		}

		// If the key already exists, return an error
		if i >= 0 && key == node.kvStore[i].key {
			return node.kvStore[i].valOff, errors.New("Key already exists")
		}

    kvPair := KV{
      key : key,
      valOff : -1,
    }

    i++
		// Insert the key into the correct position in the node
		node.kvStore = append(node.kvStore[:i], append([]KV{kvPair}, node.kvStore[i:]...)...)
    
    valOff,err := GetFreeBlock()
    if err != nil {
      return -1, errors.New("Could not allocate memory.")
    }
    node.kvStore[i].valOff = valOff

		return valOff, nil
	}

	// For internal nodes, find the correct child to insert the key
	i := len(node.kvStore) - 1
	for i >= 0  && key < node.kvStore[i].key {
		i--
	}

  i++

	// If the child is full, split it
	if len(node.children[i].kvStore) == 2*tree.minNode-1 {
		tree.splitChild(node, i)
		// After splitting, if the key is greater than the promoted key, move to the next child
		if key > node.kvStore[i].key {
			i++
		}
	}

	// Recursively call the function to fill a child
  valOff, err := tree.insertNonFull(node.children[i], key)
	return valOff, err
}


