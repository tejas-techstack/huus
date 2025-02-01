// storage.go
// Used to store data page wise
package kv

import (
  "fmt"
  "os"
)

// define the maximum size of metadata as 1000 bytes.
const metadataSize=1000


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
    return nil, fmt.Errorf("Error trying to stat the file : %w", err)
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
      return nil, fmt.Errorf("Error Writing metadata : %w", err)
    }

    err := s.initializeFreePages()
    if err != nil {
      return nil, fmt.Errorf("Error initializing free pages : %w", err)
    }

    if err := s.flush(); err != nil {
      return nil, fmt.Errorf("Error flushing the file : %w", err)
    }

    return s, nil
  }

  metadata, err := readStorageMetadata()
  if err != nil {
    return nil, fmt.Errorf("Error reading storage metadata : %w",err)
  }

  freePages,err := readFreePages()
  if err != nil {
    return nil, fmt.Errorf("Error reading free pages : %w" ,err)
  }


  // TODO fix this.
  lastPageId, err := getLastPageId()
  if err != nil {
    return nil, fmt.Errorf("Error loading lastPageId : %w", err)
  }

  return &storage{fo, pageSize, freePages, lastPageId, metadata}, nil
}

func (s *storage) initializeFreePages() error {
  return fmt.Errorf("Not yet implemented")
}

func (s *storage) writeStorageMetadata() error {

  return fmt.Errorf("Not yet implemented")
}

func readStorageMetadata() (*storageMetadata, error) {
  return nil, fmt.Errorf("Not yet implemented")
}

func readFreePages() ([]uint32, error) {
  return nil, fmt.Errorf("Not yet implemented")
}

