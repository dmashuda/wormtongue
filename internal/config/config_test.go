package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNonExistent(t *testing.T) {
	cfg, err := Load("/tmp/wormtongue-test-nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(cfg.Sources) != 0 {
		t.Fatalf("expected empty sources, got %d", len(cfg.Sources))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &Config{
		Sources: []Source{
			{Name: "test", Path: "/tmp/examples"},
		},
	}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(loaded.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(loaded.Sources))
	}
	if loaded.Sources[0].Name != "test" {
		t.Errorf("expected name 'test', got %q", loaded.Sources[0].Name)
	}
	if loaded.Sources[0].Path != "/tmp/examples" {
		t.Errorf("expected path '/tmp/examples', got %q", loaded.Sources[0].Path)
	}
}

func TestSaveCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "config.yaml")

	cfg := &Config{}
	if err := Save(path, cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
