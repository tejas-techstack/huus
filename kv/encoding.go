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
  var data [2]byte

  binary.BigEndian.PutUint16(data[:], val)

  return data[:]
}

func decodeUint16(data []byte) uint16 {
  return binary.BigEndian.Uint16(data)
}

func encodeUint32(val uint32) []byte {
  var data [4]byte

  binary.BigEndian.PutUint32(data[:], val)

  return data[:]
}

func decodeUint32(data []byte) uint32 {
  return binary.BigEndian.Uint32(data)
}

func encodeBool(val bool) []byte {
  var data [1]byte

  if val {
    data[0] = 1
  } else {
    data[0] = 0
  }

  return data[:]
}

func decodeBool(data []byte) bool {
  if data[0] == 1{
    return true
  } else {
    return false
  }
}

func encodeNode(curr *node) []byte {
  data := make([]byte, 0)

  data = append(data, encodeUint32(curr.id)...)
  data = append(data, encodeUint32(curr.parentId)...)
  data = append(data, encodeBool(curr.isLeaf)...)

  // number of keys to read.
  data = append(data, encodeUint32(uint32(len(curr.key)))...)
  for i:=0; i<len(curr.key); i++ {
    data = append(data, encodeUint32(uint32(len(curr.key[i])))...)
    data = append(data, curr.key[i]...)
  }

  data = append(data, encodeUint32(uint32(len(curr.pointers)))...)

  if curr.isLeaf {
    for i:=0; i<len(curr.pointers); i++ {
      // encode length of value and then encode value.
      data = append(data, encodeUint16(uint16(len(curr.pointers[i].asValue())))...)
      data = append(data, curr.pointers[i].asValue()...)
    }
  } else {
    for i:=0; i<len(curr.pointers); i++ {
      data = append(data, encodeUint32(curr.pointers[i].asNodeId())...)
    }
  }

  data = append(data, encodeUint32(curr.sibling)...)

  return data[:]
}

func decodeNode(data []byte) (*node, error) {
  pos := 0

  id := decodeUint32(data[pos: pos+4])
  pos += 4
  parentId := decodeUint32(data[pos: pos+4])
  pos += 4
  isLeaf := decodeBool(data[pos: pos+1])
  pos += 1

  keyNum := int(decodeUint32(data[pos:pos+4]))
  pos += 4

  key := make([][]byte, 0)
  for i:=0; i < keyNum; i++ {
    keyLen := int(decodeUint32(data[pos: pos+4]))
    pos += 4
    key = append(key, data[pos:pos+keyLen])
    pos += keyLen
  }

  pointerNum := int(decodeUint32(data[pos:pos+4]))
  pos += 4

  pointers := make([]*pointer, 0)
  if isLeaf {
    for i := 0; i < pointerNum; i++ {
      valLen := int(decodeUint16(data[pos:pos+2]))
      pos += 2

      value := data[pos : pos+valLen]
      pos += valLen

      pointers = append(pointers, &pointer{value})
    }
  } else {
    // node is not a leaf.
    for i := 0; i < pointerNum; i++ {
    
      childId := decodeUint32(data[pos: pos+4])
      pos += 4

      pointers = append(pointers, &pointer{childId})
    }
  }

  sibling := decodeUint32(data[pos: pos+4])

  newNode := &node{
    id : id,
    parentId: parentId,
    isLeaf : isLeaf,
    key : key,
    pointers : pointers,
    sibling : sibling, 
  }

  return newNode, nil
}

func encodeMetadata(metadata *treeMetaData) []byte {
  var data []byte

  data = append(data, encodeUint16(metadata.order)...)
  data = append(data, encodeUint32(metadata.rootId)...)
  data = append(data, encodeUint16(metadata.pageSize)...)

  return data
}

func decodeMetadata(data []byte) (*treeMetaData, error) {
  order := decodeUint16(data[0:2])
  rootId := decodeUint32(data[2:6])
  pageSize := decodeUint16(data[6:8])

  metadata := &treeMetaData{
    order,
    rootId,
    pageSize,
  }

  return metadata, nil
}

func encodeStorageMetadata(md *storageMetadata) []byte {
  var data []byte

  data = append(data, encodeUint16(md.pageSize)...)
  data = append(data, encodeUint32(md.lastPageId)...)
  data = append(data, md.custom...)

  return data
}

func decodeStorageMetadata(data []byte) *storageMetadata {
  pageSize := decodeUint16(data[0:2])
  lastPageId := decodeUint32(data[2:6])
  custom := data[6:]

  metadata := &storageMetadata{
    pageSize,
    lastPageId,
    custom,
  }

  return metadata
}
