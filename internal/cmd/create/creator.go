package create

import (
	"fmt"
	"os"
	"path/filepath"
)

type Creator struct{}

func (c *Creator) Create(path, fileText string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("create directories %s: %w", dir, err)
		}
	}

	err := os.WriteFile(path, []byte(fileText), 0644)
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	fmt.Printf("Created file: %s (%d bytes)\n", path, len(fileText))
	return nil
}
