package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMcpServer_createViewTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createViewTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "view", tool.Name)
	assert.Contains(t, tool.Description, "View file contents")
}

func TestMcpServer_createStrReplaceTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createStrReplaceTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "str_replace", tool.Name)
	assert.Contains(t, tool.Description, "Replace all occurrences")
}

func TestMcpServer_createCreateTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createCreateTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "create", tool.Name)
	assert.Contains(t, tool.Description, "Create a new file")
}

func TestMcpServer_createInsertTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createInsertTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "insert", tool.Name)
	assert.Contains(t, tool.Description, "Insert a new line")
}

func TestMcpServer_createUndoEditTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createUndoEditTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "undo_edit", tool.Name)
	assert.Contains(t, tool.Description, "Undo the last edit")
}

func TestMcpServer_createGlobTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createGlobTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "glob", tool.Name)
	assert.Contains(t, tool.Description, "Fast file pattern matching")
}

func TestMcpServer_createSearchTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createSearchTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "search", tool.Name)
	assert.Contains(t, tool.Description, "Search for code patterns")
}

func TestMcpServer_createBatchTool(t *testing.T) {
	m := &McpServer{}
	tool := m.createBatchTool()

	assert.NotNil(t, tool)
	assert.Equal(t, "batch", tool.Name)
	assert.Contains(t, tool.Description, "Execute multiple eddie operations")
}

func TestMcpServer_handleSearch(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")
	goContent := `package main

func hello() {
	println("Hello, World!")
}

func goodbye() {
	println("Goodbye!")
}
`

	err := os.WriteFile(goFile, []byte(goContent), 0o644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		args    map[string]any
		name    string
		wantErr bool
	}{
		{
			name: "search Go functions",
			args: map[string]any{
				"path":              goFile,
				"tree_sitter_query": "(function_declaration name: (identifier) @func)",
			},
			wantErr: false,
		},
		{
			name: "missing path",
			args: map[string]any{
				"tree_sitter_query": "(function_declaration name: (identifier) @func)",
			},
			wantErr: true,
		},
		{
			name: "missing query",
			args: map[string]any{
				"path": goFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := m.handleSearch(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)

			if tt.name == "search Go functions" {
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, textContent.Text, "hello")
				assert.Contains(t, textContent.Text, "goodbye")
			}
		})
	}
}

func TestMcpServer_handleBatch(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	initialContent := "line 1\nline 2\nline 3"
	
	err := os.WriteFile(testFile, []byte(initialContent), 0o644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
		check   func(t *testing.T, result *mcp.CallToolResult)
	}{
		{
			name: "successful batch with mixed operations",
			args: map[string]any{
				"operations": `{
					"operations": [
						{
							"type": "view",
							"path": "` + testFile + `"
						},
						{
							"type": "str_replace",
							"path": "` + testFile + `",
							"old_str": "line 2",
							"new_str": "modified line 2"
						},
						{
							"type": "view",
							"path": "` + testFile + `"
						}
					]
				}`,
			},
			wantErr: false,
			check: func(t *testing.T, result *mcp.CallToolResult) {
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				
				assert.Contains(t, textContent.Text, "success")
				assert.Contains(t, textContent.Text, "view")
				assert.Contains(t, textContent.Text, "str_replace")
			},
		},
		{
			name: "batch with some operations failing",
			args: map[string]any{
				"operations": `{
					"operations": [
						{
							"type": "view",
							"path": "` + testFile + `"
						},
						{
							"type": "view",
							"path": "/nonexistent/file.txt"
						},
						{
							"type": "ls",
							"path": "` + tmpDir + `"
						}
					]
				}`,
			},
			wantErr: false,
			check: func(t *testing.T, result *mcp.CallToolResult) {
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				
				assert.Contains(t, textContent.Text, "success")
				assert.Contains(t, textContent.Text, "error")
			},
		},
		{
			name: "missing operations parameter",
			args: map[string]any{},
			wantErr: true,
		},
		{
			name: "invalid JSON in operations",
			args: map[string]any{
				"operations": `{"invalid": json}`,
			},
			wantErr: false,
			check: func(t *testing.T, result *mcp.CallToolResult) {
				assert.True(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, textContent.Text, "Error parsing batch operations")
			},
		},
		{
			name: "empty operations array",
			args: map[string]any{
				"operations": `{"operations": []}`,
			},
			wantErr: false,
			check: func(t *testing.T, result *mcp.CallToolResult) {
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, textContent.Text, "results")
			},
		},
		{
			name: "batch with create and insert operations",
			args: map[string]any{
				"operations": `{
					"operations": [
						{
							"type": "create",
							"path": "` + filepath.Join(tmpDir, "new.txt") + `",
							"content": "new file content"
						},
						{
							"type": "insert",
							"path": "` + filepath.Join(tmpDir, "new.txt") + `",
							"insert_line": 2,
							"new_str": "inserted line"
						}
					]
				}`,
			},
			wantErr: false,
			check: func(t *testing.T, result *mcp.CallToolResult) {
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, textContent.Text, "success")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := m.handleBatch(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)

			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}
