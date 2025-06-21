package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/undo_edit"
)

var undoEditCmd = &cobra.Command{
	Use:   "undo_edit",
	Short: "Undo the last edit operation on a file by restoring from backup.",
	Long: `Undo the last edit operation on a file by restoring from backup.

This command restores a file to its previous state before the last edit operation 
(str_replace, insert, etc.). It looks for the most recent backup file and restores 
the original content.

Usage:
	undo_edit path [--show-diff] [--show-result] [--count N]

Parameters:
	path: The path to the file to restore from backup.

Flags:
	--show-diff: Show the changes made during the undo operation.
	--show-result: Show the new content after the undo operation.
	--count: Number of edits to undo (default: 1).

Example:
	eddie undo_edit /path/to/file.txt
	eddie undo_edit config.json --show-diff
	eddie undo_edit script.sh --show-result
	eddie undo_edit script.sh --count 3`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: path is required")
			return
		}
		path := args[0]
		showChanges, _ := cmd.Flags().GetBool("show-diff")
		showResult, _ := cmd.Flags().GetBool("show-result")
		count, _ := cmd.Flags().GetInt("count")

		checkErr(undo_edit.UndoEdit(path, showChanges, showResult, count))
	},
}

func init() {
	undoEditCmd.Flags().Bool("show-diff", false, "Show the changes made during the undo operation")
	undoEditCmd.Flags().Bool("show-result", false, "Show the new content after the undo operation")
	undoEditCmd.Flags().Int("count", 1, "Number of edits to undo")
	rootCmd.AddCommand(undoEditCmd)
}
