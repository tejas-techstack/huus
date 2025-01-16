package kv

import (
  "os"
  "fmt"
  "math"
  "bytes"
)

// set default order of tree.
const defaultOrder = 500

// pager constants
const maxPageSize = 4096 //change to math.MaxUint16 later on.
const minPageSize = 32

const maxOrder = math.MaxUint16

type BPTree struct {
  order uint16

  storage *storage

  metadata *treeMetaData

  // minKey = ceil(order/2) - 1
  minKeyNum int
}

type treeMetaData struct {
  order uint16
  rootId uint32
  pageSize uint16
}

// Open either opens a new tree or loads a pre existing tree.
func Open(path string) (*BPTree, error) {
  // replace defaultOrder with user selected order
  // replace pageSize with user selected pageSize

  // use to set page size
  pageSize := os.Getpagesize()

  storage, err := newStorage(path, pageSize)
  if err != nil {
    return nil, fmt.Errorf("failed to init the storage: %w", err)
  }

  // loads metadata from file header.
  // header has a fixed size and cannot be modified.
  metadata, err := storage.loadMetadata()
  if err != nil {
    return nil, fmt.Errorf("failed to init the metadata: %w", err)
  } 

  if metadata != nil && metadata.order != defaultOrder {
    return nil, fmt.Errorf("Tried to open a tree with order %w, but has order %w", metadata.order, defaultOrder)
  }

  minKeyNum := calcMinOrder(defaultOrder)

  return &BPTree{order : defaultOrder, storage : storage, metadata : metadata, minKeyNum : minKeyNum}, nil
}

type node struct {
  id uint32

  parentId uint32

  key [][]byte
  
  // pointer can either be a value or children
  // based on if the node is a leaf or not
  pointers []*pointer
  
  // leaf vs non leaf.
  isLeaf bool
  sibling uint32
}

type pointer struct {
  value interface{}
}

func (p *pointer) isValue() bool {
  _, ok := p.value.([]byte)
  return ok
}

func (p *pointer) isNodeId() bool {
  _, ok := p.value.(uint32)
  return ok
}

func (p *pointer) asValue() []byte {
  return p.value.([]byte)
}

func (p *pointer) asNodeId() uint32 {
  return p.value.(uint32)
}

// returns (value, err)
func (t *BPTree) Get(key []byte) ([]byte, error) {
  if t.metadata == nil {
    return nil, fmt.Errorf("Not initialized")
  }

  leaf, err := t.findLeaf(key)
  if err != nil {
    return nil, fmt.Errorf("Could not find leaf : %w", err)
  }

  for i := 0; i < len(cur.key); i++ {
    if compare(cur.key[i], key) == 0{
      return leaf.pointers[index].asValue(), nil
    }
  }

  return nil, fmt.Errorf("Key not found error")
}


func (t *BPTree) Put(key, value []byte) (error) {
  if t.metadata == nil {
    return fmt.Errorf("Tree not initialized")
  }

  if len(value) > maxPageSize {
    return fmt.Errorf("value greater than pageSize")
  }

  // the leaf returned here is pre processed to always have space to insert.
  leaf, err := t.findLeafToInsert(key)
  if err != nil {
    return fmt.Errorf("Put failed : %w", err)
  }

  err := t.insertIntoNode(leaf,key, pointer(value))
  if err != nil {
    return fmt.Errorf("Failed to insert into node : %w", err)
  } else {
    return nil
  }

  return fmt.Errorf("Unknow error while executing Put")
}


// insert (key, child)
func (t *BPTree) insertIntoNode(cur *node,key []byte, pointer *pointer) error {

  if len(cur.key) == t.order {
    return fmt.Errorf("Node is full cannot insert into node.")
  }

  //TODO find index to insert at using key.
  index := 0
  for index < len(cur.key) {
    if compare(key, cur.key[i]) < 0 {
      break
    }
    index++
  }

  if pointer.isValue() == 0 {
    err := insertValueAt(cur, index, pointer.asValue())
    if err != nil {
      return fmt.Errorf("Inserting value failed : %w", err)
    }
  }

  if pointer.isNodeId() == 0 {
    err := insertNodeAt(cur, index, pointer.asNodeId())
    if err != nil {
      return fmt.Errorf("Inserting node failed : %w", err)
    }
  }

  return fmt.Errorf("Unknown error while inserting into node.")
}

func (t *BPTree) findLeaf(key []byte) (*node, error) {
  
  /* 
  
    does not care if next node is full or not.
    returns error if key is not in range.

  */

  node, err := t.storage.loadNode(t.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error loading root : %w", err)
  }

  for !node.isLeaf {
    //TODO find index of child.
    index, err := t.findChildIndex(node, key)
    if err != nil {
      return nil, fmt.Errorf("Could not find child index")
    }

    childId := node.pointers[index+1].asNodeId
    node := t.storage.loadNode(childId)
  }

  return node, nil
}


