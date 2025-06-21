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

	// Get file info for permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	// Write the reversed content
	err = os.WriteFile(path, []byte(newContent), info.Mode())
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (u *UndoEditor) reverseStrReplace(content, oldStr, newStr string) (string, error) {
	// Count occurrences to verify the edit can be reversed
	newCount := strings.Count(content, newStr)
	if newCount == 0 {
		return "", fmt.Errorf("no occurrences of %q found to reverse", newStr)
	}

	// Replace back from new to old
	reversed := strings.ReplaceAll(content, newStr, oldStr)
	return reversed, nil
}

func (u *UndoEditor) reverseInsert(content string, lineNum int) (string, error) {
	lines := strings.Split(content, "\n")
	hasTrailingNewline := strings.HasSuffix(content, "\n")

	// Remove empty line at end if file ends with newline
	if hasTrailingNewline && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Validate line number
	if lineNum < 1 || lineNum > len(lines) {
		return "", fmt.Errorf("line number %d is out of range (1-%d)", lineNum, len(lines))
	}

	// Remove the line at the specified position
	var result []string
	if lineNum == 1 {
		// Remove first line
		result = lines[1:]
	} else if lineNum == len(lines) {
		// Remove last line
		result = lines[:len(lines)-1]
	} else {
		// Remove middle line
		result = make([]string, 0, len(lines)-1)
		result = append(result, lines[:lineNum-1]...)
		result = append(result, lines[lineNum:]...)
	}

	// Join with newlines and preserve original trailing newline behavior
	joined := strings.Join(result, "\n")
	if hasTrailingNewline {
		joined += "\n"
	}

	return joined, nil
}
