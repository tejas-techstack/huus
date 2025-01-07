package test

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "testing"
  // "math/rand"
)

/*
func TestInsert(t *testing.T) {
  _ = engine.CreateNewTree(10)
  fmt.Println("Hello")
  fmt.Println("Hello")
}
*/

func BenchmarkInsert(b *testing.B){
  _ = engine.CreateNewTree(10)
  fmt.Println("Hello")
}
