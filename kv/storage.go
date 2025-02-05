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
  lastPageId uint32
  // custom metadata is the tree metadata.
  custom []byte
}

func newStorage (path string, pageSize uint16, order uint16) (*storage, error){

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

    s := &storage{
      fo : fo,
      pageSize : pageSize,
      freePages: nil,
      lastPageId : 1,
      metadata : &storageMetadata{pageSize, 1, nil,},
    }
  
    /*
    rootId,_ := s.newPage()
    custom := encodeMetadata(&treeMetaData{
      order : order,
      rootId : rootId,
      pageSize : s.pageSize,
    })

    s.metadata.custom = custom
    */

    if err := s.writeStorageMetadata(s.metadata); err != nil {
      return nil, fmt.Errorf("Error Writing metadata : %w", err)
    }

    if err := s.flush(); err != nil {
      return nil, fmt.Errorf("Error flushing the file : %w", err)
    }

    return s, nil
  }

  metadata, err := readStorageMetadata(fo)
  if err != nil {
    return nil, fmt.Errorf("Error reading storage metadata : %w",err)
  }

  // TODO implement later once everything else is working.
  /*
  freePages,err := readFreePages()
  if err != nil {
    return nil, fmt.Errorf("Error reading free pages : %w" ,err)
  }
  */

  lastPageId := metadata.lastPageId
  
  // BUG: freePages initialized to nil this needs to be changed.
  return &storage{fo, pageSize, nil, lastPageId, metadata}, nil
}

func (s *storage) writeStorageMetadata(md *storageMetadata) error {
  offset := int64(0)
  dataToWrite := make([]byte, 1000)
  copy(dataToWrite, encodeStorageMetadata(md))

  n, err := s.fo.WriteAt(dataToWrite, offset)
  if err != nil {
    return fmt.Errorf("Error writing to page : %w", err)
  } 
  if n != len(dataToWrite) {
    return fmt.Errorf("Bytes written lesser than given bytes.")
  }

  err = s.flush()
  if err != nil {
    return fmt.Errorf("Error flushing to file : %w", err)
  }

  return nil

}

// func readFreePages() ([]uint32, error) { /*returns freePages by reading from file.*/}

func readStorageMetadata(fo *os.File) (*storageMetadata, error) {
  offset := int64(0)
  data := make([]byte, 1000)

  _, err := fo.ReadAt(data, offset)
  if err != nil {
    return nil,fmt.Errorf("Error reading data : %w", err)
  }

  metadata := decodeStorageMetadata(data)
  
  return metadata, nil
}

func (s *storage) loadMetadata() (*treeMetaData, error) {
  if s.metadata.custom == nil {
    return nil, nil
  }

  md, _ := decodeMetadata(s.metadata.custom)
  return md, nil
}

func (s *storage) updateMetadata(tmd *treeMetaData) error {

  s.metadata.custom = encodeMetadata(tmd)

  err := s.writeStorageMetadata(s.metadata)
  if err != nil {
    return fmt.Errorf("Error updating storage metadata : %w", err)
  }

  s.metadata, err = readStorageMetadata(s.fo)
  if err != nil {
    return fmt.Errorf("Error reading metadata : %w", err)
  }

  return nil
  // return fmt.Errorf("Not yet implemented")
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
    _, err := s.fo.ReadAt(tempData, int64(offset))
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
