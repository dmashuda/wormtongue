package examples

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestFixtures(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create test example structure
	files := map[string]string{
		"go/concurrency/worker-pool.md": "# Worker Pool\n\nDistributes work across goroutines.\n\n```go\nfunc worker() {}\n```\n",
		"go/patterns/options.md":        "# Functional Options\n\nConfigure structs with option functions.\n\n```go\ntype Option func(*Config)\n```\n",
		"csharp/async/async-await.md":   "# Async/Await\n\nAsynchronous programming in C#.\n\n```csharp\nasync Task DoWork() {}\n```\n",
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

func TestList_All(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.List(Filter{})
	if len(results) != 3 {
		t.Fatalf("expected 3 examples, got %d", len(results))
	}
}

func TestList_FilterByLanguage(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.List(Filter{Language: "go"})
	if len(results) != 2 {
		t.Fatalf("expected 2 Go examples, got %d", len(results))
	}
	for _, ex := range results {
		if ex.Language != "go" {
			t.Errorf("expected language 'go', got %q", ex.Language)
		}
	}
}

func TestList_FilterByCategory(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.List(Filter{Category: "concurrency"})
	if len(results) != 1 {
		t.Fatalf("expected 1 concurrency example, got %d", len(results))
	}
	if results[0].Name != "worker-pool" {
		t.Errorf("expected 'worker-pool', got %q", results[0].Name)
	}
}

func TestGet(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	ex, content, err := store.Get("go/concurrency/worker-pool")
	if err != nil {
		t.Fatal(err)
	}
	if ex.Name != "worker-pool" {
		t.Errorf("expected name 'worker-pool', got %q", ex.Name)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}

func TestGet_WithExtension(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	_, _, err := store.Get("go/concurrency/worker-pool.md")
	if err != nil {
		t.Fatalf("expected .md extension to be handled, got error: %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	_, _, err := store.Get("go/missing/example")
	if err == nil {
		t.Fatal("expected error for missing example")
	}
}

func TestSearch_ByName(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.Search("worker", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Example.Name != "worker-pool" {
		t.Errorf("expected 'worker-pool', got %q", results[0].Example.Name)
	}
}

func TestSearch_ByContent(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.Search("goroutines", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].MatchLine == "" {
		t.Error("expected a match line for content match")
	}
}

func TestSearch_Limit(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	// Search for something that matches all examples (they all contain markdown headers)
	results := store.Search("#", 2)
	if len(results) > 2 {
		t.Errorf("expected at most 2 results with limit, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	dir := setupTestFixtures(t)
	store := NewStore([]string{dir})

	results := store.Search("WORKER", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}
}

func TestMultipleSources(t *testing.T) {
	dir1 := setupTestFixtures(t)

	// Create a second source with different examples
	dir2 := t.TempDir()
	path := filepath.Join(dir2, "rust/memory/ownership.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# Ownership\n\nRust ownership model.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	store := NewStore([]string{dir1, dir2})

	results := store.List(Filter{})
	if len(results) != 4 {
		t.Fatalf("expected 4 examples from 2 sources, got %d", len(results))
	}

	rustResults := store.List(Filter{Language: "rust"})
	if len(rustResults) != 1 {
		t.Fatalf("expected 1 rust example, got %d", len(rustResults))
	}
}
