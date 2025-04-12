/*
* Helpers for the delete functionality.
* Contains functions that help to rebalance the tree.
*/

package kv

import (
  "fmt"
)

func (t *BPTree) removeKeyAtLeaf(cur *node, key []byte) error {

  if len(cur.key) == t.minKeyNum && cur.id != t.metadata.rootId{
    for i, v := range cur.key {
      if compare(v, key) == 0 {
        cur.key = append(cur.key[:i], cur.key[i+1:]...)
        cur.pointers = append(cur.pointers[:i], cur.pointers[i+1:]...)
      }
    }

    if cur.sibling == uint32(0) {
      // this means we need to handle left hand side only.
      parent, err := t.storage.loadNode(cur.parentId)
      if err != nil{
        return fmt.Errorf("Error loading parent : %w", err)
      }

      leftSibId := parent.pointers[len(parent.pointers) - 2].asNodeId()
      leftSib, err := t.storage.loadNode(leftSibId)
      if err != nil {
        return fmt.Errorf("Error loading left sibling : %w", err)
      }

      if len(leftSib.key) <= t.minKeyNum + 1 {
        // merge
        if err := t.mergeNode(cur); err != nil {
          return fmt.Errorf("Error merging from Node : %w", err)
        }

      } else {
        // borrow 
        err := t.borrowFromLeft(cur)
        if err != nil {
          return fmt.Errorf("Error borrowing key from left : %w", err)
        }
      }

      // left hand side was handled.
      return nil
    }
    // end if to check if we need to handle left edge case
    // borrow/ merge only from right.

    sibling, err := t.storage.loadNode(cur.sibling)
    if err != nil {
      return fmt.Errorf("Error loading sibling : %w", err)
    }

    if len(sibling.key) < t.minKeyNum + 1 {
      // merge.
      err = t.mergeNode(cur)
      if err != nil {
        return fmt.Errorf("Error merging nodes : %w", err)
      }
    } else {
      // else borrow from right
      if err := t.borrowFromRight(cur); err != nil {
        return fmt.Errorf("Error borrowing key from right : %w", err)
      }
    }

    return nil
  }

  // no need to merge or borrow.
  for i, v := range cur.key {
    if compare(v, key) == 0 {
      cur.key = append(cur.key[:i], cur.key[i+1:]...)
      cur.pointers = append(cur.pointers[:i], cur.pointers[i+1:]...)
    }
  }

  err := t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  return nil
}


func (t *BPTree) mergeNode(cur *node) error {

  if cur.sibling == uint32(0) {
    if err := t.mergeLeftNode(cur); err != nil {
      return fmt.Errorf("Error merging from left node : %w", err)
    }
  } else {
    if err := t.mergeRight(cur); err != nil {
      return fmt.Errorf("Error merging from the right node : %w", err)
    }
  }

  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil {
    return fmt.Errorf("error loading parent : %w", err)
  }

  for parent.id != t.metadata.rootId {
    if len(parent.key) < t.minKeyNum{

      if parent.sibling == uint32(0) {
        // handle merge/ borrow of left node.
        grandparent, err := t.storage.loadNode(parent.parentId)
        if err != nil{
          return fmt.Errorf("Error loading parent : %w", err)
        }

        leftSibId := grandparent.pointers[len(grandparent.pointers) - 2].asNodeId()
        leftSib, err := t.storage.loadNode(leftSibId)
        if err != nil {
          return fmt.Errorf("Error loading left sibling : %w", err)
        }

        if len(leftSib.key) <= t.minKeyNum + 1 {
          // merge
          if err := t.mergeLeftNode(parent); err != nil {
            return fmt.Errorf("Error merging from Node : %w", err)
          }

        } else {
          // borrow 
          err := t.borrowFromLeft(cur)
          if err != nil {
            return fmt.Errorf("Error borrowing key from left : %w", err)
          }

          return nil
        }
      } else {
        // handle merge/ borrow from the right.
        sibling, err := t.storage.loadNode(parent.sibling)
        if err != nil {
          return fmt.Errorf("Error loading sibling : %w", err)
        }

        if len(sibling.key) <= t.minKeyNum + 1 {
          if err := t.mergeRight(parent); err != nil {
            return fmt.Errorf("Error merging right : %w", err)
          }
        } else {
          if err := t.borrowFromRight(parent); err != nil {
            return fmt.Errorf("Error borrowing from right : %w", err)
          }

          return nil
        }
      }
  
      grandparent, err := t.storage.loadNode(parent.parentId)
      if err != nil {
        return fmt.Errorf("Error loading grandparent : %w", err)
      }

      // set parent as grandparent for next loop iteration
      parent = grandparent
    } else {
      return nil
  }
  // end for loop.
  }

  // parent is root and we need to check if height should be reduced.
  if len(parent.key) == 0 {
    newRootId := parent.pointers[0].asNodeId()
  
    t.metadata.rootId = newRootId
    if err := t.storage.updateMetadata(t.metadata); err != nil {
      return fmt.Errorf("Error updating metadata : %w", err)
    }
  }

  return nil
}

