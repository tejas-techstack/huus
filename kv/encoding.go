package kv

import (
  "encoding/binary"
)

// encode/ decode :
// treeMetaData (uint16, uint32, uint16)
//            - (order, rootId, pageSize)
// node:
// id : uint32
// parentId : uint32
// key : [][]byte
// pointers []*pointer
// isLeaf : bool
// sibling : uint32

func encodeUint16(val uint16) []byte {
  data [2]byte

  binary.BigEndian.PutUint16(data[:], val)

  return data[:]
}

func decodeUint16(data []byte) uint16 {
  return binary.BigEndian.Uint16(data)
}

func encodeUint32(val uint32) []byte {
  data [4]byte

  binary.BigEndian.PutUint32(data[:], val)

  return data[:]
}

func decodeUint32(data []byte) uint32 {
  return binary.BigEndian.Uint32(data)
}

func encodeBool(val bool) []byte {
  data [1]byte

  if val {
    data[0] = 1
  } else {
    data[0] = 0
  }

  return data[:]
}

func decodeBool(data []byte) bool {
  if data[0] {
    return true
  } else {
    return false
  }
}

func encodeNode(node *node) []byte {
  data := make([]byte, 0)

  data = append(data, encodeUint32(node.id))
  data = append(data, encodeUint32(node.parentId))
  data = append(data, encodeBool(node.isLeaf))

  // number of keys to read.
  data = append(data, encodeUint32(uint32(len(node.key))))
  for i:=0; i<len(node.key); i++ {
    data = append(data, encodeUint16(uint16(node.key[i]...)))
  }

  data = append(data, encodeUint32(uint32(len(node.pointers))))

  if node.isLeaf {
    for i:=0; i<len(node.pointers); i++ {
      // encode length of value and then encode value.
      data = append(data, encodeUint16(uint16(len(pointers[i].asValue()...))))
      data = append(data, node.pointers[i].asValue()...)
    }
  } else {
    for i:=0; i<len(node.pointers); i++ {
      data = append(data, encodeUint32(node.pointers[i].asNodeId()))
    }
  }

  data = append(data, encodeUint32(node.sibling))

  return data[:]
}

func decodeNode(data []byte) *node {
  pos := 0

  id := decodeUint32(data[pos: pos+4])
  pos += 4
  parentId := decodeUint32(data[pos: pos+4])
  pos += 4
  isLeaf := decodeBool(data[pos: pos+1])
  pos += 1

  keyNum := decodeUint32(data[pos:pos+4])
  pos += 4

  key := make([]byte, 0)
  for i:=0; i < keyNum; i++ {
    key = append(decodeUint16(data[pos:pos + 2]))
    pos += 2
  }

  pointerNum := decodeUint32(data[pos:pos+4])
  pos += 4

  pointers := make([]*pointer, 0)
  if newNode.isLeaf {
    var varValLen int
    for i := 0; i < pointerNum; i++ {
      valLen := decodeUint16(data[pos:pos+2])
      pos += 2
      pointers := append(pointers, data[])
    }
  }


}

func encodeMetaData(metadata *treeMetaData) []byte {}

func encodeMetaData(data []byte) *treeMetaData {}
