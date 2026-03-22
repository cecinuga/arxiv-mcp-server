package handler

import (
	httpclient "arxiv-mcp-server/internal/http-client"
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/tools/blog/atom"
	"arxiv-mcp-server/internal/api"
)

func fetchMetadata(ctx context.Context, client *httpclient.Client, input httpclient.QueryParams) (*atom.Feed, error) {
	feed, err := client.Get(ctx, input);
	if err != nil {
		return nil, fmt.Errorf("error failed to fetch: \n%s", err)
	}

	return feed, nil
}

func fetchPdfUrl(ctx context.Context, client *httpclient.Client, input httpclient.QueryParams) ([]api.PdfResource, error){
	feed, err := fetchMetadata(ctx, client, input);
	if err != nil {
		return nil, err;
	}

	var resources []api.PdfResource; 
	
	for _, entry := range feed.Entry {
		url := strings.Replace(entry.ID, "abs", "pdf", 1)
		resources = append(resources, api.PdfResource{
			Url: url,
			Meta: api.Metadata{ Author: entry.Author.Name, Title: entry.Title },
		})
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
        resources, err := fetchPdfUrl(ctx, client, input);
		if err != nil {
			return nil, nil, err;
		}
		
		var contents []mcp.Content
		for _, resource := range resources {
			metamap := map[string]any{}
			metamap["title"] = resource.Meta.Title
			metamap["author"] = resource.Meta.Author

			var meta mcp.Meta
			meta.SetMeta(metamap)

			contents = append(contents, &mcp.TextContent{ 
					Text: resource.Url,
					Meta: meta,
				},
			)
		}

        return &mcp.CallToolResult{Content: contents}, nil, nil
    }
}