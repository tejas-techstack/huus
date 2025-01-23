// storage.go
// Used to store data page wise
package kv

import (
  "fmt"
  "os"
)

// define the maximum size of metadata as 1000 bytes.
const metadataSize 1000


type storage struct {
  fo *os.File

  pageSize uint16

  // maintain array of freepages.
  freePages []uint32

  lastPageId uint32

  metadata *storageMetadata
}

type storageMetadata struct {
  pageSize uint16

  // custom metadata is the tree metadata.
  custom []byte
}

func newStorage (path string, pageSize uint16) (*storage, error){

  /* new  storage is supposed to  
    open the file to write.
      > if the file is opened
      a. it needs to load the free pages.
      b. read storage metadata and ensure the pagesize is equal to given pagesize
      > if file is empty
      a. check if pageSize is less than minPageSize
      b. initialize metadata for the file and write it to the file.
  */

  fo, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
  if err != nil {
    return nil, fmt.Errorf("Error opening file : %w", err)
  }

  if pageSize < minPageSize {
    return nil, fmt.Errorf("Given Page Size is lesser than minimum page size")
  }

  info, err := fo.Stat()
  if err != nil {
    return fmt.Errorf("Error trying to stat the file : %w", err)
  }

  if info.Size() == 0 {
    // file is empty, need to initialize:
    
    // metadata is basically a copy of the storage.
    // that will be written to the file.
    // the lastPageId is not passed to writeMetadata
    // since it is defined based on the contents of the file.
    
    storage := &storage{
      fo : fo,
      pageSize : pageSize,
      freePages: nil,
      lastPageId : 0,
      metadata : &storageMetadata{pageSize, nil,},
    }

    if err := s.writeMetadata(); err != nil {
      return nil, fmt.Errorf("Error Writing metadata")
    }

    freePages, err := initializeFreePages(storage)
    if err != nil {
      return nil, fmt.Errorf("Error initializing free pages")
    }

    if err := storage.flush(); err != nil {
      return nil, fmt.Errorf("Error flushing the file.")
    }
  }



  return nil, fmt.Errorf("Not yet implemented")

}

func (s *storage) writeMetadata() error {

  return fmt.Errorf("Not yet implemented")
}

func (s *storage) loadMetadata() (*treeMetaData, error) {

  // need to see how order seems to be managed.

  /* from the file load the custom metadata.
     the custom metadata contains: 
     1. order
     2. rootId
     3. pageSize
  */ 

  return nil, fmt.Errorf("Not yet implemented")
}

func (s *storage) updateMetaData(newRootId uint32) error {

  // write new metadata to the file with new root id.
  // read current metadata, decode, change the root id ,
  // encode and write back.

  return fmt.Errorf("Not yet implemented")
}

func (s *storage) loadNode(nodeId uint32) (*node, error) {

  // the nodeId itself represents the pageId.
  // calculate the offset from this pageId and the
  // read the data from that pageId.
  // in the data read, the last 4 bytes contain either 0, or
  // it contains the next pageId, if the nextPage Id is non zero,
  // read from that page Id as well and then keep chaining until the
  // data is not fully read.
  // once read decode the data as a node and return it.

  return nil, fmt.Errorf("Not yet implemented")
}

func (s *storage) updateNode(cur *node) error {

  // using the node id search for the node
  // get the offset and write to this offset.

  return fmt.Errorf("Not yet implemented")
}

func (s *storage) newNode() (uint32, error) {

  // create a node with empty data,
  // get a free page from the existing free pages.
  // load the node into this free page, return the nodeId

  return uint32(0), fmt.Errorf("Not yet implemented")
}

func (s *storage) flush() error {          i
  if err := s.fo.Flush(); err != nil {
    return fmt.Errorf("Error flushing the file : %w",err)
  }

  return nil
}
