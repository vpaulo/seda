# Seda Programming Language

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A modern scripting language combining the best of Ruby, Lua, and TypeScript—designed for simplicity, expressiveness, and developer productivity.

> **Note:** Seda uses the `.s` file extension.

## Features

- **Everything is an Object** - Ruby-inspired object orientation with prototype-based extensions
- **Clean Syntax** - Lua-like simplicity with minimal boilerplate
- **Optional Static Typing** - TypeScript-style type annotations and inference
- **Built-in Testing** - First-class testing support with `check::` and `where::` blocks
- **Module System** - Standard library, third-party packages, and local modules
- **Dynamic Extension** - Add custom properties and methods to any type at runtime
- **Functional Features** - First-class functions, closures, and higher-order functions

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git (for package management)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd seda

# Build the interpreter and tools
go build -o seda ./cmd/parser

# Add to PATH (optional)
export PATH=$PATH:$(pwd)
```

## Project Structure

```
seda/
├── ast/              # Abstract Syntax Tree definitions
├── evaluator/        # Runtime evaluation and execution
├── lexer/            # Lexical analysis (tokenization)
├── parser/           # Syntax analysis (parsing)
├── object/           # Object system and types
├── pkg/              # Package management
├── cmd/
│   └── parser/       # Main interpreter
├── examples/         # Example programs and tutorials
└── docs/             # Language specifications

Standard Library:
~/.seda/std/          # Standard library modules
~/.seda/packages/     # Third-party packages cache
```

## Package Management

```bash
# Install a third-party package
./seda install github.com/user/awesome-lib

# List installed packages
./seda list

# Update a package
./seda update awesome-lib

# Remove a package
./seda remove awesome-lib
```

## Development

### Building from Source

```bash
# Build all binaries
go build ./cmd/...

# Run tests
go test ./...

# Build specific components
go build -o seda ./cmd/parser
```

### Running Tests

```bash
# Run Go tests
go test ./evaluator
go test ./parser
go test ./lexer

# Run language tests
./seda examples/basic.s
./seda examples/custom_properties.s
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Inspired by:
- **Ruby** - Object-oriented philosophy and elegant syntax
- **Lua** and **Pyret** - Simplicity and embeddability
- **TypeScript** - Optional type system and developer experience