package main

import (
  "os"
  "fmt"
  kv "github.com/tejas-techstack/huus/pkg/kv"
)

func main(){

  defer func() {
  err := os.Remove("./example.db")
  if err != nil {
    fmt.Println("Error occured while removing example.db")
  }
  fmt.Println("Deleted example.db")
  }()

  tree, _ := kv.Open("./example.db", 10, 4096)

  err := kv.BeginQueryLoop(tree)
  if err != nil {
    fmt.Println("Error :", err)
    return
  }

}
