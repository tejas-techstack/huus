package main

import (
  "fmt"
  kv "github.com/tejas-techstack/huus/pkg/kv"
)

func main(){

  tree, _ := kv.Open("./example.db", 10, 4096)

  err := kv.BeginQueryLoop(tree)
  if err != nil {
    fmt.Println("Error :", err)
    return
  }

}
