package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmashuda/wormtongue/internal/examples"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func setupTestServer(t *testing.T) *server.MCPServer {
	t.Helper()
	s, _ := setupTestServerWithDir(t)
	return s
}

func setupTestServerWithDir(t *testing.T) (*server.MCPServer, string) {
	t.Helper()
	dir := setupFixtures(t)
	store := examples.NewStore([]string{dir})

	s := server.NewMCPServer("wormtongue-test", "0.0.1", server.WithToolCapabilities(true))
	registerTools(s, store)
	return s, dir
}

func setupFixtures(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"go/concurrency/worker-pool.md": "# Worker Pool\n\nDistributes work across goroutines.\n",
		"go/patterns/options.md":        "# Functional Options\n\nConfigure structs with option functions.\n",
		"csharp/async/async-await.md":   "# Async/Await\n\nAsynchronous programming in C#.\n",
	}
	for relPath, content := range files {
		fullPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func callTool(t *testing.T, s *server.MCPServer, name string, args map[string]any) mcp.JSONRPCMessage {
	t.Helper()
	params := map[string]any{
		"name":      name,
		"arguments": args,
	}
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params":  params,
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// Need to initialize the server first
	initReq := map[string]any{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"clientInfo":      map[string]any{"name": "test", "version": "0.0.1"},
			"capabilities":    map[string]any{},
		},
	}
	initRaw, _ := json.Marshal(initReq)
	s.HandleMessage(context.Background(), initRaw)

	return s.HandleMessage(context.Background(), raw)
}

type toolResponse struct {
	Result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	} `json:"result"`
}

func parseToolResponse(t *testing.T, msg mcp.JSONRPCMessage) toolResponse {
	t.Helper()
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var resp toolResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

func getTextContent(t *testing.T, msg mcp.JSONRPCMessage) string {
	t.Helper()
	resp := parseToolResponse(t, msg)
	if len(resp.Result.Content) == 0 {
		t.Fatal("expected content in response")
	}
	return resp.Result.Content[0].Text
}

// list_examples tests

func TestListExamples_All(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "list_examples", map[string]any{}))

	if text == "No examples found." {
		t.Fatal("expected examples to be listed")
	}
	// Should contain all 3 examples
	for _, expected := range []string{"go/concurrency/worker-pool", "go/patterns/options", "csharp/async/async-await"} {
		if !contains(text, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, text)
		}
	}
}

func TestListExamples_FilterByLanguage(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "list_examples", map[string]any{"language": "go"}))

	if contains(text, "csharp") {
		t.Errorf("expected no csharp results when filtering by go, got:\n%s", text)
	}
	if !contains(text, "go/concurrency/worker-pool") {
		t.Errorf("expected go example in output, got:\n%s", text)
	}
}

func TestListExamples_FilterByCategory(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "list_examples", map[string]any{"category": "concurrency"}))

	if !contains(text, "worker-pool") {
		t.Errorf("expected worker-pool in output, got:\n%s", text)
	}
	if contains(text, "options") {
		t.Errorf("expected no patterns results, got:\n%s", text)
	}
}

func TestListExamples_NoResults(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "list_examples", map[string]any{"language": "rust"}))

	if text != "No examples found." {
		t.Errorf("expected 'No examples found.', got: %s", text)
	}
}

// get_example tests

func TestGetExample_Valid(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "get_example", map[string]any{"path": "go/concurrency/worker-pool"}))

	if !contains(text, "Worker Pool") {
		t.Errorf("expected example content, got:\n%s", text)
	}
}

func TestGetExample_NotFound(t *testing.T) {
	s := setupTestServer(t)
	resp := parseToolResponse(t, callTool(t, s, "get_example", map[string]any{"path": "go/missing/example"}))

	if !resp.Result.IsError {
		t.Error("expected error response for missing example")
	}
}

func TestGetExample_MissingParam(t *testing.T) {
	s := setupTestServer(t)
	resp := parseToolResponse(t, callTool(t, s, "get_example", map[string]any{}))

	if !resp.Result.IsError {
		t.Error("expected error response for missing path param")
	}
}

// search_examples tests

func TestSearchExamples_Match(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "search_examples", map[string]any{"query": "worker"}))

	if !contains(text, "worker-pool") {
		t.Errorf("expected worker-pool in results, got:\n%s", text)
	}
}

func TestSearchExamples_WithLimit(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "search_examples", map[string]any{"query": "#", "limit": 1}))

	// Should only have 1 result line (path)
	lines := nonEmptyLines(text)
	if len(lines) > 2 { // path + possible match line
		t.Errorf("expected at most 1 result with limit=1, got %d lines:\n%s", len(lines), text)
	}
}

