package main

import (
  "fmt"
  "os"
  engine "github.com/tejas-techstack/huus/internal/engine"
)

func main(){

  tree, err := engine.Open("./example.db", 10, 4096)
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
  }

  err = engine.BeginQueryLoop(tree)
  if err != nil {
    fmt.Println("Error :", err)
    return
  }

}
