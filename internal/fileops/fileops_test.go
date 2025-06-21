package fileops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOps_ValidateFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	f := &FileOps{}

	tests := []struct {
		setup   func() string
		name    string
		wantErr bool
	}{
		{
			name: "existing file",
			setup: func() string {
				path := filepath.Join(tmpDir, "exists.txt")
				err := os.WriteFile(path, []byte("content"), 0o644)
				require.NoError(t, err)
				return path
			},
			wantErr: false,
		},
		{
			name: "existing directory",
			setup: func() string {
				path := filepath.Join(tmpDir, "dir")
				err := os.Mkdir(path, 0o755)
				require.NoError(t, err)
				return path
			},
			wantErr: false,
		},
		{
			name: "nonexistent file",
			setup: func() string {
				return "/nonexistent/file.txt"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			info, err := f.ValidateFileExists(path)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, info)
				assert.Contains(t, err.Error(), "stat")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, info)
			}
		})
	}
}

func TestFileOps_ValidateNotDir(t *testing.T) {
	tmpDir := t.TempDir()
	f := &FileOps{}

	filePath := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(filePath, []byte("content"), 0o644)
	require.NoError(t, err)

	dirPath := filepath.Join(tmpDir, "dir")
	err = os.Mkdir(dirPath, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name      string
		path      string
		operation string
		wantErr   bool
	}{
		{
			name:      "file - should pass",
			path:      filePath,
			operation: "test operation on",
			wantErr:   false,
		},
		{
			name:      "directory - should fail",
			path:      dirPath,
			operation: "test operation on",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := os.Stat(tt.path)
			require.NoError(t, err)

			err = f.ValidateNotDir(tt.path, info, tt.operation)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cannot test operation on directory")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileOps_ReadFileContentForOperation(t *testing.T) {
	tmpDir := t.TempDir()
	f := &FileOps{}

	tests := []struct {
		name      string
		setup     func() string
		operation string
		wantErr   string
		wantText  string
	}{
		{
			name: "successful read",
			setup: func() string {
				path := filepath.Join(tmpDir, "success.txt")
				err := os.WriteFile(path, []byte("test content"), 0o644)
				require.NoError(t, err)
				return path
			},
			operation: "test",
			wantText:  "test content",
		},
		{
			name: "directory error",
			setup: func() string {
				path := filepath.Join(tmpDir, "dir")
				err := os.Mkdir(path, 0o755)
				require.NoError(t, err)
				return path
			},
			operation: "test",
			wantErr:   "cannot test directory",
		},
		{
			name: "nonexistent file",
			setup: func() string {
				return "/nonexistent/file.txt"
			},
			operation: "test",
			wantErr:   "stat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			content, info, err := f.ReadFileContentForOperation(path, tt.operation)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Empty(t, content)
				assert.Nil(t, info)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantText, content)
				assert.NotNil(t, info)
			}
		})
	}
}

func TestFileOps_WriteFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	f := &FileOps{}

	tests := []struct {
		name    string
		path    string
		content string
		mode    os.FileMode
		wantErr bool
	}{
		{
			name:    "write to valid path",
			path:    filepath.Join(tmpDir, "write.txt"),
			content: "test content",
			mode:    0o644,
			wantErr: false,
		},
		{
			name:    "write to invalid path",
			path:    "/invalid/path/file.txt",
			content: "content",
			mode:    0o644,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := f.WriteFileContent(tt.path, tt.content, tt.mode)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "write file")
			} else {
				assert.NoError(t, err)

				content, err := os.ReadFile(tt.path)
				require.NoError(t, err)
				assert.Equal(t, tt.content, string(content))

				info, err := os.Stat(tt.path)
				require.NoError(t, err)
				assert.Equal(t, tt.mode, info.Mode().Perm())
			}
		})
	}
}

func TestFileOps_CreateFile(t *testing.T) {
	tmpDir := t.TempDir()
	f := &FileOps{}

	tests := []struct {
		name    string
		setup   func() string
		content string
		wantErr string
	}{
		{
			name: "create new file",
			setup: func() string {
				return filepath.Join(tmpDir, "new.txt")
			},
			content: "new content",
		},
		{
			name: "create file with directories",
			setup: func() string {
				return filepath.Join(tmpDir, "subdir", "deep", "new.txt")
			},
			content: "deep content",
		},
		{
			name: "file already exists",
			setup: func() string {
				path := filepath.Join(tmpDir, "exists.txt")
				err := os.WriteFile(path, []byte("old"), 0o644)
				require.NoError(t, err)
				return path
			},
			content: "new content",
			wantErr: "file already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := f.CreateFile(path, tt.content)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)

				content, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, tt.content, string(content))

				info, err := os.Stat(path)
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
			}
		})
	}
}
