package view

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkViewer_ViewFile(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files of different sizes
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
			// Create file with specified number of lines
			content := strings.Repeat("this is a test line with some content\n", size.size)
			testFile := filepath.Join(tmpDir, "bench_"+size.name+".txt")
			err := os.WriteFile(testFile, []byte(content), 0644)
			if err != nil {
				b.Fatal(err)
			}

			v := &Viewer{}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := v.viewFile(testFile, "")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkViewer_ViewFileWithRange(b *testing.B) {
	tmpDir := b.TempDir()

	// Create large test file
	content := strings.Repeat("this is a test line with some content\n", 100000)
	testFile := filepath.Join(tmpDir, "bench_range.txt")
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	ranges := []struct {
		name   string
		range_ string
	}{
		{"first_100", "1,100"},
		{"middle_100", "50000,50100"},
		{"last_100", "99900,100000"},
		{"to_end", "99900,-1"},
	}

	for _, r := range ranges {
		b.Run(r.name, func(b *testing.B) {
			v := &Viewer{}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := v.viewFile(testFile, r.range_)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkViewer_ViewDir(b *testing.B) {
	tmpDir := b.TempDir()

	// Create directories with different numbers of files
	sizes := []struct {
		name  string
		files int
	}{
		{"small", 10},
		{"medium", 100},
		{"large", 1000},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			testDir := filepath.Join(tmpDir, "bench_"+size.name)
			err := os.Mkdir(testDir, 0755)
			if err != nil {
				b.Fatal(err)
			}

			// Create files
			for i := 0; i < size.files; i++ {
				fileName := filepath.Join(testDir, "file"+string(rune(i))+".txt")
				err := os.WriteFile(fileName, []byte("content"), 0644)
				if err != nil {
					b.Fatal(err)
				}
			}

			v := &Viewer{}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := v.viewDir(testDir)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkViewer_ParseRange(b *testing.B) {
	v := &Viewer{}

	ranges := []string{
		"",
		"1,10",
		"1000,2000",
		"50,-1",
		" 100 , 200 ",
	}

	for _, r := range ranges {
		b.Run("range_"+r, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := v.parseRange(r)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
