# arxiv-mcp-server

An MCP (Model Context Protocol) server written in Go to interact with the [arXiv.org](https://arxiv.org) API. It allows any MCP-compatible client (Claude Desktop, MCP-enabled IDEs, AI agents) to search scientific articles, retrieve metadata, and get PDF links directly from arXiv.

> *Thank you to arXiv for use of its open access interoperability.*

---

## Architecture

```
┌─────────────────────┐       stdio        ┌─────────────────────────┐
│    MCP Client        │◄──────────────────►│     MCP Server          │
│  (Claude, IDE, ...)  │                    │  cmd/server/main.go     │
└─────────────────────┘                    └──────────┬──────────────┘
                                                      │
                                           ┌──────────▼──────────────┐
                                           │    Tool Handlers         │
                                           │  internal/handler/       │
                                           │                          │
                                           │  - export-metadata       │
                                           │  - export-pdf-url        │
                                           └──────────┬──────────────┘
                                                      │
                                           ┌──────────▼──────────────┐
                                           │    HTTP Client           │
                                           │  internal/http-client/   │
                                           │                          │
                                           │  - Rate limiting (3s)    │
                                           │  - Retry + backoff       │
                                           │  - Query builder         │
                                           └──────────┬──────────────┘
                                                      │
                                           ┌──────────▼──────────────┐
                                           │   arXiv API              │
                                           │   export.arxiv.org/api   │
                                           └──────────────────────────┘
```

### Project structure

```
arxiv-mcp-server/
├── cmd/
│   ├── server/main.go             # MCP server entry point
│   └── client/main.go             # Test client
├── internal/
│   ├── handler/
│   │   └── export.go              # Tool handler implementations
│   └── http-client/
│       ├── http-client.go         # HTTP client with rate limiting and retry
│       ├── request.go             # Query parameter parsing (reflection-based)
│       └── response.go            # Response handling
├── scripts/
│   └── build.sh                   # Build script
├── go.mod
└── go.sum
```

### Core components

**MCP Server** (`cmd/server/main.go`) — Initializes the server using `go-sdk/mcp`, registers the tools, and starts the stdio transport.

**Tool Handlers** (`internal/handler/export.go`) — Implement tool logic as closure factories. Each handler receives the HTTP client and returns an MCP-compatible function.

**HTTP Client** (`internal/http-client/`) — Configured to respect arXiv API guidelines:
- **Rate limiting**: 3-second ticker between requests
- **Retry with exponential backoff**: up to 3 attempts, delay = `2^attempt * 3s`
- **Connection pooling**: max 1 idle connection, 90s timeout
- **Timeout**: 10 seconds per request

**Query Builder** (`internal/http-client/request.go`) — Reflection-based system that converts Go structs into query strings for the arXiv API, using `query` and `queryschema` struct tags.

---

## Available tools

### `export-metadata`

Searches arXiv and returns complete metadata as an Atom feed (title, authors, abstract, categories, dates, links).

### `export-pdf-url`

Searches arXiv and returns direct PDF URLs for matching articles.

### Input parameters

Both tools accept the same parameter schema:

| Parameter | Type | Description |
|-----------|------|-------------|
| `id_list` | `[]int` | List of specific article IDs |
| `start` | `int` | Start index for pagination |
| `max_results` | `int` | Maximum number of results returned |
| `search_query` | `object` | Search filters (see below) |

**Search filters** (`search_query`):

| Field | API tag | Description |
|-------|---------|-------------|
| `Title` | `ti` | Article title |
| `Author` | `au` | Author name |
| `Abstract` | `abs` | Abstract content |
| `Comment` | `co` | Article comments |
| `JournalReference` | `jr` | Journal reference |
| `SubjectCategory` | `cat` | Subject category (e.g. `cs.AI`, `math.CO`, `physics.optics`) |
| `ReportNumber` | `rn` | Report number |
| `All` | `all` | Search across all fields |

When multiple filters are specified, they are combined with the `AND` operator.

---

## Requirements

- Go >= 1.24.0

## Build

```bash
# Build the server
./scripts/build.sh server

# Build the test client
./scripts/build.sh client
```

This generates a `server.exe` (or `client.exe`) executable in the project root.

## Claude Desktop configuration

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "arxiv": {
      "command": "/absolute/path/to/server.exe"
    }
  }
}
```

Or, without compiling, run directly from source:

```json
{
  "mcpServers": {
    "arxiv": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "cwd": "/absolute/path/to/arxiv-mcp-server"
    }
  }
}
```

## Using with other MCP clients

Any MCP-compatible client that supports **stdio** transport can use this server. The server communicates via stdin/stdout following the MCP standard.

```bash
# Start the server (communicates via stdio)
./server.exe
```

---

## Usage examples

Below are practical examples of how an LLM or AI agent can leverage this MCP server.

### 1. Search by title

> *"Search for the article 'Attention Is All You Need' on arXiv and show me its metadata"*

The client invokes `export-metadata` with:
```json
{
  "max_results": 1,
  "search_query": { "Title": "Attention is all you need" }
}
```

### 2. Explore a scientific category

> *"Give me the 10 most recent articles in the cs.AI category"*

```json
{
  "max_results": 10,
  "search_query": { "SubjectCategory": "cs.AI" }
}
```

### 3. Search for a specific author's work

> *"Find all articles by Yann LeCun on arXiv"*

```json
{
  "max_results": 20,
  "search_query": { "Author": "Yann LeCun" }
}
```

### 4. Get the PDF of an article

> *"Get me the PDF link for the BERT paper"*

The client invokes `export-pdf-url` with:
```json
{
  "max_results": 1,
  "search_query": { "Title": "BERT Pre-training of Deep Bidirectional Transformers" }
}
```

### 5. Combined search (multi-filter)

> *"Search for machine learning articles in the bioinformatics field"*

```json
{
  "max_results": 5,
  "search_query": {
    "SubjectCategory": "q-bio",
    "All": "machine learning"
  }
}
```

### 6. Paginated results

> *"Show me results 10 through 20 for 'quantum computing'"*

```json
{
  "start": 10,
  "max_results": 10,
  "search_query": { "All": "quantum computing" }
}
```

### 7. AI-assisted literature review

> *"Do a literature review on transformers applied to computer vision, show me titles and abstracts of the first 15 articles"*

The AI agent invokes `export-metadata` with the appropriate parameters and then synthesizes the results for the user.

### 8. Journal monitoring

> *"Search for the most recent articles published in Physical Review Letters"*

```json
{
  "max_results": 10,
  "search_query": { "JournalReference": "Physical Review Letters" }
}
```

---

## TODO — Improvements and future features


- [ ] **Fix `logging in stderr`**: currently all logs are printed to stdio, which can interfere with MCP communication. Implement logging to stderr or a separate log file to avoid this issue.

- [ ] **Unit and integration tests**: add a test suite for handlers, the HTTP client, and parameter parsing. Use HTTP mocks for unit tests and integration tests against the real API.

- [ ] **`export-bibtex` tool**: add a tool that returns citations in BibTeX format, useful for LaTeX integration and bibliography management tools.

- [ ] **`get-article-by-id` tool**: a dedicated tool to retrieve a single article by its arXiv ID (e.g. `2106.09685`), without going through search.

- [ ] **Support OR and ANDNOT operators in queries**: currently search filters are combined with `AND` only. Add support for more complex boolean operators.

- [ ] **Response caching**: implement an in-memory (or on-disk) cache to avoid duplicate API requests and reduce latency.

- [ ] **Structured logging**: replace `log.Fatal` with a structured logger (e.g. `slog`) with configurable levels, useful for debugging and production monitoring.

- [ ] **External configuration**: allow configuring the base URL, timeout, rate limit, and retry via environment variables or a configuration file.

- [ ] **Dockerfile**: add a Dockerfile to simplify deployment and integration in containerized environments.

- [ ] **`search-similar` tool**: given an article ID, find related articles based on category and abstract keywords.

- [ ] **OpenAPI/JSON Schema documentation**: automatically generate tool parameter documentation from the `queryschema` struct tags.

- [ ] **Graceful shutdown handling**: intercept OS signals (SIGINT, SIGTERM) to properly close connections and release resources.

---

## License

See [LICENSE](LICENSE) file.

## Credits

This project uses the [arXiv API](https://info.arxiv.org/help/api/index.html) and the [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).