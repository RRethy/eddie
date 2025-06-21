package ls

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func BenchmarkLister_Ls(b *testing.B) {
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
			testDir := b.TempDir()

			for i := 0; i < size.files; i++ {
				fileName := filepath.Join(testDir, "file"+strconv.Itoa(i)+".txt")
				err := os.WriteFile(fileName, []byte("content"), 0644)
				if err != nil {
					b.Fatal(err)
				}
			}

			l := &Lister{}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := l.Ls(testDir)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
