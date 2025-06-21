package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/ls"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List directory contents",
	Long: `List directory contents.

Usage:
	ls [path]

Parameters:
	[path]: (Optional) The path to the directory to list. Defaults to current directory if not provided.

Example:
	eddie ls
	eddie ls /path/to/directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			path = "."
		}

		checkErr(ls.Ls(path))
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
