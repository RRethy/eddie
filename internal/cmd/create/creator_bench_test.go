package create

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkCreator_Create(b *testing.B) {
	tmpDir := b.TempDir()

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
			content := strings.Repeat("x", size.size)
			c := &Creator{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				filePath := filepath.Join(tmpDir, "bench_"+size.name+"_"+string(rune(i))+".txt")
				b.StartTimer()

				err := c.Create(filePath, content)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCreator_CreateWithDirectories(b *testing.B) {
	tmpDir := b.TempDir()
	content := "benchmark content"

	// Test different directory depths
	depths := []struct {
		name  string
		depth int
	}{
		{"flat", 0},
		{"shallow", 3},
		{"deep", 10},
	}

	for _, depth := range depths {
		b.Run(depth.name, func(b *testing.B) {
			c := &Creator{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// Build path with specified depth
				pathParts := []string{tmpDir}
				for j := 0; j < depth.depth; j++ {
					pathParts = append(pathParts, "dir"+string(rune(j)))
				}
				pathParts = append(pathParts, "file_"+string(rune(i))+".txt")
				filePath := filepath.Join(pathParts...)
				b.StartTimer()

				err := c.Create(filePath, content)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCreator_CreateDifferentContentTypes(b *testing.B) {
	tmpDir := b.TempDir()

	contentTypes := []struct {
		name    string
		content string
	}{
		{"text", "This is a simple text file with some content"},
		{"json", `{"key": "value", "number": 42, "array": [1, 2, 3], "nested": {"inner": "data"}}`},
		{"multiline", "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"},
		{"binary", "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xfe\xfd\xfc\xfb\xfa"},
		{"unicode", "Hello 世界 🌍 Здравствуй мир مرحبا بالعالم"},
	}

	for _, ct := range contentTypes {
		b.Run(ct.name, func(b *testing.B) {
			c := &Creator{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				filePath := filepath.Join(tmpDir, "bench_"+ct.name+"_"+string(rune(i))+".txt")
				b.StartTimer()

				err := c.Create(filePath, ct.content)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkFileOperations(b *testing.B) {
	tmpDir := b.TempDir()
	content := "benchmark file operations content"

	b.Run("os_WriteFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			filePath := filepath.Join(tmpDir, "writefile_"+string(rune(i))+".txt")
			b.StartTimer()

			err := os.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("os_MkdirAll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			dirPath := filepath.Join(tmpDir, "mkdir", "deep", "path", string(rune(i)))
			b.StartTimer()

			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("filepath_Dir", func(b *testing.B) {
		testPath := "/very/long/path/to/some/deeply/nested/file.txt"
		for i := 0; i < b.N; i++ {
			_ = filepath.Dir(testPath)
		}
	})

	b.Run("os_Stat", func(b *testing.B) {
		// Create a file to stat
		testFile := filepath.Join(tmpDir, "stattest.txt")
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := os.Stat(testFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCreator_CreateConcurrent(b *testing.B) {
	tmpDir := b.TempDir()
	content := "concurrent creation test content"
	c := &Creator{}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			filePath := filepath.Join(tmpDir, "concurrent_"+string(rune(i))+".txt")
			err := c.Create(filePath, content)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}
