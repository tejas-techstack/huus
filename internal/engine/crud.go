/*
* This file Contains all the operations exposed to outside the library.
* Everything else is intended to be internal to the library.
*/

package engine

import (
  "fmt"
)

// Implements the Get function but takes input as int.
func (t *BPTree) GetInt(key int) ([]byte, bool, error) {
  return t.Get(encodeUint64(key))
}

// Returns (value, exists, error)
func (t *BPTree) Get(key []byte) ([]byte, bool ,error) {
  if t.metadata == nil {
    return nil, false, nil
  }

  leaf, err := t.findLeaf(key)
  if err != nil {
    return nil, false, fmt.Errorf("Could not find leaf : %w", err)
  }

  // find the value in the given leaf.
  for i := 0; i < len(leaf.key); i++ {
    if compare(leaf.key[i], key) == 0{
      return leaf.pointers[i].asValue(), true, nil
    }
  }

  return nil, false, nil 
}

// Implements the Put function but takes input as int.
func (t *BPTree) PutInt(key, value int) error {
  return t.Put(encodeUint64(key), encodeUint64(value))
}

// Returns error if any.
func (t *BPTree) Put(key, value []byte) error {

  // TODO : BAD CODE, THIS WILL CAUSE A READ EVERYSINGlE TIME A WRITE HAPPENS
  // Ideally the write should propogate to the end of the tree and there if the key is found
  // to already exist then it should say key already exists and return.
  _, exists, err := t.Get(key)
  if err != nil {
    return fmt.Errorf("Error Checking if key exists : %w", err)
  }
  if exists {
    fmt.Println("Key", decodeUint64(key), "already exists")
    return nil
  }

  if len(value) > maxPageSize {
    return fmt.Errorf("value greater than pageSize")
  }

  if t.metadata == nil {
    err := t.initializeRoot(key, value)
    if err != nil {
      return fmt.Errorf("Error initializing root : %w", err)
    }

    return nil
  }

  // the leaf returned here is pre processed to always have space to insert.
  leaf, err := t.findLeafToInsert(key)
  if err != nil {
    return fmt.Errorf("Put failed : %w", err)
  }

  err = t.insertIntoNode(leaf, key, &pointer{value})
  if err != nil {
    return fmt.Errorf("Failed to insert into node : %w", err)
  }

  return nil
}

// Implements the Delete function but takes input as int.
func (t *BPTree) DeleteInt(key int) (bool, error) {
  return t.Delete(encodeUint64(key))
}

// Delete returns true, nil if deletion was successful
// returns false, nil if key did not exist
func (t *BPTree) Delete(key []byte) (bool, error) {

  _, exists, err := t.Get(key)
  if err != nil {
    return false,fmt.Errorf("Error getting key %w", err)
  }
  if !exists {
    return true, nil
  }

  leaf, err := t.findLeafToDelete(key)
  if err != nil {
    return false,fmt.Errorf("Error searching for leaf %w", err)
  }

  err = t.removeKeyAtLeaf(leaf, key)
  if err != nil {
    return false,fmt.Errorf("Error removing key at leaf : %w", err)
  }

  // Empty stack after deletion.
  // If done right should already be empty?
  t.stack.emptyStack()

  return true, nil
}
