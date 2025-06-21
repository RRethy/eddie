package str_replace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkReplacer_StrReplace(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files of different sizes
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
			// Create file with specified number of lines containing target string
			line := "this is a test line with hello in it\n"
			content := strings.Repeat(line, size.size)

			r := &Replacer{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_"+size.name+"_"+string(rune(i))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = r.StrReplace(testFile, "hello", "hi")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkReplacer_StrReplacePatterns(b *testing.B) {
	tmpDir := b.TempDir()

	// Create a medium-sized file
	line := "this is a test line with various patterns to replace\n"
	content := strings.Repeat(line, 10000)

	patterns := []struct {
		name   string
		oldStr string
		newStr string
	}{
		{"short_to_short", "is", "was"},
		{"short_to_long", "a", "absolutely"},
		{"long_to_short", "various patterns", "stuff"},
		{"long_to_long", "test line with", "example string containing"},
		{"no_match", "xyz", "abc"},
	}

	for _, pattern := range patterns {
		b.Run(pattern.name, func(b *testing.B) {
			r := &Replacer{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_pattern_"+pattern.name+"_"+string(rune(i))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = r.StrReplace(testFile, pattern.oldStr, pattern.newStr)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkReplacer_StrReplaceFrequency(b *testing.B) {
	tmpDir := b.TempDir()

	// Test different frequencies of target string
	frequencies := []struct {
		name       string
		targetFreq int // every Nth word is the target
	}{
		{"rare", 1000},    // 1 in 1000 words
		{"uncommon", 100}, // 1 in 100 words
		{"common", 10},    // 1 in 10 words
		{"frequent", 3},   // 1 in 3 words
	}

	for _, freq := range frequencies {
		b.Run(freq.name, func(b *testing.B) {
			// Create content with target string at specified frequency
			var contentBuilder strings.Builder
			for i := 0; i < 10000; i++ {
				if i%freq.targetFreq == 0 {
					contentBuilder.WriteString("target ")
				} else {
					contentBuilder.WriteString("word ")
				}
			}
			content := contentBuilder.String()

			r := &Replacer{}

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testFile := filepath.Join(tmpDir, "bench_freq_"+freq.name+"_"+string(rune(i))+".txt")
				err := os.WriteFile(testFile, []byte(content), 0644)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()

				err = r.StrReplace(testFile, "target", "replacement")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkStringOperations(b *testing.B) {
	// Benchmark the core string operations used in str_replace
	content := strings.Repeat("hello world hello universe hello galaxy\n", 10000)

	b.Run("strings_Count", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.Count(content, "hello")
		}
	})

	b.Run("strings_ReplaceAll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.ReplaceAll(content, "hello", "hi")
		}
	})

	b.Run("combined_operations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			modified := strings.ReplaceAll(content, "hello", "hi")
			_ = strings.Count(content, "hello")
			_ = modified
		}
	})
}
