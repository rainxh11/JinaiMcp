package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	defaultReaderEndpoint = "http://reader-container:3000"
	defaultPort           = "8080"
)

type ReaderServer struct {
	client  *http.Client
	baseURL string
}

func NewReaderServer(baseURL string) *ReaderServer {
	if baseURL == "" {
		baseURL = defaultReaderEndpoint
	}
	return &ReaderServer{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (s *ReaderServer) fetchURL(ctx context.Context, url string, responseType string) (string, error) {
	reqURL := s.baseURL + "/" + url
	log.Printf("[fetchURL] Starting request: url=%q responseType=%q fullURL=%q", url, responseType, reqURL)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		log.Printf("[fetchURL] ERROR creating request: %v", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Respond-With", responseType)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[fetchURL] ERROR fetching URL: %v", err)
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[fetchURL] Response received: statusCode=%d contentType=%q", resp.StatusCode, resp.Header.Get("Content-Type"))

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[fetchURL] ERROR non-200 status: body=%q", string(body))
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[fetchURL] ERROR reading response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("[fetchURL] SUCCESS: responseSize=%d bytes", len(body))
	return string(body), nil
}

// Tool request/response types

type FetchParams struct {
	URL string `json:"url"`
}

type FetchOutput struct {
	Result string `json:"result"`
}

func (s *ReaderServer) FetchMarkdown(ctx context.Context, req *mcp.CallToolRequest, args FetchParams) (any, *FetchOutput, error) {
	result, err := s.fetchURL(ctx, args.URL, "markdown")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, &FetchOutput{Result: result}, nil
}

func (s *ReaderServer) FetchHTML(ctx context.Context, req *mcp.CallToolRequest, args FetchParams) (any, *FetchOutput, error) {
	result, err := s.fetchURL(ctx, args.URL, "html")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, &FetchOutput{Result: result}, nil
}

func (s *ReaderServer) FetchText(ctx context.Context, req *mcp.CallToolRequest, args FetchParams) (any, *FetchOutput, error) {
	result, err := s.fetchURL(ctx, args.URL, "text")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, &FetchOutput{Result: result}, nil
}

func (s *ReaderServer) FetchScreenshot(ctx context.Context, req *mcp.CallToolRequest, args FetchParams) (any, *FetchOutput, error) {
	result, err := s.fetchURL(ctx, args.URL, "screenshot")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, &FetchOutput{Result: result}, nil
}

func (s *ReaderServer) FetchPageshot(ctx context.Context, req *mcp.CallToolRequest, args FetchParams) (any, *FetchOutput, error) {
	result, err := s.fetchURL(ctx, args.URL, "pageshot")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, &FetchOutput{Result: result}, nil
}

// Manual handler for tools that need custom handling
func (s *ReaderServer) manualFetch(responseType string) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		toolName := req.Params.Name
		log.Printf("[MCP Tool] Called: tool=%s", toolName)

		var args FetchParams
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			log.Printf("[MCP Tool] ERROR invalid arguments: tool=%s error=%v", toolName, err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid arguments: %v", err)},
				},
				IsError: true,
			}, nil
		}

		log.Printf("[MCP Tool] Arguments parsed: tool=%s url=%q responseType=%s", toolName, args.URL, responseType)

		result, err := s.fetchURL(ctx, args.URL, responseType)
		if err != nil {
			log.Printf("[MCP Tool] ERROR fetch failed: tool=%s error=%v", toolName, err)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
				},
				IsError: true,
			}, nil
		}

		log.Printf("[MCP Tool] SUCCESS: tool=%s responseSize=%d bytes", toolName, len(result))
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: result},
			},
		}, nil
	}
}

func main() {
	readerEndpoint := os.Getenv("READER_ENDPOINT")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("=== Reader MCP Server Starting ===")
	log.Printf("[Config] READER_ENDPOINT: %s", readerEndpoint)
	log.Printf("[Config] PORT: %s", port)

	server := NewReaderServer(readerEndpoint)
	log.Printf("[Init] Reader server initialized with baseURL: %s", server.baseURL)

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "reader-mcp",
		Version: "1.0.0",
	}, nil)
	log.Println("[Init] MCP server created")

	// Define the input schema
	inputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"url": {
				Type:        "string",
				Description: "The URL to fetch",
			},
		},
		Required: []string{"url"},
	}

	// Register tools using manual handlers
	tools := []struct {
		name        string
		description string
		responseType string
	}{
		{"fetch_markdown", "Fetch a webpage and return its content as Markdown (bypasses readability processing)", "markdown"},
		{"fetch_html", "Fetch a webpage and return its HTML (documentElement.outerHTML)", "html"},
		{"fetch_text", "Fetch a webpage and return its text content (document.body.innerText)", "text"},
		{"fetch_screenshot", "Fetch a screen-size screenshot of a webpage (returns the URL of the screenshot)", "screenshot"},
		{"fetch_pageshot", "Fetch a full-page screenshot of a webpage (returns the URL of the screenshot)", "pageshot"},
	}

	for _, tool := range tools {
		mcpServer.AddTool(&mcp.Tool{
			Name:        tool.name,
			Description: tool.description,
			InputSchema: inputSchema,
		}, server.manualFetch(tool.responseType))
		log.Printf("[Register] Tool registered: %s", tool.name)
	}

	// Create HTTP Streamable handler with Stateless mode to avoid session management
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, &mcp.StreamableHTTPOptions{Stateless: true})

	addr := ":" + port
	log.Printf("[Server] Starting HTTP server on %s", addr)
	log.Printf("[Server] MCP endpoint: http://localhost%s/mcp/http", addr)
	log.Printf("[Server] Total tools registered: %d", len(tools))
	log.Println("=== Ready to accept requests ===")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("[Server] FATAL: %v", err)
	}
}
