// storage.go
// Used to store data page wise
package kv

type storage struct {
  file *file

  pageSize uint16

  // maintain array of freepages.
  freePages []uint32

  lastPageId uint32
}

func newStorage (path string, pageSize uint16) (*storage, error){}

func loadMetadata() (*treeMetadata, error) {}
