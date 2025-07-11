package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/insert"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a new line at the specified line number in a file.",
	Long: `Insert a new line at the specified line number in a file.

Usage:
	insert path insert_line new_str [--show-diff] [--show-result]

Parameters:
	path: The path to the file to modify.
	insert_line: The line number where the new line should be inserted (1-based).
	new_str: The content of the new line to insert.

Flags:
	--show-diff: Show the changes made to the file.
	--show-result: Show the new content after the edit operation.

Example:
	eddie insert /path/to/file.txt 5 "This is a new line"
	eddie insert config.json 10 "  \"newKey\": \"newValue\"," --show-diff
	eddie insert script.sh 1 "#!/bin/bash" --show-result`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Error: path, insert_line, and new_str are required")
			return
		}
		path := args[0]
		insertLine := args[1]
		newStr := args[2]
		showChanges, _ := cmd.Flags().GetBool("show-diff")
		showResult, _ := cmd.Flags().GetBool("show-result")

		checkErr(insert.Insert(path, insertLine, newStr, showChanges, showResult))
	},
}

func init() {
	insertCmd.Flags().Bool("show-diff", false, "Show the changes made to the file")
	insertCmd.Flags().Bool("show-result", false, "Show the new content after the edit operation")
	rootCmd.AddCommand(insertCmd)
}
