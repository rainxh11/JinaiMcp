# Reader MCP Server

An MCP (Model Context Protocol) server that provides tools for fetching web content in various formats (Markdown, HTML, text, and screenshots) using the [Reader](https://github.com/intergalacticalvariable/reader) service.

## Features

- **fetch_markdown**: Fetch a webpage and return its content as Markdown (bypasses readability processing)
- **fetch_html**: Fetch a webpage and return its HTML (documentElement.outerHTML)
- **fetch_text**: Fetch a webpage and return its text content (document.body.innerText)
- **fetch_screenshot**: Fetch a screen-size screenshot of a webpage
- **fetch_pageshot**: Fetch a full-page screenshot of a webpage

## Quick Start with Docker

The easiest way to run both the Reader service and MCP server is using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start the Reader service on port 3000
- Start the MCP server on port 8080 with SSE transport

## Setup

### 1. Start the Services

```bash
docker-compose up -d
```

### 2. Configure Claude Desktop

To use this MCP server with Claude Desktop, add the following to your Claude Desktop configuration file:

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "reader": {
      "transport": {
        "type": "sse",
        "url": "http://localhost:8080/sse"
      }
    }
  }
}
```

### 3. Restart Claude Desktop

After updating the configuration, restart Claude Desktop to load the MCP server.

## Usage

Once configured, you can use the tools in Claude Desktop:

- "Fetch https://example.com as markdown"
- "Get the HTML of https://example.com"
- "Take a screenshot of https://example.com"

## API Endpoints

When running via Docker Compose, the MCP server exposes:

- **SSE Endpoint**: `http://localhost:8080/sse` - Server-Sent Events endpoint for MCP communication

## Development

### Dependencies

- Go 1.23+
- Docker (for the Reader service)

### Building

```bash
go build -o reader-mcp .
```

### Running Locally

1. Start the Reader service:
```bash
docker-compose up reader -d
```

2. Run the MCP server:
```bash
READER_ENDPOINT=http://localhost:3000 PORT=8080 go run main.go
```

### Environment Variables

- `READER_ENDPOINT`: URL of the Reader service (default: `http://reader-container:3000`)
- `PORT`: Port for the MCP server (default: `8080`)

## License

MIT
