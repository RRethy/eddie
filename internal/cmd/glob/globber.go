package glob

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Globber struct{}

type fileInfo struct {
	modTime time.Time
	path    string
}

func (g *Globber) Glob(pattern, path string) error {
	if path == "" {
		path = "."
	}

	var matches []string
	var err error

	if strings.Contains(pattern, "**") {
		if path == "." {
			matches, err = g.recursiveGlobFromPattern(pattern)
		} else {
			matches, err = g.recursiveGlob(pattern, path)
		}
	} else {
		fullPattern := filepath.Join(path, pattern)
		matches, err = filepath.Glob(fullPattern)
	}

	if err != nil {
		return fmt.Errorf("glob %s: %w", pattern, err)
	}

	if len(matches) == 0 {
		return nil
	}

	files := make([]fileInfo, 0, len(matches))
	for _, match := range matches {
		info, err := os.Lstat(match)
		if err != nil {
			continue
		}
		files = append(files, fileInfo{
			path:    match,
			modTime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})

	for _, f := range files {
		fmt.Println(f.path)
	}

	return nil
}

func (g *Globber) recursiveGlobFromPattern(pattern string) ([]string, error) {
	parts := strings.Split(pattern, "**")
	if len(parts) != 2 {
		return g.recursiveGlob(pattern, ".")
	}

	basePart := strings.TrimSuffix(parts[0], "/")
	suffixPart := strings.TrimPrefix(parts[1], "/")

	if basePart == "" {
		return g.recursiveGlob(pattern, ".")
	}

	var matches []string
	dirsOnly := strings.HasSuffix(pattern, "/")

	err := filepath.WalkDir(basePart, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if dirsOnly && !d.IsDir() {
			return nil
		}

		var testPattern string
		if suffixPart == "" {
			testPattern = basePart + "/**"
		} else {
			testPattern = basePart + "/**/" + suffixPart
		}

		matched, err := g.matchDoubleStarPattern(testPattern, path)
		if err != nil {
			return nil
		}

		if matched {
			if path == basePart && d.IsDir() && !dirsOnly {
				matches = append(matches, path+"/")
			} else if dirsOnly {
				matches = append(matches, path+"/")
			} else {
				matches = append(matches, path)
			}
		}

		return nil
	})

	return matches, err
}

func (g *Globber) recursiveGlob(pattern, basePath string) ([]string, error) {
	var matches []string
	dirsOnly := strings.HasSuffix(pattern, "/")

	err := filepath.WalkDir(basePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if dirsOnly && !d.IsDir() {
			return nil
		}

		matched, err := g.matchDoubleStarPattern(pattern, path)
		if err != nil {
			return nil
		}

		if matched {
			if path == basePath && d.IsDir() && !dirsOnly {
				matches = append(matches, path+"/")
			} else if dirsOnly {
				matches = append(matches, path+"/")
			} else {
				matches = append(matches, path)
			}
		}

		return nil
	})

	return matches, err
}

func (g *Globber) matchDoubleStarPattern(pattern, path string) (bool, error) {
	if !strings.Contains(pattern, "**") {
		return filepath.Match(pattern, path)
	}

	parts := strings.Split(pattern, "**")
	if len(parts) != 2 {
		return filepath.Match(strings.ReplaceAll(pattern, "**", "*"), path)
	}

	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := strings.TrimPrefix(parts[1], "/")

	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false, nil
	}

	if suffix == "" {
		if prefix == "" {
			return true, nil
		}
		return strings.HasPrefix(path, prefix), nil
	}

	if prefix == "" {
		pathParts := strings.Split(path, "/")
		for _, part := range pathParts {
			matched, err := filepath.Match(suffix, part)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	}

	if !strings.HasPrefix(path, prefix) {
		return false, nil
	}

	remaining := strings.TrimPrefix(path, prefix)
	remaining = strings.TrimPrefix(remaining, "/")

	if remaining == "" {
		return false, nil
	}

	pathParts := strings.Split(remaining, "/")
	for _, part := range pathParts {
		matched, err := filepath.Match(suffix, part)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}
