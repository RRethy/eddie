package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/create"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new file with the specified content.",
	Long: `Create a new file with the specified content.

Usage:
	create path file_text [--show-diff] [--show-result]

Parameters:
	path: The path where the new file should be created.
	file_text: The content to write to the new file.

Flags:
	--show-diff: Show the content of the created file.
	--show-result: Show the new content after the file creation.

Example:
	eddie create /path/to/newfile.txt "Hello, World!"
	eddie create config.json '{"key": "value"}' --show-diff
	eddie create script.sh "#!/bin/bash\necho 'Hello'" --show-result`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: path and file_text are required")
			return
		}
		path := args[0]
		fileText := args[1]
		showChanges, _ := cmd.Flags().GetBool("show-diff")
		showResult, _ := cmd.Flags().GetBool("show-result")

		checkErr(create.Create(path, fileText, showChanges, showResult))
	},
}

func init() {
	createCmd.Flags().Bool("show-diff", false, "Show the content of the created file")
	createCmd.Flags().Bool("show-result", false, "Show the new content after the file creation")
	rootCmd.AddCommand(createCmd)
}
