package test

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
)

func TestSuite(){
  tree := engine.CreateNewTree(10)

  for i:=0; i<1000000000; i++{
    start := time.Now()
    err := tree.Insert(i)
    if err != nil{
      fmt.Printf("error inserting %v because %v\n", i, err)
    }
   fmt.Printf("Time to execute: %V\n", time.Since(start))
  }

  err := tree.Insert(12381)
  if err != nil{
    fmt.Printf("Error inserting")
  }
  // tree.PrintTree()

  node, index, err := tree.SearchKey(12381)
  if err != nil {
    fmt.Println("error occured: ", err)
  } else {
    fmt.Println(node)
    fmt.Println(index)
  }
}
