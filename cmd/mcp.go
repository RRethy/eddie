package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/mcp"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start eddie as an MCP (Model Context Protocol) server",
	Long: `Start eddie as an MCP (Model Context Protocol) server that exposes all eddie commands as tools for LLM integration.

Usage:
	mcp

Example:
	eddie mcp`,
	Run: func(cmd *cobra.Command, args []string) {
		checkErr(mcp.Mcp())
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
