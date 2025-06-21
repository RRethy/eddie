package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eddie",
	Short: "A text editor designed for AI Agents (e.g. `claude` code), not humans.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error:"), err)
		os.Exit(1)
	}
}
