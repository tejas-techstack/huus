package test

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
)

func TestSuite(){
  tree := engine.CreateNewTree(10)
  placeHolder := make([]byte, 0)
  start := time.Now()

  for i:=0; i<10000000; i++{
    err := tree.Insert(i, placeHolder)
    if err != nil{
    }
  }

  err := tree.Insert(12, placeHolder)
  if err != nil{
    fmt.Println(err)
    fmt.Println("Error inserting")
  }

  // tree.Print()

  /*
  node, index, err := tree.SearchKey(999)
  if err != nil {
    fmt.Println("error occured: ", err)
  } else {
    fmt.Println(node)
    fmt.Println(index)
  }
  */

  fmt.Printf("Time to execute: %v\n", time.Since(start))
}
