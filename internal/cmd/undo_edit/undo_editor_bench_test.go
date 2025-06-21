package undo_edit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkUndoEditor_RecordEdit(b *testing.B) {
	tmpDir := b.TempDir()
	u := &UndoEditor{}

	// Set XDG_CACHE_HOME to tmpDir for test isolation
	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	// Test different file sizes
	sizes := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 10000},
		{"large", 1000000},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			content := strings.Repeat("test line\n", size.size)

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, fmt.Sprintf("bench_%s_%d.txt", size.name, n))
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = u.RecordEdit(testFile, "str_replace", "old", "new", -1)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUndoEditor_UndoEdit(b *testing.B) {
	tmpDir := b.TempDir()
	u := &UndoEditor{}

	// Set XDG_CACHE_HOME to tmpDir for test isolation
	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	// Test different file sizes
	sizes := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 10000},
		{"large", 100000},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			originalContent := strings.Repeat("original line\n", size.size)
			modifiedContent := strings.Repeat("modified line\n", size.size)

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, fmt.Sprintf("bench_undo_%s_%d.txt", size.name, n))

				// Create original file
				err := os.WriteFile(testFile, []byte(originalContent), 0644)
				if err != nil {
					b.Fatal(err)
				}

				// Modify file
				err = os.WriteFile(testFile, []byte(modifiedContent), 0644)
				if err != nil {
					b.Fatal(err)
				}

				// Record edit after modification
				err = u.RecordEdit(testFile, "str_replace", "original", "modified", -1)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				// Undo edit
				err = u.UndoEdit(testFile, false, false, 1)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUndoEditor_MultipleEdits(b *testing.B) {
	tmpDir := b.TempDir()
	u := &UndoEditor{}

	// Set XDG_CACHE_HOME to tmpDir for test isolation
	oldCacheHome := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		if oldCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	}()
	os.Setenv("XDG_CACHE_HOME", tmpDir)

	// Test different numbers of edits
	editCounts := []struct {
		name  string
		count int
	}{
		{"few", 5},
		{"many", 50},
		{"lots", 200},
	}

	for _, ec := range editCounts {
		b.Run(ec.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, fmt.Sprintf("multi_edit_%s_%d.txt", ec.name, n))
				content := "original content"

				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}

				// Create multiple edits
				for i := 0; i < ec.count; i++ {
					newContent := fmt.Sprintf("content_%d", i)
					err = os.WriteFile(testFile, []byte(newContent), 0644)
					if err != nil {
						b.Fatal(err)
					}

					err = u.RecordEdit(testFile, "str_replace", "content", fmt.Sprintf("content_%d", i), -1)
					if err != nil {
						b.Fatal(err)
					}
				}

				b.StartTimer()

				// Undo one edit
				err = u.UndoEdit(testFile, false, false, 1)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFileOperations_UndoEdit(b *testing.B) {
	tmpDir := b.TempDir()
	content := strings.Repeat("test line\n", 1000)

	b.Run("os_ReadFile", func(b *testing.B) {
		testFile := filepath.Join(tmpDir, "read_test.txt")
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_, err := os.ReadFile(testFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("os_WriteFile", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			testFile := filepath.Join(tmpDir, fmt.Sprintf("write_test_%d.txt", n))
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("os_Remove", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			testFile := filepath.Join(tmpDir, fmt.Sprintf("remove_test_%d.txt", n))
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				b.Fatal(err)
			}
			b.StartTimer()

			err = os.Remove(testFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkStringOperations_UndoEdit(b *testing.B) {
	// Test string operations used in edit file management

	b.Run("strings_ReplaceAll", func(b *testing.B) {
		content := strings.Repeat("hello world test content\n", 1000)

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = strings.ReplaceAll(content, "hello", "hi")
		}
	})

	b.Run("createSafeFilename", func(b *testing.B) {
		u := &UndoEditor{}
		paths := []string{
			"/home/user/file.txt",
			"/very/long/path/to/some/file/with/many/segments.txt",
			"C:\\Users\\user\\Documents\\file with spaces.txt",
			"/tmp/test.txt",
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for _, path := range paths {
				_ = u.createSafeFilename(path)
			}
		}
	})

	b.Run("filepath_operations", func(b *testing.B) {
		path := "/path/to/test.txt"

		for n := 0; n < b.N; n++ {
			_ = filepath.Dir(path)
			_ = filepath.Base(path)
			_ = filepath.Join("/tmp", "edit.json")
		}
	})
}
