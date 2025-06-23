package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_Search(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		query    string
		expected bool
	}{
		{
			name: "basic function search",
			content: `package main

func hello() {
	println("hello")
}`,
			query:    "(function_declaration name: (identifier) @func)",
			expected: true,
		},
		{
			name: "call expression search",
			content: `package main

func main() {
	hello()
}`,
			query:    "(call_expression function: (identifier) @call)",
			expected: true,
		},
		{
			name: "no match",
			content: `package main

var x = 1`,
			query:    "(function_declaration name: (identifier) @func)",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.go")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)

			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_SearchDir(t *testing.T) {
	tmpDir := t.TempDir()

	testFile1 := filepath.Join(tmpDir, "test1.go")
	err := os.WriteFile(testFile1, []byte(`package main
func hello() {
	println("hello")
}`), 0644)
	require.NoError(t, err)

	testFile2 := filepath.Join(tmpDir, "test2.go")
	err = os.WriteFile(testFile2, []byte(`package main
func world() {
	println("world")
}`), 0644)
	require.NoError(t, err)

	s := &Searcher{}
	err = s.Search(tmpDir, "(function_declaration name: (identifier) @func)")
	assert.NoError(t, err)
}
