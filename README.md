# eddie

An experimental text editor designed for AI Agents, not humans.

## Features

- **File Operations**: View, create, edit, and manage files with full undo support
- **Code Search**: Powerful tree-sitter-based search across multiple programming languages
- **Edit History**: Automatic backup and undo functionality for all file modifications
- **MCP Server**: Built-in Model Context Protocol server support

## Installation

```bash
go install github.com/RRethy/eddie@latest
```

## MCP Setup for Claude

To use eddie with Claude Code via MCP:

1. Add eddie to your Claude MCP configuration:

```json
{
  "mcpServers": {
    "eddie": {
      "command": "eddie",
      "args": ["mcp"],
      "type": "stdio"
    }
  }
}
```

2. Restart Claude Code to load the MCP server
3. Eddie commands will be available as tools in Claude Code

## Commands

### view

Examine file contents or list directory contents.

```bash
eddie view <path> [range]

# Examples
eddie view file.txt                # View entire file
eddie view file.txt 10,20         # View lines 10-20
eddie view file.txt 15,-1         # View from line 15 to end
eddie view /path/to/directory      # List directory contents
```

### str_replace

Replace all occurrences of a string in a file.

```bash
eddie str_replace <path> <old_str> <new_str> [flags]

# Examples
eddie str_replace config.json "localhost" "example.com"
eddie str_replace app.py "old_function" "new_function" --show-diff
eddie str_replace main.go "TODO" "DONE" --show-result

# Flags
--show-diff     Show changes made to the file
--show-result   Show file content after modification
```

### create

Create a new file with specified content.

```bash
eddie create <path> <content> [flags]

# Examples
eddie create hello.txt "Hello, World!"
eddie create config.json '{"port": 8080}' --show-result
eddie create script.sh "#!/bin/bash\necho 'Hello'" --show-diff

# Flags
--show-diff     Show the created file content
--show-result   Show file content after creation
```

### insert

Insert a new line at a specific line number.

```bash
eddie insert <path> <line_number> <content> [flags]

# Examples
eddie insert app.py 10 "import os"
eddie insert config.json 5 '  "debug": true,' --show-diff
eddie insert README.md 1 "# My Project" --show-result

# Flags
--show-diff     Show changes made to the file
--show-result   Show file content after insertion
```

### undo_edit

Undo previous file modifications.

```bash
eddie undo_edit <path> [flags]

# Examples
eddie undo_edit app.py                # Undo last edit
eddie undo_edit config.json --count 3 # Undo last 3 edits
eddie undo_edit main.go --show-diff   # Show what was undone

# Flags
--show-diff     Show changes made during undo
--show-result   Show file content after undo
--count N       Number of edits to undo (default: 1)
```

### ls

List directory contents.

```bash
eddie ls [path]

# Examples
eddie ls                    # List current directory
eddie ls /path/to/directory # List specific directory
```

### search

Search for code patterns using tree-sitter queries.

```bash
eddie search <path> --tree-sitter-query "<query>"

# Language Support
Go          .go                    (function_declaration name: (identifier) @func)
JavaScript  .js, .mjs, .jsx        (function_declaration name: (identifier) @func)
TypeScript  .ts, .tsx              (function_declaration name: (identifier) @func)
Python      .py, .pyi              (function_definition name: (identifier) @func)
Rust        .rs                    (function_item name: (identifier) @func)
Java        .java                  (class_declaration name: (identifier) @class)
C           .c, .h                 (function_definition declarator: (function_declarator declarator: (identifier) @func))
C++         .cpp, .hpp, .cxx       (function_definition declarator: (function_declarator declarator: (identifier) @func))

# Examples
# Find all Go functions
eddie search . --tree-sitter-query "(function_declaration name: (identifier) @func)"

# Find Python functions
eddie search src/ --tree-sitter-query "(function_definition name: (identifier) @func)"

# Find Java classes
eddie search . --tree-sitter-query "(class_declaration name: (identifier) @class)"

# Find Go method calls
eddie search main.go --tree-sitter-query "(call_expression function: (identifier) @call)"

# Find Go struct types
eddie search . --tree-sitter-query "(type_declaration (type_spec name: (type_identifier) @struct type: (struct_type)))"
```

## MCP Server

Eddie includes a built-in MCP (Model Context Protocol) server that exposes all commands as tools for AI assistants.

```bash
eddie mcp
```

The MCP server provides structured access to all eddie commands with proper parameter validation and error handling.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
