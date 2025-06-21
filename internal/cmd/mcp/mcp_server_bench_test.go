package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func BenchmarkMcpServer_handleView(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5\n"

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	m := &McpServer{}
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"path": testFile,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.handleView(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMcpServer_handleCreate(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 10000},
		{"large", 100000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			tmpDir := b.TempDir()
			m := &McpServer{}
			content := fmt.Sprintf("content_%s", string(make([]byte, tt.size)))

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				req := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Arguments: map[string]interface{}{
							"path":    filepath.Join(tmpDir, fmt.Sprintf("file_%d.txt", i)),
							"content": content,
						},
					},
				}

				_, err := m.handleCreate(context.Background(), req)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkMcpServer_handleStrReplace(b *testing.B) {
	tmpDir := b.TempDir()
	m := &McpServer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("test_%d.txt", i))
		content := "Hello World Hello World Hello World"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			b.Fatal(err)
		}

		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"path":    testFile,
					"old_str": "World",
					"new_str": "Universe",
				},
			},
		}

		_, err = m.handleStrReplace(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMcpServer_handleInsert(b *testing.B) {
	tmpDir := b.TempDir()
	m := &McpServer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("test_%d.txt", i))
		content := "line1\nline2\nline3\nline4\nline5\n"

		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			b.Fatal(err)
		}

		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"path":    testFile,
					"line":    float64(3),
					"content": "inserted line",
				},
			},
		}

		_, err = m.handleInsert(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMcpServer_handleGlob(b *testing.B) {
	tmpDir := b.TempDir()
	
	for i := 0; i < 100; i++ {
		err := os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("test%d.txt", i)), []byte("content"), 0644)
		if err != nil {
			b.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.go", i)), []byte("content"), 0644)
		if err != nil {
			b.Fatal(err)
		}
	}

	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		err := os.WriteFile(filepath.Join(subDir, fmt.Sprintf("nested%d.txt", i)), []byte("content"), 0644)
		if err != nil {
			b.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		pattern string
	}{
		{"simple", "*.txt"},
		{"recursive", "**/*.txt"},
		{"all_recursive", "**"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			m := &McpServer{}
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"pattern": tt.pattern,
						"path":    tmpDir,
					},
				},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := m.handleGlob(context.Background(), req)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
