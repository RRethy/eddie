package insert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInserter_parseLineNumber(t *testing.T) {
	i := &Inserter{}

	tests := []struct {
		name       string
		insertLine string
		want       int
		wantErr    bool
	}{
		{"valid line 1", "1", 1, false},
		{"valid line 10", "10", 10, false},
		{"with spaces", " 5 ", 5, false},
		{"zero line", "0", 0, true},
		{"negative line", "-1", 0, true},
		{"invalid string", "abc", 0, true},
		{"empty string", "", 0, true},
		{"float number", "1.5", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.parseLineNumber(tt.insertLine)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInserter_insertLine(t *testing.T) {
	i := &Inserter{}

	tests := []struct {
		name    string
		content string
		lineNum int
		newStr  string
		want    string
		wantErr bool
	}{
		{
			name:    "insert at beginning",
			content: "line1\nline2\nline3\n",
			lineNum: 1,
			newStr:  "new first line",
			want:    "new first line\nline1\nline2\nline3\n",
			wantErr: false,
		},
		{
			name:    "insert in middle",
			content: "line1\nline2\nline3\n",
			lineNum: 2,
			newStr:  "inserted line",
			want:    "line1\ninserted line\nline2\nline3\n",
			wantErr: false,
		},
		{
			name:    "insert at end",
			content: "line1\nline2\nline3\n",
			lineNum: 4,
			newStr:  "new last line",
			want:    "line1\nline2\nline3\nnew last line\n",
			wantErr: false,
		},
		{
			name:    "insert in empty file",
			content: "",
			lineNum: 1,
			newStr:  "first line",
			want:    "first line\n",
			wantErr: false,
		},
		{
			name:    "insert in single line file with newline",
			content: "only line\n",
			lineNum: 1,
			newStr:  "new first",
			want:    "new first\nonly line\n",
			wantErr: false,
		},
		{
			name:    "insert in single line file without newline",
			content: "only line",
			lineNum: 2,
			newStr:  "second line",
			want:    "only line\nsecond line",
			wantErr: false,
		},
		{
			name:    "line number too high",
			content: "line1\nline2\n",
			lineNum: 5,
			newStr:  "too far",
			want:    "",
			wantErr: true,
		},
		{
			name:    "insert with empty string",
			content: "line1\nline2\n",
			lineNum: 2,
			newStr:  "",
			want:    "line1\n\nline2\n",
			wantErr: false,
		},
		{
			name:    "insert with multiline content",
			content: "line1\nline2\n",
			lineNum: 2,
			newStr:  "multi\nline\ninsert",
			want:    "line1\nmulti\nline\ninsert\nline2\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.insertLine(tt.content, tt.lineNum, tt.newStr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInserter_Insert(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		initialContent string
		insertLine     string
		newStr         string
		wantContent    string
		wantErr        bool
	}{
		{
			name:           "basic insertion",
			initialContent: "line1\nline2\nline3\n",
			insertLine:     "2",
			newStr:         "inserted line",
			wantContent:    "line1\ninserted line\nline2\nline3\n",
			wantErr:        false,
		},
		{
			name:           "insert at beginning",
			initialContent: "existing content\n",
			insertLine:     "1",
			newStr:         "new first line",
			wantContent:    "new first line\nexisting content\n",
			wantErr:        false,
		},
		{
			name:           "insert JSON line",
			initialContent: "{\n  \"key1\": \"value1\"\n}\n",
			insertLine:     "2",
			newStr:         "  \"key2\": \"value2\",",
			wantContent:    "{\n  \"key2\": \"value2\",\n  \"key1\": \"value1\"\n}\n",
			wantErr:        false,
		},
		{
			name:           "invalid line number",
			initialContent: "line1\nline2\n",
			insertLine:     "abc",
			newStr:         "new line",
			wantContent:    "",
			wantErr:        true,
		},
		{
			name:           "line number too high",
			initialContent: "line1\nline2\n",
			insertLine:     "10",
			newStr:         "new line",
			wantContent:    "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test_"+tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.initialContent), 0644))

			i := &Inserter{}
			err := i.Insert(testFile, tt.insertLine, tt.newStr, false, false)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			result, err := os.ReadFile(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.wantContent, string(result))
		})
	}
}

func TestInserter_Insert_Errors(t *testing.T) {
	tmpDir := t.TempDir()
	i := &Inserter{}

	tests := []struct {
		name    string
		setup   func() string
		line    string
		content string
		wantErr string
	}{
		{
			name: "nonexistent file",
			setup: func() string {
				return "/nonexistent/file.txt"
			},
			line:    "1",
			content: "content",
			wantErr: "stat",
		},
		{
			name: "directory instead of file",
			setup: func() string {
				return tmpDir
			},
			line:    "1",
			content: "content",
			wantErr: "cannot insert line in directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := i.Insert(path, tt.line, tt.content, false, false)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestInserter_Insert_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	i := &Inserter{}

	tests := []struct {
		name           string
		initialContent string
		insertLine     string
		newStr         string
		wantContent    string
	}{
		{
			name:           "file without trailing newline",
			initialContent: "line1\nline2",
			insertLine:     "2",
			newStr:         "inserted",
			wantContent:    "line1\ninserted\nline2",
		},
		{
			name:           "empty file",
			initialContent: "",
			insertLine:     "1",
			newStr:         "first line",
			wantContent:    "first line\n",
		},
		{
			name:           "single line without newline",
			initialContent: "single",
			insertLine:     "1",
			newStr:         "before",
			wantContent:    "before\nsingle",
		},
		{
			name:           "insert very long line",
			initialContent: "short\n",
			insertLine:     "1",
			newStr:         string(make([]byte, 10000)),
			wantContent:    string(make([]byte, 10000)) + "\nshort\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "edge_"+tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.initialContent), 0644))

			err := i.Insert(testFile, tt.insertLine, tt.newStr, false, false)
			require.NoError(t, err)

			result, err := os.ReadFile(testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.wantContent, string(result))
		})
	}
}
