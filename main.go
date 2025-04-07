package main

import (
  "fmt"
  kv "github.com/tejas-techstack/huus/kv"
)

func main(){
  tree, _ := kv.Open("./example.db", 10, 4096)
  for i := 1; i < 5; i++{
    _ = tree.PutInt(i, i)
  }

  err := kv.BeginQueryLoop()
  if err != nil {
    fmt.Println("Error :", err)
    return
  }
}
