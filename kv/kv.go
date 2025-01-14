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
  rootID uint32
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
  // keyNum represents number of keys present in the node.
  keyNum int
  
  // pointer can either be a value or children
  // based on if the node is a leaf or not
  pointers []*pointer
  
  // leaf vs non leaf.
  isLeaf bool
  sibling [1]uint32 // holds id of next node if child, else holds nil
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

  index, err := findIndex(leaf, key)
  if err != nil {
    return nil, fmt.Errorf("Could not Get value, finding index failed : %w", err)
  } else {
    return leaf.pointers[index].asValue(), nil
  }

  return nil, fmt.Errorf("Key not found error")
}

func (t *BPTree) findLeaf(key []byte) (*node, error) {}

func (t *BPTree) Put(key, value []byte) (error) {
  if t.metadata == nil {
    return fmt.Errorf("Tree not initialized")
  }

  if len(value) > maxPageSize {
    return fmt.Errorf("value greater than pageSize")
  }

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

func (t *BPTree) findLeafToInsert(key []byte) error {
  return fmt.Errorf("Unknown error while executing findLeafToInsert")
}

// insert (key, child)
func (t *BPTree) insertIntoNode(curr *node,key []byte, pointer *pointer) error {
  // find index to insert at using key.
  index, err := findIndex(curr, key)
  if err != nil {
    return fmt.Errorf("Inserting into node failed : %w", err)
  }

  if pointer.isValue() == 0 {
    err := insertValueAt(curr, index, pointer.asValue())
    if err != nil {
      return fmt.Errorf("Inserting value failed : %w", err)
    }
  }

  if pointer.isNodeId() == 0 {
    err := insertNodeAt(curr, index, pointer.asNodeId())
    if err != nil {
      return fmt.Errorf("Inserting node failed : %w", err)
    }
  }

  return fmt.Errorf("Unknown error while inserting into node.")
}

func findIndex(curr *node, key []byte) (int, error) {}

func insertValueAt(curr *node, index int, value []byte) error {}

func insertNodeAt(curr *node, index int, child uint32) error {}


func compare(byteA ,byteB []byte) int {
  return bytes.Compare(byteA, byteB)
}

func calcMinOrder(order uint16) {}
