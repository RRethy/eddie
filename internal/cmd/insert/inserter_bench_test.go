package insert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkInserter_Insert(b *testing.B) {
	tmpDir := b.TempDir()

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
			// Create file with specified number of lines
			content := strings.Repeat("this is a test line\n", size.size)

			i := &Inserter{}

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_"+size.name+"_"+string(rune(n))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = i.Insert(testFile, "50", "inserted line")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkInserter_InsertPosition(b *testing.B) {
	tmpDir := b.TempDir()

	// Create medium-sized file
	content := strings.Repeat("this is a test line\n", 10000)

	positions := []struct {
		name string
		line string
	}{
		{"beginning", "1"},
		{"early", "100"},
		{"middle", "5000"},
		{"late", "9000"},
		{"end", "10001"},
	}

	for _, pos := range positions {
		b.Run(pos.name, func(b *testing.B) {
			i := &Inserter{}

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_pos_"+pos.name+"_"+string(rune(n))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = i.Insert(testFile, pos.line, "inserted line")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkInserter_InsertContentLength(b *testing.B) {
	tmpDir := b.TempDir()

	// Create base file
	content := strings.Repeat("line\n", 1000)

	contentLengths := []struct {
		name    string
		content string
	}{
		{"short", "x"},
		{"medium", strings.Repeat("x", 100)},
		{"long", strings.Repeat("x", 10000)},
		{"multiline", "line1\nline2\nline3\nline4\nline5"},
	}

	for _, cl := range contentLengths {
		b.Run(cl.name, func(b *testing.B) {
			i := &Inserter{}

			for n := 0; n < b.N; n++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_content_"+cl.name+"_"+string(rune(n))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = i.Insert(testFile, "500", cl.content)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkInserter_insertLine(b *testing.B) {
	// Test the core insertLine function without file I/O
	i := &Inserter{}

	sizes := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 10000},
		{"large", 100000},
	}

	for _, size := range sizes {
		content := strings.Repeat("test line\n", size.size)

		b.Run(size.name+"_beginning", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := i.insertLine(content, 1, "new line")
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(size.name+"_middle", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := i.insertLine(content, size.size/2, "new line")
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(size.name+"_end", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := i.insertLine(content, size.size+1, "new line")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkStringOperations_Insert(b *testing.B) {
	// Benchmark the core operations used in insert
	content := strings.Repeat("test line\n", 10000)
	lines := strings.Split(content, "\n")
	newLine := "inserted line"

	b.Run("strings_Split", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = strings.Split(content, "\n")
		}
	})

	b.Run("strings_Join", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = strings.Join(lines, "\n")
		}
	})

	b.Run("slice_append_beginning", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			result := append([]string{newLine}, lines...)
			_ = result
		}
	})

	b.Run("slice_append_middle", func(b *testing.B) {
		middle := len(lines) / 2
		for n := 0; n < b.N; n++ {
			result := make([]string, 0, len(lines)+1)
			result = append(result, lines[:middle]...)
			result = append(result, newLine)
			result = append(result, lines[middle:]...)
			_ = result
		}
	})

	b.Run("slice_append_end", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			result := append(lines, newLine)
			_ = result
		}
	})

	b.Run("strings_HasSuffix", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = strings.HasSuffix(content, "\n")
		}
	})
}

func BenchmarkInserter_parseLineNumber(b *testing.B) {
	i := &Inserter{}

	lineNumbers := []string{
		"1",
		"100",
		"10000",
		" 5000 ",
		"999999",
	}

	for _, lineNum := range lineNumbers {
		b.Run("line_"+lineNum, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := i.parseLineNumber(lineNum)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
