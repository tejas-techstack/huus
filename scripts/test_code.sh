#!/bin/bash

mkdir -p logs

name="${1:-logs_$(date +%Y%m%d_%H%M%S)}"
log="logs/${name}.txt"

# Run 5000
go run benchmark/main.go >"$log" 2>&1
cat "$log"
