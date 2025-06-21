package undo_edit

import (
	"fmt"
	"os"
	"strings"
)

func (u *UndoEditor) applyReverseEdit(path string, record *EditRecord) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	var newContent string
	switch record.EditType {
	case "str_replace":
		newContent, err = u.reverseStrReplace(string(content), record.OldContent, record.NewContent)
		if err != nil {
			return fmt.Errorf("reverse str_replace: %w", err)
		}
	case "insert":
		newContent, err = u.reverseInsert(string(content), record.Position)
		if err != nil {
			return fmt.Errorf("reverse insert: %w", err)
		}
	default:
		return fmt.Errorf("unknown edit type: %s", record.EditType)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	err = os.WriteFile(path, []byte(newContent), info.Mode())
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (u *UndoEditor) reverseStrReplace(content, oldStr, newStr string) (string, error) {
	newCount := strings.Count(content, newStr)
	if newCount == 0 {
		return "", fmt.Errorf("no occurrences of %q found to reverse", newStr)
	}

	reversed := strings.ReplaceAll(content, newStr, oldStr)
	return reversed, nil
}

func (u *UndoEditor) reverseInsert(content string, lineNum int) (string, error) {
	lines := strings.Split(content, "\n")
	hasTrailingNewline := strings.HasSuffix(content, "\n")

	if hasTrailingNewline && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if lineNum < 1 || lineNum > len(lines) {
		return "", fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(lines))
	}

	var result []string
	if lineNum == 1 {
		result = lines[1:]
	} else if lineNum == len(lines) {
		result = lines[:len(lines)-1]
	} else {
		result = make([]string, 0, len(lines)-1)
		result = append(result, lines[:lineNum-1]...)
		result = append(result, lines[lineNum:]...)
	}

	joined := strings.Join(result, "\n")
	if hasTrailingNewline {
		joined += "\n"
	}

	return joined, nil
}
