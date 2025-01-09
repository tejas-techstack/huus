package test

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "testing"
  // "math/rand"
)

func TestInsert(t *testing.T) {
  testString := "hello world"
  tree := engine.CreateNewTree(5)
  err := tree.Insert(1, []byte(testString))
  if err != nil {
    t.Errorf("Insert Error : %v", err)
  }
}

func BenchmarkInsert(b *testing.B){
  testString := "Hello world"
  tree := engine.CreateNewTree(5)
  for i := range b.N {
    err := tree.Insert(i, []byte(testString))
    if err != nil{
      b.Errorf("Insert Error : %v", err)
    }
  }

  // b.Log("Time to insert:", b.Elapsed())

  for i := range b.N{
    _, _, _, err := tree.Get(i)
    if err != nil{
      b.Errorf("Get Error : %v", err)
    }
    // b.Logf("Key is :%v, Value is :%v", key, value)
  }
}

