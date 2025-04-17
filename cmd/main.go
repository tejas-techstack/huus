package main

import (
  "fmt"
  engine "github.com/tejas-techstack/huus/internal/engine"
)

func main(){

  tree, _ := engine.Open("./example.db", 10, 4096)

  err := engine.BeginQueryLoop(tree)
  if err != nil {
    fmt.Println("Error :", err)
    return
  }

}
