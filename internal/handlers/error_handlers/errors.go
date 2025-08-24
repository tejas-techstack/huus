package error_handlers

import "errors"

var (
  ErrNotFound = errors.New("Not found")
  ErrInvalidInput = errors.New("Invalid Input")
)
