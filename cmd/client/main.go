package main

import (
	httpclient "arxiv-mcp-server/internal/http-client"
	"context"
	"fmt"
	"log"
	"os/exec"

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
		log.Fatalf("error creating client: %s", err)
	}
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:	  "export-pdfurl",
		Arguments: httpclient.QueryParams{ 
			Search: httpclient.SearchQuery{ 
				Title: "attention",
			},
		},
	})
	if err != nil {
		log.Fatal(err);
	}

	for _, content := range res.Content {
		fmt.Println(content.(*mcp.TextContent).Text)
		//meta := content.(*mcp.TextContent).Meta.GetMeta()
	}
}