package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrReplaceCommand(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		initialContent string
		args           []string
		wantContent    string
		wantOutput     string
		wantErr        bool
	}{
		{
			name:           "basic replacement",
			initialContent: "hello world hello",
			args:           []string{"hello", "hi"},
			wantContent:    "hi world hi",
			wantOutput:     "Replaced 2 occurrence(s) of \"hello\" with \"hi\"",
			wantErr:        false,
		},
		{
			name:           "no matches found",
			initialContent: "hello world",
			args:           []string{"foo", "bar"},
			wantContent:    "hello world",
			wantOutput:     "No occurrences of \"foo\" found",
			wantErr:        false,
		},
		{
			name:           "replace with empty string",
			initialContent: "hello world hello",
			args:           []string{"hello ", ""},
			wantContent:    "world hello",
			wantOutput:     "Replaced 1 occurrence(s) of \"hello \" with \"\"",
			wantErr:        false,
		},
		{
			name:           "multiline replacement",
			initialContent: "line1\nhello\nline3\nhello\nline5",
			args:           []string{"hello", "hi"},
			wantContent:    "line1\nhi\nline3\nhi\nline5",
			wantOutput:     "Replaced 2 occurrence(s) of \"hello\" with \"hi\"",
			wantErr:        false,
		},
		{
			name:           "special characters",
			initialContent: "func(x) { return x + 1; }",
			args:           []string{"func(x)", "function(x)"},
			wantContent:    "function(x) { return x + 1; }",
			wantOutput:     "Replaced 1 occurrence(s) of \"func(x)\" with \"function(x)\"",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.initialContent), 0644))

			args := []string{"str_replace", testFile}
			args = append(args, tt.args...)

			stdout, stderr, err := runEddie(t, args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
				return
			}

			require.NoError(t, err, "stderr: %s", stderr)
			assert.Contains(t, stdout, tt.wantOutput)

			result, err := os.ReadFile(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.wantContent, string(result))
		})
	}
}

func TestStrReplaceCommandErrors(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing arguments",
			args:    []string{"str_replace"},
			wantErr: "path, old_str, and new_str are required",
		},
		{
			name:    "missing new_str",
			args:    []string{"str_replace", "file.txt", "old"},
			wantErr: "path, old_str, and new_str are required",
		},
		{
			name:    "nonexistent file",
			args:    []string{"str_replace", "/nonexistent/file.txt", "old", "new"},
			wantErr: "stat",
		},
		{
			name:    "directory instead of file",
			args:    []string{"str_replace", tmpDir, "old", "new"},
			wantErr: "cannot replace strings in directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.name == "missing arguments" || tt.name == "missing new_str" {
					assert.Contains(t, stdout, tt.wantErr)
				return
			}

			assert.True(t, err != nil || stderr != "", "Expected error but got none")
			if err != nil {
				assert.Contains(t, err.Error(), "exit status")
			}
		})
	}
}
