package kv

import (
  "fmt"
  lychee "github.com/tejas-techstack/huus/parser"
)

func BeginQueryLoop() error {
  err := lychee.StartQueryLoop()
  if err != nil {
    fmt.Println("Query loop exited with error: ", err)
    return err
  }

  return nil
}
