package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

type Searcher struct{}

type Match struct {
	File    string
	Content string
	Capture string
	Line    int
	Column  int
}

var languageMap = map[string]func() *tree_sitter.Language{
	".go":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_go.Language()) },
	".js":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_javascript.Language()) },
	".mjs": func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_javascript.Language()) },
	".jsx": func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_javascript.Language()) },
	".ts": func() *tree_sitter.Language {
		return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	},
	".tsx":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTSX()) },
	".py":   func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_python.Language()) },
	".pyi":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_python.Language()) },
	".rs":   func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_rust.Language()) },
	".java": func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_java.Language()) },
	".c":    func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_c.Language()) },
	".h":    func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_c.Language()) },
	".cc":   func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_cpp.Language()) },
	".cpp":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_cpp.Language()) },
	".cxx":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_cpp.Language()) },
	".hpp":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_cpp.Language()) },
	".hxx":  func() *tree_sitter.Language { return tree_sitter.NewLanguage(tree_sitter_cpp.Language()) },
}

func getLanguageFromFile(filename string) *tree_sitter.Language {
	ext := strings.ToLower(filepath.Ext(filename))
	if langFunc, exists := languageMap[ext]; exists {
		return langFunc()
	}
	return nil
}

func (s *Searcher) Search(path, queryStr string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return s.searchDir(path, queryStr)
	}
	return s.searchFile(path, queryStr)
}

func (s *Searcher) searchDir(dir, queryStr string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if getLanguageFromFile(path) == nil {
			return nil
		}

		return s.searchFile(path, queryStr)
	})
}

func (s *Searcher) searchFile(filename, queryStr string) error {
	lang := getLanguageFromFile(filename)
	if lang == nil {
		return nil
	}

	parser := tree_sitter.NewParser()
	defer parser.Close()

	err := parser.SetLanguage(lang)
	if err != nil {
		return fmt.Errorf("set language for %s: %w", filename, err)
	}

	query, queryErr := tree_sitter.NewQuery(lang, queryStr)
	if queryErr != nil {
		return fmt.Errorf("invalid query for %s: %s", filename, queryErr.Message)
	}
	defer query.Close()

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
			if startPos.Row < uint(len(lines)) {
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
