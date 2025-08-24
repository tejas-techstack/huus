// This is the heart of the storage
// contains the b-tree algorithm
package index

import (
  types "github.com/tejas-techstack/huus/internal/engine/types"
  errs "github.com/tejas-techstack/huus/internal/handlers/error_handlers"
)

func Insert (key Key, val Value) error {}

func Update (key Key, val Value) error {}

func Read (key Key) (Value, error) {}

func Delete (key Key) (Value, error) {}
