package main

import (
  "strings"
  "regexp"
  "fmt"
  "bufio"
  "os"
)


var insert_regex = regexp.MustCompile(`(?i)INSERT\s*\(([a-zA-Z0-9]+),\s*([a-zA-Z0-9]+)\)`)
var update_regex = regexp.MustCompile(`(?i)UPDATE\s*\(([a-zA-Z0-9]+),\s*([a-zA-Z0-9]+)\)`)
var read_regex   = regexp.MustCompile(`(?i)READ\s*\(([a-zA-Z0-9]+)\s*\)`)
var delete_regex = regexp.MustCompile(`(?i)DELETE\s*\(([a-zA-Z0-9]+)\s*\)`)

func parseQuery(line string) (int, error){

  switch {
  case line == insert_regex.FindString(line):
    fmt.Println("insert query")
    matches := insert_regex.FindStringSubmatch(line)
    fmt.Printf("entered key:%v Value:%v\n", matches[1], matches[2])
  case line == update_regex.FindString(line):
    fmt.Println("Update query")
    matches := update_regex.FindStringSubmatch(line)
    fmt.Printf("entered key:%v, Value:%v\n", matches[1], matches[2])
  case line == read_regex.FindString(line):
    fmt.Println("Read query")
    matches := read_regex.FindStringSubmatch(line)
    fmt.Printf("entered key:%v", matches[1])
  case line == delete_regex.FindString(line):
    fmt.Println("Delete query")
    matches := delete_regex.FindStringSubmatch(line)
    fmt.Printf("entered key:%v", matches[1])
  default:
    fmt.Println("Not a valid query")
  }

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
      return fmt.Errorf("Invalid Query")
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

