package glob

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkGlobber(b *testing.B) {
	tmpDir := b.TempDir()
	
	for i := 0; i < 100; i++ {
		file, err := os.Create(tmpDir + "/" + "test" + string(rune('0'+i%10)) + ".txt")
		if err != nil {
			b.Fatal(err)
		}
		file.Close()
	}

	g := &GlobberSilent{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Glob("*.txt", tmpDir)
	}
}

type GlobberSilent struct{}

func (g *GlobberSilent) Glob(pattern, path string) error {
	if path == "" {
		path = "."
	}

	fullPattern := filepath.Join(path, pattern)
	_, err := filepath.Glob(fullPattern)
	if err != nil {
		return fmt.Errorf("glob %s: %w", fullPattern, err)
	}

	return nil
}