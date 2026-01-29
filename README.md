# Reader MCP Server

A FastMCP server that provides tools to convert URLs to LLM-friendly formats using the Reader service.

## Features

This MCP server provides 5 tools:

1. **get_markdown** - Convert a URL to markdown format (bypasses readability processing)
2. **get_html** - Convert a URL to HTML format (returns documentElement.outerHTML)
3. **get_text** - Convert a URL to plain text format (returns document.body.innerText)
4. **get_screenshot** - Take a screen-size screenshot of a URL
5. **get_pageshot** - Take a full-page screenshot of a URL

## Quick Start with Docker

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

### Using Docker directly

```bash
docker build -t reader-mcp .
docker run -d -p 3000:3000 -v $(pwd)/screenshots:/app/local-storage --name reader-mcp reader-mcp
```

## Configuration

The following environment variables can be configured:

- `READER_URL` - URL of the Reader service (default: `http://localhost:3000`)
- `PUPPETEER_SKIP_CHROMIUM_DOWNLOAD` - Skip Puppeteer chromium download (default: `true`)
- `PUPPETEER_EXECUTABLE_PATH` - Path to Chrome executable (default: `/usr/bin/google-chrome-stable`)

## Usage with MCP Clients

Add to your Claude Desktop config (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "reader": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "reader-mcp"]
    }
  }
}
```

## API Tools

### get_markdown

Converts a URL to markdown format.

**Parameters:**
- `url` (string, required): The URL to fetch and convert

**Example:**
```
Call get_markdown with url="https://example.com"
```

### get_html

Converts a URL to HTML format.

**Parameters:**
- `url` (string, required): The URL to fetch and convert

### get_text

Converts a URL to plain text format.

**Parameters:**
- `url` (string, required): The URL to fetch and convert

### get_screenshot

Takes a screen-size screenshot of a URL.

**Parameters:**
- `url` (string, required): The URL to screenshot

**Returns:** URL of the screenshot image

### get_pageshot

Takes a full-page screenshot of a URL.

**Parameters:**
- `url` (string, required): The URL to screenshot

**Returns:** URL of the full-page screenshot image

## Based On

This project is based on [Jina AI's Reader](https://github.com/jina-ai/reader) and the [Docker deployment version](https://github.com/intergalacticalvariable/reader) by intergalacticalvariable.

## License

Apache-2.0
