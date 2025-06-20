# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI application named "eddie" built using the Cobra framework. The project is currently in its initial state with basic Cobra scaffolding.

## Development Commands

### Building the Application
```bash
go build -o eddie
```

### Running the Application
```bash
# Run directly with go
go run main.go

# Or build and run the binary
go build -o eddie && ./eddie
```

### Standard Go Commands
```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Tidy dependencies
go mod tidy
```

## Project Structure

- `main.go` - Entry point that calls `cmd.Execute()`
- `cmd/root.go` - Root Cobra command definition with basic CLI setup
- `go.mod` - Go module file defining dependencies (Cobra CLI framework)

## Architecture

The application follows the standard Cobra CLI pattern:
- `main.go` serves as the entry point
- `cmd/` package contains command definitions
- Root command is defined in `cmd/root.go` with placeholder descriptions

The application is licensed under GNU AGPL v3.