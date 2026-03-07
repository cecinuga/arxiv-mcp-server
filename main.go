// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main(){
	ctx := context.Background();
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:		"arxiv.org server",
			Version:  	"v1.0.0",
		},
		nil,
	);

	mcp.AddTool(server, &mcp.Tool{
		Name:		"greet",
		Description:"Saluta una persona per nome",
	}, handleGreet);

	clientTransport, serverTransport := mcp.NewInMemoryTransports();
	serverSession, err := server.Connect(ctx, serverTransport, nil);
	if err != nil{
		log.Fatal(err);
	}

	client := mcp.NewClient(&mcp.Implementation{ Name: "test-client" }, nil);
	clientSession, err := client.Connect(ctx, clientTransport, nil);

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:	  "greet",
		Arguments: map[string]any{ "name": "Cecinuga" },
	})

	if err != nil{
		log.Fatal(err);
	}

	fmt.Println(res.Content[0].(*mcp.TextContent).Text)

	clientSession.Close()
	serverSession.Wait()
}