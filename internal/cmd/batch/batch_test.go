package batch

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    *BatchRequest
		wantErr bool
	}{
		{
			name: "valid JSON",
			json: `{"operations":[{"type":"view","path":"test.txt"}]}`,
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
				},
			},
		},
		{
			name:    "invalid JSON",
			json:    `{"operations":`,
			wantErr: true,
		},
		{
			name: "empty operations",
			json: `{"operations":[]}`,
			want: &BatchRequest{Operations: []Operation{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromJSON(tt.json)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    *BatchRequest
		wantErr bool
	}{
		{
			name:    "valid file",
			content: `{"operations":[{"type":"view","path":"test.txt"}]}`,
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
				},
			},
		},
		{
			name:    "invalid JSON in file",
			content: `{"operations":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tmpDir + "/test.json"
			err := os.WriteFile(path, []byte(tt.content), 0644)
			require.NoError(t, err)

			got, err := ParseFromFile(path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ParseFromFile("/nonexistent/file.json")
		assert.Error(t, err)
	})
}

func TestParseFromOps(t *testing.T) {
	tests := []struct {
		name    string
		ops     []string
		want    *BatchRequest
		wantErr bool
	}{
		{
			name: "view operation",
			ops:  []string{"view,test.txt"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
				},
			},
		},
		{
			name: "view with range",
			ops:  []string{"view,test.txt,1,10"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt", ViewRange: "1,10"},
				},
			},
		},
		{
			name: "str_replace operation",
			ops:  []string{"str_replace,test.txt,old,new"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "str_replace", Path: "test.txt", OldStr: "old", NewStr: "new"},
				},
			},
		},
		{
			name: "create operation",
			ops:  []string{"create,new.txt,content"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "create", Path: "new.txt", Content: "content"},
				},
			},
		},
		{
			name: "insert operation",
			ops:  []string{"insert,test.txt,5,new line"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "insert", Path: "test.txt", InsertLine: 5, NewStr: "new line"},
				},
			},
		},
		{
			name: "search operation",
			ops:  []string{"search,test.txt,query"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "search", Path: "test.txt", TreeQuery: "query"},
				},
			},
		},
		{
			name: "multiple operations",
			ops:  []string{"view,test.txt", "create,new.txt,content"},
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
					{Type: "create", Path: "new.txt", Content: "content"},
				},
			},
		},
		{
			name:    "invalid operation format",
			ops:     []string{"view"},
			wantErr: true,
		},
		{
			name:    "str_replace missing args",
			ops:     []string{"str_replace,test.txt,old"},
			wantErr: true,
		},
		{
			name:    "create missing content",
			ops:     []string{"create,test.txt"},
			wantErr: true,
		},
		{
			name:    "insert missing args",
			ops:     []string{"insert,test.txt,5"},
			wantErr: true,
		},
		{
			name:    "insert invalid line number",
			ops:     []string{"insert,test.txt,abc,content"},
			wantErr: true,
		},
		{
			name:    "search missing query",
			ops:     []string{"search,test.txt"},
			wantErr: true,
		},
		{
			name:    "unknown operation type",
			ops:     []string{"unknown,test.txt"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromOps(tt.ops)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessor_ProcessBatch(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"
	err := os.WriteFile(testFile, []byte("line1\nline2\nline3\n"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name string
		req  *BatchRequest
		want func(t *testing.T, resp *BatchResponse)
	}{
		{
			name: "successful view operation",
			req: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: testFile},
				},
			},
			want: func(t *testing.T, resp *BatchResponse) {
				require.Len(t, resp.Results, 1)
				result := resp.Results[0]
				assert.True(t, result.Success)
				assert.Nil(t, result.Error)
				assert.Equal(t, "view", result.Operation.Type)
				assert.Equal(t, testFile, result.Operation.Path)
			},
		},
		{
			name: "failed view operation",
			req: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "/nonexistent/file.txt"},
				},
			},
			want: func(t *testing.T, resp *BatchResponse) {
				require.Len(t, resp.Results, 1)
				result := resp.Results[0]
				assert.False(t, result.Success)
				assert.NotNil(t, result.Error)
				assert.Contains(t, *result.Error, "no such file or directory")
			},
		},
		{
			name: "unknown operation type",
			req: &BatchRequest{
				Operations: []Operation{
					{Type: "unknown", Path: testFile},
				},
			},
			want: func(t *testing.T, resp *BatchResponse) {
				require.Len(t, resp.Results, 1)
				result := resp.Results[0]
				assert.False(t, result.Success)
				assert.NotNil(t, result.Error)
				assert.Contains(t, *result.Error, "unknown operation type")
			},
		},
		{
			name: "mixed success and failure",
			req: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: testFile},
					{Type: "view", Path: "/nonexistent/file.txt"},
					{Type: "ls", Path: tmpDir},
				},
			},
			want: func(t *testing.T, resp *BatchResponse) {
				require.Len(t, resp.Results, 3)
				assert.True(t, resp.Results[0].Success)
				assert.False(t, resp.Results[1].Success)
				assert.True(t, resp.Results[2].Success)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			processor := NewProcessor(&buf)
			resp, err := processor.ProcessBatch(tt.req)
			require.NoError(t, err)
			tt.want(t, resp)
		})
	}
}

func TestParseFromStdin(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *BatchRequest
		wantErr bool
	}{
		{
			name:  "valid JSON from stdin",
			input: `{"operations":[{"type":"view","path":"test.txt"}]}`,
			want: &BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
				},
			},
		},
		{
			name:    "invalid JSON from stdin",
			input:   `{"operations":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			require.NoError(t, err)
			os.Stdin = r

			go func() {
				defer w.Close()
				_, _ = w.Write([]byte(tt.input))
			}()

			got, err := ParseFromStdin()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
