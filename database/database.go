package main

import (
  "os"
  "bufio"
  "github.com/tejas-techstack/storageEngine/shared"
)

// my job is to write one line to a file.
// and then read that line.

func WriteVal([]byte, valOff shared.ValueOffset) (error){

  f,err := os.Openfile(shared.testFile, os.O_RDWR|os.O_CREATE, 0644)
  if err != nil {
    return err
  }

  // seek to the offset.

  offset := valOff * shared.blockSize

  if err = f.Seek(valOff*s, 0) ; err != nil {
    return err
  }



}

func ReadVal(valOff shared.ValueOffset) ([]byte){
  return []byte{}
}
