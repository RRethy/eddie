package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/RRethy/eddie/internal/cmd/batch"
	"github.com/RRethy/eddie/internal/cmd/create"
	"github.com/RRethy/eddie/internal/cmd/glob"
	"github.com/RRethy/eddie/internal/cmd/insert"
	"github.com/RRethy/eddie/internal/cmd/ls"
	"github.com/RRethy/eddie/internal/cmd/search"
	"github.com/RRethy/eddie/internal/cmd/str_replace"
	"github.com/RRethy/eddie/internal/cmd/undo_edit"
	"github.com/RRethy/eddie/internal/cmd/view"
)

type McpServer struct{}

func (m *McpServer) Mcp() error {
	s := server.NewMCPServer(
		"Eddie MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(*m.createViewTool(), m.handleView)
	s.AddTool(*m.createStrReplaceTool(), m.handleStrReplace)
	s.AddTool(*m.createCreateTool(), m.handleCreate)
	s.AddTool(*m.createInsertTool(), m.handleInsert)
	s.AddTool(*m.createUndoEditTool(), m.handleUndoEdit)
	s.AddTool(*m.createGlobTool(), m.handleGlob)
	s.AddTool(*m.createLsTool(), m.handleLs)
	s.AddTool(*m.createSearchTool(), m.handleSearch)
	s.AddTool(*m.createBatchTool(), m.handleBatch)

	return server.ServeStdio(s)
}

func (m *McpServer) createViewTool() *mcp.Tool {
	tool := mcp.NewTool("view",
		mcp.WithDescription("View file contents or list directory contents"),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path to the file or directory to view")),
		mcp.WithString("range", mcp.Description("Range of lines to view in format \"start,end\". If \"end\" is -1, reads to end of file. Ignored for directories.")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	return &tool
}

func (m *McpServer) createStrReplaceTool() *mcp.Tool {
	tool := mcp.NewTool("str_replace",
		mcp.WithDescription("Replace all occurrences of a string in a file"),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path to the file to modify")),
		mcp.WithString("old_str", mcp.Required(), mcp.Description("The string to search for and replace")),
		mcp.WithString("new_str", mcp.Required(), mcp.Description("The string to replace old_str with")),
		mcp.WithBoolean("show_changes", mcp.Description("Show the changes made to the file")),
		mcp.WithBoolean("show_result", mcp.Description("Show the new content after the edit operation")),
	)
	return &tool
}

func (m *McpServer) createCreateTool() *mcp.Tool {
	tool := mcp.NewTool("create",
		mcp.WithDescription("Create a new file with specified content"),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path where the new file should be created")),
		mcp.WithString("content", mcp.Required(), mcp.Description("The content to write to the new file")),
		mcp.WithBoolean("show_changes", mcp.Description("Show the content of the created file")),
		mcp.WithBoolean("show_result", mcp.Description("Show the new content after the file creation")),
	)
	return &tool
}

func (m *McpServer) createInsertTool() *mcp.Tool {
	tool := mcp.NewTool("insert",
		mcp.WithDescription("Insert a new line at specified line number"),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path to the file to modify")),
		mcp.WithNumber("line", mcp.Required(), mcp.Description("The line number where the new line should be inserted (1-based)")),
		mcp.WithString("content", mcp.Required(), mcp.Description("The content of the new line to insert")),
		mcp.WithBoolean("show_changes", mcp.Description("Show the changes made to the file")),
		mcp.WithBoolean("show_result", mcp.Description("Show the new content after the edit operation")),
	)
	return &tool
}

func (m *McpServer) createUndoEditTool() *mcp.Tool {
	tool := mcp.NewTool("undo_edit",
		mcp.WithDescription("Undo the last edit operation on a file"),
		mcp.WithString("path", mcp.Required(), mcp.Description("The path to the file to restore from backup")),
		mcp.WithBoolean("show_changes", mcp.Description("Show the changes made during the undo operation")),
		mcp.WithBoolean("show_result", mcp.Description("Show the new content after the undo operation")),
	)
	return &tool
}

func (m *McpServer) createGlobTool() *mcp.Tool {
	tool := mcp.NewTool("glob",
		mcp.WithDescription("Fast file pattern matching tool that works with any codebase size"),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("The glob pattern to match files against")),
		mcp.WithString("path", mcp.Description("The directory to search in. If not specified, the current working directory will be used. IMPORTANT: Omit this field to use the default directory. DO NOT enter \"undefined\" or \"null\" - simply omit it for the default behavior. Must be a valid directory path if provided.")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	return &tool
}

func (m *McpServer) handleView(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}

	rangeStr := ""
	if r, ok := args["range"].(string); ok {
		rangeStr = r
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := view.View(path, rangeStr)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(output),
		},
	}, nil
}

func (m *McpServer) handleStrReplace(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}
	oldStr, ok := args["old_str"].(string)
	if !ok {
		return nil, fmt.Errorf("old_str parameter required")
	}
	newStr, ok := args["new_str"].(string)
	if !ok {
		return nil, fmt.Errorf("new_str parameter required")
	}

	showChanges := false
	if sc, ok := args["show_changes"].(bool); ok {
		showChanges = sc
	}

	showResult := false
	if sr, ok := args["show_result"].(bool); ok {
		showResult = sr
	}

	err := str_replace.StrReplace(path, oldStr, newStr, showChanges, showResult)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("String replacement completed successfully"),
		},
	}, nil
}

