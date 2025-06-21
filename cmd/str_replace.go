package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/str_replace"
)

var strReplaceCmd = &cobra.Command{
	Use:   "str_replace",
	Short: "Replace all occurrences of a string in a file with another string.",
	Long: `Replace all occurrences of a string in a file with another string.

Usage:
	str_replace path old_str new_str [--show-diff]

Parameters:
	path: The path to the file to modify.
	old_str: The string to search for and replace.
	new_str: The string to replace old_str with.

Flags:
	--show-diff: Show the changes made to the file.

Example:
	eddie str_replace /path/to/file.txt "old text" "new text"
	eddie str_replace config.json "localhost" "example.com" --show-diff`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Error: path, old_str, and new_str are required")
			return
		}
		path := args[0]
		oldStr := args[1]
		newStr := args[2]
		showChanges, _ := cmd.Flags().GetBool("show-diff")

		checkErr(str_replace.StrReplace(path, oldStr, newStr, showChanges))
	},
}

func init() {
	strReplaceCmd.Flags().Bool("show-diff", false, "Show the changes made to the file")
	rootCmd.AddCommand(strReplaceCmd)
}
