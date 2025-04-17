/*
* Contains the functions that are required to setup the engine.
*/

package engine

import (
  "fmt"
  "path/filepath"
)

// Open either opens a new tree or loads a pre existing tree.
func Open(filePath string, order uint16, pageSize uint16) (*BPTree, error) {
  // Load the storage.
  storage, err := newStorage(filePath, pageSize, order)
  if err != nil {
    return nil, fmt.Errorf("failed to init the storage: %w", err)
  }

  // loads metadata from file header.
  // header has a fixed size and cannot be modified.
  metadata, err := storage.loadMetadata()
  if err != nil {
    return nil, fmt.Errorf("failed to init the metadata: %w", err)
  }

  // metdata != nil takes care of the case 
  // where the tree is not yet initialized.
  if metadata != nil && metadata.order != order {
    return nil, fmt.Errorf("Tried to open a tree with order %v, but has order %v", metadata.order, defaultOrder)
  }
  minKeyNum := calcMinOrder(order)

  fmt.Println("Using tree:", filepath.Base(filePath))

  return &BPTree{order : order, storage : storage, metadata : metadata, minKeyNum : minKeyNum}, nil
}

// Initializes the root node on first insert.
func (t *BPTree) initializeRoot(key, value []byte) error {
  rootId, err := t.storage.newNode()
  if err != nil {
    return fmt.Errorf("Error creating newNode : %w", err)
  }

  root, err := t.storage.loadNode(rootId)
  if err != nil {
    return fmt.Errorf("Error reading root : %w", err)
  }

  root.isLeaf = true
  root.key = append(root.key, key)
  root.pointers = append(root.pointers, &pointer{value})

  err = t.storage.updateNode(root)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  t.metadata = &treeMetaData{
    t.order,
    rootId,
    t.storage.pageSize,
  }

  err = t.storage.updateMetadata(t.metadata)
  if err != nil {
    return fmt.Errorf("Error updating metadata : %w", err)
  }

  t.stack = InitializeStack(rootId)

  return nil
}
