package insert

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
	"github.com/RRethy/eddie/internal/display"
	"github.com/RRethy/eddie/internal/fileops"
)

type Inserter struct {
	fileOps *fileops.FileOps
	display *display.Display
}

func NewInserter(w io.Writer) *Inserter {
	return &Inserter{
		fileOps: &fileops.FileOps{},
		display: display.New(w),
	}
}

func (i *Inserter) Insert(path, insertLine, newStr string, showChanges, showResult bool) error {
	original, info, err := i.fileOps.ReadFileContentForOperation(path, "insert line in")
	if err != nil {
		return err
	}

	lineNum, err := i.parseLineNumber(insertLine)
	if err != nil {
		return fmt.Errorf("parse line number: %w", err)
	}

	modified, err := i.insertLine(original, lineNum, newStr)
	if err != nil {
		return fmt.Errorf("insert line: %w", err)
	}

	if showChanges {
		i.display.ShowInsertDiff(path, original, modified, lineNum)
	}

	err = i.fileOps.WriteFileContent(path, modified, info.Mode())
	if err != nil {
		return err
	}

	if showResult {
		i.display.ShowResult(path, modified)
	}

	undoEditor := undo_edit.NewUndoEditor(os.Stdout)
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
