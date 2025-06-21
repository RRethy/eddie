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
# Run tests (unit and integration)
go test ./...

# Run only unit tests
go test ./internal/...

# Run only e2e tests
go test ./test/e2e/...

# Run benchmarks
go test -bench=. ./internal/...

# Run specific benchmark
go test -bench=BenchmarkViewer ./internal/cmd/view/...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./internal/...

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
- `cmd/view.go` - View command definition for file/directory inspection
- `cmd/str_replace.go` - String replacement command definition
- `cmd/create.go` - File creation command definition
- `cmd/insert.go` - Line insertion command definition
- `cmd/undo_edit.go` - Undo edit command definition
- `internal/cmd/view/` - Business logic for the view command
- `internal/cmd/str_replace/` - Business logic for the str_replace command
- `internal/cmd/create/` - Business logic for the create command
- `internal/cmd/insert/` - Business logic for the insert command
- `internal/cmd/undo_edit/` - Business logic for the undo_edit command
- `test/e2e/` - End-to-end tests that test the CLI binary
- `go.mod` - Go module file defining dependencies (Cobra CLI framework)

## Architecture

The application follows the standard Cobra CLI pattern:
- `main.go` serves as the entry point
- `cmd/` package contains command definitions, only parses flags and arguments and calls internal logic
- `internal/` package contains business logic
- `internal/cmd/` contains business logic for specific commands

## Current Commands

### view
Examine the contents of a file or list the contents of a directory. It can read the entire file or a specific range of lines.

**Usage:**
```bash
eddie view path [view_range]
```

**Parameters:**
- `path`: The path to the file or directory to view
- `[view_range]`: (Optional) Range of lines to view in format "start,end". If "end" is -1, reads to end of file. Ignored for directories.

**Examples:**
```bash
eddie view /path/to/file.txt
eddie view /path/to/directory  
eddie view /path/to/file.txt 10,20
```

### str_replace
Replace all occurrences of a string in a file with another string.

**Usage:**
```bash
eddie str_replace path old_str new_str
```

**Parameters:**
- `path`: The path to the file to modify
- `old_str`: The string to search for and replace
- `new_str`: The string to replace old_str with

**Examples:**
```bash
eddie str_replace /path/to/file.txt "old text" "new text"
eddie str_replace config.json "localhost" "example.com"
```

### create
Create a new file with the specified content.

**Usage:**
```bash
eddie create path file_text
```

**Parameters:**
- `path`: The path where the new file should be created
- `file_text`: The content to write to the new file

**Examples:**
```bash
eddie create /path/to/newfile.txt "Hello, World!"
eddie create config.json '{"key": "value"}'
eddie create script.sh "#!/bin/bash\necho 'Hello'"
```

### insert
Insert a new line at the specified line number in a file.

**Usage:**
```bash
eddie insert path insert_line new_str
```

**Parameters:**
- `path`: The path to the file to modify
- `insert_line`: The line number where the new line should be inserted (1-based)
- `new_str`: The content of the new line to insert

**Examples:**
```bash
eddie insert /path/to/file.txt 5 "This is a new line"
eddie insert config.json 10 "  \"newKey\": \"newValue\","
eddie insert script.sh 1 "#!/bin/bash"
```

### undo_edit
Undo the last edit operation on a file by restoring from backup.

**Usage:**
```bash
eddie undo_edit path
```

**Parameters:**
- `path`: The path to the file to restore from backup

**Examples:**
```bash
eddie undo_edit /path/to/file.txt
eddie undo_edit config.json
eddie undo_edit script.sh
```

**Note:** This command automatically records edit operations when using `str_replace` or `insert` commands. Edit records are stored in `$XDG_CACHE_HOME/eddie/edits` (or `~/.cache/eddie/edits` if `XDG_CACHE_HOME` is not set). It reverses the most recent edit and validates that the file hasn't been modified by other means since the last tracked edit.

# Development Guidelines

## Go Code Style - Rob Pike Style (MANDATORY)

### Rob Pike's Philosophy
**"Simplicity is the ultimate sophistication"** - Write code as if the person maintaining it is a violent psychopath who knows where you live.

### Core Principles (NON-NEGOTIABLE)
- **Clarity above all** - If it's not immediately obvious, rewrite it
- **No clever code** - Clever code is bug-prone code
- **Shorter is better** - But not at the expense of clarity
- **Gofmt is gospel** - Never commit unformatted code
- **One thing per function** - Functions should do exactly one thing well
- **Fail fast and explicitly** - No hidden control flow, no magic
- **NO COMMENTS** - Code should be self-explanatory. Only add comments for exported functions/types or truly complex algorithms that cannot be simplified