func TestSearchExamples_NoResults(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "search_examples", map[string]any{"query": "zzzznonexistent"}))

	if text != "No matching examples found." {
		t.Errorf("expected 'No matching examples found.', got: %s", text)
	}
}

func TestSearchExamples_MissingQuery(t *testing.T) {
	s := setupTestServer(t)
	resp := parseToolResponse(t, callTool(t, s, "search_examples", map[string]any{}))

	if !resp.Result.IsError {
		t.Error("expected error response for missing query param")
	}
}

// add_example tests

func TestAddExample_Success(t *testing.T) {
	s, dir := setupTestServerWithDir(t)
	text := getTextContent(t, callTool(t, s, "add_example", map[string]any{
		"language": "go",
		"category": "testing",
		"name":     "table-tests",
		"content":  "# Table Tests\n\nContent here.\n",
	}))

	if !contains(text, "Added example: go/testing/table-tests") {
		t.Errorf("expected success message, got: %s", text)
	}

	// Verify file on disk
	data, err := os.ReadFile(filepath.Join(dir, "go", "testing", "table-tests.md"))
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if string(data) != "# Table Tests\n\nContent here.\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestAddExample_MissingParams(t *testing.T) {
	s := setupTestServer(t)

	tests := []struct {
		name string
		args map[string]any
	}{
		{"missing language", map[string]any{"category": "cat", "name": "n", "content": "c"}},
		{"missing category", map[string]any{"language": "go", "name": "n", "content": "c"}},
		{"missing name", map[string]any{"language": "go", "category": "cat", "content": "c"}},
		{"missing content", map[string]any{"language": "go", "category": "cat", "name": "n"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := parseToolResponse(t, callTool(t, s, "add_example", tc.args))
			if !resp.Result.IsError {
				t.Error("expected error response")
			}
		})
	}
}

func TestAddExample_InvalidName(t *testing.T) {
	s := setupTestServer(t)
	resp := parseToolResponse(t, callTool(t, s, "add_example", map[string]any{
		"language": "Go",
		"category": "testing",
		"name":     "table-tests",
		"content":  "# Content\n",
	}))

	if !resp.Result.IsError {
		t.Error("expected error for invalid language")
	}
}

func TestAddExample_Duplicate(t *testing.T) {
	s := setupTestServer(t)
	resp := parseToolResponse(t, callTool(t, s, "add_example", map[string]any{
		"language": "go",
		"category": "concurrency",
		"name":     "worker-pool",
		"content":  "# Duplicate\n",
	}))

	if !resp.Result.IsError {
		t.Error("expected error for duplicate example")
	}
}

func TestAddExample_Force(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "add_example", map[string]any{
		"language": "go",
		"category": "concurrency",
		"name":     "worker-pool",
		"content":  "# Updated\n",
		"force":    true,
	}))

	if !contains(text, "Added example") {
		t.Errorf("expected success with force, got: %s", text)
	}
}

func TestAddExample_ThenGet(t *testing.T) {
	s := setupTestServer(t)

	// Add
	getTextContent(t, callTool(t, s, "add_example", map[string]any{
		"language": "python",
		"category": "web",
		"name":     "flask-routes",
		"content":  "# Flask Routes\n\nRoute examples.\n",
	}))

	// Get
	text := getTextContent(t, callTool(t, s, "get_example", map[string]any{
		"path": "python/web/flask-routes",
	}))

	if !contains(text, "Flask Routes") {
		t.Errorf("expected added content from get, got: %s", text)
	}
}

// list_languages tests

func TestListLanguages_All(t *testing.T) {
	s := setupTestServer(t)
	text := getTextContent(t, callTool(t, s, "list_languages", map[string]any{}))

	if text == "No languages found." {
		t.Fatal("expected languages to be listed")
	}
	for _, expected := range []string{"go", "csharp"} {
		if !contains(text, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, text)
		}
	}
}

func TestListLanguages_Empty(t *testing.T) {
	dir := t.TempDir()
	store := examples.NewStore([]string{dir})
	s := server.NewMCPServer("wormtongue-test", "0.0.1", server.WithToolCapabilities(true))
	registerTools(s, store)

	text := getTextContent(t, callTool(t, s, "list_languages", map[string]any{}))
	if text != "No languages found." {
		t.Errorf("expected 'No languages found.', got: %s", text)
	}
}

// helpers

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func nonEmptyLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	return lines
}
