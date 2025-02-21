package kv

import (
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
func Open(path string, order uint16, pageSize uint16) (*BPTree, error) {
  // replace defaultOrder with user selected order
  // replace pageSize with user selected pageSize

  // use to set page size
  // pageSizeOfSystem := os.Getpagesize()

  storage, err := newStorage(path, pageSize, order)
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

  return &BPTree{order : order, storage : storage, metadata : metadata, minKeyNum : minKeyNum}, nil
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

// returns (value, exists, err)
func (t *BPTree) Get(key []byte) ([]byte, bool ,error) {
  if t.metadata == nil {
    // succeeds even if tree doesnt exist and returns not found.
    return nil, false, nil
  }

  leaf, err := t.findLeaf(key)
  if err != nil {
    return nil, false, fmt.Errorf("Could not find leaf : %w", err)
  }



  for i := 0; i < len(leaf.key); i++ {
    if compare(leaf.key[i], key) == 0{
      return leaf.pointers[i].asValue(), true, nil
    }
  }

  return nil, false, nil 
}


func (t *BPTree) Put(key, value []byte) (error) {
  if t.metadata == nil {
    err := t.initializeRoot(key, value)
    if err != nil {
      return fmt.Errorf("Error initializing root : %w", err)
    }

    return nil
  }

  if len(value) > maxPageSize {
    return fmt.Errorf("value greater than pageSize")
  }

  // the leaf returned here is pre processed to always have space to insert.
  leaf, err := t.findLeafToInsert(key)
  if err != nil {
    return fmt.Errorf("Put failed : %w", err)
  }

  err = t.insertIntoNode(leaf,key, &pointer{value})
  if err != nil {
    return fmt.Errorf("Failed to insert into node : %w", err)
  }
  
  return nil
}

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

  return nil
}

// insert (key, child)
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

func (t *BPTree) findLeaf(key []byte) (*node, error) {
  /* 
    does not care if next node is full or not.
    returns error if key is not in range.
  */
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


func (t *BPTree) findLeafToInsert(key []byte) (*node, error) {

  // load root node, if it has t.Order
  // split the root.
  // load child, split if it is full and insert seperator
  // into the parent.
  // keep both parent and child nodes in memory to make adding seperator easier.
  // keep switching the pairs.
  
  root,err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error inserting into leaf : %w", err)
  }


  if len(root.key) == int(t.order) - 1{
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
      return nil, fmt.Errorf("Error splitting child : %w", err)
    }
  }


  for !child.isLeaf{
    parent = child;
    child, err = t.findChild(parent, key)

    if len(child.key) >= int(t.order) - 1{
      err := t.splitNode(child, parent)
      if err != nil {
        return nil, fmt.Errorf("Error splitting child : %w", err)
      }
    }

    child, err = t.findChild(parent,key)
  }

  child, err = t.findChild(parent, key)

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

func (t *BPTree) splitRoot() (*node, error) {
  // change metadata of tree.

  newRootId, err := t.storage.newNode()
  if err != nil {
    return nil, fmt.Errorf("Error loading new node : %w",err)
  }

  newRoot := &node {
    id : newRootId,
    parentId : 0,
    key : [][]byte{},
    pointers : []*pointer{&pointer{t.metadata.rootId}},
    isLeaf : false,
    sibling : 0,
  }

  curRoot, err := t.storage.loadNode(t.metadata.rootId)
  if err != nil {
    return nil, fmt.Errorf("Error loading root node : %w", err)
  }

  // change properties of newRoot and current Root.
 
  err = t.splitNode(curRoot, newRoot)
  if err != nil {
    // revert changes to current root
    curRoot.parentId = 0
    return nil, fmt.Errorf("Error splitting node.")
  }

  // update metaData of tree with rootId as newRootId.
  t.metadata = &treeMetaData{
    t.order,
    newRoot.id,
    t.storage.pageSize,
  }
  err = t.storage.updateMetadata(t.metadata)
  if err != nil {
    return nil, fmt.Errorf("Failed to update metadata : %w", err)
  }

  return newRoot, nil

}

func (t *BPTree) splitNode(cur *node, parent *node) error {
  // splitNode assumes cur and parent nodes have space present already.
  // create new node.
  // copy contents of cur node to new node
  // clear contents of cur node that were copied.
  // find seperator, if leaf, copy up, if not leaf promote up.
  // put seperator in parent.
  // write both parent and cur nodes.
  var err error
  if parent == nil {
    parent, err = t.storage.loadNode(cur.parentId)
    if err != nil {
      return fmt.Errorf("Failed to load parent : %w",err)
    }
  }

  newNodeId, err := t.storage.newNode()
  if err != nil {
    return fmt.Errorf("Error creating newNode : %w", err)
  }

  // need to update current's parent id.
  // in case it may be a root node that is being split
  cur.parentId = parent.id

  newNode := &node {
    id : newNodeId,
    parentId : parent.id,
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

func (t *BPTree) insertKeyAt(cur *node, index int, key []byte) error {
  if len(cur.key) == int(t.order){
    return fmt.Errorf("Cannot insert into node as it is full.")
  }

  // increase size of keys in cur node by 1.
  // use copy function to copy the keys into correct place.
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
  if len(cur.pointers) == int(t.order) {
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

func (t *BPTree) Delete(key []byte) (bool, error) {

  _, exists, err := t.Get(key)
  if err != nil {
    return false,fmt.Errorf("Error getting key %w", err)
  }
  if !exists {
    return true, nil
  }


  leaf, err := t.findLeaf(key)
  if err != nil {
    return false,fmt.Errorf("Error searching for leaf %w", err)
  }


  err = t.removeKeyAtLeaf(leaf, key)
  if err != nil {
    return false,fmt.Errorf("Error removing key at leaf : %w", err)
  }

  return true, nil
}

func (t *BPTree) removeKeyAtLeaf(cur *node, key []byte) error {

  if len(cur.key) == t.minKeyNum && cur.id != t.metadata.rootId{
    for i, v := range cur.key {
      if compare(v, key) == 0 {
        cur.key = append(cur.key[:i], cur.key[i+1:]...)
      }
    }

    if cur.sibling == uint32(0) {
      return fmt.Errorf("Not yet implemented edge case.")
    }

    sibling, err := t.storage.loadNode(cur.sibling)
    if err != nil {
      return fmt.Errorf("Error loading sibling : %w", err)
    }

    if len(sibling.key) <= t.minKeyNum + 1 {
      // merge.
      err = t.mergeNode(sibling, cur)
      if err != nil {
        return fmt.Errorf("Error merging nodes : %w", err)
      }

      return nil
    }

    // else borrow.
    err = t.borrowKey(sibling, cur)
    if err != nil {
      return fmt.Errorf("Error borrowing key : %w", err)
    }

    return nil
  }

  // no need to merge or borrow.
  for i, v := range cur.key {
    if compare(v, key) == 0 {
      cur.key = append(cur.key[:i], cur.key[i+1:]...)
    }
  }

  err := t.storage.updateNode(cur)
  if err != nil {
    return fmt.Errorf("Error updating node : %w", err)
  }

  return nil
}

func (t *BPTree) borrowKey(sibling, cur *node) error {
  // NOTE: borrowKey is tested only with leaf nodes.
  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil {
    return fmt.Errorf("Error loading parent.")
  }

  bKey := sibling.key[0]
  bVal := sibling.pointers[0]

  sibling.key = sibling.key[1:]
  sibling.pointers = sibling.pointers[1:]

  cur.key = append(cur.key, bKey)
  cur.pointers = append(cur.pointers, bVal)

  index := 0
  // find the postion in parent.
  for i, v := range parent.pointers {
    if v.asNodeId() == cur.id{
      // found index.
      index = i
    }
  }

  // pointers[index] represents the node that has keys lesser than key[index]
  // we need to change this one based on the borrow.

  if compare(parent.key[index], bKey) == 0 {
    parent.key[index] = sibling.key[0]
    err = t.storage.updateNode(parent)
    if err != nil {
      return fmt.Errorf("Error updating parent : %w", err)
    }
  }

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating current : %w", err)
  }
  
  if err := t.storage.updateNode(sibling); err != nil {
    return fmt.Errorf("Error updating sibling : %w", err)
  }

  return nil
}

func (t *BPTree) mergeNode(sibling, cur *node) error {
  
  // merge effectively causes deletion of a node.
  // this means the parent will have 1 key less.
  // this means that the parent can go below the threshold as well.
  // this means we have may have to merge/borrow from a non leaf node as well.
 
  parent, err := t.storage.loadNode(cur.parentId)
  if err != nil {
    return fmt.Errorf("Error loading parent.")
  }

  cur.key = append(cur.key, sibling.key...)
  cur.pointers = append(cur.pointers, sibling.pointers...)
  cur.sibling = sibling.sibling

  if err = t.storage.deleteNode(sibling.id); err != nil {
    return fmt.Errorf("Error deleting node : %w", err)
  }

  index := 0
  for i, v := range parent.pointers {
    if v.asNodeId() == cur.id{
      // found index.
      index = i
    }
  }
  
  parent.key = append(parent.key[:index], parent.key[index+1:]...)
  parent.pointers = append(parent.pointers[:index+1], parent.pointers[index+2:]...)

  if err := t.storage.updateNode(cur); err != nil {
    return fmt.Errorf("Error updating current : %w", err)
  }


  for parent.id != t.metadata.rootId {
    if len(parent.key) < t.minKeyNum && parent.id != t.metadata.rootId{
      // parent is loaded.
      // we deleted a key from the parent.
      // we now have to borrow or merge the parent.
      // load the sibling of the parent.
      
      sibling, err := t.storage.loadNode(parent.sibling)
      if err != nil {
        return fmt.Errorf("Error loading sibling : %w", err)
      }

      if len(sibling.key) < t.minKeyNum + 1 {

        // merge parent and sibling and demote key of grandparent.
        return fmt.Errorf("Not yet finished merging non leaf node.")


      } else {

        // demote from parent and promote from sibling.
        // exit since we are not reducing number of keys in parent.

        grandparent, err := t.storage.loadNode(parent.parentId)
        if err != nil {
          return fmt.Errorf("Error Loading grandparent : %w", err)
        }

        //find key in grandparent.
        index := 0
        for i, v := range grandparent.pointers {
          if v.asNodeId() == parent.id{
            // found index.
            index = i
          }
        }

        parent.key = append(parent.key, grandparent.key[index])
        parent.pointers = append(parent.pointers, sibling.pointers[0])

        grandparent.key[index] = sibling.key[0]

        sibling.key = sibling.key[1:]
        sibling.pointers = sibling.pointers[1:]

        // update parent grandparent and sibling.

        if err := t.storage.updateNode(parent); err != nil {
          return fmt.Errorf("Error updaing parent while borrowing: %w", err)
        }

        if err := t.storage.updateNode(grandparent); err != nil {
          return fmt.Errorf("Error updating grandparent while borrowing : %w", err)
        }

        if err := t.storage.updateNode(sibling); err != nil {
          return fmt.Errorf("Error updating parent's sibling while borrowing : %w", err)
        }

        return nil
      }

    // end if.
    } else {
      return fmt.Errorf("Not yet invented reducing height i guess.")
    }
  // end for loop.
  }

  // if it comes out of loop.
  // parent is root and we need to check if height should be reduced.
  if len(parent.key) < t.minKeyNum {
    return fmt.Errorf("Reducing height has not yet been invented.")
  }

  return nil
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
