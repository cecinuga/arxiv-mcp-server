package main

import (
	"fmt"
	"os/exec"
	"context"
	httpclient "arxiv-mcp-server/internal/http-client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main(){
	ctx := context.Background();
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:		"arxiv.org client",
			Version:  	"v1.0.0",
		},
		nil,
	)

	transport := &mcp.CommandTransport{ 
		Command: exec.Command("go", "run", "./cmd/server"),
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:	  "export",
		Arguments: httpclient.QueryParams{ 
			MaxResults: 1, 
			Search: httpclient.SearchQuery{ Title:"Attention is all you need" },
		},
	})

	fmt.Println(res.Content[0].(*mcp.TextContent).Text)
}