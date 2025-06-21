package display

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplay_ShowResult(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		content  string
		expected []string
	}{
		{
			name:    "simple content",
			path:    "/test/file.txt",
			content: "hello world",
			expected: []string{
				"Result of /test/file.txt:",
				"hello world",
			},
		},
		{
			name:    "multiline content",
			path:    "config.json",
			content: "{\n  \"key\": \"value\"\n}",
			expected: []string{
				"Result of config.json:",
				"{\n  \"key\": \"value\"\n}",
			},
		},
		{
			name:    "empty content",
			path:    "empty.txt",
			content: "",
			expected: []string{
				"Result of empty.txt:",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			d := New(&buf)
			d.ShowResult(tt.path, tt.content)

			output := buf.String()
			for _, exp := range tt.expected {
				assert.Contains(t, output, exp)
			}
		})
	}
}

func TestDisplay_ShowDiff(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		before   string
		after    string
		expected []string
	}{
		{
			name:   "single line change",
			path:   "test.txt",
			before: "hello world",
			after:  "hello there",
			expected: []string{
				"Changes in test.txt:",
				"--- Before",
				"+++ After",
				"-hello world",
				"+hello there",
			},
		},
		{
			name:   "multiline changes",
			path:   "config.txt",
			before: "line1\nline2\nline3",
			after:  "line1\nmodified\nline3",
			expected: []string{
				"Changes in config.txt:",
				"--- Before",
				"+++ After",
				"-line2",
				"+modified",
			},
		},
		{
			name:   "addition",
			path:   "add.txt",
			before: "existing",
			after:  "existing\nnew line",
			expected: []string{
				"Changes in add.txt:",
				"+new line",
			},
		},
		{
			name:   "deletion",
			path:   "del.txt",
			before: "keep\nremove",
			after:  "keep",
			expected: []string{
				"Changes in del.txt:",
				"-remove",
			},
		},
		{
			name:   "no changes",
			path:   "same.txt",
			before: "unchanged",
			after:  "unchanged",
			expected: []string{
				"Changes in same.txt:",
				"--- Before",
				"+++ After",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			d := New(&buf)
			d.ShowDiff(tt.path, tt.before, tt.after)

			output := buf.String()
			for _, exp := range tt.expected {
				assert.Contains(t, output, exp)
			}
		})
	}
}

func TestDisplay_ShowNewFileContent(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		content  string
		expected []string
	}{
		{
			name:    "single line",
			path:    "new.txt",
			content: "hello world",
			expected: []string{
				"Content of new.txt:",
				"--- New file",
				"+hello world",
			},
		},
		{
			name:    "multiline content",
			path:    "multi.txt",
			content: "line1\nline2\nline3",
			expected: []string{
				"Content of multi.txt:",
				"--- New file",
				"+line1",
				"+line2",
				"+line3",
			},
		},
		{
			name:    "empty file",
			path:    "empty.txt",
			content: "",
			expected: []string{
				"Content of empty.txt:",
				"--- New file",
				"+",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			d := New(&buf)
			d.ShowNewFileContent(tt.path, tt.content)

			output := buf.String()
			for _, exp := range tt.expected {
				assert.Contains(t, output, exp)
			}
		})
	}
}

func TestDisplay_ShowInsertDiff(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		original string
		modified string
		expected []string
		lineNum  int
	}{
		{
			name:     "insert at beginning",
			path:     "test.txt",
			original: "line2\nline3\nline4",
			modified: "inserted\nline2\nline3\nline4",
			lineNum:  1,
			expected: []string{
				"Changes in test.txt:",
				"--- Before",
				"+++ After",
				"+inserted",
				" line2",
				" line3",
			},
		},
		{
			name:     "insert in middle",
			path:     "mid.txt",
			original: "line1\nline2\nline3\nline4\nline5",
			modified: "line1\nline2\ninserted\nline3\nline4\nline5",
			lineNum:  3,
			expected: []string{
				"Changes in mid.txt:",
				"+inserted",
				" line2",
				" line3",
			},
		},
		{
			name:     "insert at end",
			path:     "end.txt",
			original: "line1\nline2",
			modified: "line1\nline2\ninserted",
			lineNum:  3,
			expected: []string{
				"Changes in end.txt:",
				"+inserted",
				" line2",
			},
		},
		{
			name:     "insert with context boundaries",
			path:     "boundary.txt",
			original: "line1",
			modified: "inserted\nline1",
			lineNum:  1,
			expected: []string{
				"Changes in boundary.txt:",
				"+inserted",
				" line1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			d := New(&buf)
			d.ShowInsertDiff(tt.path, tt.original, tt.modified, tt.lineNum)

			output := buf.String()
			for _, exp := range tt.expected {
				assert.Contains(t, output, exp)
			}
		})
	}
}

func TestDisplay_ShowInsertDiff_ContextWindow(t *testing.T) {
	original := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9"
	modified := "line1\nline2\nline3\nline4\ninserted\nline5\nline6\nline7\nline8\nline9"
	lineNum := 5

	var buf bytes.Buffer
	d := New(&buf)
	d.ShowInsertDiff("context.txt", original, modified, lineNum)
	output := buf.String()

	assert.Contains(t, output, " line2")
	assert.Contains(t, output, " line3")
	assert.Contains(t, output, " line4")
	assert.Contains(t, output, "+inserted")
	assert.Contains(t, output, " line5")
	assert.Contains(t, output, " line6")
	assert.Contains(t, output, " line7")

	assert.NotContains(t, output, " line1")
	assert.NotContains(t, output, " line8")
}
