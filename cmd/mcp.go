package cmd

import (
	"fmt"
	"os"

	"github.com/ghchinoy/calctl/mcp"
)

// ExecuteMCP starts the MCP server.
func ExecuteMCP(httpAddr string) {
	if err := mcp.Start(rootCmd, httpAddr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
