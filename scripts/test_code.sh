#!/bin/bash

mkdir -p logs

log="logs/logs_$(date +%Y%m%d_%H%M%S).txt"

# Run 5000
go run benchmark/main.go >"$log" 2>&1
cat "$log"
