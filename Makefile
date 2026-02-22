GO = go
MAIN_GO = cmd/main.go
DB_FILE = example.db

run:
	rm -rf $(DB_FILE)
	$(GO) run $(MAIN_GO)

clean:
	rm -rf $(DB_FILE)
