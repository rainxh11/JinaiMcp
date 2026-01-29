package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	readerURL = os.Getenv("READER_URL")
	mcpPort   = getEnv("MCP_PORT", "8000")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// MCPRequest represents a JSON-RPC request
type MCPRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Method  string        `json:"method"`
	Params  interface{}   `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents a JSON-RPC error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolCallParams represents parameters for tool calls
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Content represents MCP content
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// fetchURL fetches a URL through the Reader service
func fetchURL(url, respondWith string) (string, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	readerURL := fmt.Sprintf("%s/%s", readerURL, url)

	req, err := http.NewRequest("GET", readerURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Respond-With", respondWith)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// listTools returns all available tools
func listTools() []Tool {
	return []Tool{
		{
			Name:        "get_markdown",
			Description: "Convert a URL to markdown format. This bypasses readability processing and returns the raw content as markdown.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to fetch and convert to markdown",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "get_html",
			Description: "Convert a URL to HTML format. Returns documentElement.outerHTML.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to fetch and convert to HTML",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "get_text",
			Description: "Convert a URL to plain text format. Returns document.body.innerText.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to fetch and convert to text",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "get_screenshot",
			Description: "Take a screen-size screenshot of a URL. Returns the URL of the screenshot image.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to screenshot",
					},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "get_pageshot",
			Description: "Take a full-page screenshot of a URL. Returns the URL of the full-page screenshot image.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "The URL to take a full-page screenshot of",
					},
				},
				"required": []string{"url"},
			},
		},
	}
}

// callTool executes a tool call
func callTool(name string, arguments map[string]interface{}) (*MCPResponse, error) {
	url, ok := arguments["url"].(string)
	if !ok || url == "" {
		return &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params: url is required",
			},
		}, nil
	}

	toolToResponse := map[string]string{
		"get_markdown":   "markdown",
		"get_html":       "html",
		"get_text":       "text",
		"get_screenshot": "screenshot",
		"get_pageshot":   "pageshot",
	}

	respondWith, ok := toolToResponse[name]
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Unknown tool: %s", name),
			},
		}, nil
	}

	result, err := fetchURL(url, respondWith)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32603,
				Message: fmt.Sprintf("Error fetching URL: %v", err),
			},
		}, nil
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"content": []Content{
				{
					Type: "text",
					Text: result,
				},
			},
		},
	}, nil
}

// handleMCP processes MCP requests
func handleMCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":        "reader-mcp",
			"version":     "1.0.0",
			"description": "MCP server for Reader URL to LLM-friendly conversion",
		})
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "reader-mcp",
				"version": "1.0.0",
			},
		}

	case "tools/list":
		tools := listTools()
		toolsData := make([]map[string]interface{}, len(tools))
		for i, tool := range tools {
			toolsData[i] = map[string]interface{}{
				"name":        tool.Name,
				"description": tool.Description,
				"inputSchema":  tool.InputSchema,
			}
		}
		resp.Result = map[string]interface{}{
			"tools": toolsData,
		}

	case "tools/call":
		paramsData, _ := req.Params.(map[string]interface{})
		name, _ := paramsData["name"].(string)
		arguments, _ := paramsData["arguments"].(map[string]interface{})

		result, err := callTool(name, arguments)
		if err != nil {
			resp.Error = &MCPError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			resp.Result = result.Result
			resp.Error = result.Error
		}

	default:
		resp.Error = &MCPError{
			Code:    -32601,
			Message: fmt.Sprintf("Method not found: %s", req.Method),
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func handleStreamableHTTP(w http.ResponseWriter, r *http.Request) {
	// For now, use the same handler
	handleMCP(w, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"reader": readerURL,
	})
}

func main() {
	if readerURL == "" {
		readerURL = "http://localhost:3000"
	}

	port, _ := strconv.Atoi(mcpPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStreamableHTTP)
	mux.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf(":%d", port)

	log.Printf("Starting MCP server on port %s...", mcpPort)
	log.Printf("Reader service URL: %s", readerURL)
	log.Printf("MCP Streamable HTTP endpoint: http://0.0.0.0%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
