package handler

import (
	httpclient "arxiv-mcp-server/internal/http-client"
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/tools/blog/atom"
)

func fetchMetadata(ctx context.Context, client *httpclient.Client, input httpclient.QueryParams) (*atom.Feed, error) {
	feed, err := client.Get(ctx, input);
	if err != nil {
		return nil, fmt.Errorf("error failed to fetch: \n%s", err)
	}

	return feed, nil
}

func fetchPdfUrl(ctx context.Context, client *httpclient.Client, input httpclient.QueryParams) ([]string, error){
	feed, err := fetchMetadata(ctx, client, input);
	if err != nil {
		return nil, err;
	}

	var resources []string; 
	
	for _, entry := range feed.Entry {
		resource := strings.Replace(entry.ID, "abs", "pdf", 1)
		resources = append(resources, resource)
	}

	return resources, nil
}

func NewExportMetadata(client *httpclient.Client) func(context.Context, *mcp.CallToolRequest, httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, _ *mcp.CallToolRequest, input httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
        res, err := fetchMetadata(ctx, client, input);
		if err != nil {
			return nil, nil, err;
		}
		
        return &mcp.CallToolResult{}, res, nil
    }
}


func NewExportPdfUrl(client *httpclient.Client) func(context.Context, *mcp.CallToolRequest, httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, _ *mcp.CallToolRequest, input httpclient.QueryParams) (*mcp.CallToolResult, any, error) {
        urls, err := fetchPdfUrl(ctx, client, input);
		if err != nil {
			return nil, nil, err;
		}
		
		var contents []mcp.Content
		for _, url := range urls {
			var meta mcp.Meta
			//meta.SetMeta("title")

			contents = append(contents, &mcp.TextContent{ 
					Text: url,
					Meta: meta,
				},
			)
		}

        return &mcp.CallToolResult{Content: contents}, nil, nil
    }
}