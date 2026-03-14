package handler

import (
	"log"
	"context"
	atomparser "github.com/wbernest/atom-parser"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	httpclient "arxiv-mcp-server/internal/http-client"
)

func NewExport(client *httpclient.Client) func(context.Context, *mcp.CallToolRequest, httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, _ *mcp.CallToolRequest, input httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
        res, err := client.Get(ctx, input);
		if err != nil {
			log.Fatal(err);
		}

		data, err := httpclient.ReadBody(res);
		if err != nil {
			log.Fatal(err)
		}

		feed, err := atomparser.ParseString(string(data));
		if err != nil {
			log.Fatal(err)
		}
		
        return &mcp.CallToolResult{}, feed, nil
    }
}