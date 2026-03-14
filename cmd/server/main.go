// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"context"
	"arxiv-mcp-server/internal/handler"
	"arxiv-mcp-server/internal/http-client"
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

	client := httpclient.New()

	mcp.AddTool(server, &mcp.Tool{
		Name:        "export",
		Description: "Export article from arxiv.org",
	}, handler.NewExport(client))

	transport := &mcp.StdioTransport{};

	if err := server.Run(ctx, transport); err != nil {
		log.Fatal(err);	
	}

	return;
}