func (t *BPTree) borrowFromLeft(cur *node) error {

  if cur.sibling != uint32(0) {
    return fmt.Errorf("Should not be borrowing from left, right sibling is not nil.")
  }

  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil{
    return fmt.Errorf("Error loading parent : %w", err)
  }


  leftSibId := parent.pointers[len(parent.pointers) - 2].asNodeId()
  leftSib, err := t.storage.loadNode(leftSibId)
  if err != nil {
    return fmt.Errorf("Error loading left sibling : %w", err)
  }

  if len(leftSib.key) < t.minKeyNum + 1 {
    if err := t.mergeNode(cur); err != nil {
      return fmt.Errorf("Error merging node : %w", err)
    }
  }

  index := len(parent.key) - 1
  lastElem := len(leftSib.key) - 1

  if !leftSib.isLeaf{
    child, err := t.storage.loadNode(leftSib.pointers[lastElem+1].asNodeId())
    if err != nil {
      return fmt.Errorf("Error loading child : %w", err)
    }
    child.parentId = cur.id
    if err := t.storage.updateNode(child); err != nil {
      return fmt.Errorf("Error updating child's parent id : %w", err)
    }
  }

  cur.key = append(parent.key[index:index+1], cur.key...)
  cur.pointers = append(leftSib.pointers[lastElem+1:], cur.pointers...)

  parent.key[index] = leftSib.key[lastElem]

  leftSib.key = leftSib.key[0:lastElem]
  leftSib.pointers = leftSib.pointers[0:lastElem+1]

  if cur.isLeaf {
    cur.key[0] = parent.key[index]
  }

  if err := t.storage.updateNode(parent); err != nil {
    return fmt.Errorf("Error updating parent : %w", err)
  }

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating current node : %w", err)
  }

  if err := t.storage.updateNode(leftSib); err != nil {
    return fmt.Errorf("Error updating sibling : %w", err)
  }

  return nil
}

func (t *BPTree) borrowFromRight(cur *node) error {

  if cur.sibling == uint32(0) {
    return fmt.Errorf("Right sibling is nil.")
  }

  sibling, err := t.storage.loadNode(cur.sibling)
  if err != nil {
    return fmt.Errorf("Error loading sibling : %w", err)
  }



  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil {
    return fmt.Errorf("Error loading parent : %w", err)
  }

  index := 0
  // find the postion in parent.
  for i, v := range parent.pointers {
    if v.asNodeId() == cur.id{
      // found index.
      index = i
    }
  }

  if !sibling.isLeaf{
    child, err := t.storage.loadNode(sibling.pointers[0].asNodeId())
    if err != nil {
      return fmt.Errorf("Error loading child : %w", err)
    }
    child.parentId = cur.parentId
    if err := t.storage.updateNode(child); err != nil {
      return fmt.Errorf("Error updating child's parent id : %w", err)
    }
  }

  cur.key = append(cur.key, parent.key[index])
  cur.pointers = append(cur.pointers, sibling.pointers[0])

  parent.key[index] = sibling.key[0]

  sibling.key = sibling.key[1:]
  sibling.pointers = sibling.pointers[1:]

  if cur.isLeaf {
    parent.key[index] = sibling.key[0]
  }

  if err := t.storage.updateNode(parent); err != nil {
    return fmt.Errorf("Error updating parent : %w", err)
  }

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating current node : %w", err)
  }

  if err := t.storage.updateNode(sibling); err != nil {
    return fmt.Errorf("Error updating sibling : %w", err)
  }

  return nil
}

