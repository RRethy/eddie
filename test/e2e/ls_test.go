package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLsCommand(t *testing.T) {
	tmpDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("package main"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "subdir"), 0o755))

	tests := []struct {
		name       string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "ls help",
			args:       []string{"ls", "--help"},
			wantOutput: []string{"List directory contents"},
			wantErr:    false,
		},
		{
			name:       "ls current directory",
			args:       []string{"ls"},
			wantOutput: []string{},
			wantErr:    false,
		},
		{
			name:       "ls specific directory",
			args:       []string{"ls", tmpDir},
			wantOutput: []string{"file1.txt", "file2.go", "subdir"},
			wantErr:    false,
		},
		{
			name:       "ls nonexistent directory",
			args:       []string{"ls", "/nonexistent"},
			wantOutput: []string{"Error:"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
				if len(tt.wantOutput) > 0 {
					output := stdout + stderr
					for _, want := range tt.wantOutput {
						assert.Contains(t, output, want)
					}
				}
				return
			}

			assert.NoError(t, err, "stderr: %s", stderr)
			for _, want := range tt.wantOutput {
				assert.Contains(t, stdout, want)
			}
		})
	}
}

func TestLsCommandOutput(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"alpha.txt", "beta.go", "gamma.md"}
	for _, file := range files {
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, file), []byte("content"), 0o644))
	}

	stdout, stderr, err := runEddie(t, "ls", tmpDir)
	assert.NoError(t, err, "stderr: %s", stderr)

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	assert.Len(t, lines, len(files), "Expected %d files but got %d", len(files), len(lines))

	for _, file := range files {
		assert.Contains(t, stdout, file)
	}
}
