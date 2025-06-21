package str_replace

import (
	"fmt"
	"os"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
)

type Replacer struct{}

func (r *Replacer) StrReplace(path, oldStr, newStr string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("cannot replace strings in directory: %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	original := string(content)
	modified := strings.ReplaceAll(original, oldStr, newStr)

	if original == modified {
		fmt.Printf("No occurrences of %q found in %s\n", oldStr, path)
		return nil
	}

	err = os.WriteFile(path, []byte(modified), info.Mode())
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	// Record edit for undo after file is written
	undoEditor := &undo_edit.UndoEditor{}
	err = undoEditor.RecordEdit(path, "str_replace", oldStr, newStr, -1)
	if err != nil {
		return fmt.Errorf("record edit: %w", err)
	}

	count := strings.Count(original, oldStr)
	fmt.Printf("Replaced %d occurrence(s) of %q with %q in %s\n", count, oldStr, newStr, path)
	return nil
}
