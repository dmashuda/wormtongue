package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupCmdFixtures(t *testing.T) (examplesDir string, configPath string) {
	t.Helper()
	dir := t.TempDir()
	examplesDir = filepath.Join(dir, "examples")

	files := map[string]string{
		"go/concurrency/worker-pool.md": "# Worker Pool\n\nDistributes work across goroutines.\n",
		"go/patterns/options.md":        "# Functional Options\n\nConfigure structs with option functions.\n",
		"csharp/async/async-await.md":   "# Async/Await\n\nAsynchronous programming in C#.\n",
	}
	for relPath, content := range files {
		fullPath := filepath.Join(examplesDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	configPath = filepath.Join(dir, "config.yaml")
	return examplesDir, configPath
}

func executeCmd(t *testing.T, examplesDir string, configPath string, args ...string) (string, error) {
	t.Helper()

	// Reset package-level state between tests
	store = nil

	t.Setenv("WORMTONGUE_EXAMPLES", examplesDir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"--config", configPath}, args...))

	err := rootCmd.Execute()
	return buf.String(), err
}

// list tests

func TestCmdList_All(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "worker-pool") {
		t.Errorf("expected worker-pool in output, got:\n%s", out)
	}
	if !strings.Contains(out, "async-await") {
		t.Errorf("expected async-await in output, got:\n%s", out)
	}
}

func TestCmdList_FilterLanguage(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "list", "--language", "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "csharp") {
		t.Errorf("expected no csharp when filtering by go, got:\n%s", out)
	}
	if !strings.Contains(out, "worker-pool") {
		t.Errorf("expected worker-pool in output, got:\n%s", out)
	}
}

func TestCmdList_FilterCategory(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "list", "--category", "concurrency")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "worker-pool") {
		t.Errorf("expected worker-pool, got:\n%s", out)
	}
	if strings.Contains(out, "options") {
		t.Errorf("expected no patterns results, got:\n%s", out)
	}
}

func TestCmdList_NoResults(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "list", "--language", "rust")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No examples found") {
		t.Errorf("expected 'No examples found', got:\n%s", out)
	}
}

// show tests

func TestCmdShow_Valid(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "show", "go/concurrency/worker-pool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Worker Pool") {
		t.Errorf("expected example content, got:\n%s", out)
	}
}

func TestCmdShow_NotFound(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	_, err := executeCmd(t, exDir, cfgPath, "show", "go/missing/example")
	if err == nil {
		t.Error("expected error for missing example")
	}
}

func TestCmdShow_MissingArg(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	_, err := executeCmd(t, exDir, cfgPath, "show")
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

// search tests

func TestCmdSearch_Match(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "search", "worker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "worker-pool") {
		t.Errorf("expected worker-pool in results, got:\n%s", out)
	}
}

func TestCmdSearch_NoMatch(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "search", "zzzznonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No matching examples found") {
		t.Errorf("expected 'No matching examples found', got:\n%s", out)
	}
}

func TestCmdSearch_Limit(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	out, err := executeCmd(t, exDir, cfgPath, "search", "#", "--limit", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Count result lines (paths contain "/")
	count := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "/") && !strings.HasPrefix(strings.TrimSpace(line), "...") {
			count++
		}
	}
	if count > 1 {
		t.Errorf("expected at most 1 result with --limit 1, got %d:\n%s", count, out)
	}
}

// source tests

func TestCmdSource_Lifecycle(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)

	// Add a source
	_, err := executeCmd(t, exDir, cfgPath, "source", "add", "test-source", exDir)
	if err != nil {
		t.Fatalf("source add failed: %v", err)
	}

	// List sources
	out, err := executeCmd(t, exDir, cfgPath, "source", "list")
	if err != nil {
		t.Fatalf("source list failed: %v", err)
	}
	if !strings.Contains(out, "test-source") {
		t.Errorf("expected test-source in list, got:\n%s", out)
	}

	// Remove source
	_, err = executeCmd(t, exDir, cfgPath, "source", "remove", "test-source")
	if err != nil {
		t.Fatalf("source remove failed: %v", err)
	}

	// Verify removed
	out, err = executeCmd(t, exDir, cfgPath, "source", "list")
	if err != nil {
		t.Fatalf("source list after remove failed: %v", err)
	}
	if !strings.Contains(out, "No external sources registered") {
		t.Errorf("expected no sources after remove, got:\n%s", out)
	}
}

func TestCmdSource_AddDuplicate(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)

	_, err := executeCmd(t, exDir, cfgPath, "source", "add", "dup", exDir)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	_, err = executeCmd(t, exDir, cfgPath, "source", "add", "dup", exDir)
	if err == nil {
		t.Error("expected error for duplicate source name")
	}
}

func TestCmdSource_RemoveNotFound(t *testing.T) {
	exDir, cfgPath := setupCmdFixtures(t)
	_, err := executeCmd(t, exDir, cfgPath, "source", "remove", "nonexistent")
	if err == nil {
		t.Error("expected error for removing nonexistent source")
	}
}
