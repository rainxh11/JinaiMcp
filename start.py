#!/usr/bin/env python3
"""
Startup script that runs both the Reader service and MCP server.
"""
import subprocess
import threading
import time
import os
import signal
import sys

# Reader service URL from environment
READER_URL = os.environ.get("READER_URL", "http://localhost:3000")
MCP_PORT = os.environ.get("MCP_PORT", "8000")

def run_reader_service():
    """Start the Reader service (Node.js)"""
    try:
        process = subprocess.Popen(
            ["node", "build/server.js"],
            cwd="/app/backend/functions",
            env={**os.environ, "PORT": "3000"}
        )
        process.wait()
    except Exception as e:
        print(f"Reader service error: {e}", flush=True)

def run_mcp_server():
    """Start the MCP server (Python)"""
    # Wait a bit for reader service to start
    time.sleep(3)
    try:
        # Run the MCP server directly
        import asyncio
        from mcp_server.server import main as mcp_main
        asyncio.run(mcp_main())
    except Exception as e:
        print(f"MCP server error: {e}", flush=True)

def handler(signum, frame):
    """Handle shutdown signals"""
    print("Shutting down...", flush=True)
    sys.exit(0)

# Register signal handlers
signal.signal(signal.SIGTERM, handler)
signal.signal(signal.SIGINT, handler)

if __name__ == "__main__":
    print("Starting Reader MCP Server...", flush=True)
    print(f"Reader URL: {READER_URL}", flush=True)
    print(f"MCP Port: {MCP_PORT}", flush=True)

    # Start reader service in a thread
    reader_thread = threading.Thread(target=run_reader_service, daemon=True)
    reader_thread.start()

    # Start MCP server in main thread
    run_mcp_server()
