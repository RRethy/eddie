package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/view"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Examine the contents of a file or list the contents of a directory. It can read the entire file or a specific range of lines.",
	Long: `Examine the contents of a file or list the contents of a directory. It can read the entire file or a specific range of lines.

Usage:
	view path [view_range]

Parameters:
	path: The path to the file or directory to view.
	[view_range]: (Optional) An optional parameter specifying the range of lines to view in a file, formatted as "start,end". If "end" is -1, it means read to the end of the file. This parameter is ignored when viewing directories.

Example:
	eddie view /path/to/file.txt
	eddie view /path/to/directory
	eddie view /path/to/file.txt 10,20`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: path is required")
			return
		}
		path := args[0]
		var viewRange string
		if len(args) > 1 {
			viewRange = args[1]
		}

		checkErr(view.View(path, viewRange))
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
