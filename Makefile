GO = go
MAIN_GO = cmd/main.go
DB_FILE = example.db
WAL_FILE = wal.txt

run:
	$(GO) run $(MAIN_GO)

clean:
	rm -rf $(DB_FILE)
	rm -rf $(WAL_FILE)
