// storage.go
// Used to store data page wise
package kv

import (
  "fmt"
  "os"
)

type storage struct {
  fo *os.File

  pageSize uint16

  // maintain array of freepages.
  freePages []uint32

  lastPageId uint32
}

func newStorage (path string, pageSize uint16) (*storage, error){
  return nil, fmt.Errorf("Not yet implemented")
}

func (s *storage) loadMetadata() (*treeMetaData, error) {
  return nil, fmt.Errorf("Not yet implemented")
}

func (s *storage) updateMetaData(nodeId uint32) error {
  return fmt.Errorf("Not yet implemented")
}

func (s *storage) loadNode(nodeId uint32) (*node, error) {
  return nil, fmt.Errorf("Not yet implemented")
}

func (s *storage) updateNode(nodeId uint32) error {
  return fmt.Errorf("Not yet implemented")
}

func (s *storage) newNode() (uint32, error) {
  return uint32(0), fmt.Errorf("Not yet implemented")
}
