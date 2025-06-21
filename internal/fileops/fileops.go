package fileops

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileOps struct{}

func (f *FileOps) ValidateFileExists(path string) (*os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	return &info, nil
}

func (f *FileOps) ValidateNotDir(path string, info os.FileInfo, operation string) error {
	if info.IsDir() {
		return fmt.Errorf("cannot %s directory: %s", operation, path)
	}
	return nil
}

func (f *FileOps) ReadFileContentForOperation(path, operation string) (string, os.FileInfo, error) {
	info, err := f.ValidateFileExists(path)
	if err != nil {
		return "", nil, err
	}

	if err := f.ValidateNotDir(path, *info, operation); err != nil {
		return "", nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read file %s: %w", path, err)
	}

	return string(content), *info, nil
}

func (f *FileOps) WriteFileContent(path, content string, mode os.FileMode) error {
	err := os.WriteFile(path, []byte(content), mode)
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

func (f *FileOps) CreateFile(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("create directories %s: %w", dir, err)
		}
	}

	return f.WriteFileContent(path, content, 0o644)
}
