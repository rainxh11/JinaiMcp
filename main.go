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
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Respond-With", responseType)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

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
		var args FetchParams
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid arguments: %v", err)},
				},
				IsError: true,
			}, nil
		}

		result, err := s.fetchURL(ctx, args.URL, responseType)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
				},
				IsError: true,
			}, nil
		}

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

	server := NewReaderServer(readerEndpoint)

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "reader-mcp",
		Version: "1.0.0",
	}, nil)

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
	mcpServer.AddTool(&mcp.Tool{
		Name:        "fetch_markdown",
		Description: "Fetch a webpage and return its content as Markdown (bypasses readability processing)",
		InputSchema: inputSchema,
	}, server.manualFetch("markdown"))

	mcpServer.AddTool(&mcp.Tool{
		Name:        "fetch_html",
		Description: "Fetch a webpage and return its HTML (documentElement.outerHTML)",
		InputSchema: inputSchema,
	}, server.manualFetch("html"))

	mcpServer.AddTool(&mcp.Tool{
		Name:        "fetch_text",
		Description: "Fetch a webpage and return its text content (document.body.innerText)",
		InputSchema: inputSchema,
	}, server.manualFetch("text"))

	mcpServer.AddTool(&mcp.Tool{
		Name:        "fetch_screenshot",
		Description: "Fetch a screen-size screenshot of a webpage (returns the URL of the screenshot)",
		InputSchema: inputSchema,
	}, server.manualFetch("screenshot"))

	mcpServer.AddTool(&mcp.Tool{
		Name:        "fetch_pageshot",
		Description: "Fetch a full-page screenshot of a webpage (returns the URL of the screenshot)",
		InputSchema: inputSchema,
	}, server.manualFetch("pageshot"))

	// Create HTTP Streamable handler
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	addr := ":" + port
	log.Printf("Reader MCP Server starting on %s, connecting to: %s", addr, server.baseURL)
	log.Printf("HTTP MCP endpoint available at http://localhost%s/mcp/http", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
