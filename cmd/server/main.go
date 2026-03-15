// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"arxiv-mcp-server/internal/handler"
	"arxiv-mcp-server/internal/http-client"
	"context"
	"fmt"

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
		Name:        "export-metadata",
		Description: "Export article's feed from arxiv.org",
	}, handler.NewExportMetadata(client))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "export-pdf-url",
		Description: "Export the article's pdf URL from arxiv.org",
	}, handler.NewExportPdfUrl(client))

	transport := &mcp.StdioTransport{};

	if err := server.Run(ctx, transport); err != nil {
		fmt.Println(err)
	}
}