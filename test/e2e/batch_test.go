package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RRethy/eddie/internal/cmd/batch"
)

func TestBatchCommand(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("line1\nline2\nline3\n"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name  string
		args  []string
		input string
		want  func(t *testing.T, output string, err error)
	}{
		{
			name:  "batch from stdin",
			args:  []string{"batch"},
			input: `{"operations":[{"type":"view","path":"` + testFile + `"}]}`,
			want: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)

				var resp batch.BatchResponse
				unmarshalErr := json.Unmarshal([]byte(output), &resp)
				require.NoError(t, unmarshalErr)

				require.Len(t, resp.Results, 1)
				assert.True(t, resp.Results[0].Success)
				assert.Nil(t, resp.Results[0].Error)
				assert.Equal(t, "view", resp.Results[0].Operation.Type)
			},
		},
		{
			name: "batch from JSON flag",
			args: []string{"batch", "--json", `{"operations":[{"type":"ls","path":"` + tmpDir + `"}]}`},
			want: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)

				var resp batch.BatchResponse
				unmarshalErr := json.Unmarshal([]byte(output), &resp)
				require.NoError(t, unmarshalErr)

				require.Len(t, resp.Results, 1)
				assert.True(t, resp.Results[0].Success)
				assert.Nil(t, resp.Results[0].Error)
				assert.Equal(t, "ls", resp.Results[0].Operation.Type)
			},
		},
		{
			name: "batch from op flags",
			args: []string{"batch", "--op", "view," + testFile, "--op", "ls," + tmpDir},
			want: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)

				var resp batch.BatchResponse
				unmarshalErr := json.Unmarshal([]byte(output), &resp)
				require.NoError(t, unmarshalErr)

				require.Len(t, resp.Results, 2)
				assert.True(t, resp.Results[0].Success)
				assert.True(t, resp.Results[1].Success)
				assert.Equal(t, "view", resp.Results[0].Operation.Type)
				assert.Equal(t, "ls", resp.Results[1].Operation.Type)
			},
		},
		{
			name: "batch with error continues execution",
			args: []string{"batch", "--op", "view,/nonexistent/file.txt", "--op", "ls," + tmpDir},
			want: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)

				var resp batch.BatchResponse
				unmarshalErr := json.Unmarshal([]byte(output), &resp)
				require.NoError(t, unmarshalErr)

				require.Len(t, resp.Results, 2)
				assert.False(t, resp.Results[0].Success)
				assert.NotNil(t, resp.Results[0].Error)
				assert.True(t, resp.Results[1].Success)
				assert.Nil(t, resp.Results[1].Error)
			},
		},
		{
			name:  "invalid JSON input",
			args:  []string{"batch"},
			input: `{"operations":`,
			want: func(t *testing.T, output string, err error) {
				assert.Error(t, err)
				assert.Contains(t, output, "Error parsing input")
			},
		},
		{
			name: "multiple input methods error",
			args: []string{"batch", "--json", `{"operations":[]}`, "--op", "view,test.txt"},
			want: func(t *testing.T, output string, err error) {
				assert.Error(t, err)
				assert.Contains(t, output, "only one input method allowed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runEddieWithInput(t, tt.input, tt.args...)
			tt.want(t, output, err)
		})
	}
}

func TestBatchCommandFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("line1\nline2\nline3\n"), 0644)
	require.NoError(t, err)

	batchFile := filepath.Join(tmpDir, "batch.json")
	batchContent := `{"operations":[{"type":"view","path":"` + testFile + `"},{"type":"ls","path":"` + tmpDir + `"}]}`
	err = os.WriteFile(batchFile, []byte(batchContent), 0644)
	require.NoError(t, err)

	stdout, stderr, err := runEddie(t, "batch", "--file", batchFile)
	assert.NoError(t, err)
	output := stdout + stderr

	var resp batch.BatchResponse
	unmarshalErr := json.Unmarshal([]byte(output), &resp)
	require.NoError(t, unmarshalErr)

	require.Len(t, resp.Results, 2)
	assert.True(t, resp.Results[0].Success)
	assert.True(t, resp.Results[1].Success)
	assert.Equal(t, "view", resp.Results[0].Operation.Type)
	assert.Equal(t, "ls", resp.Results[1].Operation.Type)
}

func TestBatchCommandFileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("old content\nline2\n"), 0644)
	require.NoError(t, err)

	newFile := filepath.Join(tmpDir, "new.txt")

	operations := []string{
		"view," + testFile,
		"str_replace," + testFile + ",old content,new content",
		"create," + newFile + ",hello world",
		"insert," + testFile + ",1,inserted line",
		"view," + testFile,
	}

	args := []string{"batch"}
	for _, op := range operations {
		args = append(args, "--op", op)
	}

	stdout, stderr, err := runEddie(t, args...)
	assert.NoError(t, err)
	output := stdout + stderr

	var resp batch.BatchResponse
	unmarshalErr := json.Unmarshal([]byte(output), &resp)
	require.NoError(t, unmarshalErr)

	require.Len(t, resp.Results, 5)

	for i, result := range resp.Results {
		assert.True(t, result.Success, "Operation %d should succeed", i)
		assert.Nil(t, result.Error, "Operation %d should not have error", i)
	}

	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "new content")
	assert.Contains(t, string(content), "inserted line")

	content, err = os.ReadFile(newFile)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(content))
}

func runEddieWithInput(t *testing.T, input string, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command("../../eddie", args...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}
