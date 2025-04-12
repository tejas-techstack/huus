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
  stack.push(14)
  stack.push(16)
  fmt.Println(stack)

  x,err := stack.pop()
  fmt.Println(x, stack)

  _,err = stack.pop()
  _,err = stack.pop()
  _,err = stack.pop()
  _,err = stack.pop()

  if err != nil {
    t.Fatalf("Error occured : %s", err)
  }

  fmt.Println(stack)
  y := stack.showTop()
  fmt.Println(y, stack)
}
