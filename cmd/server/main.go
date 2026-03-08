// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"context"

	"arxiv-mcp-server/internal/handler"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main(){
	ctx := context.Background();
	server := mcp.NewServer(
		&mcp.Implementation{
			Version:  	"v1.0.0",
			Name:		"arxiv.org server",
			Title: 		"Thank you to arXiv for use of its open access interoperability.",
		},
		nil,
	);

	mcp.AddTool(server, &mcp.Tool{
		Name:		"greet",
		Description:"Saluta una persona per nome",
	}, handler.HandleGreet);

	transport := &mcp.StdioTransport{};

	if err := server.Run(ctx, transport); err != nil {
		log.Fatal(err);	
	}
}