package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type Searcher struct{}

type Match struct {
	File    string
	Line    int
	Column  int
	Content string
	Capture string
}

func (s *Searcher) Search(path, queryStr string) error {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	lang := tree_sitter.NewLanguage(tree_sitter_go.Language())
	parser.SetLanguage(lang)

	query, queryErr := tree_sitter.NewQuery(lang, queryStr)
	if queryErr != nil {
		return fmt.Errorf("invalid query: %s", queryErr.Message)
	}
	defer query.Close()

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return s.searchDir(path, parser, query)
	}
	return s.searchFile(path, parser, query)
}

func (s *Searcher) searchDir(dir string, parser *tree_sitter.Parser, query *tree_sitter.Query) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		return s.searchFile(path, parser, query)
	})
}

func (s *Searcher) searchFile(filename string, parser *tree_sitter.Parser, query *tree_sitter.Query) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %s: %w", filename, err)
	}

	tree := parser.Parse(content, nil)
	if tree == nil {
		return fmt.Errorf("failed to parse %s", filename)
	}
	defer tree.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(query, tree.RootNode(), content)
	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			node := capture.Node
			startPos := node.StartPosition()

			lines := strings.Split(string(content), "\n")
			var lineContent string
			if int(startPos.Row) < len(lines) {
				lineContent = strings.TrimSpace(lines[startPos.Row])
			}

			captureName := query.CaptureNames()[capture.Index]

			fmt.Printf("%s:%d:%d: @%s: %s\n",
				filename,
				startPos.Row+1,
				startPos.Column+1,
				captureName,
				lineContent)
		}
	}

	return nil
}
