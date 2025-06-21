package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

	// Create test directory with files
	testSubDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.Mkdir(testSubDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(testSubDir, "file1.txt"), []byte("content"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(testSubDir, "file2.txt"), []byte("content"), 0644))

	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "view entire file",
			args:       []string{"view", testFile},
			wantOutput: "line1\nline2\nline3\nline4\nline5\n",
			wantErr:    false,
		},
		{
			name:       "view file with range",
			args:       []string{"view", testFile, "2,4"},
			wantOutput: "line2\nline3\nline4\n",
			wantErr:    false,
		},
		{
			name:       "view file from line to end",
			args:       []string{"view", testFile, "3,-1"},
			wantOutput: "line3\nline4\nline5\n",
			wantErr:    false,
		},
		{
			name:       "view directory",
			args:       []string{"view", testSubDir},
			wantOutput: "file1.txt\nfile2.txt\n",
			wantErr:    false,
		},
		{
			name:    "view nonexistent file",
			args:    []string{"view", "/nonexistent/file.txt"},
			wantErr: true,
		},
		{
			name:    "view with invalid range",
			args:    []string{"view", testFile, "invalid"},
			wantErr: true,
		},
		{
			name:    "view without arguments",
			args:    []string{"view"},
			wantErr: false, // Should print error message but not exit with error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
				return
			}

			if tt.name == "view without arguments" {
				assert.Contains(t, stdout, "Error: path is required")
				return
			}

			require.NoError(t, err, "stderr: %s", stderr)
			assert.Equal(t, tt.wantOutput, stdout)
		})
	}
}