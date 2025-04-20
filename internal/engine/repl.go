/*
* Interacts with the parser to evaluate queries
*/

package engine

import (
  "fmt"
  lychee "github.com/tejas-techstack/huus/internal/parser"
  wal "github.com/tejas-techstack/huus/internal/wal"
)

// invalid  : -1
// exit   : 0
// insert : 1
// update : 2
// read   : 3
// delete : 4
// print  : 99 (debugging only will remove.)

func BeginQueryLoop(tree *BPTree) error {

  fmt.Println(tree)

  // create a logger object for wal.
  lo, err := wal.OpenLogger()
  if err != nil {
    return fmt.Errorf("Error opening log file : %w", err)
  }

  for true {
    line, err := lychee.ReadLine()
    if err != nil {
      return fmt.Errorf("Error in query loop: %w", err)
    }

    query, err := lychee.ParseLine(line)
    if err != nil {
      return fmt.Errorf("Error parsing line: %w", err)
    }

    if err = lo.LogQuery(line); err != nil {
      return fmt.Errorf("Error logging query : %w", err)
    }

    switch query {
    case -1 :
      fmt.Println("Invalid Query")
    case 0:
      fmt.Println("Exiting")
      return nil
    case 1:
      key, value, err := lychee.ParseInsert(line)
      if err != nil {
        return fmt.Errorf("Error parsing insert query : %w", err)
      }

      err = tree.PutInt(key, value)
      if err != nil {
        return fmt.Errorf("Error putting value : %w", err)
      }

      fmt.Println("Inserted successfully")
    case 2:
      _, _, err := lychee.ParseUpdate(line)
      if err != nil {
        return fmt.Errorf("Error parsing update query : %w", err)
      }
      fmt.Println("WHAT DO YOU WANT FROM ME ;;")
    case 3:
      key, err := lychee.ParseRead(line)
      if err != nil {
        return fmt.Errorf("Error parsing read query : %w", err)
      }

      value, exists, err := tree.GetInt(key)
      if err != nil {
        return fmt.Errorf("Error getting value : %w", err)
      }

      if !exists {
        fmt.Println("Key does not exist")
      } else {
        fmt.Printf("Key:%v, Value:%v\n" ,key, value)
      }
    case 4:
      key, err := lychee.ParseDelete(line)
      if err != nil {
        return fmt.Errorf("Error parsing delete query : %w", err)
      }

      deleted, err := tree.DeleteInt(key)
      if err != nil {
        return fmt.Errorf("Error deleting key : %w", err)
      }

      if !deleted {
        fmt.Println("Key does not exist")
      }
      fmt.Println("Deleted successfully")
    case 99:
      printTree(tree)
    }
  }

  return nil
}
