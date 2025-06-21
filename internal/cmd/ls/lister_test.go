package ls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLister_Ls(t *testing.T) {
	tmpDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("package main"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, "subdir"), 0o755))

	l := &Lister{}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"list directory", tmpDir, false},
		{"list current directory", ".", false},
		{"nonexistent directory", "/nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := l.Ls(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
