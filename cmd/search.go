package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/search"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for code patterns using tree-sitter queries across files.",
	Long: `Search for code patterns using tree-sitter queries across files.

Usage:
	search <file|dir> --tree-sitter-query "<tree-sitter-query>"

Parameters:
	<file|dir>: Path to file or directory to search.

Flags:
	--tree-sitter-query: Tree-sitter query pattern (required).

Example:
	eddie search ./src --tree-sitter-query "(function_declaration name: (identifier) @func)"
	eddie search main.go --tree-sitter-query "(call_expression function: (identifier) @call)"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: file or directory path is required")
			return
		}
		path := args[0]
		query, _ := cmd.Flags().GetString("tree-sitter-query")
		if query == "" {
			fmt.Println("Error: --tree-sitter-query flag is required")
			return
		}

		checkErr(search.Search(path, query))
	},
}

func init() {
	searchCmd.Flags().StringP("tree-sitter-query", "q", "", "Tree-sitter query pattern (required)")
	rootCmd.AddCommand(searchCmd)
}
