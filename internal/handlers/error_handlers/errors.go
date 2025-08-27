package error_handlers

import "errors"

var (
  ErrNotFound = errors.New("Not found")
  ErrInvalidInput = errors.New("Invalid Input")
  ErrPageNotFound = errors.New("Page not found")
  ErrKeyDNE = errors.New("Key does not exist")
)
