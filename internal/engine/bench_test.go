/*
*
* This is to test an UNBOUNDED cache and notice its benefits
*
*/
package engine

import (
  "path/filepath"
  "testing"
  "time"
)

func TestCache(t *testing.T) {
  dir := t.TempDir()
  path := filepath.Join(dir, "bench.db")

  tree, err := Open(path, 10, 4096)
  if err != nil {
    t.Fatalf("open: %v", err)
  }

  const n = 5000
  start := time.Now()
  for i := 0; i < n; i++ {
    if err := tree.PutInt(i, i*2); err != nil {
      t.Fatalf("put %d: %v", i, err)
    }
  }
  elapsed := time.Since(start)
  ops := float64(n) / elapsed.Seconds()
  t.Logf("inserted %d keys in %s => %.0f ops/sec", n, elapsed, ops)

  for i := 0; i < n; i++ {
    val, ok, err := tree.GetInt(i)
    if err != nil {
      t.Fatalf("get %d: %v", i, err)
    }
    if !ok {
      t.Fatalf("key %d missing after insert", i)
    }
    if got := decodeUint64(val); got != i*2 {
      t.Fatalf("key %d: got %d want %d", i, got, i*2)
    }
  }

  if err := tree.Close(); err != nil {
    t.Fatalf("close: %v", err)
  }

  reopened, err := Open(path, 10, 4096)
  if err != nil {
    t.Fatalf("reopen: %v", err)
  }
  defer reopened.Close()

  for i := 0; i < n; i++ {
    val, ok, err := reopened.GetInt(i)
    if err != nil {
      t.Fatalf("get %d after reopen: %v", i, err)
    }
    if !ok {
      t.Fatalf("key %d missing after reopen", i)
    }
    if got := decodeUint64(val); got != i*2 {
      t.Fatalf("key %d after reopen: got %d want %d", i, got, i*2)
    }
  }
}
