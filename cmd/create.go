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
	create path file_text [--show-changes]

Parameters:
	path: The path where the new file should be created.
	file_text: The content to write to the new file.

Flags:
	--show-changes: Show the content of the created file.

Example:
	eddie create /path/to/newfile.txt "Hello, World!"
	eddie create config.json '{"key": "value"}' --show-changes
	eddie create script.sh "#!/bin/bash\necho 'Hello'"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: path and file_text are required")
			return
		}
		path := args[0]
		fileText := args[1]
		showChanges, _ := cmd.Flags().GetBool("show-changes")

		checkErr(create.Create(path, fileText, showChanges))
	},
}

func init() {
	createCmd.Flags().Bool("show-changes", false, "Show the content of the created file")
	rootCmd.AddCommand(createCmd)
}
