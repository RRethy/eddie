package str_replace

import (
	"fmt"
	"os"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
)

type Replacer struct{}

func (r *Replacer) StrReplace(path, oldStr, newStr string, showChanges, showResult bool) error {
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

	if showChanges {
		r.showDiff(path, original, modified)
	}

	err = os.WriteFile(path, []byte(modified), info.Mode())
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	if showResult {
		r.showResult(path, modified)
	}

	undoEditor := &undo_edit.UndoEditor{}
	err = undoEditor.RecordEdit(path, "str_replace", oldStr, newStr, -1)
	if err != nil {
		return fmt.Errorf("record edit: %w", err)
	}

	count := strings.Count(original, oldStr)
	fmt.Printf("Replaced %d occurrence(s) of %q with %q in %s\n", count, oldStr, newStr, path)
	return nil
}

func (r *Replacer) showDiff(path, original, modified string) {
	fmt.Printf("\nChanges in %s:\n", path)
	fmt.Println("--- Before")
	fmt.Println("+++ After")
	
	origLines := strings.Split(original, "\n")
	modLines := strings.Split(modified, "\n")
	
	maxLines := len(origLines)
	if len(modLines) > maxLines {
		maxLines = len(modLines)
	}
	
	for i := 0; i < maxLines; i++ {
		origLine := ""
		modLine := ""
		
		if i < len(origLines) {
			origLine = origLines[i]
		}
		if i < len(modLines) {
			modLine = modLines[i]
		}
		
		if origLine != modLine {
			if origLine != "" {
				fmt.Printf("-%s\n", origLine)
			}
			if modLine != "" {
				fmt.Printf("+%s\n", modLine)
			}
		}
	}
	fmt.Println()
}

func (r *Replacer) showResult(path, content string) {
	fmt.Printf("\nResult of %s:\n", path)
	fmt.Println(content)
}
