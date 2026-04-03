package cmd

import (
	"github.com/dmashuda/wormtongue/internal/mcpserver"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  "Start a Model Context Protocol server over stdio for LLM tool use.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpserver.Run(cmd.Context(), store)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
