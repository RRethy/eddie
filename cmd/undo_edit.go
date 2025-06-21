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
	undo_edit path [--show-diff]

Parameters:
	path: The path to the file to restore from backup.

Flags:
	--show-diff: Show the changes made during the undo operation.

Example:
	eddie undo_edit /path/to/file.txt
	eddie undo_edit config.json --show-diff
	eddie undo_edit script.sh`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: path is required")
			return
		}
		path := args[0]
		showChanges, _ := cmd.Flags().GetBool("show-diff")

		checkErr(undo_edit.UndoEdit(path, showChanges))
	},
}

func init() {
	undoEditCmd.Flags().Bool("show-diff", false, "Show the changes made during the undo operation")
	rootCmd.AddCommand(undoEditCmd)
}
