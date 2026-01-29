#!/usr/bin/env python3
"""
FastMCP Server for Reader - URL to LLM-friendly conversion.

This server provides 5 tools based on the Reader service:
1. get_markdown - Returns content as markdown (bypasses readability processing)
2. get_html - Returns documentElement.outerHTML
3. get_text - Returns document.body.innerText
4. get_screenshot - Returns URL of screen-size screenshot
5. get_pageshot - Returns URL of full-page screenshot
"""

import os
import asyncio
from typing import Any
import httpx

from mcp.server.models import InitializationOptions
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import (
    Tool,
    TextContent,
)

# Get Reader service URL from environment
READER_URL = os.environ.get("READER_URL", "http://localhost:3000")

# Create server instance
server = Server("reader-mcp")

# HTTP client for making requests to Reader service
client = httpx.AsyncClient(timeout=60.0)


async def fetch_url(url: str, respond_with: str) -> str:
    """
    Fetch a URL through the Reader service with specified response type.

    Args:
        url: The URL to fetch
        respond_with: The response type (markdown, html, text, screenshot, pageshot)

    Returns:
        The response content from the Reader service
    """
    reader_url = f"{READER_URL}/{url}"

    headers = {
        "X-Respond-With": respond_with
    }

    try:
        response = await client.get(reader_url, headers=headers, follow_redirects=True)
        response.raise_for_status()
        return response.text
    except httpx.HTTPError as e:
        return f"Error fetching URL: {str(e)}"


@server.list_tools()
async def handle_list_tools() -> list[Tool]:
    """List available tools."""
    return [
        Tool(
            name="get_markdown",
            description="Convert a URL to markdown format. This bypasses readability processing and returns the raw content as markdown.",
            inputSchema={
                "type": "object",
                "properties": {
                    "url": {
                        "type": "string",
                        "description": "The URL to fetch and convert to markdown"
                    }
                },
                "required": ["url"]
            }
        ),
        Tool(
            name="get_html",
            description="Convert a URL to HTML format. Returns documentElement.outerHTML.",
            inputSchema={
                "type": "object",
                "properties": {
                    "url": {
                        "type": "string",
                        "description": "The URL to fetch and convert to HTML"
                    }
                },
                "required": ["url"]
            }
        ),
        Tool(
            name="get_text",
            description="Convert a URL to plain text format. Returns document.body.innerText.",
            inputSchema={
                "type": "object",
                "properties": {
                    "url": {
                        "type": "string",
                        "description": "The URL to fetch and convert to text"
                    }
                },
                "required": ["url"]
            }
        ),
        Tool(
            name="get_screenshot",
            description="Take a screen-size screenshot of a URL. Returns the URL of the screenshot image.",
            inputSchema={
                "type": "object",
                "properties": {
                    "url": {
                        "type": "string",
                        "description": "The URL to screenshot"
                    }
                },
                "required": ["url"]
            }
        ),
        Tool(
            name="get_pageshot",
            description="Take a full-page screenshot of a URL. Returns the URL of the full-page screenshot image.",
            inputSchema={
                "type": "object",
                "properties": {
                    "url": {
                        "type": "string",
                        "description": "The URL to take a full-page screenshot of"
                    }
                },
                "required": ["url"]
            }
        )
    ]


@server.call_tool()
async def handle_call_tool(name: str, arguments: dict[str, Any]) -> list[TextContent]:
    """Handle tool calls."""
    url = arguments.get("url", "")

    if not url:
        return [TextContent(type="text", text="Error: URL is required")]

    # Map tool names to response types
    tool_to_response = {
        "get_markdown": "markdown",
        "get_html": "html",
        "get_text": "text",
        "get_screenshot": "screenshot",
        "get_pageshot": "pageshot"
    }

    respond_with = tool_to_response.get(name)
    if not respond_with:
        return [TextContent(type="text", text=f"Unknown tool: {name}")]

    result = await fetch_url(url, respond_with)
    return [TextContent(type="text", text=result)]


async def main():
    """Main entry point for the MCP server."""
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="reader-mcp",
                server_version="1.0.0",
                capabilities=server.get_capabilities(
                    notification_options=None,
                    experimental_capabilities={},
                ),
            ),
        )


if __name__ == "__main__":
    asyncio.run(main())
