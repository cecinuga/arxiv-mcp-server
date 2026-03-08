package handler

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"arxiv-mcp-server/internal/api"
)

func HandleGreet(_ context.Context, _ *mcp.CallToolRequest, input api.GreetParams) (
	*mcp.CallToolResult, any, error,
) {
	msg := fmt.Sprintf("Ciao, %s! Benvenuto nel server MCP.", input.Name)
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{ Text: msg },
		},
	}, nil, nil
}