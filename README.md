# Reader MCP Server

A FastMCP server that provides tools to convert URLs to LLM-friendly formats using the Reader service. Built with **Deno + Hono** for the Reader service and **Go** for the MCP server with Streamable HTTP transport!

## Features

This MCP server provides 5 tools:

1. **get_markdown** - Convert a URL to markdown format (bypasses readability processing)
2. **get_html** - Convert a URL to HTML format (returns documentElement.outerHTML)
3. **get_text** - Convert a URL to plain text format (returns document.body.innerText)
4. **get_screenshot** - Take a screen-size screenshot of a URL
5. **get_pageshot** - Take a full-page screenshot of a URL

## Tech Stack

- **Reader Service**: Deno + Hono + Puppeteer
- **MCP Server**: Go with Streamable HTTP transport
- **Base Image**: browserless/chrome (Chrome/Puppeteer pre-installed)
- **Transport**: Streamable HTTP (bidirectional streaming over HTTP)

## Quick Start with Docker

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

### Using Docker directly

```bash
docker build -t reader-mcp .
docker run -d -p 3000:3000 -p 8000:8000 -v $(pwd)/screenshots:/app/local-storage --cap-add=SYS_ADMIN --shm-size=2g --name reader-mcp reader-mcp
```

## MCP Endpoints

Once running, the MCP server will be available at:

```
MCP Streamable HTTP: http://localhost:8000/
```

The Reader service runs on:
```
Reader Service: http://localhost:3000
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `READER_URL` | `http://localhost:3000` | URL of the Reader service |
| `MCP_PORT` | `8000` | Port for the MCP server |
| `PUPPETEER_EXECUTABLE_PATH` | `/usr/bin/google-chrome-stable` | Path to Chrome executable |

## Usage with MCP Clients

### Claude Desktop Configuration

Add to your Claude Desktop config (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "reader": {
      "transport": {
        "type": "http",
        "url": "http://localhost:8000",
        "streaming": true
      }
    }
  }
}
```

### Direct HTTP Access

You can also test the Reader service directly:

```bash
# Get markdown from a URL
curl -H "X-Respond-With: markdown" http://localhost:3000/https://example.com

# Get HTML from a URL
curl -H "X-Respond-With: html" http://localhost:3000/https://example.com

# Get text from a URL
curl -H "X-Respond-With: text" http://localhost:3000/https://example.com

# Get screenshot URL
curl -H "X-Respond-With: screenshot" http://localhost:3000/https://example.com

# Get full-page screenshot URL
curl -H "X-Respond-With: pageshot" http://localhost:3000/https://example.com
```

## API Tools

### get_markdown
Converts a URL to markdown format.
- **Parameters:** `url` (string, required)
- **Returns:** Markdown content

### get_html
Converts a URL to HTML format.
- **Parameters:** `url` (string, required)
- **Returns:** HTML content

### get_text
Converts a URL to plain text format.
- **Parameters:** `url` (string, required)
- **Returns:** Plain text content

### get_screenshot
Takes a screen-size screenshot of a URL.
- **Parameters:** `url` (string, required)
- **Returns:** URL of the screenshot image

### get_pageshot
Takes a full-page screenshot of a URL.
- **Parameters:** `url` (string, required)
- **Returns:** URL of the full-page screenshot image

## Architecture

```
┌─────────────────┐
│   MCP Client    │
│  (Claude, etc.) │
└────────┬────────┘
         │ Streamable HTTP
         ▼
┌─────────────────┐      ┌─────────────────┐
│   MCP Server    │──────│  Reader Service │
│   (Go)          │      │   (Deno+Hono)   │
│   Port: 8000    │      │    Port: 3000    │
└─────────────────┘      └────────┬────────┘
                                  │
                                  ▼
                          ┌─────────────────┐
                          │ browserless/    │
                          │    chrome       │
                          └─────────────────┘
```

## Project Structure

```
ReaderMcp/
├── Dockerfile              # Multi-stage build with Go + Deno
├── docker-compose.yaml     # Docker Compose configuration
├── go.mod                  # Go module definition
├── main.go                 # Go MCP server entry point
├── start.sh                # Startup script
└── deno/
    ├── deno.json          # Deno configuration
    ├── main.ts            # Hono server entry point
    └── services/
        ├── puppeteer.ts   # Puppeteer service
        └── storage.ts     # File storage utilities
```

## Development

### Running locally with Deno

```bash
# Install Deno
curl -fsSL https://deno.land/install.sh | sh

# Run Reader service
cd deno
deno task dev
```

### Running MCP server locally (Go)

```bash
# Run Go MCP server
go run main.go
```

## Based On

This project is based on [Jina AI's Reader](https://github.com/jina-ai/reader) and the [Docker deployment version](https://github.com/intergalacticalvariable/reader) by intergalacticalvariable.

## License

Apache-2.0
