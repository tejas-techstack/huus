package test

import (
  "github.com/tejas-techstack/storageEngine/engine"
  "fmt"
  "time"
  "math/rand"
)

func getRandomNumber(n int) int {
  return rand.Intn(n)
}

func TestSuite(n, minNode, rs int){
  tree := engine.CreateNewTree(minNode)
  placeHolder := make([]byte, 0)
  start := time.Now()

  for i:=1; i<n; i++{
    num := i
    if rs == 0 {
      num = getRandomNumber(10000)
    }
    // fmt.Println("Inserting :", num)
    err := tree.Insert(num, placeHolder)
    if err != nil{
      i--
      continue
      // fmt.Println("error occured:", err)
    }
  }

  fmt.Printf("Time to execute: %v\n", time.Since(start))
  if n < 1000{
    tree.Print()
  }
  
  findNum := 9
  _, _, _, err := tree.Get(findNum)
  if err != nil {
    fmt.Println("error occured: ", err)
  } else {
    fmt.Println("Found :",findNum)
  }

  err = tree.Delete(9)
  if err != nil{
    fmt.Println("Error occured: ", err)
  }

  err = tree.Delete(10)
  err = tree.Delete(11)
  err = tree.Delete(12)

  err = tree.Insert(10, []byte{})
  err = tree.Insert(10, []byte{})
  if err != nil{
    fmt.Println("Error occured :",err)
  }
  
  if n < 1000{
    tree.Print()
  }

  findNum = 10
  _, _, _, err = tree.Get(findNum)
  if err != nil{
    fmt.Println("Error occured: ", err)
  } else {
    fmt.Println("Found ", findNum)
  }
}
