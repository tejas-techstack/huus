// storage.go
// Used to store data page wise
package kv

type storage struct {
  fo *file

  pageSize uint16

  // maintain array of freepages.
  freePages []uint32

  lastPageId uint32
}

func newStorage (path string, pageSize uint16) (*storage, error){}

func loadMetadata() (*treeMetadata, error) {
  // read from storage.fo header

}

func (s *storage) loadNode(nodeId uint32) (*node, error) {}

func (s *storage) updateNode() (error) {}
