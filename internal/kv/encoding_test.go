package kv

import (
  "testing"
)


func TestEncodeNode(t *testing.T) {
  newNode := &node{
    id : uint32(42),
    isLeaf : false,
    parentId : uint32(12),
    key : [][]byte{
      {1,2,3,4},
      {5,6,7,8},
    },
    pointers : []*pointer{
      {uint32(52)},
      {uint32(61)},
      {uint32(91)},
    },
    sibling : uint32(16),
  }

  encoded := encodeNode(newNode)
  decoded, err := decodeNode(encoded)
  if err != nil {
    t.Fatalf("failed to decode node : %s", err)
  }

  t.Log("Decoded Node:", decoded)
}

func TestEncodeMetaData(t *testing.T) {
  metadata := &treeMetaData{
    order : uint16(12),
    rootId : uint32(16),
    pageSize : uint16(61),
  }

  encoded := encodeMetadata(metadata)
  decoded, err := decodeMetadata(encoded)
  if err != nil {
    t.Fatalf("Failed to decode metadata : %s", err)
  }

  t.Log("Decoded metadata: ", decoded) 
}
