package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func handleGreet(ctx context.Context, req *mcp.CallToolRequest, input GreetInput) (
	*mcp.CallToolResult, any, error,
) {
	msg := fmt.Sprintf("Ciao, %s! Benvenuto nel server MCP.", input.Name)
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{ Text: msg },
		},
	}, nil, nil
}