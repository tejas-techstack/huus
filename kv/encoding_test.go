package kv

import (
  "testing"
  "fmt"
)

func TestEncodeNode(t *testing.T) {
  newNode := &node{
    id : 42,
    isLeaf : true,
    parentId : 12,
    key : [][]byte{
      {1,2,3,4},
      {5,6,7,8},
      nil,
    },
    pointers : []*pointer{
      {uint32(52)},
      {uint32(61)},
    },
    sibling : uint32(16),
  }

  fmt.Println(newNode)

  decoded, err := decodeNode(encodeNode(newNode))
  if err != nil {
    t.Fatalf("failed to decode node : %s", err)
  }

  fmt.Println(decoded)
}
