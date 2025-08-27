// This is the heart of the storage
// contains the b-tree algorithm
package index

import (
  types "github.com/tejas-techstack/huus/internal/engine/types"
  errs "github.com/tejas-techstack/huus/internal/handlers/error_handlers"
  page "github.com/tejas-techstack/huus/internal/engine/pager"
  cache "github.com/tejas-techstack/huus/internal/engine/storage/cache"
)

func Insert (key Key, val Value) error {
  // 1. find target page
  //    a. check if page in cache
  //    b. If not fault it and bring it into cache
  //    c. Once in cache get the reference here.
  pageToInsert, err := findTargetPage(key);
  if err == errs.ErrPageNotFound {
    continue
  } else {
    return err
  }
  
  // 2. Update the page
  // this updates the page in the cache only
  // Two cases:
  //    a. If key exists in page -> upsert
  //    b. If key does not exist in page -> insert
  err = pageToInsert.UpdatePage(val)
  if err != nil {
    return err
  }

  return nil
}

func Update (key Key, val Value) error {
  // 1. find target page 
  // if it DNE return error since update means update only existing
  pageToUpdate, err := findTargetPage(key);
  if err == errs.ErrPageNotFound {
    return errs.ErrKeyDNE
  } else {
    return err
  }

  // 2. update the page
  err = pageToUpdate.UpdatePage(val)
  if err != nil {
    return err
  }

  return nil
}

func Read (key Key) (Value, error) {
  // find target page 
  // if it DNE return error and Value as -1

  pageToRead, err := findTargetPage(key);
  if err != nil {
    return (-1), err
  }

  val, err := pageToRead.getVal(key);
  if err != nil {
    return (-1), err
  }

  return val, nil
}

func Delete (key Key) (Value, error) {
  // find target page 
  pageToDeleteFrom, err := findTargetPage(key);
  if err != nil {
    return (-1), err
  }

  val, err := pageToDeleteFrom.DeleteKey(key)
  if err != nil {
    return (-1), err
  }

  return val, nil
}