func getLastPageId() (uint32, error) {
  return uint32(0), fmt.Errorf("Not yet implemented")
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

func (s *storage) loadNodeRaw(nodeId uint32) ([]byte, error) {

  offset := (int(nodeId) * int(s.pageSize)) + metadataSize

  data := make([]byte, int(s.pageSize))
  _, err := s.fo.ReadAt(data, int64(offset))
  if err != nil {
    return nil, fmt.Errorf("Error reading file")
  }


  
  nextPageId := decodeUint32(data[4:8])
  data = data[8:]
  for nextPageId != uint32(0) {
    tempData := make([]byte, int(s.pageSize))
    offset := (int(nextPageId) * int(s.pageSize)) + metadataSize
    _, err := s.fo.ReadAt(data, int64(offset))
    if err != nil {
      return nil, fmt.Errorf("error reading file : %w", err)
    }
    nextPageId = decodeUint32(tempData[4:8])
    tempData = tempData[8:]

    data = append(data, tempData...)
  }



  return data, nil
}

func (s *storage) loadNode(nodeId uint32) (*node, error) {

  if s.isFreePage(nodeId) == true {
    return nil, fmt.Errorf("Node does not exist.")
  }

  data, err := s.loadNodeRaw(nodeId)
  if err != nil {
    return nil, fmt.Errorf("Error loading raw node data : %w", err)
  }

  node, _ := decodeNode(data)
  return node, nil
}

func (s *storage) updateNode(cur *node) error {

  if s.isFreePage(cur.id) == true {
    return fmt.Errorf("Node does node exist.")
  }

  curPageId := cur.id
  data := encodeNode(cur)

  oldData, err := s.loadNodeRaw(curPageId)
  if err != nil{
    return fmt.Errorf("Error loading raw node data : %w", err)
  }

  if len(oldData) <= len(data) {
    // new pages may be required.
    nextPageId, err := s.writePages(curPageId, data)
    if err != nil {
      return fmt.Errorf("Error writing pages : %w", err)
    }

    // should never occur but is a safety net.
    if nextPageId != uint32(0) {
      return fmt.Errorf("Unknown error while writing pages, all pages not initialized properly")
    }


  } else {
    // no new pages required.
    // need to free pages.

    nextPageId, err := s.writePages(curPageId, data)
    if err != nil {
      return fmt.Errorf("Error writing pages : %w", err)
    }

    pageToFree := nextPageId
    for pageToFree != uint32(0) {
      nextPageData, err := s.readPageData(pageToFree)
      if err != nil {
        return fmt.Errorf("Error reading page data : %w", err);
      }

      err = s.freeThePage(pageToFree)
      if err != nil {
        return fmt.Errorf("Error freeing the page : %w", err)
      }

      pageToFree = decodeUint32(nextPageData[4:8])
    }
  }

  err = s.flush()
  if err != nil {
    return fmt.Errorf("Error flushing to file : %w", err)
  }

  return nil
}

// writes all data, creates new pages if required.
// returns nextPageId if its not 0.
// (all data is written but nextPageId isnt 0 means we need to free pages.)
func (s *storage) writePages(curPageId uint32, data []byte) (uint32, error) {

  nextPageId := uint32(0)
  for len(data) != 0 {
    nextPageData, err := s.readPageData(curPageId)
    if err != nil {
      return uint32(0), fmt.Errorf("Error reading page data : %w", err)
    }

    nextPageId = decodeUint32(nextPageData[4:8])

    // if nextPageId is 0 and data is still there to write
    // generate a new page.
    // else last page to write anyways and keep it 0
    if nextPageId == uint32(0) && (len(data) > int(s.pageSize-8)) {
      nextPageId, err = s.newPage()
      if err != nil {
        return uint32(0), fmt.Errorf("Error generating new page : %w", err)
      }
    }

    dataToWrite := make([]byte, s.pageSize)
    copy(dataToWrite[0:4], encodeUint32(curPageId))
    copy(dataToWrite[4:8], encodeUint32(nextPageId))

    if len(data) <= int(s.pageSize-8) {
      copy(dataToWrite[8:], data)
      data = nil
    } else {
      copy(dataToWrite[8:], data[:s.pageSize-8])
      data = data[s.pageSize-8:]
    }

    err = s.writePage(curPageId, dataToWrite)
    if err != nil {
      return uint32(0), fmt.Errorf("Error writing page : %w", err)
    }

    curPageId = nextPageId
  }

  return nextPageId, nil
}

func (s *storage) readPageData(pageId uint32) ([]byte, error) {
  if s.isFreePage(pageId) {
    return nil, fmt.Errorf("Page Id : %d is empty.", pageId)
  }

  offset := int64(int(pageId) * int(s.pageSize) + metadataSize)
  pageData := make([]byte, 8)
  _, err := s.fo.ReadAt(pageData, offset)
  if err != nil {
    return nil, fmt.Errorf("Error reading page data : %w", err)
  }

  return pageData, nil
}

func (s *storage) writePage(pageId uint32, data []byte) error {
  if s.isFreePage(pageId) {
    return fmt.Errorf("Page id : %d is empty", pageId)
  }

  if pageId == uint32(0) {
    return fmt.Errorf("Writing to page index 0 is not allowed.")
  }

  page := make([]byte, s.pageSize)
  copy(page, data)
  
  offset := int64(int(pageId) * int(s.pageSize) + metadataSize)
  n, err := s.fo.WriteAt(page, offset)
  if err != nil {
    return fmt.Errorf("Error writing to file : %w", err)
  } else {
    if n != len(page) {
      return fmt.Errorf("Wanted to write: %d, wrote : %d", len(page), n)
    }
  }

  return nil
}

func (s *storage) isFreePage(pageId uint32) bool {
  for i:=0; i<len(s.freePages); i++ {
    if pageId == s.freePages[i] {
      return true
    }
  }

  return false
}

func (s *storage) freeThePage(pageId uint32) error {
  if pageId == uint32(0) {
    return fmt.Errorf("Page cannot be freed, reserved.")
  }

  emptyData := make([]byte, s.pageSize)

  err := s.writePage(pageId, emptyData)
  if err != nil {
    return fmt.Errorf("Error writing to page : %w", err)
  }

  err = s.flush()
  if err != nil {
    return fmt.Errorf("Error flushing : %w", err)
  }

  s.freePages = append(s.freePages, pageId)

  return nil
}

func (s *storage) newPage() (uint32, error) {
  // generates a new page with nextPageId = 0 by default.
  
  if len(s.freePages) == 0 {
    pageId := s.lastPageId
    s.lastPageId++

    emptyData := make([]byte, int(s.pageSize))
    err := s.writePage(pageId, emptyData)
    if err != nil {
      return uint32(0),fmt.Errorf("Error writing to page : %w", err)
    }

    err = s.flush()
    if err != nil {
      return uint32(0), fmt.Errorf("Error flushing : %w", err)
    }

    return pageId, nil
  }

  pageId := s.freePages[len(s.freePages)-1]
  s.freePages = s.freePages[:len(s.freePages)-1]

  if pageId == uint32(0) {
    return uint32(0), fmt.Errorf("Error generating new page id for some reason.")
  }

  return pageId, nil
}

func (s *storage) newNode() (uint32, error) {

  // as of right now a new node is always assigned from the end of the file.
  // later on shift it to making it link with exisiting freePageIds.
  newNodeId := s.lastPageId
  s.lastPageId++
  nextNodeId := uint32(0);

  newNode := &node{
    id : newNodeId,
    parentId : uint32(0),
    key : nil,
    pointers : nil,
    isLeaf : false,
    sibling : 0,
  }

  data := make([]byte, s.pageSize)
  copy(data[0:4], encodeUint32(newNodeId))
  copy(data[4:8], encodeUint32(nextNodeId))
  copy(data[8:], encodeNode(newNode))

  offset := (int(newNodeId) * int(s.pageSize)) + metadataSize
  n, err := s.fo.WriteAt(data, int64(offset))
  if err != nil {
    return uint32(0), fmt.Errorf("Error writing to file :%w",err)
  } else {
    if n != len(data) {
      return uint32(0), fmt.Errorf("Had to write %d, only wrote %d", len(data), n)
    }
  }

  err = s.flush()
  if err != nil {
    return uint32(0), fmt.Errorf("Error flushing : %w", err)
  }

  return newNodeId, nil
}

func (s *storage) flush() error {
  if err := s.fo.Sync(); err != nil {
    return fmt.Errorf("Error flushing the file : %w",err)
  }

  return nil
}
