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

  t.Log("Test Node:")
  t.Log(newNode)

  encoded := encodeNode(newNode)

  t.Log("Encoded Node:")
  t.Log(encoded)

  decoded, err := decodeNode(encoded)
  if err != nil {
    t.Fatalf("failed to decode node : %s", err)
  }

  t.Log("Decoded Node:")
  t.Log(decoded)
}

func TestEncodeMetaData(t *testing.T) {
  metadata := &treeMetaData{
    order : uint16(12),
    rootId : uint32(16),
    pageSize : uint16(61),
  }

  t.Log("Test metaData:")
  t.Log(metadata)

  encoded := encodeMetaData(metadata)

  t.Log("Encoded metadata:")
  t.Log(encoded)

  decoded, err := decodeMetaData(encoded)
  if err != nil {
    t.Fatalf("Failed to decode metadata : %s", err)
  }

  t.Log("Decoded metadata")
  t.Log(decoded)
}
