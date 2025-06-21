package ls

import (
	"fmt"
	"os"
)

type Lister struct{}

func (l *Lister) Ls(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", path, err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
	return nil
}
