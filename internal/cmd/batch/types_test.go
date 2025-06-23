package batch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchRequest_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		req  BatchRequest
		want string
	}{
		{
			name: "empty operations",
			req:  BatchRequest{Operations: []Operation{}},
			want: `{"operations":[]}`,
		},
		{
			name: "single view operation",
			req: BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt", ViewRange: "1,10"},
				},
			},
			want: `{"operations":[{"type":"view","path":"test.txt","view_range":"1,10"}]}`,
		},
		{
			name: "multiple operations",
			req: BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
					{Type: "str_replace", Path: "test.txt", OldStr: "old", NewStr: "new", ShowChanges: true},
					{Type: "create", Path: "new.txt", Content: "hello"},
				},
			},
			want: `{"operations":[{"type":"view","path":"test.txt"},{"type":"str_replace","path":"test.txt","old_str":"old","new_str":"new","show_changes":true},{"type":"create","path":"new.txt","content":"hello"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestBatchRequest_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    BatchRequest
		wantErr bool
	}{
		{
			name: "empty operations",
			json: `{"operations":[]}`,
			want: BatchRequest{Operations: []Operation{}},
		},
		{
			name: "single view operation",
			json: `{"operations":[{"type":"view","path":"test.txt","view_range":"1,10"}]}`,
			want: BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt", ViewRange: "1,10"},
				},
			},
		},
		{
			name: "multiple operations",
			json: `{"operations":[{"type":"view","path":"test.txt"},{"type":"str_replace","path":"test.txt","old_str":"old","new_str":"new","show_changes":true},{"type":"create","path":"new.txt","content":"hello"}]}`,
			want: BatchRequest{
				Operations: []Operation{
					{Type: "view", Path: "test.txt"},
					{Type: "str_replace", Path: "test.txt", OldStr: "old", NewStr: "new", ShowChanges: true},
					{Type: "create", Path: "new.txt", Content: "hello"},
				},
			},
		},
		{
			name:    "invalid JSON",
			json:    `{"operations":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BatchRequest
			err := json.Unmarshal([]byte(tt.json), &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBatchResponse_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		resp BatchResponse
		want string
	}{
		{
			name: "empty results",
			resp: BatchResponse{Results: []OperationResult{}},
			want: `{"results":[]}`,
		},
		{
			name: "successful operation",
			resp: BatchResponse{
				Results: []OperationResult{
					{
						Operation: Operation{Type: "view", Path: "test.txt"},
						Success:   true,
						Output:    "file content",
						Error:     nil,
					},
				},
			},
			want: `{"results":[{"operation":{"type":"view","path":"test.txt"},"success":true,"output":"file content","error":null}]}`,
		},
		{
			name: "failed operation",
			resp: BatchResponse{
				Results: []OperationResult{
					{
						Operation: Operation{Type: "view", Path: "missing.txt"},
						Success:   false,
						Output:    "",
						Error:     stringPtr("file not found"),
					},
				},
			},
			want: `{"results":[{"operation":{"type":"view","path":"missing.txt"},"success":false,"output":"","error":"file not found"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.resp)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestBatchResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    BatchResponse
		wantErr bool
	}{
		{
			name: "empty results",
			json: `{"results":[]}`,
			want: BatchResponse{Results: []OperationResult{}},
		},
		{
			name: "successful operation",
			json: `{"results":[{"operation":{"type":"view","path":"test.txt"},"success":true,"output":"file content","error":null}]}`,
			want: BatchResponse{
				Results: []OperationResult{
					{
						Operation: Operation{Type: "view", Path: "test.txt"},
						Success:   true,
						Output:    "file content",
						Error:     nil,
					},
				},
			},
		},
		{
			name: "failed operation",
			json: `{"results":[{"operation":{"type":"view","path":"missing.txt"},"success":false,"output":"","error":"file not found"}]}`,
			want: BatchResponse{
				Results: []OperationResult{
					{
						Operation: Operation{Type: "view", Path: "missing.txt"},
						Success:   false,
						Output:    "",
						Error:     stringPtr("file not found"),
					},
				},
			},
		},
		{
			name:    "invalid JSON",
			json:    `{"results":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BatchResponse
			err := json.Unmarshal([]byte(tt.json), &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
