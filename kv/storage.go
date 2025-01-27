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

  // last page id is basically used to store
  // the last offset of the file.
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
    
    s := &storage{
      fo : fo,
      pageSize : pageSize,
      freePages: nil,
      lastPageId : 0,
      metadata : &storageMetadata{pageSize, nil,},
    }
   
    if err := s.writeStorageMetadata(); err != nil {
      return nil, fmt.Errorf("Error Writing metadata")
    }

    freePages, err := s.initializeFreePages()
    if err != nil {
      return nil, fmt.Errorf("Error initializing free pages")
    }

    if err := s.flush(); err != nil {
      return nil, fmt.Errorf("Error flushing the file.")
    }
  }

  metadata, err := s.readStorageMetadata()
  if err != nil {
    return fmt.Errorf("Error reading storage metadata : %w",err)
  }

  freePages,err := s.readFreePages()
  if err != nil {
    return fmt.Errorf("Error reading free pages : %w" ,err)
  }


  // TODO fix this. 
  lastPageId, err := s.getLastPageId()
  if err != nil {
    return fmt.Errorf("Error loading lastPageId : %w", err)
  }

  return &storage{fo, pageSize, freePages, lastPageId, metadata}, nil
}

func (s *storage) writeStorageMetadata() error {

  return fmt.Errorf("Not yet implemented")
}

func (s *storage) readStorageMetadata() error {
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
  // calculate the offset from this pageId and then
  // read the data from that pageId.
  // in the data read, the last 4 bytes contain either 0, or
  // it contains the next pageId, if the nextPage Id is non zero,
  // read from that page Id as well and then keep chaining until the
  // data is not fully read.
  // once read decode the data as a node and return it.

  // pageId 1 : always will start at 1001 byte 
  // since 1000 bytes are reserved for metadata.
  offset := (int(nodeId) * int(s.pageSize)) + metadataSize

  data := make([]byte, int(s.pageSize))
  _, err := s.fo.ReadAt(data, offset)
  if err != nil {
    return fmt.Errorf("Error reading file")
  }

  dataLen := len(data)

  nextPageId, err := decodeUint32(data[dataLen-4:])
  data = data[:dataLen-4]
  for nextPageId != uint32(0) {
    tempData := make([]byte, int(s.pageSize))
    _, err := s.fo.ReadAt(data, nextPageId)
    if err != nil {
      return nil, fmt.Errorf("error reading file : %w", err)
    }
    nextPageId = decodeUint32(tempData[dataLen-4:])
    tempData = tempData[:dataLen-4]

    data = append(data, tempData)
  }

  node, err := decodeNode(data)
  if err != nil {
    return nil, fmt.Errorf("Error decoding node : %w", err)
  }

  return node, nil
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
