package create

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreator_Create(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		setup          func() string
		fileText       string
		wantErr        string
		cleanupCurrent bool
	}{
		{
			name:     "basic file creation",
			setup:    func() string { return filepath.Join(tmpDir, "test.txt") },
			fileText: "Hello, World!",
			wantErr:  "",
		},
		{
			name:     "empty file",
			setup:    func() string { return filepath.Join(tmpDir, "empty.txt") },
			fileText: "",
			wantErr:  "",
		},
		{
			name:     "multiline content",
			setup:    func() string { return filepath.Join(tmpDir, "multi.txt") },
			fileText: "line1\nline2\nline3",
			wantErr:  "",
		},
		{
			name:     "json content",
			setup:    func() string { return filepath.Join(tmpDir, "config.json") },
			fileText: `{"key": "value", "number": 42}`,
			wantErr:  "",
		},
		{
			name:     "create with nested directories",
			setup:    func() string { return filepath.Join(tmpDir, "nested", "deep", "file.txt") },
			fileText: "nested content",
			wantErr:  "",
		},
		{
			name:     "special characters",
			setup:    func() string { return filepath.Join(tmpDir, "special.txt") },
			fileText: "Special chars: äöü ñ 中文 🚀",
			wantErr:  "",
		},
		{
			name: "file already exists",
			setup: func() string {
				testFile := filepath.Join(tmpDir, "existing.txt")
				require.NoError(t, os.WriteFile(testFile, []byte("existing"), 0o644))
				return testFile
			},
			fileText: "new content",
			wantErr:  "file already exists",
		},
		{
			name: "invalid path characters",
			setup: func() string {
				return filepath.Join(tmpDir, "invalid\x00file.txt")
			},
			fileText: "content",
			wantErr:  "write file",
		},
		{
			name:     "nested directory creation",
			setup:    func() string { return filepath.Join(tmpDir, "level1", "level2", "level3", "file.txt") },
			fileText: "deep content",
			wantErr:  "",
		},
		{
			name:     "very long content",
			setup:    func() string { return filepath.Join(tmpDir, "long.txt") },
			fileText: string(make([]byte, 10000)),
			wantErr:  "",
		},
		{
			name:     "binary-like content",
			setup:    func() string { return filepath.Join(tmpDir, "binary.dat") },
			fileText: "\x00\x01\x02\xff\xfe\xfd",
			wantErr:  "",
		},
		{
			name:           "file in current directory",
			setup:          func() string { return "temp_test_file.txt" },
			fileText:       "current dir",
			wantErr:        "",
			cleanupCurrent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Creator{}
			path := tt.setup()
			err := c.Create(path, tt.fileText, false, false)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)

			content, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, tt.fileText, string(content))

			info, err := os.Stat(path)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())

			parentDir := filepath.Dir(path)
			if parentDir != tmpDir && parentDir != "." && parentDir != "/" {
				assert.DirExists(t, parentDir)

				dirInfo, err := os.Stat(parentDir)
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0o755), dirInfo.Mode().Perm())
			}

			if tt.cleanupCurrent {
				defer os.Remove(path)
			}
		})
	}
}
