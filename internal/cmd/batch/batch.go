package batch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/RRethy/eddie/internal/cmd/create"
	"github.com/RRethy/eddie/internal/cmd/insert"
	"github.com/RRethy/eddie/internal/cmd/ls"
	"github.com/RRethy/eddie/internal/cmd/search"
	"github.com/RRethy/eddie/internal/cmd/str_replace"
	"github.com/RRethy/eddie/internal/cmd/undo_edit"
	"github.com/RRethy/eddie/internal/cmd/view"
)

type Processor struct {
	out io.Writer
}

func NewProcessor(out io.Writer) *Processor {
	return &Processor{out: out}
}

func (p *Processor) ProcessBatch(req *BatchRequest) (*BatchResponse, error) {
	resp := &BatchResponse{
		Results: make([]OperationResult, len(req.Operations)),
	}

	for i, op := range req.Operations {
		result := p.processOperation(op)
		resp.Results[i] = result
	}

	return resp, nil
}

func (p *Processor) processOperation(op Operation) OperationResult {
	var buf bytes.Buffer
	var err error

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputDone := make(chan string)
	go func() {
		var output bytes.Buffer
		output.ReadFrom(r)
		outputDone <- output.String()
	}()

	switch op.Type {
	case "view":
		err = view.View(op.Path, op.ViewRange)
	case "str_replace":
		err = str_replace.NewReplacer(&buf).StrReplace(op.Path, op.OldStr, op.NewStr, op.ShowChanges, op.ShowResult)
	case "create":
		err = create.NewCreator(&buf).Create(op.Path, op.Content, op.ShowChanges, op.ShowResult)
	case "insert":
		insertLine := strconv.Itoa(op.InsertLine)
		err = insert.NewInserter(&buf).Insert(op.Path, insertLine, op.NewStr, op.ShowChanges, op.ShowResult)
	case "undo_edit":
		err = undo_edit.NewUndoEditor(&buf).UndoEdit(op.Path, op.ShowChanges, op.ShowResult, op.Count)
	case "ls":
		err = ls.Ls(op.Path)
	case "search":
		err = search.Search(op.Path, op.TreeQuery)
	default:
		err = fmt.Errorf("unknown operation type: %s", op.Type)
	}

	w.Close()
	os.Stdout = oldStdout
	capturedOutput := <-outputDone

	var output string
	if buf.Len() > 0 {
		output = buf.String()
	} else {
		output = capturedOutput
	}

	result := OperationResult{
		Operation: op,
		Success:   err == nil,
		Output:    output,
	}

	if err != nil {
		errStr := err.Error()
		result.Error = &errStr
	}

	return result
}

func ParseFromStdin() (*BatchRequest, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}

	var req BatchRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	return &req, nil
}

func ParseFromFile(path string) (*BatchRequest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var req BatchRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("parse JSON from %s: %w", path, err)
	}

	return &req, nil
}

func ParseFromJSON(jsonStr string) (*BatchRequest, error) {
	var req BatchRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return nil, fmt.Errorf("parse JSON string: %w", err)
	}

	return &req, nil
}

func ParseFromOps(ops []string) (*BatchRequest, error) {
	req := &BatchRequest{
		Operations: make([]Operation, len(ops)),
	}

	for i, op := range ops {
		parts := strings.Split(op, ",")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid operation format: %s", op)
		}

		operation := Operation{
			Type: parts[0],
			Path: parts[1],
		}

		switch parts[0] {
		case "view":
			if len(parts) > 2 {
				operation.ViewRange = strings.Join(parts[2:], ",")
			}
		case "str_replace":
			if len(parts) < 4 {
				return nil, fmt.Errorf("str_replace requires old_str and new_str: %s", op)
			}
			operation.OldStr = parts[2]
			operation.NewStr = parts[3]
		case "create":
			if len(parts) < 3 {
				return nil, fmt.Errorf("create requires content: %s", op)
			}
			operation.Content = parts[2]
		case "insert":
			if len(parts) < 4 {
				return nil, fmt.Errorf("insert requires line number and content: %s", op)
			}
			var line int
			if _, err := fmt.Sscanf(parts[2], "%d", &line); err != nil {
				return nil, fmt.Errorf("invalid line number in insert: %s", op)
			}
			operation.InsertLine = line
			operation.NewStr = parts[3]
		case "undo_edit":
		case "ls":
		case "search":
			if len(parts) < 3 {
				return nil, fmt.Errorf("search requires tree-sitter query: %s", op)
			}
			operation.TreeQuery = parts[2]
		default:
			return nil, fmt.Errorf("unknown operation type: %s", parts[0])
		}

		req.Operations[i] = operation
	}

	return req, nil
}
