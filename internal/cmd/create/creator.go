package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Creator struct{}

func (c *Creator) Create(path, fileText string, showChanges, showResult bool) error {
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

	if showChanges {
		c.showContent(path, fileText)
	}

	if showResult {
		c.showResult(path, fileText)
	}

	fmt.Printf("Created file: %s (%d bytes)\n", path, len(fileText))
	return nil
}

func (c *Creator) showContent(path, content string) {
	fmt.Printf("\nContent of %s:\n", path)
	fmt.Println("--- New file")
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fmt.Printf("+%s\n", line)
	}
	fmt.Println()
}

func (c *Creator) showResult(path, content string) {
	fmt.Printf("\nResult of %s:\n", path)
	fmt.Println(content)
}
