package main

import (
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
		log.Fatal(err)
	}
	defer session.Close()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:	  "greet",
		Arguments: map[string]any{ "name": "Cecinuga" },
	})

	fmt.Println(res.Content[0].(*mcp.TextContent).Text)
}