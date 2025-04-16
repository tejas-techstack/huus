/*
* Contain helper functions for insert functionality.
*/

package kv

import (
  "fmt"
)

func (t *BPTree) insertIntoNode(cur *node,key []byte, pointer *pointer) error {

  if len(cur.key) == int(t.order) {
    return fmt.Errorf("Node is full cannot insert into node.")
  }

  index := 0
  for index < len(cur.key) {
    if compare(key, cur.key[index]) < 0 {
      break
    }
    index++
  }

  err := t.insertKeyAt(cur, index, key)
  if err != nil {
    return fmt.Errorf("Error inserting key : %w", err)
  }

  if pointer.isValue() {
    err := t.insertValueAt(cur, index, pointer.asValue())
    if err != nil {
      return fmt.Errorf("Inserting value failed : %w", err)
    }
  }

  if pointer.isNodeId() {
    err := t.insertNodeAt(cur, index, pointer.asNodeId())
    if err != nil {
      return fmt.Errorf("Inserting node failed : %w", err)
    }
  }

  return nil
}


func (t *BPTree) findLeafToInsert(key []byte) (*node, error) {
  
  root,err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error inserting into leaf : %w", err)
  }

  if len(root.key) == int(t.order){
    root, err = t.splitRoot()
    if err != nil {
      return nil, fmt.Errorf("Error splitting root : %w", err)
    }
  }

  if root.isLeaf{
    return root, nil
  }

  parent := root

  // find child returns a loaded child node and error if any.
  child, err := t.findChild(parent, key)
  if err != nil {
    return nil, fmt.Errorf("Error finding child : %w", err)
  }

  if len(child.key) >= int(t.order) - 1 {
    err := t.splitNode(child, parent)
    if err != nil {
      return nil, fmt.Errorf("Error splitting child in findLeafToInsert: %w", err)
    }
  }

  for !child.isLeaf{
    parent = child;
    child, err = t.findChild(parent, key)

    if len(child.key) >= int(t.order) - 1{
      err := t.splitNode(child, parent)
      if err != nil {
        return nil, fmt.Errorf("Error splitting child in findLeafToInsert: %w", err)
      }
    }

    // This should be replaced with functionallity that directly checks if the child is going to be
    // the new child or if it is going to be the old child itself after the split.
    // The below line does that but searches from the root which may not be required.
    child, err = t.findChild(parent,key)
  }

  child, err = t.findLeaf(key)
  if err != nil {
    return nil, fmt.Errorf("Error finding leaf : %w", err)
  }

  return child, nil
}


func (t *BPTree) insertKeyAt(cur *node, index int, key []byte) error {
  if len(cur.key) == int(t.order){
    return fmt.Errorf("Cannot insert into node as it is full.")
  }

  cur.key = append(cur.key, []byte{0})
  copy(cur.key[index+1:], cur.key[index:])
  cur.key[index] = key

  err := t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  return nil
}

func (t *BPTree) insertValueAt(cur *node, index int, value []byte) error {
  if len(cur.pointers) == int(t.order) {
    return fmt.Errorf("Cannot insert value, node is full")
  }

  cur.pointers = append(cur.pointers, &pointer{0})
  copy(cur.pointers[index+1:], cur.pointers[index:])
  cur.pointers[index] = &pointer{value}

  err := t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  return nil
}

func (t *BPTree) insertNodeAt(cur *node, index int, child uint32) error {
  if len(cur.pointers) == int(t.order)+1 {
    return fmt.Errorf("Cannot insert node, node is full")
  }

  cur.pointers = append(cur.pointers, &pointer{0})
  copy(cur.pointers[index+1:], cur.pointers[index:])
  cur.pointers[index] = &pointer{child}

  err := t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  return nil
}
