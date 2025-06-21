package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/glob"
)

var globCmd = &cobra.Command{
	Use:   "glob",
	Short: "Find files matching a glob pattern",
	Long: `Find files matching a glob pattern.

Usage:
	glob pattern [path]

Parameters:
	pattern: The glob pattern to match files against
	[path]: (Optional) The directory to search in. Defaults to current directory.

Example:
	eddie glob "*.go"
	eddie glob "**/*.js" src/
	eddie glob "test_*.py" tests/`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: pattern is required")
			return
		}
		pattern := args[0]
		var path string
		if len(args) > 1 {
			path = args[1]
		}

		checkErr(glob.Glob(pattern, path))
	},
}

func init() {
	rootCmd.AddCommand(globCmd)
}
