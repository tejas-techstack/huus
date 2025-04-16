package kv

import (
  "testing"
  "fmt"
)

func TestBTStack (t *testing.T) {
  // Initialize stack with stack pointing to root first.
  stack := InitializeStack(10)
  fmt.Println(stack)

  stack.push(12)
  fmt.Println(stack)

  x, _ := stack.pop()
  fmt.Println(x)

  fmt.Println(stack.showTop())

  stack.push(12)
  fmt.Println(stack.getParent(12))
  
  fmt.Println(stack.getParent(10))
}
