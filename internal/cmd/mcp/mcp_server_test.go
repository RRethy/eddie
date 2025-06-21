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

func TestMcpServer_handleView(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid file view",
			args: map[string]interface{}{
				"path": testFile,
			},
			wantErr: false,
		},
		{
			name: "view with range",
			args: map[string]interface{}{
				"path":  testFile,
				"range": "1,2",
			},
			wantErr: false,
		},
		{
			name:    "missing path",
			args:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "nonexistent file",
			args: map[string]interface{}{
				"path": "/nonexistent/file.txt",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := m.handleView(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)
		})
	}
}

func TestMcpServer_handleCreate(t *testing.T) {
	tmpDir := t.TempDir()

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "create new file",
			args: map[string]interface{}{
				"path":    filepath.Join(tmpDir, "new.txt"),
				"content": "Hello, World!",
			},
			wantErr: false,
		},
		{
			name: "missing path",
			args: map[string]interface{}{
				"content": "Hello",
			},
			wantErr: true,
		},
		{
			name: "missing content",
			args: map[string]interface{}{
				"path": filepath.Join(tmpDir, "test.txt"),
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

			result, err := m.handleCreate(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsError)
			assert.NotEmpty(t, result.Content)
		})
	}
}

func TestMcpServer_handleStrReplace(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello World"

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid replacement",
			args: map[string]interface{}{
				"path":    testFile,
				"old_str": "World",
				"new_str": "Universe",
			},
			wantErr: false,
		},
		{
			name: "missing path",
			args: map[string]interface{}{
				"old_str": "old",
				"new_str": "new",
			},
			wantErr: true,
		},
		{
			name: "missing old_str",
			args: map[string]interface{}{
				"path":    testFile,
				"new_str": "new",
			},
			wantErr: true,
		},
		{
			name: "missing new_str",
			args: map[string]interface{}{
				"path":    testFile,
				"old_str": "old",
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

			result, err := m.handleStrReplace(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsError)
		})
	}
}

func TestMcpServer_handleInsert(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid insertion",
			args: map[string]interface{}{
				"path":    testFile,
				"line":    float64(2),
				"content": "inserted line",
			},
			wantErr: false,
		},
		{
			name: "missing path",
			args: map[string]interface{}{
				"line":    float64(1),
				"content": "text",
			},
			wantErr: true,
		},
		{
			name: "missing line",
			args: map[string]interface{}{
				"path":    testFile,
				"content": "text",
			},
			wantErr: true,
		},
		{
			name: "missing content",
			args: map[string]interface{}{
				"path": testFile,
				"line": float64(1),
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

			result, err := m.handleInsert(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsError)
		})
	}
}

func TestMcpServer_handleUndoEdit(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "original content"

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "missing path",
			args:    map[string]interface{}{},
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

			result, err := m.handleUndoEdit(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestMcpServer_handleGlob(t *testing.T) {
	tmpDir := t.TempDir()
	
	files := []string{"test1.txt", "test2.txt", "main.go"}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(tmpDir, f), []byte("content"), 0644)
		require.NoError(t, err)
	}

	m := &McpServer{}

	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "glob txt files",
			args: map[string]interface{}{
				"pattern": "*.txt",
				"path":    tmpDir,
			},
			wantErr: false,
		},
		{
			name: "glob without path",
			args: map[string]interface{}{
				"pattern": "*.go",
			},
			wantErr: false,
		},
		{
			name: "recursive glob",
			args: map[string]interface{}{
				"pattern": "**/*.txt",
				"path":    tmpDir,
			},
			wantErr: false,
		},
		{
			name:    "missing pattern",
			args:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "invalid pattern",
			args: map[string]interface{}{
				"pattern": "[",
				"path":    tmpDir,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := m.handleGlob(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)
		})
	}
}
