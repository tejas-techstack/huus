/*
*
* Test Scripts written by AI and audited.
*
*/

// benchmark/main.go
//
// Standalone throughput / memory / disk benchmark for the huus engine.
// Run from the repo root:
//
//   go run ./benchmark
//   go run ./benchmark -n 100000        # bigger run
//   go run ./benchmark -order 200       # try a saner tree order
//
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	engine "github.com/tejas-techstack/huus/internal/engine"
)

func main() {
	n := flag.Int("n", 5000, "number of keys to insert/delete")
	order := flag.Int("order", 10, "tree order (keys per node)")
	pageSize := flag.Int("page", 4096, "page size in bytes")
	path := flag.String("db", "./benchmark.db", "database file path")
	flag.Parse()

	// start clean so disk numbers are honest.
	os.Remove(*path)
	defer os.Remove(*path)

	fmt.Printf("== huus benchmark ==\n")
	fmt.Printf("keys=%d  order=%d  pageSize=%d  db=%s\n\n", *n, *order, *pageSize, *path)

	tree, err := engine.Open(*path, uint16(*order), uint16(*pageSize))
	must(err)

	baselineRAM := heapMB()

	// ---- INSERT ----
	start := time.Now()
	for i := 0; i < *n; i++ {
		must(tree.PutInt(i, i*2))
	}
	must(tree.Sync())
	insElapsed := time.Since(start)
	report("INSERT", *n, insElapsed)
	fmt.Printf("  heap=%.1f MB (Δ %.1f MB)   disk=%.1f MB   bytes/key=%.0f\n\n",
		heapMB(), heapMB()-baselineRAM, fileMB(*path), fileMB(*path)*1e6/float64(*n))

	// ---- READ-BACK VERIFY ----
	start = time.Now()
	for i := 0; i < *n; i++ {
		val, ok, err := tree.GetInt(i)
		must(err)
		if !ok {
			fatalf("key %d missing after insert", i)
		}
		if got := int(binary.BigEndian.Uint64(val)); got != i*2 {
			fatalf("key %d: got %d want %d", i, got, i*2)
		}
	}
	report("READ", *n, time.Since(start))
	fmt.Println()

	// ---- DELETE ----
	start = time.Now()
	for i := 0; i < *n; i++ {
		ok, err := tree.DeleteInt(i)
		must(err)
		if !ok {
			fatalf("delete reported key %d missing", i)
		}
	}
	must(tree.Sync())
	report("DELETE", *n, time.Since(start))
	fmt.Printf("  heap=%.1f MB   disk=%.1f MB (note: deletes don't shrink the file)\n\n",
		heapMB(), fileMB(*path))

	// verify they're actually gone.
	for i := 0; i < *n; i++ {
		_, ok, err := tree.GetInt(i)
		must(err)
		if ok {
			fatalf("key %d still present after delete", i)
		}
	}
	fmt.Println("DELETE verify: all keys gone  ✓")

	// ---- PERSISTENCE (close / reopen) ----
	for i := 0; i < *n; i++ {
		must(tree.PutInt(i, i*2))
	}
	must(tree.Close())

	reopened, err := engine.Open(*path, uint16(*order), uint16(*pageSize))
	must(err)
	defer reopened.Close()

	missing := 0
	for i := 0; i < *n; i++ {
		val, ok, err := reopened.GetInt(i)
		must(err)
		if !ok || int(binary.BigEndian.Uint64(val)) != i*2 {
			missing++
		}
	}
	if missing == 0 {
		fmt.Printf("PERSISTENCE: re-inserted %d keys, closed, reopened, all readable  ✓\n", *n)
	} else {
		fatalf("PERSISTENCE: %d/%d keys wrong after reopen", missing, *n)
	}
}

func report(label string, n int, d time.Duration) {
	fmt.Printf("%-7s %d ops in %-12s => %.0f ops/sec  (%.1f µs/op)\n",
		label, n, d.Round(time.Microsecond), float64(n)/d.Seconds(),
		float64(d.Microseconds())/float64(n))
}

// heapMB returns current Go heap usage in MB (forces a GC first so the
// number reflects live data, e.g. the engine's in-memory page cache).
func heapMB() float64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.HeapAlloc) / 1e6
}

func fileMB(path string) float64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return float64(fi.Size()) / 1e6
}

func must(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", a...)
	os.Exit(1)
}
