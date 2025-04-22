package wal

import (
  "fmt"
  "os"
)

const logFile = "./wal.txt"

type logger struct {
  fo *os.File
}

func OpenLogger() (*logger, error){
  fo, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
  if err != nil {
    return nil,fmt.Errorf("Error opening log file")
  }

  return &logger{fo}, nil
}

func (l *logger) LogQuery(query string) error {
  if string(query[0]) == "." {
    return nil
  }
  query += "\n"
  n, err := l.fo.WriteString(query)
  if err != nil {
    return fmt.Errorf("Error appending to log file : %w", err)
  }
  if n != len(query) {
    return fmt.Errorf("Had to write %v, Wrote %v to the log file", len(query), n)
  }

  return nil
}
