package main

import (
  "strings"
  "fmt"
  "bufio"
  "os"
)

func parseQuery(line string) (int, error){
  // place holder code
  fmt.Println(line)
  return 1, nil
}

func readLine() (string, error) {
  fmt.Printf(">")
  reader := bufio.NewReader(os.Stdin)
  line, err := reader.ReadString('\n')
  if err != nil {
    return "", fmt.Errorf("Error occured while taking input : %w", err)
  }

  line = strings.TrimSuffix(line, "\n")
  return line, nil
}

func StartQueryLoop() error {
  line, err := readLine()
  if err != nil {
    return fmt.Errorf("Error occured while reading line : %w", err)
  }
  for line != "exit" {
    valid, err := parseQuery(line)
    if err != nil {
      return fmt.Errorf("Error occured while parsing query : %w", err)
    }
    if valid == -1 {
      fmt.Println("Invalid query")
    }

    line, err = readLine()
    if err != nil {
      return fmt.Errorf("Error occured while reading line : %w", err)
    }
  }

  // Print an exit statement to show end of query loop
  fmt.Println("Exiting")
  return nil
}

func main() {
  // begin the query loop
  err := StartQueryLoop()
  if err != nil {
    fmt.Println("Query loop exited with error: ", err)
  }

  return
}

