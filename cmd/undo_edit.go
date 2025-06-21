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
	undo_edit path

Parameters:
	path: The path to the file to restore from backup.

Example:
	eddie undo_edit /path/to/file.txt
	eddie undo_edit config.json
	eddie undo_edit script.sh`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: path is required")
			return
		}
		path := args[0]

		checkErr(undo_edit.UndoEdit(path))
	},
}

func init() {
	rootCmd.AddCommand(undoEditCmd)
}
