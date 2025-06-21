package e2e

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MCPRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]string      `json:"clientInfo"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

func TestMCPServerStartup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "eddie", "mcp")
	cmd.Dir = "../.."

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)

	initReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]interface{}{},
			ClientInfo:      map[string]string{"name": "test-client", "version": "1.0.0"},
		},
	}

	reqData, err := json.Marshal(initReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "2.0", resp.Jsonrpc)
		assert.Equal(t, 1, resp.ID)
		assert.NotNil(t, resp.Result)
		assert.Nil(t, resp.Error)
	} else {
		t.Fatal("No response received from MCP server")
	}
}

func TestMCPListTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "eddie", "mcp")
	cmd.Dir = "../.."

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)

	initReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]interface{}{},
			ClientInfo:      map[string]string{"name": "test-client", "version": "1.0.0"},
		},
	}

	reqData, err := json.Marshal(initReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	}

	listToolsReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	reqData, err = json.Marshal(listToolsReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "2.0", resp.Jsonrpc)
		assert.Equal(t, 2, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		resultBytes, err := json.Marshal(resp.Result)
		require.NoError(t, err)

		var result ListToolsResult
		err = json.Unmarshal(resultBytes, &result)
		require.NoError(t, err)

		expectedTools := []string{"view", "str_replace", "create", "insert", "undo_edit", "glob"}
		actualToolNames := make([]string, len(result.Tools))
		for i, tool := range result.Tools {
			actualToolNames[i] = tool.Name
		}

		for _, expected := range expectedTools {
			assert.Contains(t, actualToolNames, expected, "Tool %s should be available", expected)
		}

		for _, tool := range result.Tools {
			assert.NotEmpty(t, tool.Name)
			assert.NotEmpty(t, tool.Description)
			assert.NotNil(t, tool.InputSchema)
		}
	}
}

func TestMCPToolExecution(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello World"

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "eddie", "mcp")
	cmd.Dir = "../.."

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)

	initReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]interface{}{},
			ClientInfo:      map[string]string{"name": "test-client", "version": "1.0.0"},
		},
	}

	reqData, err := json.Marshal(initReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	}

	viewReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "view",
			"arguments": map[string]interface{}{
				"path": testFile,
			},
		},
	}

	reqData, err = json.Marshal(viewReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "2.0", resp.Jsonrpc)
		assert.Equal(t, 2, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		resultMap, ok := resp.Result.(map[string]interface{})
		require.True(t, ok)

		contentArray, ok := resultMap["content"].([]interface{})
		require.True(t, ok)
		require.Len(t, contentArray, 1)

		contentObj, ok := contentArray[0].(map[string]interface{})
		require.True(t, ok)

		text, ok := contentObj["text"].(string)
		require.True(t, ok)
		assert.Contains(t, text, "Hello World")
	}
}

func TestMCPGlobTool(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"test1.txt", "test2.txt", "main.go"}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(tmpDir, f), []byte("content"), 0644)
		require.NoError(t, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "eddie", "mcp")
	cmd.Dir = "../.."

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)

	initReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]interface{}{},
			ClientInfo:      map[string]string{"name": "test-client", "version": "1.0.0"},
		},
	}

	reqData, err := json.Marshal(initReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	}

	globReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "glob",
			"arguments": map[string]interface{}{
				"pattern": "*.txt",
				"path":    tmpDir,
			},
		},
	}

	reqData, err = json.Marshal(globReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "2.0", resp.Jsonrpc)
		assert.Equal(t, 2, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		resultMap, ok := resp.Result.(map[string]interface{})
		require.True(t, ok)

		contentArray, ok := resultMap["content"].([]interface{})
		require.True(t, ok)
		require.Len(t, contentArray, 1)

		contentObj, ok := contentArray[0].(map[string]interface{})
		require.True(t, ok)

		text, ok := contentObj["text"].(string)
		require.True(t, ok)

		lines := strings.Split(strings.TrimSpace(text), "\n")
		txtFiles := make([]string, 0)
		for _, line := range lines {
			if strings.HasSuffix(line, ".txt") {
				txtFiles = append(txtFiles, filepath.Base(line))
			}
		}

		assert.ElementsMatch(t, []string{"test1.txt", "test2.txt"}, txtFiles)
	}
}

func TestMCPErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "eddie", "mcp")
	cmd.Dir = "../.."

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)

	initReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]interface{}{},
			ClientInfo:      map[string]string{"name": "test-client", "version": "1.0.0"},
		},
	}

	reqData, err := json.Marshal(initReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)
		assert.Nil(t, resp.Error)
	}

	viewReq := MCPRequest{
		Jsonrpc: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "view",
			"arguments": map[string]interface{}{
				"path": "/nonexistent/file.txt",
			},
		},
	}

	reqData, err = json.Marshal(viewReq)
	require.NoError(t, err)

	_, err = stdin.Write(append(reqData, '\n'))
	require.NoError(t, err)

	if scanner.Scan() {
		var resp MCPResponse
		err = json.Unmarshal(scanner.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "2.0", resp.Jsonrpc)
		assert.Equal(t, 2, resp.ID)
		assert.Nil(t, resp.Error)
		assert.NotNil(t, resp.Result)

		resultMap, ok := resp.Result.(map[string]interface{})
		require.True(t, ok)

		isError, ok := resultMap["isError"].(bool)
		require.True(t, ok)
		assert.True(t, isError)
	}
}
