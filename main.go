package main

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
)

func main(){
  // TODO write a test suite

  start := time.Now()

  tree := engine.CreateNewTree(10)

  for i:=0; i<1000000000; i++{
    err := tree.Insert(i)
    if err != nil{
      fmt.Printf("error inserting %v because %v\n", i, err)
    }
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

  fmt.Printf("Time to execute: ", time.Since(start))
}
