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
		name     string
		path     string
		fileText string
		wantErr  bool
	}{
		{
			name:     "basic file creation",
			path:     filepath.Join(tmpDir, "test.txt"),
			fileText: "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "empty file",
			path:     filepath.Join(tmpDir, "empty.txt"),
			fileText: "",
			wantErr:  false,
		},
		{
			name:     "multiline content",
			path:     filepath.Join(tmpDir, "multi.txt"),
			fileText: "line1\nline2\nline3",
			wantErr:  false,
		},
		{
			name:     "json content",
			path:     filepath.Join(tmpDir, "config.json"),
			fileText: `{"key": "value", "number": 42}`,
			wantErr:  false,
		},
		{
			name:     "create with nested directories",
			path:     filepath.Join(tmpDir, "nested", "deep", "file.txt"),
			fileText: "nested content",
			wantErr:  false,
		},
		{
			name:     "special characters",
			path:     filepath.Join(tmpDir, "special.txt"),
			fileText: "Special chars: äöü ñ 中文 🚀",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Creator{}
			err := c.Create(tt.path, tt.fileText)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file was created with correct content
			content, err := os.ReadFile(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.fileText, string(content))

			// Verify file permissions
			info, err := os.Stat(tt.path)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
		})
	}
}

func TestCreator_Create_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	c := &Creator{}

	tests := []struct {
		name    string
		setup   func() string
		content string
		wantErr string
	}{
		{
			name: "file already exists",
			setup: func() string {
				testFile := filepath.Join(tmpDir, "existing.txt")
				require.NoError(t, os.WriteFile(testFile, []byte("existing"), 0644))
				return testFile
			},
			content: "new content",
			wantErr: "file already exists",
		},
		{
			name: "invalid path characters",
			setup: func() string {
				return filepath.Join(tmpDir, "invalid\x00file.txt")
			},
			content: "content",
			wantErr: "write file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := c.Create(path, tt.content)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCreator_Create_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	c := &Creator{}

	// Test creating file in non-existent nested directories
	deepPath := filepath.Join(tmpDir, "level1", "level2", "level3", "file.txt")
	err := c.Create(deepPath, "deep content")
	require.NoError(t, err)

	// Verify all directories were created
	assert.DirExists(t, filepath.Join(tmpDir, "level1"))
	assert.DirExists(t, filepath.Join(tmpDir, "level1", "level2"))
	assert.DirExists(t, filepath.Join(tmpDir, "level1", "level2", "level3"))

	// Verify file content
	content, err := os.ReadFile(deepPath)
	require.NoError(t, err)
	assert.Equal(t, "deep content", string(content))

	// Verify directory permissions
	info, err := os.Stat(filepath.Join(tmpDir, "level1"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
}

func TestCreator_Create_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	c := &Creator{}

	tests := []struct {
		name     string
		path     string
		fileText string
	}{
		{
			name:     "very long content",
			path:     filepath.Join(tmpDir, "long.txt"),
			fileText: string(make([]byte, 10000)),
		},
		{
			name:     "binary-like content",
			path:     filepath.Join(tmpDir, "binary.dat"),
			fileText: "\x00\x01\x02\xff\xfe\xfd",
		},
		{
			name:     "file in current directory",
			path:     "temp_test_file.txt",
			fileText: "current dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Create(tt.path, tt.fileText)
			require.NoError(t, err)

			// Clean up current directory file
			if tt.name == "file in current directory" {
				defer os.Remove(tt.path)
			}

			content, err := os.ReadFile(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.fileText, string(content))
		})
	}
}
