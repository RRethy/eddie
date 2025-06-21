package insert

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
)

type Inserter struct{}

func (i *Inserter) Insert(path, insertLine, newStr string, showChanges bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("cannot insert line in directory: %s", path)
	}

	lineNum, err := i.parseLineNumber(insertLine)
	if err != nil {
		return fmt.Errorf("parse line number: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	original := string(content)
	modified, err := i.insertLine(original, lineNum, newStr)
	if err != nil {
		return fmt.Errorf("insert line: %w", err)
	}

	if showChanges {
		i.showDiff(path, original, modified, lineNum)
	}

	err = os.WriteFile(path, []byte(modified), info.Mode())
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	undoEditor := &undo_edit.UndoEditor{}
	err = undoEditor.RecordEdit(path, "insert", "", newStr, lineNum)
	if err != nil {
		return fmt.Errorf("record edit: %w", err)
	}

	fmt.Printf("Inserted line at position %d in %s\n", lineNum, path)
	return nil
}

func (i *Inserter) parseLineNumber(insertLine string) (int, error) {
	lineNum, err := strconv.Atoi(strings.TrimSpace(insertLine))
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %w", err)
	}

	if lineNum < 1 {
		return 0, fmt.Errorf("line number must be >= 1, got %d", lineNum)
	}

	return lineNum, nil
}

func (i *Inserter) insertLine(content string, lineNum int, newStr string) (string, error) {
	if content == "" {
		return newStr + "\n", nil
	}

	lines := strings.Split(content, "\n")
	hasTrailingNewline := strings.HasSuffix(content, "\n")

	if hasTrailingNewline && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if lineNum > len(lines)+1 {
		return "", fmt.Errorf("line number %d exceeds file length (%d lines)", lineNum, len(lines))
	}

	var result []string

	if lineNum == 1 {
		result = append([]string{newStr}, lines...)
	} else if lineNum > len(lines) {
		result = append(lines, newStr)
	} else {
		result = make([]string, 0, len(lines)+1)
		result = append(result, lines[:lineNum-1]...)
		result = append(result, newStr)
		result = append(result, lines[lineNum-1:]...)
	}

	joined := strings.Join(result, "\n")
	if hasTrailingNewline {
		joined += "\n"
	}

	return joined, nil
}

func (i *Inserter) showDiff(path, original, modified string, lineNum int) {
	fmt.Printf("\nChanges in %s:\n", path)
	fmt.Println("--- Before")
	fmt.Println("+++ After")
	
	origLines := strings.Split(original, "\n")
	modLines := strings.Split(modified, "\n")
	
	// Show context around the insertion point
	start := lineNum - 3
	if start < 1 {
		start = 1
	}
	end := lineNum + 3
	if end > len(modLines) {
		end = len(modLines)
	}
	
	for i := start; i <= end; i++ {
		if i == lineNum {
			// Show the inserted line
			if i <= len(modLines) {
				fmt.Printf("+%s\n", modLines[i-1])
			}
		} else {
			// Show existing lines for context
			origIdx := i
			if i > lineNum {
				origIdx = i - 1 // Adjust for insertion
			}
			if origIdx <= len(origLines) && origIdx > 0 {
				fmt.Printf(" %s\n", origLines[origIdx-1])
			}
		}
	}
	fmt.Println()
}
