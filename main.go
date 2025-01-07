package main

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
  "strconv"
)

func main(){
  tree := engine.CreateNewTree(3)
  start := time.Now()

  for i := 1; i<10; i++{
    testString := "This is test string " + strconv.Itoa(i)
    placeHolder := []byte(testString)
    err := tree.Insert(i,placeHolder)
    if err != nil {
      fmt.Println(err)
    }
  }

  for i := 1; i<10; i++{
    _, key, value, err := tree.Get(i)
    if err != nil {
      fmt.Println(err)
    }

    fmt.Println("Key :", key, "Value:", string(value))
  }

  tree.Print()

  fmt.Println("Time taken : ", time.Since(start))
}
