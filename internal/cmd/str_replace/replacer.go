package str_replace

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
	"github.com/RRethy/eddie/internal/display"
	"github.com/RRethy/eddie/internal/fileops"
)

type Replacer struct {
	fileOps *fileops.FileOps
	display *display.Display
}

func NewReplacer(w io.Writer) *Replacer {
	return &Replacer{
		fileOps: &fileops.FileOps{},
		display: display.New(w),
	}
}

func (r *Replacer) StrReplace(path, oldStr, newStr string, showChanges, showResult bool) error {
	original, info, err := r.fileOps.ReadFileContentForOperation(path, "replace strings in")
	if err != nil {
		return err
	}

	modified := strings.ReplaceAll(original, oldStr, newStr)

	if original == modified {
		fmt.Printf("No occurrences of %q found in %s\n", oldStr, path)
		return nil
	}

	if showChanges {
		r.display.ShowDiff(path, original, modified)
	}

	err = r.fileOps.WriteFileContent(path, modified, info.Mode())
	if err != nil {
		return err
	}

	if showResult {
		r.display.ShowResult(path, modified)
	}

	undoEditor := undo_edit.NewUndoEditor(os.Stdout)
	err = undoEditor.RecordEdit(path, "str_replace", oldStr, newStr, -1)
	if err != nil {
		return fmt.Errorf("record edit: %w", err)
	}

	count := strings.Count(original, oldStr)
	fmt.Printf("Replaced %d occurrence(s) of %q with %q in %s\n", count, oldStr, newStr, path)
	return nil
}
