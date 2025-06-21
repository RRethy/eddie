package glob

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobber_Glob(t *testing.T) {
	tests := []struct {
		setup    func(t *testing.T) string
		validate func(t *testing.T, output string)
		name     string
		pattern  string
		path     string
		wantErr  bool
	}{
		{
			name:    "empty pattern",
			pattern: "",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "no matches",
			pattern: "*.nonexistent",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "go files",
			pattern: "*.go",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "recursive all files",
			pattern: "**",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "invalid pattern",
			pattern: "[",
			path:    ".",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Globber{}
			err := g.Glob(tt.pattern, tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGlobber_GlobWithTempFiles(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"test1.txt", "test2.txt", "other.go"}
	for _, f := range files {
		file, err := os.Create(tmpDir + "/" + f)
		require.NoError(t, err)
		file.Close()
	}

	g := &Globber{}
	err := g.Glob("*.txt", tmpDir)
	assert.NoError(t, err)
}
