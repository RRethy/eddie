package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobCommand(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{
		"test1.txt",
		"test2.txt",
		"config.json",
		"main.go",
		"helper.go",
	}

	for _, f := range files {
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, f), []byte("content"), 0o644))
	}

	subDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "nested.go"), []byte("content"), 0o644))

	tests := []struct {
		name      string
		args      []string
		wantFiles []string
		wantErr   bool
	}{
		{
			name:      "glob txt files",
			args:      []string{"glob", "*.txt", tmpDir},
			wantFiles: []string{"test1.txt", "test2.txt"},
			wantErr:   false,
		},
		{
			name:      "glob go files",
			args:      []string{"glob", "*.go", tmpDir},
			wantFiles: []string{"main.go", "helper.go"},
			wantErr:   false,
		},
		{
			name:      "glob json files",
			args:      []string{"glob", "*.json", tmpDir},
			wantFiles: []string{"config.json"},
			wantErr:   false,
		},
		{
			name:      "glob all files",
			args:      []string{"glob", "*", tmpDir},
			wantFiles: append(files, "subdir"),
			wantErr:   false,
		},
		{
			name:      "glob no matches",
			args:      []string{"glob", "*.xyz", tmpDir},
			wantFiles: []string{},
			wantErr:   false,
		},
		{
			name:    "glob invalid pattern",
			args:    []string{"glob", "[", tmpDir},
			wantErr: true,
		},
		{
			name:      "glob recursive txt files",
			args:      []string{"glob", "**/*.txt", tmpDir},
			wantFiles: []string{"test1.txt", "test2.txt"},
			wantErr:   false,
		},
		{
			name:    "glob without arguments",
			args:    []string{"glob"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
				return
			}

			if tt.name == "glob without arguments" {
				assert.Contains(t, stdout, "Error: pattern is required")
				return
			}

			require.NoError(t, err, "stderr: %s", stderr)

			if len(tt.wantFiles) == 0 {
				assert.Empty(t, strings.TrimSpace(stdout))
				return
			}

			lines := strings.Split(strings.TrimSpace(stdout), "\n")
			actualFiles := make([]string, len(lines))
			for i, line := range lines {
				actualFiles[i] = filepath.Base(line)
			}

			assert.ElementsMatch(t, tt.wantFiles, actualFiles)
		})
	}
}
