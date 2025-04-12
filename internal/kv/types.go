package kv

import (
  "math"
)

const defaultOrder = 500

const maxPageSize = 4096
const minPageSize = 32

const maxOrder = math.MaxUint16

type BPTree struct {
  order uint16

  storage *storage

  metadata *treeMetaData

  minKeyNum int
}

type treeMetaData struct {
  order uint16
  rootId uint32
  pageSize uint16
}

type node struct {
  id uint32

  parentId uint32

  key [][]byte
  
  // pointer can either be a value or child Node Id.
  pointers []*pointer
  
  isLeaf bool

  // each level of the tree are connected in a linked list type of order.
  sibling uint32
}

type pointer struct {
  // value can either be a value of the kv pair.
  // or hold the value of the child node id.
  value interface{}
}



