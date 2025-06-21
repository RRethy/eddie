package str_replace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplacer_StrReplace(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name         string
		content      string
		oldStr       string
		newStr       string
		wantContent  string
		wantErr      bool
		expectOutput string
	}{
		{
			name:         "basic replacement",
			content:      "hello world hello",
			oldStr:       "hello",
			newStr:       "hi",
			wantContent:  "hi world hi",
			wantErr:      false,
			expectOutput: "Replaced 2 occurrence(s)",
		},
		{
			name:         "no matches",
			content:      "hello world",
			oldStr:       "foo",
			newStr:       "bar",
			wantContent:  "hello world",
			wantErr:      false,
			expectOutput: "No occurrences",
		},
		{
			name:         "empty old string",
			content:      "hello",
			oldStr:       "",
			newStr:       "x",
			wantContent:  "xhxexlxlxox",
			wantErr:      false,
			expectOutput: "Replaced 6 occurrence(s)",
		},
		{
			name:         "replace with empty string",
			content:      "hello world hello",
			oldStr:       "hello ",
			newStr:       "",
			wantContent:  "world hello",
			wantErr:      false,
			expectOutput: "Replaced 1 occurrence(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.content), 0644))

			r := &Replacer{}
			err := r.StrReplace(testFile, tt.oldStr, tt.newStr)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			result, err := os.ReadFile(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.wantContent, string(result))
		})
	}
}

func TestReplacer_StrReplace_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Replacer{}

	tests := []struct {
		name    string
		setup   func() string
		wantErr string
	}{
		{
			name: "nonexistent file",
			setup: func() string {
				return "/nonexistent/file.txt"
			},
			wantErr: "stat",
		},
		{
			name: "directory instead of file",
			setup: func() string {
				return tmpDir
			},
			wantErr: "cannot replace strings in directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := r.StrReplace(path, "old", "new")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}