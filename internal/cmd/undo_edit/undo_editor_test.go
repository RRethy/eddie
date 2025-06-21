package undo_edit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUndoEditor_RecordEdit(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	tests := []struct {
		name        string
		editType    string
		oldContent  string
		newContent  string
		position    int
		fileContent string
	}{
		{"str_replace edit", "str_replace", "old", "new", -1, "some old content"},
		{"insert edit", "insert", "", "new line", 2, "line1\nline2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.fileContent), 0644))

			err := u.RecordEdit(testFile, tt.editType, tt.oldContent, tt.newContent, tt.position)
			require.NoError(t, err)

			editPath, err := u.getEditFilePath(testFile)
			require.NoError(t, err)

			history, err := u.readEditHistory(editPath)
			require.NoError(t, err)

			assert.Equal(t, testFile, history.FilePath)
			assert.Len(t, history.Edits, 1)

			edit := history.Edits[0]
			assert.Equal(t, tt.editType, edit.EditType)
			assert.Equal(t, tt.oldContent, edit.OldContent)
			assert.Equal(t, tt.newContent, edit.NewContent)
			assert.Equal(t, tt.position, edit.Position)
		})
	}
}

func TestUndoEditor_RecordEdit_Multiple(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("original content"), 0644))

	err := u.RecordEdit(testFile, "str_replace", "original", "first", -1)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(testFile, []byte("first content"), 0644))

	err = u.RecordEdit(testFile, "str_replace", "content", "text", -1)
	require.NoError(t, err)

	editPath, err := u.getEditFilePath(testFile)
	require.NoError(t, err)

	history, err := u.readEditHistory(editPath)
	require.NoError(t, err)

	assert.Len(t, history.Edits, 2)
	assert.Equal(t, "str_replace", history.Edits[0].EditType)
	assert.Equal(t, "original", history.Edits[0].OldContent)
	assert.Equal(t, "first", history.Edits[0].NewContent)

	assert.Equal(t, "str_replace", history.Edits[1].EditType)
	assert.Equal(t, "content", history.Edits[1].OldContent)
	assert.Equal(t, "text", history.Edits[1].NewContent)
}

func TestUndoEditor_UndoEdit_StrReplace(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := "hello world\nline2\nline3\n"
	require.NoError(t, os.WriteFile(testFile, []byte(originalContent), 0644))

	modifiedContent := "hi world\nline2\nline3\n"
	require.NoError(t, os.WriteFile(testFile, []byte(modifiedContent), 0644))

	err := u.RecordEdit(testFile, "str_replace", "hello", "hi", -1)
	require.NoError(t, err)

	currentContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, modifiedContent, string(currentContent))

	err = u.UndoEdit(testFile)
	require.NoError(t, err)

	restoredContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(restoredContent))

	editPath, err := u.getEditFilePath(testFile)
	require.NoError(t, err)
	_, err = os.Stat(editPath)
	assert.True(t, os.IsNotExist(err), "Edit history file should be removed after all edits are undone")
}

func TestUndoEditor_UndoEdit_Insert(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := "line1\nline3\n"
	require.NoError(t, os.WriteFile(testFile, []byte(originalContent), 0644))

	modifiedContent := "line1\nline2\nline3\n"
	require.NoError(t, os.WriteFile(testFile, []byte(modifiedContent), 0644))

	err := u.RecordEdit(testFile, "insert", "", "line2", 2)
	require.NoError(t, err)

	currentContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, modifiedContent, string(currentContent))

	err = u.UndoEdit(testFile)
	require.NoError(t, err)

	restoredContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(restoredContent))
}

func TestUndoEditor_UndoEdit_Multiple(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")

	content1 := "version 1\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content1), 0644))

	content2 := "version 2\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content2), 0644))
	err := u.RecordEdit(testFile, "str_replace", "1", "2", -1)
	require.NoError(t, err)

	content3 := "version 3\n"
	require.NoError(t, os.WriteFile(testFile, []byte(content3), 0644))
	err = u.RecordEdit(testFile, "str_replace", "2", "3", -1)
	require.NoError(t, err)

	err = u.UndoEdit(testFile)
	require.NoError(t, err)

	restoredContent, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, content2, string(restoredContent))

	err = u.UndoEdit(testFile)
	require.NoError(t, err)

	restoredContent, err = os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, content1, string(restoredContent))

	editPath, err := u.getEditFilePath(testFile)
	require.NoError(t, err)
	_, err = os.Stat(editPath)
	assert.True(t, os.IsNotExist(err), "Edit history file should be removed after all edits are undone")
}

func TestUndoEditor_UndoEdit_ModificationTimeValidation(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := "hello world\n"
	require.NoError(t, os.WriteFile(testFile, []byte(originalContent), 0644))

	// Modify file and record edit
	modifiedContent := "hi world\n"
	require.NoError(t, os.WriteFile(testFile, []byte(modifiedContent), 0644))
	err := u.RecordEdit(testFile, "str_replace", "hello", "hi", -1)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	externalContent := "external change\n"
	require.NoError(t, os.WriteFile(testFile, []byte(externalContent), 0644))

	err = u.UndoEdit(testFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file has been modified since last tracked edit")
}

func TestUndoEditor_UndoEdit_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	u := &UndoEditor{}

	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	tests := []struct {
		name    string
		setup   func() string
		wantErr string
	}{
		{
			name: "file does not exist",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent.txt")
			},
			wantErr: "file does not exist",
		},
		{
			name: "no edit history",
			setup: func() string {
				testFile := filepath.Join(tmpDir, "noedits.txt")
				require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))
				return testFile
			},
			wantErr: "read edit history",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := u.UndoEdit(path)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestUndoEditor_createSafeFilename(t *testing.T) {
	u := &UndoEditor{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple path",
			input:    "/home/user/file.txt",
			expected: "_home_user_file_txt",
		},
		{
			name:     "path with spaces",
			input:    "/home/user/my file.txt",
			expected: "_home_user_my_file_txt",
		},
		{
			name:     "windows path",
			input:    "C:\\Users\\user\\file.txt",
			expected: "C_\\Users\\user\\file_txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := u.createSafeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
