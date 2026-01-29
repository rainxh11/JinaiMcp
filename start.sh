#!/bin/sh
set -e

echo "Starting Reader MCP Server..."
echo "Reader URL: ${READER_URL:-http://localhost:3000}"
echo "MCP Port: ${MCP_PORT:-8000}"

# Start Deno Reader service in background
deno run --allow-net --allow-read --allow-write --allow-env --allow-run --allow-sys /app/deno/main.ts &

# Wait for Reader service to start
sleep 3

# Start Go MCP server in foreground
exec /app/mcp-server
