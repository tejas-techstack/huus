/*
* Helper functions that control node splits.
*/

package engine

import (
  "fmt"
)

func (t *BPTree) splitRoot() (*node, error) {
  // change metadata of tree.
  newRootId, err := t.storage.newNode()
  if err != nil {
    return nil, fmt.Errorf("Error loading new node : %w",err)
  }

  newRoot := &node {
    id : newRootId,
    key : [][]byte{},
    pointers : []*pointer{&pointer{t.metadata.rootId}},
    isLeaf : false,
    sibling : 0,
  }

  curRoot, err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error loading root node : %w", err)
  }

  err = t.splitNode(curRoot, newRoot)
  if err != nil {
    return nil, fmt.Errorf("Error splitting node.")
  }

  t.metadata = &treeMetaData{
    t.order,
    newRoot.id,
    t.storage.pageSize,
  }
  err = t.storage.updateMetadata(t.metadata)
  if err != nil {
    return nil, fmt.Errorf("Failed to update metadata : %w", err)
  }

  t.stack.updateZeroPointer(newRoot.id)

  return newRoot, nil
}

func (t *BPTree) splitNode(cur *node, parent *node) error {

  newNodeId, err := t.storage.newNode()
  if err != nil {
    return fmt.Errorf("Error creating newNode : %w", err)
  }

  // update current's parent id in case of root node.
  // otherwise does not matter.

  newNode := &node {
    id : newNodeId,
    key : [][]byte{},
    pointers : []*pointer{},
    isLeaf : cur.isLeaf,
    sibling : 0,
  }

  seperator := cur.key[t.minKeyNum]

  // sibling is updated regardless of parent or child.
  newNode.sibling = cur.sibling
  cur.sibling = newNode.id

  if newNode.isLeaf {
    // update keys.
    newNode.key = append(newNode.key, cur.key[t.minKeyNum : ]...)
    cur.key = cur.key[:t.minKeyNum]

    // update values.
    newNode.pointers = append(newNode.pointers, cur.pointers[t.minKeyNum : ]...)
    cur.pointers = cur.pointers[:t.minKeyNum]

  } else {
    // node is not a leaf.
    // update child node.
    newNode.key = append(newNode.key, cur.key[t.minKeyNum + 1: ]...)
    cur.key = cur.key[:t.minKeyNum]

    // update child pointers
    newNode.pointers = append(newNode.pointers,cur.pointers[t.minKeyNum+1 :]...)
    cur.pointers = cur.pointers[:t.minKeyNum + 1]
  }

  // move seperator and newNode to parent.
  index, err := t.findChildIndex(parent, seperator)
  if err != nil {
    return fmt.Errorf("Error finding child index : %w", err)
  }

  err = t.insertKeyAt(parent, index, seperator)
  if err != nil {
    return fmt.Errorf("Error inserting seperator into parent : %w", err) 
  }

  err = t.insertNodeAt(parent,index+1, newNode.id)
  if err != nil {
    return fmt.Errorf("Error inserting newNode as a child into parent : %w", err)
  }

  err = t.storage.updateNode(newNode)
  if err != nil {
    return fmt.Errorf("Error updating newly created Node : %w", err)
  }

  err = t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating child node : %w", err)
  }

  return nil
}
