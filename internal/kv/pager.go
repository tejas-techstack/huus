/*
* This is an extension of storage functionallity.
*/
package kv

import (
  "fmt"
)


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

func (s *storage) readPage(pageId uint32) ([]byte, error) {
  if s.isFreePage(pageId) {
    return nil, fmt.Errorf("Error reading a free page.")
  }

  page := make([]byte, s.pageSize)

  offset := int64(int(pageId) * int(s.pageSize) + metadataSize)
  _, err := s.fo.ReadAt(page, offset)
  if err != nil {
    return nil, fmt.Errorf("Error reading from file : %w", err)
  }

  return page,nil
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

