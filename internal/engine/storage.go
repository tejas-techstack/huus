// storage.go
// Used to store data page wise
package engine

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

	// Page cache that is in memory.
	// Key is pageId.
	// BUG: CURRENT ITERATION CACHE IS UNBOUNDED
  cache map[uint32][]byte
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
      cache : make(map[uint32][]byte),
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

    if err := s.sync(); err != nil {
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
  return &storage{
    fo : fo,
    pageSize : pageSize,
    freePages : nil,
    lastPageId : lastPageId,
    metadata : metadata,
    cache : make(map[uint32][]byte),
  }, nil
}


// Return cached buffer, If it is not in memory then read it from disk
func (s *storage) getPage(pageId uint32) ([]byte, error) {
  if buf, ok := s.cache[pageId]; ok {
    return buf, nil
  }

  buf := make([]byte, s.pageSize)
  offset := int64(int(pageId) * int(s.pageSize) + metadataSize)
  if _, err := s.fo.ReadAt(buf, offset); err != nil {
    return nil, fmt.Errorf("Error reading page %d : %w", pageId, err)
  }

  s.cache[pageId] = buf
  return buf, nil
}

func (s *storage) setPage(pageId uint32, data []byte) error {
  buf, ok := s.cache[pageId]
  if !ok {
    buf = make([]byte, s.pageSize)
    s.cache[pageId] = buf
  }

  // data may be shorter than a full page, so clear before copying.
  clear(buf)
  copy(buf, data)

  offset := int64(int(pageId) * int(s.pageSize) + metadataSize)
  if _, err := s.fo.WriteAt(buf, offset); err != nil {
    return fmt.Errorf("Error writing page %d : %w", pageId, err)
  }

  return nil
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

  page, err := s.getPage(nodeId)
  if err != nil {
    return nil, fmt.Errorf("Error reading file : %w", err)
  }

	// Copy to avoid modifying the cache Directly
  data := make([]byte, int(s.pageSize))
  copy(data, page)

  nextPageId := decodeUint32(data[4:8])
  data = data[8:]
  for nextPageId != uint32(0) {
    tempPage, err := s.getPage(nextPageId)
    if err != nil {
      return nil, fmt.Errorf("error reading file : %w", err)
    }
    tempData := make([]byte, int(s.pageSize))
    copy(tempData, tempPage)

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

  return nil
}



func (s *storage) deleteNode(nodeId uint32) error {

  if nodeId == uint32(0) {
    return fmt.Errorf("Error trying to free the wrong page (ie page 0).")
  }
  // need to load all pageIds and free them.

  pageId := nodeId
  for pageId != uint32(0) {

    pageData, err := s.readPageData(pageId)
    if err != nil {
      return fmt.Errorf("Error reading pageData : %w", err)
    }

    s.freePages = append(s.freePages, pageId)

    pageId = decodeUint32(pageData[4:8])
  }

  return nil
}

func (s *storage) newNode() (uint32, error) {

  // as of right now a new node is always assigned from the end of the file.
  // later on shift it to making it link with exisiting freePageIds.
  newNodeId := s.lastPageId
  s.lastPageId++
  nextNodeId := uint32(0);

  newNode := &node{
    id : newNodeId,
    key : nil,
    pointers : nil,
    isLeaf : false,
    sibling : 0,
  }

  data := make([]byte, s.pageSize)
  copy(data[0:4], encodeUint32(newNodeId))
  copy(data[4:8], encodeUint32(nextNodeId))
  copy(data[8:], encodeNode(newNode))

  if err := s.setPage(newNodeId, data); err != nil {
    return uint32(0), fmt.Errorf("Error writing to file :%w", err)
  }

  return newNodeId, nil
}

func (s *storage) sync() error {
  if err := s.fo.Sync(); err != nil {
    return fmt.Errorf("Error flushing the file : %w",err)
  }

  return nil
}

func (s *storage) close() error {
  if err := s.sync(); err != nil {
    return err
  }

  return s.fo.Close()
}
