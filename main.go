package main

import (
	"os"
	"strings"

	"github.com/ghchinoy/calctl/cmd"
)

// main is the entry point of the application.
// It checks for the --mcp and --mcp-http flags to determine whether to run in MCP mode.
func main() {
	mcp := false
	mcpHTTP := ""

	for i, arg := range os.Args[1:] {
		if arg == "--mcp" {
			mcp = true
			break
		}
		if strings.HasPrefix(arg, "--mcp-http") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				mcpHTTP = parts[1]
			} else if i+1 < len(os.Args[1:]) {
				mcpHTTP = os.Args[i+2]
			}
			break
		}
	}

	if mcp || mcpHTTP != "" {
		cmd.ExecuteMCP(mcpHTTP)
	} else {
		cmd.Execute()
	}
}