### Rob Pike Naming Rules (STRICT)
- **Package names**: Single word, lowercase, no plurals (`net` not `networks`)
- **Variables**: SHORT. `i` not `index`, `n` not `numberOfItems`, `s` not `inputString`
- **Functions**: Descriptive verbs. `Get`, `Set`, `Read`, `Write` - not `Retrieve` or `Obtain`
- **No stuttering**: `log.Print()` not `log.LogPrint()`
- **Constants**: CamelCase, not SCREAMING_SNAKE_CASE
- **Receivers**: Single letter that makes sense (`c *Client`, `r *Reader`)

### Rob Pike Function Rules
- **Functions ≤ 30 lines** - If longer, split it up
- **No deep nesting** - Use early returns religiously
- **Error handling first** - Check errors immediately, don't defer
- **No side effects** - Functions should be predictable
- **Return multiple values** - Don't create structs just to return two things

```go
// GOOD - Rob Pike style
func read(r io.Reader) ([]byte, error) {
    buf := make([]byte, 1024)
    n, err := r.Read(buf)
    if err != nil {
        return nil, err
    }
    return buf[:n], nil
}

// BAD - Not Pike style
func readWithComplexity(reader io.Reader) *Result {
    result := &Result{}
    if reader != nil {
        buffer := make([]byte, 1024)
        if bytesRead, readError := reader.Read(buffer); readError == nil {
            result.Data = buffer[:bytesRead]
            result.Success = true
        } else {
            result.Error = readError
        }
    }
    return result
}
```

### Rob Pike Error Handling (NO EXCEPTIONS)
**"Errors are values"** - Treat them as such, don't hide them

```go
// ALWAYS do this immediately after function calls
f, err := os.Open(name)
if err != nil {
    return err
}
defer f.Close()

// NEVER do this
f, _ := os.Open(name) // FORBIDDEN

// NEVER say "error" when wrapping errors
if err != nil {
    return fmt.Errorf("error: %w", err) // FORBIDDEN
}

// Add context when wrapping
if err != nil {
    return fmt.Errorf("open %s: %w", name, err)
}
```

**Rules:**
- Check every error - no `_` assignments
- Handle errors at the call site, don't pass them up blindly
- Add meaningful context when wrapping
- Use `%w` for error wrapping, not `%v`

### Rob Pike Testing Approach
**"Test what matters, not what's easy"**

**MANDATORY: Always use testify for assertions and test structure**

```go
// Table-driven tests with testify - Required pattern
func TestSplit(t *testing.T) {
    tests := []struct {
        name  string
        input string
        sep   string
        want  []string
    }{
        {"basic split", "a,b,c", ",", []string{"a", "b", "c"}},
        {"empty string", "", ",", []string{""}},
        {"no separator", "a", ",", []string{"a"}},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := strings.Split(tt.input, tt.sep)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**Mandatory:**
- Always use testify/assert for assertions
- Use testify/require for assertions that should stop test execution
- Test the exported interface, not internals
- Table-driven tests with t.Run() for multiple cases
- Clear test names describing what's being tested
- Use testify/mock for mocking when necessary

## Development Workflow

### CLAUDE.md Guidelines
- Update CLAUDE.md with any new patterns or practices

### Commit Strategy
- **IMPORTANT** Commit small chunks of work frequently
- **IMPORTANT** NEVER mention claude code in commit messages
- Write succinct commit messages

### Testing Requirements
- **Always test any code that is written**
- Run tests before committing: `go test ./...`
- Include both positive and negative test cases
- Test edge cases and error conditions
- Maintain high test coverage for critical paths

### Rob Pike Code Review Checklist (ENFORCE STRICTLY)
- [ ] **Gofmt applied** - Reject if not formatted
- [ ] **Variable names ≤ 8 characters** - Longer names need justification  
- [ ] **Functions ≤ 30 lines** - Split if longer
- [ ] **No nested if statements > 3 levels** - Use early returns
- [ ] **Every error checked** - No `_` assignments
- [ ] **No clever tricks** - Code should be boring
- [ ] **NO COMMENTS** - Code must be self-explanatory, comments indicate unclear code
- [ ] **Exported functions documented** - Brief and clear (only exception to no-comment rule)
- [ ] **Tests cover the interface** - Not implementation details

### Forbidden Patterns
```go
// NEVER write code like this:
if condition {
    if anotherCondition {
        if yetAnother {
            doSomething()
        }
    }
}

// ALWAYS write like this:
if !condition {
    return
}
if !anotherCondition {
    return  
}
if !yetAnother {
    return
}
doSomething()
```

### The Pike Mantra
**"Clear is better than clever"**
**"Simple is better than complex"** 
**"Readable is better than terse"**

When in doubt, choose the most obvious solution.
