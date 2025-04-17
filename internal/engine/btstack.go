/*
* The stack is only required for Deletion since
* Insertion does not require backward traversal.
*/
package engine

import "fmt"

type BTStack struct {
  stack []uint32
  top int
  rootId uint32
}

func InitializeStack(rootId uint32) (*BTStack) {
  retStack := &BTStack{
    stack : []uint32{rootId},
    top : 0,
    rootId : rootId,
  }

  return retStack
}

func (b *BTStack) updateZeroPointer(rootId uint32) {
  b.rootId = rootId
  b.stack[0] = rootId
  return
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

// Returns parent id, ParentIsRoot.
// ParentIsRoot is true if parent is root or the stack top has only root.
// it is false if parent is not root.
func (b *BTStack) getParent(nodeId uint32) (uint32, error) {
  for i,v := range b.stack {
    if v == nodeId {

      if i-1 < 0 {
        return b.rootId, nil
      }
      return b.stack[i-1], nil
    }
  }
  return b.rootId, fmt.Errorf("The given node ID does not exist on stack.")
}

func (b *BTStack) emptyStack() {
  b.stack = []uint32{b.rootId}
  b.top = 0

  return
}
