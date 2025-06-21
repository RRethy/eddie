package view

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewer_parseRange(t *testing.T) {
	v := &Viewer{}

	tests := []struct {
		name      string
		viewRange string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{"empty range", "", 0, 0, false},
		{"valid range", "1,10", 1, 10, false},
		{"range to end", "5,-1", 5, -1, false},
		{"invalid format", "1,2,3", 0, 0, true},
		{"invalid start", "abc,10", 0, 0, true},
		{"invalid end", "1,xyz", 0, 0, true},
		{"start greater than end", "10,5", 0, 0, true},
		{"with spaces", " 1 , 10 ", 1, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := v.parseRange(tt.viewRange)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantStart, start)
			assert.Equal(t, tt.wantEnd, end)
		})
	}
}

func TestViewer_View(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0o644))

	testSubDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.Mkdir(testSubDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(testSubDir, "file1.txt"), []byte("content"), 0o644))

	v := &Viewer{}

	tests := []struct {
		name    string
		path    string
		range_  string
		wantErr bool
	}{
		{"view file", testFile, "", false},
		{"view file with range", testFile, "2,4", false},
		{"view file to end", testFile, "3,-1", false},
		{"view directory", testSubDir, "", false},
		{"nonexistent path", "/nonexistent", "", true},
		{"invalid range", testFile, "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.View(tt.path, tt.range_)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
