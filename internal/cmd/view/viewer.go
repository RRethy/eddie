package view

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Viewer struct{}

func (v *Viewer) View(path, viewRange string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return v.viewDir(path)
	}
	return v.viewFile(path, viewRange)
}

func (v *Viewer) viewDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", path, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		fmt.Println(name)
	}
	return nil
}

func (v *Viewer) viewFile(path, viewRange string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	start, end, err := v.parseRange(viewRange)
	if err != nil {
		return fmt.Errorf("parse range: %w", err)
	}

	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if start > 0 && line < start {
			line++
			continue
		}
		if end > 0 && line > end {
			break
		}
		fmt.Println(scanner.Text())
		line++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan file: %w", err)
	}
	return nil
}

func (v *Viewer) parseRange(viewRange string) (int, int, error) {
	if viewRange == "" {
		return 0, 0, nil
	}

	parts := strings.Split(viewRange, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format, expected start,end")
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start: %w", err)
	}

	endStr := strings.TrimSpace(parts[1])
	if endStr == "-1" {
		return start, -1, nil
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end: %w", err)
	}

	if start > end {
		return 0, 0, fmt.Errorf("start cannot be greater than end")
	}

	return start, end, nil
}
