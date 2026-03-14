package handler

import (
	httpclient "arxiv-mcp-server/internal/http-client"
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	atomparser "github.com/wbernest/atom-parser"
	"golang.org/x/tools/blog/atom"
)

func fetchMetadata(ctx context.Context, client *httpclient.Client, input httpclient.QueryParams) (*atom.Feed, error) {
	res, err := client.Get(ctx, input);
	if err != nil {
		return nil, fmt.Errorf("error failed to fetch: %s", err)
	}

	data, err := httpclient.ReadBody(res);
	if err != nil {
		return nil, fmt.Errorf("error reading body response: %s", err)
	}

	feed, err := atomparser.ParseString(string(data));
	if err != nil {
		return nil, fmt.Errorf("error parsing atom xml: %s", err)
	}
	
	return feed, nil
}

func NewExportMetadata(client *httpclient.Client) func(context.Context, *mcp.CallToolRequest, httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, _ *mcp.CallToolRequest, input httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
        res, err := fetchMetadata(ctx, client, input);
		if err != nil {
			return nil, nil, fmt.Errorf("error fetching metadata: %s", err);
		}
		
        return &mcp.CallToolResult{}, res, nil
    }
}

func NewExportRaw(client *httpclient.Client) func(context.Context, *mcp.CallToolRequest, httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, _ *mcp.CallToolRequest, input httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
        

        return &mcp.CallToolResult{}, nil, nil
    }
}