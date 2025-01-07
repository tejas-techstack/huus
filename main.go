package main

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
)

func main(){
  tree := engine.CreateNewTree(3)
  testString := "HELLO WORLD"
  start := time.Now()
  placeHolder := []byte(testString)

  for i := 1; i<10; i++{
    err := tree.Insert(i,placeHolder)
    if err != nil {
      fmt.Println(err)
    }
  }

  for i := 1; i<10; i++{
    _, key, valoff, err := tree.Get(i)
    if err != nil {
      fmt.Println(err)
    }

    fmt.Println("Key :", key, "Valoff:", valoff)
  }

  tree.Print()

  err := tree.Delete(1)
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println("Time taken : ", time.Since(start))
}