// merges left sibling to current node.
func (t *BPTree) mergeLeftNode(cur *node) error {

  if cur.sibling != uint32(0) {
    return fmt.Errorf("Should not be merging from left, right sibling is not nil.")
  }

  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil{
    return fmt.Errorf("Error loading parent : %w", err)
  }

  leftSibId := parent.pointers[len(parent.pointers) - 2].asNodeId()
  leftSib, err := t.storage.loadNode(leftSibId)
  if err != nil {
    return fmt.Errorf("Error loading left sibling : %w", err)
  }

  if !cur.isLeaf{
    cur.key = append(parent.key[len(parent.key)-1:], cur.key...)
  }

  if !cur.isLeaf{
    for _, v := range leftSib.pointers {
      child, err := t.storage.loadNode(v.asNodeId())
      if err != nil {
        return fmt.Errorf("Error loading child : %w", err)
      }

      child.parentId = cur.id
      if err := t.storage.updateNode(child); err != nil {
        return fmt.Errorf("Error updating Child : %w", err)
      }
    }
  }

  cur.key = append(leftSib.key, cur.key...)
  cur.pointers = append(leftSib.pointers, cur.pointers...)

  parent.key = parent.key[:len(parent.key)-1]
  parent.pointers = append(parent.pointers[:len(parent.pointers)-2], parent.pointers[len(parent.pointers)-1])

  if err := t.storage.deleteNode(leftSib.id); err != nil {
    return fmt.Errorf("Error deleting node : %w", err)
  }

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  if err := t.storage.updateNode(parent); err != nil {
    return fmt.Errorf("Error updating parent : %w", err)
  }

  return nil
}


// merges right sibling to current node.
func (t *BPTree) mergeRight(cur *node) error {

  if cur.sibling == uint32(0) {
    return fmt.Errorf("Right sibling is nil.")
  }

  sibling, err := t.storage.loadNode(cur.sibling)
  if err != nil {
    return fmt.Errorf("Error loading sibling : %w", err)
  }


  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil {
    return fmt.Errorf("Error loading parent : %w", err)
  }


  index := 0
  // find the postion in parent.
  for i, v := range parent.pointers {
    if v.asNodeId() == cur.id{
      // found index.
      index = i
    }
  }

  if !cur.isLeaf {
    cur.key = append(cur.key, parent.key[index])
  }

  cur.key = append(cur.key, sibling.key...)
  cur.pointers = append(cur.pointers, sibling.pointers...)
  cur.sibling = sibling.sibling

  if !cur.isLeaf {
    for _, v := range cur.pointers {
      child, err := t.storage.loadNode(v.asNodeId())
      if err != nil {
        return fmt.Errorf("Error loading child : %w", err)
      }

      child.parentId = cur.id 
      if err := t.storage.updateNode(child); err != nil {
        return fmt.Errorf("Error updating child : %w", err)
      }
    }
  }

  if err = t.storage.deleteNode(sibling.id); err != nil {
    return fmt.Errorf("Error deleting node : %w", err)
  }

  parent.key = append(parent.key[:index], parent.key[index+1:]...)
  parent.pointers = append(parent.pointers[:index+1], parent.pointers[index+2:]...)

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating current : %w", err)
  }

  if err := t.storage.updateNode(parent); err != nil {
    return fmt.Errorf("Error updating parent : %w", err)
  }

  return nil
}


