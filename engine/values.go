package engine

import (
  "os"
  "errors"
  // "bufio"
  "fmt"
)

const TestFile = "values.txt"

// my job is to write one line to a file.
// and then read that line.

func GetFreeBlock() (valueOffset, error) {

  // BUG: VERY VERY BAD METHOD OF DOING IT, fix this immediately

  f, err := os.OpenFile(TestFile, os.O_RDWR|os.O_CREATE, 0644)
  if err != nil {
    return -1, err
  }

  info, err := f.Stat()
  if err != nil{
    return -1, errors.New("Stat read error")
  }

  ret := info.Size()

  /*
  _, err = f.Seek(-1, 2) // Move to the last byte
	if err != nil {
		fmt.Println("Seek error:", err)
		return -1, err
	}

	// Read the last byte
	buf := make([]byte, 1)
	_, err = f.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return -1, err
	}
  */

	// fmt.Printf("Last byte: %q\n", buf[0])
  return valueOffset(int(int(ret)/blockSize)), nil
}

func WriteVal(data []byte, valOff valueOffset) (error){

  dataBlock := make([]byte, blockSize)
  copy(dataBlock[:], data)

  f,err := os.OpenFile(TestFile, os.O_RDWR|os.O_CREATE, 0644)
  if err != nil {
    return err
  }

  // seek to the offset.

  offset := int64(valOff * blockSize)
  // fmt.Println(offset)

  if _, err = f.Seek(offset, 0) ; err != nil {
    fmt.Println("here")
    return err
  }

  numberOfBytesWritten, err := f.Write(dataBlock) ; 
  if err != nil{
    return err
  }
  _ = numberOfBytesWritten

  return nil
}

func ReadVal(valOff valueOffset) ([]byte){
  return []byte{}
}
