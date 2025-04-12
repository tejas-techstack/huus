/*
* Generic helper functions that can be used everywhere.
*/

package kv

import (
  "fmt"
  "bytes"
)

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

func compare(byteA ,byteB []byte) int {
  return bytes.Compare(byteA, byteB)
}

func calcMinOrder(order uint16) int {
  // minKeyNum is given by ceil(order/2) - 1
  d := (order / 2)
	if order%2 == 0 {
		return int(d - 1)
	}

	return int(d)
}

func (t *BPTree) findLeaf(key []byte) (*node, error) {
  // does not care if next node is full or not.
  // returns error if key is not in range.
  cur, err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error loading root : %w", err)
  }

  for !cur.isLeaf {
    index := 0
    for index < len(cur.key) {
      if compare(key, cur.key[index]) < 0 {
        break
      }
      index++
    }

    childId := cur.pointers[index].asNodeId()
    cur, err = t.storage.loadNode(childId)
    if err != nil {
      return nil, fmt.Errorf("Error loading child : %w", err)
    }
  }

  return cur, nil
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
  childIndex, err := t.findChildIndex(parent, key)
  if err != nil {
    return nil, fmt.Errorf("error finding child index")
  }

  childId := parent.pointers[childIndex].asNodeId()

  child,err := t.storage.loadNode(childId)
  if err != nil {
    return nil, fmt.Errorf("error loading child : %w", err)
  }

  return child, nil
}
