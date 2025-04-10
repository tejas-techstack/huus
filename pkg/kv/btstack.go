package kv

import "fmt"

type BTStack struct {
  stack []uint32
  top int
  rootId uint32
}

func InitializeStack(rootId uint32) (*BTStack) {
  retStack := &BTStack{
    stack : []uint32{rootId},
    top : -1,
    rootId : rootId,
  }

  return retStack
}

// Takes nodeID as input and pushes to top of stack.
func (b *BTStack) push(nodeID uint32) {
  b.stack = append(b.stack, nodeID)
  b.top++
  return
}

// Removes the top of the stack and returns it.
// RETURNS ERROR IF TRYING TO POP ROOT.
func (b *BTStack) pop() (uint32, error) {
  if b.top == -1 {
    return b.rootId, fmt.Errorf("Trying to pop zeroth stack index")
  }

  retVal := b.stack[b.top]
  b.stack = append([]uint32{}, b.stack[:b.top]...)
  b.top--

  return retVal, nil
}

// Returns the top of the stack without removing.
func (b *BTStack) showTop() uint32 {
  return b.stack[b.top]
}

