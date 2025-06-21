package create

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkCreator(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 10000},
		{"large", 1000000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			tmpDir := b.TempDir()
			content := strings.Repeat("x", tt.size)
			c := &Creator{}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				path := filepath.Join(tmpDir, fmt.Sprintf("file_%d.txt", i))
				if err := c.Create(path, content); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCreatorConcurrent(b *testing.B) {
	tmpDir := b.TempDir()
	c := &Creator{}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			path := filepath.Join(tmpDir, fmt.Sprintf("concurrent_%d.txt", i))
			if err := c.Create(path, "test content"); err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}