func (t *BPTree) findLeafToInsert(key []byte) *node, error {

  // load root node, if it has t.Order - 1
  // split the root.
  // load child, split if it is full and insert seperator
  // into the parent.
  // keep both parent and child nodes in memory to make adding seperator easier.
  // keep switching the pairs.
  
  root,err := t.storage.loadNode(t.rootId)
  if err != nil {
    return fmt.Errorf("Error inserting into leaf : %w", err)
  }

  if len(root.key) == t.order - 1{
    err := t.splitRoot()
    if err != nil {
      return fmt.Errorf("Error splitting root : %w", err)
    }
  }

  parent := root

  // find child returns a loaded child node and error if any.
  child, err := findChild(parent, key)
  if err != nil {
    return fmt.Errorf("Error finding child : %w", err)
  }

  for !child.isLeaf{
    parent = child;
    child, err = findChild(parent, key)

    if len(child.key) == t.order - 1 {
      err := t.splitNode(parent, child)
      if err != nil {
        return fmt.Errorf("Error splitting child : %w", err)
      }
    }
  }

  if child.isLeaf {
    return child, nil
  }

  return nil, fmt.Errorf("Unknown error while executing findLeafToInsert")
}


func (t *BPTree) findChildIndex(cur *node, key []byte) (int, error) {
  
  index := 0
  for index < len(cur.key) {
    if compare(key, cur.key[index]) < 0 {
      break
    }
    index++
  }

  return index, nil

}

func (t *BPTree) findChild(parent *node, key []byte) (*node, error) {
  childIndex, err := findChildIndex(parent, key)
  if err != nil {
    return err
  }

  childId := parent.pointers[childIndex].asNodeId()

  child,err := t.storage.loadNode(childId)
  if err != nil {
    return nil, fmt.Errorf("error loading child : %w", err)
  }

  return child, nil
}

func (t *BPTree) splitRoot() error {
  // change metadata of tree.

  newRootId, err := t.storage.newNode()
  if err != nil {
    return fmt.Errorf("Error loading new node : %w",err)
  }

  newRoot := &node {
    id : newRootId,
    parentId : 0,
    key : []byte{},
    pointers : []*pointer{},
    isLeaf : false,
    sibling : nil,
  }

  curRoot, err := t.storage.loadNode(t.rootId)
  if err != nil {
    return fmt.Errorf("Error loading root node : %w", err)
  }

  // change properties of newRoot and current Root.
  
  err := splitNode(curRoot, newRoot)
  if err != nil {
    // revert changes to current root
    curRoot.parentId = 0
    return fmt.Errorf("Error splitting node.")
  }

  // update metaData of tree with rootId as newRootId.
  err := t.updateMetaData(newRoot.id)
  if err != nil {
    return fmt.Errorf("Failed to update metadata : %w", err)
  } else {
    return nil
  }

  return fmt.Errorf("Unknown error in splitting root.")
}

func (t *BPTree) splitNode(cur *node, parent *node) error {

  // splitNode assumes cur and parent nodes have space present already.

  // create new node.
  // copy contents of cur node to new node
  // clear contents of cur node that were copied.
  // find seperator, if leaf, copy up, if not leaf promote up.
  // put seperator in parent.
  // write both parent and cur nodes.

  if parent == nil {
    parent, err := t.storage.loadNode(cur.parentId)
    if err != nil {
      return fmt.Errorf("Failed to load parent : %w",err)
    }
  }

  newNodeId, err := t.storage.newNode()
  if err != nil {
    return fmt.Errorf("Error creating newNode : %w", err)
  }

  newNode := &node {
    id : newRootId,
    parentId : parent.id,
    key : []byte{},
    pointers : []*pointer{},
    isLeaf : cur.isLeaf,
    sibling : 0,
  }

  seperator := cur.key[t.minKeyNum]

  if newNode.isLeaf {
    // update sibling 
    newNode.sibling = cur.sibling
    cur.sibling = newNode.id

    // update keys.
    newNode.key = append(newNode.key, child.key[t.minKeyNum : ]...)
    child.key = child.key[:t.minKeyNum]

    // update values.
    newNode.pointers = append(newNode.pointers, child.pointers[t.minKeyNum : ]...)
    child.pointers = child.pointers[:t.minKeyNum]
  } else {
    // node is not a leaf.
    // update child node.
    newNode.key = append(newNode.key, child.key[t.minKeyNum + 1: ]...)
    child.key = child.key[:t.minKeyNum]

    // update child pointers
    newNode.pointers = append(newNode.pointers, child.pointers[t.minKeyNum :]...)
    child.pointers = child.pointesr[:t.minKeyNum]
  }

  // have to insert newNode to parent as well innit.

  // move seperator to parent.
  index, err := findChildIndex(parent)
  if err != nil {
    return fmt.Errorf("Error finding child index : %w", err)
  }

  err := insertNodeAt(parent, Index, seperator)
  if err != nil {
    return fmt.Errorf("Error inserting sepeartor into parent")
  }

}

func (t *BPTree) insertValueAt(cur *node, index int, value []byte) error {
  if len(cur.key) == t.order - 1 {
    return fmt.Errorf("Cannot insert value, node is full")
  }

  // write the value to storage.
}

func (t *BPTree) insertNodeAt(cur *node, index int, child uint32) error {
  if len(cur.key) == t.order - 1 {
    return fmt.Errorf("Cannot insert value, node is full")
  }

  // write the node to stoarge.
}

func compare(byteA ,byteB []byte) int {
  return bytes.Compare(byteA, byteB)
}

func calcMinOrder(order uint16) {}
