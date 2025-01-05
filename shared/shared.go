package shared

const (
  testFile = "test.txt"  
  blockSize = 4096
)

type valueOffset int

// struct to hold kv pair
type KV struct {
  key int
  valOff valueOffset
}

// struct for each node of the b+tree
// TODO: instead of nextLeaf, make the last element of <Node>.children as the nextLeaf
type Node struct{
  kvStore []KV
  children []*Node
  isLeaf bool
  nextLeaf *Node 
}

// B+ tree struct
type BPtree struct {
  root *Node
  minNode int
}
