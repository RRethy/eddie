# Eddie CLI - Feature Roadmap

This document outlines planned features to make eddie a comprehensive MCP server for Claude Code's file editing needs.

## High Priority - Essential File Management

### Core File Operations
- [ ] **`delete`** - Remove files and directories
  - Support recursive directory deletion with safety checks
  - Add `--force` flag for non-interactive deletion
  - Integrate with undo system for file recovery

- [ ] **`copy`** - Copy files and directories
  - Support recursive directory copying
  - Handle file permissions and metadata preservation
  - Add progress indication for large operations

- [ ] **`move`** - Move/rename files and directories
  - Cross-filesystem move support
  - Atomic operations where possible
  - Update edit history tracking for moved files

### Advanced Content Operations
- [ ] **`search`** - Content search with regex support
  - Multi-file search across directory trees
  - Context lines display (before/after matches)
  - File type filtering and exclusion patterns
  - Case-sensitive/insensitive options

- [ ] **`replace_lines`** - Replace specific line ranges
  - Replace lines N-M with new content
  - Support for inserting different number of lines
  - Integration with diff display

- [ ] **`append`** - Add content to file end
  - Append text or entire files
  - Optional newline handling
  - Batch append to multiple files

## Medium Priority - Enhanced Editing

### Directory Management
- [ ] **`mkdir`** - Create directories
  - Create parent directories automatically (-p flag)
  - Set permissions during creation
  - Batch directory creation

### Batch Operations
- [ ] **`multi_replace`** - Replace across multiple files
  - Use glob patterns to select files
  - Dry-run mode to preview changes
  - Rollback capability for batch operations

- [ ] **`delete_lines`** - Remove specific lines or ranges
  - Delete by line numbers or patterns
  - Support for multiple ranges in single operation
  - Integration with undo system

### File Metadata
- [ ] **`chmod`** - Set file permissions
  - Recursive permission setting
  - Symbolic and octal notation support
  - Restore original permissions in undo

- [ ] **`find_files`** - Find files by name/pattern with content filtering
  - Combine filename patterns with content search
  - Output format options (list, detailed, JSON)
  - Integration with other commands (pipe to replace, delete)

## Lower Priority - Advanced Features

### Project-Level Operations
- [ ] **`backup_dir`** - Backup entire directories
  - Create timestamped backups
  - Compression options (tar.gz, zip)
  - Incremental backup support
  - Restore from backup functionality

### Content Transformation
- [ ] **`format`** - Code formatting integration
  - Language-specific formatting (gofmt, prettier, black)
  - Configuration file support (.editorconfig, etc.)
  - Format-on-save integration

- [ ] **`sort_lines`** - Sort file contents
  - Alphabetical, numerical, and custom sorting
  - Sort specific sections of files
  - Preserve file structure (headers, comments)

- [ ] **`template`** - Template-based file generation
  - Variable substitution in templates
  - Template library management
  - Project scaffolding from templates

## Implementation Guidelines

### Architecture Consistency
- Follow existing pattern: `cmd/` → `internal/cmd/[command]/` → MCP registration
- Reuse `fileops` and `display` packages
- Maintain consistent flag naming (`--show-diff`, `--show-result`)
- Add comprehensive error handling and validation

### Testing Requirements
- Unit tests for all business logic
- End-to-end tests for CLI interface
- MCP integration tests
- Performance benchmarks for large operations

### Documentation
- Update CLAUDE.md with new command documentation
- Add usage examples for each command
- Include common workflow patterns
- Performance considerations and limitations

## Future Considerations

### Advanced Features (Long-term)
- **Symlink operations** - Create and manage symbolic links
- **Archive operations** - Create/extract tar/zip files
- **Encoding conversion** - Handle different text encodings
- **Conflict resolution** - Handle merge conflicts in collaborative editing
- **Plugin system** - Extensible command architecture
- **Configuration management** - User preferences and defaults
- **Remote file operations** - SSH/SFTP support for remote editing
- **Version control integration** - Git operations through eddie
- **File watching** - Monitor files for changes
- **Syntax highlighting** - Enhanced view command with syntax coloring

### Performance Optimizations
- **Streaming operations** - Handle large files without loading entirely in memory
- **Parallel processing** - Concurrent operations for batch commands
- **Caching** - Cache file metadata for faster operations
- **Progress indicators** - Show progress for long-running operations

### Security Enhancements
- **Permission validation** - Check user permissions before operations
- **Safe mode** - Restricted operations in untrusted environments
- **Audit logging** - Track all file operations for security
- **Sandboxing** - Limit operations to specific directories

---

## Contributing

When implementing new features:

1. Create feature branch from `master`
2. Follow Rob Pike Go style guidelines in CLAUDE.md
3. Add comprehensive tests (unit + e2e)
4. Update documentation
5. Ensure golangci-lint passes
6. Add MCP tool registration
7. Update this TODO.md to track progress

## Notes

- Prioritize features that Claude Code uses most frequently
- Maintain backward compatibility with existing MCP tools
- Consider performance impact of new operations
- Security and data safety are paramount for file operations