func (m *McpServer) handleCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}
	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter required")
	}

	showChanges := false
	if sc, ok := args["show_changes"].(bool); ok {
		showChanges = sc
	}

	showResult := false
	if sr, ok := args["show_result"].(bool); ok {
		showResult = sr
	}

	err := create.Create(path, content, showChanges, showResult)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("File created successfully"),
		},
	}, nil
}

func (m *McpServer) handleInsert(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}
	lineFloat, ok := args["line"].(float64)
	if !ok {
		return nil, fmt.Errorf("line parameter required")
	}
	line := strconv.Itoa(int(lineFloat))
	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter required")
	}

	showChanges := false
	if sc, ok := args["show_changes"].(bool); ok {
		showChanges = sc
	}

	showResult := false
	if sr, ok := args["show_result"].(bool); ok {
		showResult = sr
	}

	err := insert.Insert(path, line, content, showChanges, showResult)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("Line inserted successfully"),
		},
	}, nil
}

func (m *McpServer) handleUndoEdit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}

	showChanges := false
	if sc, ok := args["show_changes"].(bool); ok {
		showChanges = sc
	}

	showResult := false
	if sr, ok := args["show_result"].(bool); ok {
		showResult = sr
	}

	err := undo_edit.UndoEdit(path, showChanges, showResult, 1)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("Edit undone successfully"),
		},
	}, nil
}

func (m *McpServer) handleGlob(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	pattern, ok := args["pattern"].(string)
	if !ok {
		return nil, fmt.Errorf("pattern parameter required")
	}

	path := ""
	if p, ok := args["path"].(string); ok {
		path = p
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := glob.Glob(pattern, path)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(output),
		},
	}, nil
}

func (m *McpServer) handleSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter required")
	}

	query, ok := args["tree_sitter_query"].(string)
	if !ok {
		return nil, fmt.Errorf("tree_sitter_query parameter required")
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := search.Search(path, query)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(output),
		},
	}, nil
}

func (m *McpServer) createLsTool() *mcp.Tool {
	tool := mcp.NewTool("ls",
		mcp.WithDescription("List directory contents"),
		mcp.WithString("path", mcp.Description("The directory to search in. If not specified, the current working directory will be used. IMPORTANT: Omit this field to use the default directory. DO NOT enter \"undefined\" or \"null\" - simply omit it for the default behavior. Must be a valid directory path if provided.")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	return &tool
}

func (m *McpServer) createSearchTool() *mcp.Tool {
	tool := mcp.NewTool("search",
		mcp.WithDescription("Search for code patterns using tree-sitter queries across files"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to file or directory to search")),
		mcp.WithString("tree_sitter_query", mcp.Required(), mcp.Description("Tree-sitter query pattern")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	return &tool
}

func (m *McpServer) createBatchTool() *mcp.Tool {
	tool := mcp.NewTool("batch",
		mcp.WithDescription("Execute multiple eddie operations in sequence from JSON input"),
		mcp.WithString("operations", mcp.Required(), mcp.Description("JSON string containing operations array: {\"operations\": [{\"type\": \"view\", \"path\": \"file.txt\"}, ...]}")),
	)
	return &tool
}

func (m *McpServer) handleLs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	path := "."
	if p, ok := args["path"].(string); ok {
		path = p
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := ls.Ls(path)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(output),
		},
	}, nil
}

func (m *McpServer) handleBatch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}

	operationsStr, ok := args["operations"].(string)
	if !ok {
		return nil, fmt.Errorf("operations parameter required")
	}

	batchReq, err := batch.ParseFromJSON(operationsStr)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error parsing batch operations: %v", err)),
			},
		}, nil
	}

	var buf bytes.Buffer
	processor := batch.NewProcessor(&buf)
	batchResp, err := processor.ProcessBatch(batchReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error processing batch: %v", err)),
			},
		}, nil
	}

	output, err := json.Marshal(batchResp)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error marshaling response: %v", err)),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(output)),
		},
	}, nil
}
