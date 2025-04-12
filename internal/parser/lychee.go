package lychee

import (
  "strings"
  "strconv"
  "regexp"
  "fmt"
  "bufio"
  "os"
)

var insert_regex = regexp.MustCompile(`(?i)INSERT\s*\(([a-zA-Z0-9]+),\s*([a-zA-Z0-9]+)\)`)
var update_regex = regexp.MustCompile(`(?i)UPDATE\s*\(([a-zA-Z0-9]+),\s*([a-zA-Z0-9]+)\)`)
var read_regex   = regexp.MustCompile(`(?i)READ\s*\(([a-zA-Z0-9]+)\s*\)`)
var delete_regex = regexp.MustCompile(`(?i)DELETE\s*\(([a-zA-Z0-9]+)\s*\)`)
var exit_regex   = regexp.MustCompile(`(?i)EXIT`)
var print_regex  = regexp.MustCompile(`(?i)PRINT`)

func ParseInsert(line string) (int, int, error) {
  matches := insert_regex.FindStringSubmatch(line)
  key, err := strconv.Atoi(matches[1])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1,-1,nil
  }
  value,err := strconv.Atoi(matches[2])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1,-1,nil
  }
  return key, value, nil
}

func ParseUpdate(line string) (int, int, error) {
  matches := update_regex.FindStringSubmatch(line)
  key, err := strconv.Atoi(matches[1])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1,-1,nil
  }
  value,err := strconv.Atoi(matches[2])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1,-1, nil
  }
  return key, value, nil

}

func ParseRead(line string) (int, error) {
  matches := read_regex.FindStringSubmatch(line)
  key, err := strconv.Atoi(matches[1])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1, nil
  }
  return key, nil
}

func ParseDelete(line string) (int, error) {
  matches := delete_regex.FindStringSubmatch(line)
  key, err := strconv.Atoi(matches[1])
  if err != nil {
    fmt.Println("Incorrect input(s)")
    return -1, nil
  }
  return key, nil

}


func ParseLine(line string) (int, error){
  switch {
  case line == insert_regex.FindString(line):
    return 1, nil
  case line == update_regex.FindString(line):
    return 2, nil
  case line == read_regex.FindString(line):
    return 3, nil
  case line == delete_regex.FindString(line):
    return 4, nil
  case line == exit_regex.FindString(line):
    return 0, nil
  case line == print_regex.FindString(line):
    return 99, nil
  default:
    return -1, nil
  }

  return -1, nil
}

func ReadLine() (string, error) {
  fmt.Printf("lychee> ")
  reader := bufio.NewReader(os.Stdin)
  line, err := reader.ReadString('\n')
  if err != nil {
    return "", fmt.Errorf("Error occured while taking input : %w", err)
  }

  line = strings.TrimSuffix(line, "\n")
  return line, nil
}
