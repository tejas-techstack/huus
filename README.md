# Huus

[![Go Version](https://img.shields.io/badge/Go-1.22.5+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A persistent, B+ tree based key-value store written in Go. Huus demonstrates modern B+ tree implementation with optimizations using Go routines, providing a foundation for understanding database internals.

## Features

- **Persistent Storage**: Data is stored on disk and survives restarts
- **B+ Tree Implementation**: Efficient search, insertion, and deletion operations
- **REPL Interface**: Interactive command-line interface for database operations
- **Configurable**: Adjustable tree order and page size for performance tuning
- **Thread-Safe Operations**: Optimized with Go routines
- **Comprehensive Testing**: Well-tested with unit and integration tests

## Table of Contents

- [Background](#background)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Background

### What are B+ Trees?

B-trees were first developed to enable memory storage devices like hard disks to handle search queries efficiently without needing to read from multiple locations. The architecture of a B-tree allows storage in a block-oriented format, which is ideal for disk-based systems.

**B+ trees** are a variation of B-trees that optimize for:
- **Range Queries**: All values are stored only in leaf nodes
- **Sequential Access**: Leaf nodes are linked together forming a chain
- **Better Disk I/O**: Non-leaf nodes only store keys, allowing more keys per node

This makes B+ trees particularly well-suited for database systems and file systems where range scans and sequential access are common operations.

## Installation

### Prerequisites

- Go 1.22.5 or higher
- Make (optional, for using Makefile commands)

### Install

```bash
# Clone the repository
git clone https://github.com/tejas-techstack/huus.git
cd huus

# Build the project
go build -o huus cmd/main.go
```

## Quick Start

```bash
# Run the database REPL
make run

# Or run directly with go
go run cmd/main.go
```

Once the REPL starts, you can begin executing queries:

```
> INSERT (1, 100)
Inserted successfully
> READ (1)
Key: 1, Value: 100
> UPDATE (1, 200)
Updated successfully
> DELETE (1)
Deleted successfully
> EXIT
Exiting
```

## Usage

### Starting the Database

```bash
make run
```

This will:
1. Create or open a database file named `example.db`
2. Initialize a B+ tree with default order 10 and page size 4096 bytes
3. Start an interactive REPL for query execution

### Available Commands

All commands are **case-insensitive**:

| Command | Syntax | Description | Example |
|---------|--------|-------------|---------|
| **INSERT** | `INSERT (key, value)` | Insert a new key-value pair | `INSERT (42, 100)` |
| **UPDATE** | `UPDATE (key, value)` | Update an existing key's value | `UPDATE (42, 200)` |
| **READ** | `READ (key)` | Retrieve the value for a key | `READ (42)` |
| **DELETE** | `DELETE (key)` | Remove a key-value pair | `DELETE (42)` |
| **EXIT** | `EXIT` | Exit the REPL | `EXIT` |

**Note**: Currently, both keys and values must be integers.

### Programmatic Usage

You can also use Huus as a library in your Go programs:

```go
package main

import (
    "fmt"
    "github.com/tejas-techstack/huus/internal/engine"
)

func main() {
    // Open a database with order=10 and pageSize=4096
    tree, err := engine.Open("./mydata.db", 10, 4096)
    if err != nil {
        panic(err)
    }

    // Insert a key-value pair
    err = tree.PutInt(1, 100)
    if err != nil {
        panic(err)
    }

    // Read a value
    value, exists, err := tree.GetInt(1)
    if err != nil {
        panic(err)
    }
    
    if exists {
        fmt.Printf("Found value: %v\n", value)
    }

    // Delete a key
    deleted, err := tree.DeleteInt(1)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Deleted: %v\n", deleted)
}
```

### Configuration

When opening a database, you can configure:

- **Order** (default: 10): Maximum number of children per node. Higher values mean fewer levels but larger nodes.
- **Page Size** (default: 4096 bytes): Size of each disk page. Must be between 32 and 4096 bytes.

```go
tree, err := engine.Open("./data.db", order, pageSize)
```

### Cleaning the Database

```bash
make clean
```

This removes the `example.db` file, effectively deleting all stored data.

## Architecture

Huus is built with a clean architecture separating concerns:

### Core Components

1. **B+ Tree Engine** (`internal/engine`)
   - Core data structure implementation
   - Node splitting and merging logic
   - CRUD operations (Create, Read, Update, Delete)
   - Persistence layer

2. **Storage Layer** (`internal/engine/storage.go`)
   - Page-based disk I/O
   - Node serialization/deserialization
   - Metadata management

3. **Parser** (`internal/parser`)
   - Query parsing (Lychee parser)
   - Command-line input handling
   - REPL implementation

### Key Design Decisions

- **Page-Based Storage**: Data is organized into fixed-size pages for efficient disk I/O
- **Encoding**: Custom binary encoding for nodes and metadata
- **Stack-Based Traversal**: Uses a stack to track the path during tree operations
- **Linked Leaf Nodes**: Leaf nodes are linked for efficient range scans

## Project Structure

```
huus/
├── cmd/
│   └── main.go              # Entry point for the application
├── internal/
│   ├── engine/              # Core B+ tree implementation
│   │   ├── btstack.go       # Stack for tree traversal
│   │   ├── config.go        # Configuration and initialization
│   │   ├── crud.go          # CRUD operations (public API)
│   │   ├── delete_helpers.go # Deletion logic
│   │   ├── encoding.go      # Serialization/deserialization
│   │   ├── helpers.go       # Utility functions
│   │   ├── insert_helpers.go # Insertion logic
│   │   ├── pager.go         # Page management
│   │   ├── printTree.go     # Debugging utilities
│   │   ├── repl.go          # REPL interface
│   │   ├── split.go         # Node splitting logic
│   │   ├── storage.go       # Persistent storage layer
│   │   ├── types.go         # Data structures
│   │   └── *_test.go        # Unit tests
│   └── parser/
│       └── lychee.go        # Query parser
├── go.mod                   # Go module definition
├── Makefile                 # Build and run commands
├── LICENSE                  # MIT License
└── README.md                # This file
```

## Development

### Building

```bash
# Build the binary
go build -o huus cmd/main.go

# Run the binary
./huus
```

### Code Organization

- **Public API**: Exposed through `crud.go` - `Get`, `Put`, `Delete` methods
- **Internal Logic**: Helper functions for tree operations are kept internal
- **Tests**: Co-located with implementation files using `*_test.go` naming

## Testing

Huus includes comprehensive tests covering:

- Node encoding/decoding
- Tree operations (insert, search, delete)
- Node splitting and merging
- Storage and persistence
- Stack operations

### Running Tests

```bash
# Run all tests
go test ./...

```

### Test Coverage

The project includes tests for:
- B+ Tree stack operations
- Node encoding/decoding
- Insert, update, and delete operations
- Tree splitting on overflow
- Storage layer operations
- Metadata persistence


## Contributing

Contributions are welcome! This is a learning project, and improvements are encouraged.

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines

- Write tests for new features
- Follow Go conventions and idioms
- Update documentation for API changes
- Keep commits focused and descriptive

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- B+ Tree design principles from database systems literature
- Inspired by real-world database implementations
- Built with Go's excellent standard library

## References

- [B-tree - Wikipedia](https://en.wikipedia.org/wiki/B-tree)
- [B+ Tree - Wikipedia](https://en.wikipedia.org/wiki/B%2B_tree)
- [Database Internals by Alex Petrov](https://www.oreilly.com/library/view/database-internals/9781492040347/)

---

**Note**: This is an educational project designed to demonstrate B+ tree data structures and database internals.
