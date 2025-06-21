package display

import (
	"fmt"
	"io"
	"strings"
)

type Display struct {
	w io.Writer
}

func New(w io.Writer) *Display {
	return &Display{w: w}
}

func (d *Display) ShowResult(path, content string) {
	fmt.Fprintf(d.w, "\nResult of %s:\n", path)
	fmt.Fprintln(d.w, content)
}

func (d *Display) ShowDiff(path, before, after string) {
	fmt.Fprintf(d.w, "\nChanges in %s:\n", path)
	fmt.Fprintln(d.w, "--- Before")
	fmt.Fprintln(d.w, "+++ After")
	
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")
	
	maxLines := len(beforeLines)
	if len(afterLines) > maxLines {
		maxLines = len(afterLines)
	}
	
	for i := 0; i < maxLines; i++ {
		beforeLine := ""
		afterLine := ""
		
		if i < len(beforeLines) {
			beforeLine = beforeLines[i]
		}
		if i < len(afterLines) {
			afterLine = afterLines[i]
		}
		
		if beforeLine != afterLine {
			if beforeLine != "" {
				fmt.Fprintf(d.w, "-%s\n", beforeLine)
			}
			if afterLine != "" {
				fmt.Fprintf(d.w, "+%s\n", afterLine)
			}
		}
	}
	fmt.Fprintln(d.w)
}

func (d *Display) ShowNewFileContent(path, content string) {
	fmt.Fprintf(d.w, "\nContent of %s:\n", path)
	fmt.Fprintln(d.w, "--- New file")
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fmt.Fprintf(d.w, "+%s\n", line)
	}
	fmt.Fprintln(d.w)
}

func (d *Display) ShowInsertDiff(path, original, modified string, lineNum int) {
	fmt.Fprintf(d.w, "\nChanges in %s:\n", path)
	fmt.Fprintln(d.w, "--- Before")
	fmt.Fprintln(d.w, "+++ After")
	
	origLines := strings.Split(original, "\n")
	modLines := strings.Split(modified, "\n")
	
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
			if i <= len(modLines) {
				fmt.Fprintf(d.w, "+%s\n", modLines[i-1])
			}
		} else {
			origIdx := i
			if i > lineNum {
				origIdx = i - 1
			}
			if origIdx <= len(origLines) && origIdx > 0 {
				fmt.Fprintf(d.w, " %s\n", origLines[origIdx-1])
			}
		}
	}
	fmt.Fprintln(d.w)
}