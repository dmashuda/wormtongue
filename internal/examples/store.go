package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddOptions configures behavior for adding an example.
type AddOptions struct {
	Force bool // overwrite existing file if true
}

// ExampleStore manages a collection of examples from one or more source directories.
type ExampleStore struct {
	sources []string
	index   []Example
	loaded  bool
}

// NewStore creates a new ExampleStore that will scan the given source directories.
func NewStore(sources []string) *ExampleStore {
	return &ExampleStore{sources: sources}
}

// List returns all examples matching the given filter.
func (s *ExampleStore) List(filter Filter) []Example {
	s.ensureLoaded()

	var results []Example
	for _, ex := range s.index {
		if filter.Language != "" && !strings.EqualFold(ex.Language, filter.Language) {
			continue
		}
		if filter.Category != "" && !strings.EqualFold(ex.Category, filter.Category) {
			continue
		}
		results = append(results, ex)
	}
	return results
}

// Get retrieves a specific example by its relative path (e.g. "go/concurrency/worker-pool").
// Returns the example metadata, file content, and any error.
func (s *ExampleStore) Get(path string) (Example, string, error) {
	s.ensureLoaded()

	// Normalize: strip .md extension if provided
	path = strings.TrimSuffix(path, ".md")

	for _, ex := range s.index {
		if ex.Path == path {
			content, err := os.ReadFile(ex.FullPath)
			if err != nil {
				return Example{}, "", fmt.Errorf("reading example: %w", err)
			}
			return ex, string(content), nil
		}
	}
	return Example{}, "", fmt.Errorf("example not found: %s", path)
}

// Search finds examples matching the query string (case-insensitive) in filenames,
// paths, and file contents. Returns at most limit results.
func (s *ExampleStore) Search(query string, limit int) []SearchResult {
	s.ensureLoaded()

	if limit <= 0 {
		limit = 10
	}

	q := strings.ToLower(query)
	var results []SearchResult

	for _, ex := range s.index {
		if len(results) >= limit {
			break
		}

		// Check path/name match
		if strings.Contains(strings.ToLower(ex.Path), q) ||
			strings.Contains(strings.ToLower(ex.Name), q) {
			results = append(results, SearchResult{Example: ex})
			continue
		}

		// Check file content
		content, err := os.ReadFile(ex.FullPath)
		if err != nil {
			continue
		}
		if matchLine := findMatchingLine(string(content), q); matchLine != "" {
			results = append(results, SearchResult{Example: ex, MatchLine: matchLine})
		}
	}

	return results
}

// Add writes a new example to the first source directory.
// Creates language/category directories if they don't exist.
func (s *ExampleStore) Add(language, category, name, content string, opts AddOptions) (Example, error) {
	if err := validateComponent("language", language); err != nil {
		return Example{}, err
	}
	if err := validateComponent("category", category); err != nil {
		return Example{}, err
	}
	if err := validateComponent("name", name); err != nil {
		return Example{}, err
	}
	if strings.TrimSpace(content) == "" {
		return Example{}, fmt.Errorf("content must not be empty")
	}
	if len(s.sources) == 0 {
		return Example{}, fmt.Errorf("no example sources configured")
	}

	root := s.sources[0]
	dir := filepath.Join(root, language, category)
	fullPath := filepath.Join(dir, name+".md")

	if !opts.Force {
		if _, err := os.Stat(fullPath); err == nil {
			return Example{}, fmt.Errorf("example already exists: %s/%s/%s (use force to overwrite)", language, category, name)
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Example{}, fmt.Errorf("creating directories: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return Example{}, fmt.Errorf("writing example: %w", err)
	}

	s.loaded = false

	return Example{
		Name:     name,
		Language: language,
		Category: category,
		Path:     fmt.Sprintf("%s/%s/%s", language, category, name),
		FullPath: fullPath,
	}, nil
}

func validateComponent(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	for _, r := range value {
		isLower := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if !isLower && !isDigit && r != '-' {
			return fmt.Errorf("%s contains invalid character %q: must be lowercase alphanumeric or hyphens", field, r)
		}
	}
	if value[0] == '-' || value[len(value)-1] == '-' {
		return fmt.Errorf("%s must not start or end with a hyphen", field)
	}
	return nil
}

func (s *ExampleStore) ensureLoaded() {
	if s.loaded {
		return
	}
	s.loaded = true
	s.index = nil

	for _, source := range s.sources {
		s.scanSource(source)
	}
}

func (s *ExampleStore) scanSource(root string) {
	// Walk looking for pattern: <root>/<language>/<category>/<name>.md
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) != 3 {
			return nil // must be <language>/<category>/<name>.md
		}

		language := parts[0]
		category := parts[1]
		name := strings.TrimSuffix(parts[2], ".md")

		s.index = append(s.index, Example{
			Name:     name,
			Language: language,
			Category: category,
			Path:     fmt.Sprintf("%s/%s/%s", language, category, name),
			FullPath: path,
		})
		return nil
	})
}

// findMatchingLine returns the first line in content that contains the query (case-insensitive).
func findMatchingLine(content, lowerQuery string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(strings.ToLower(line), lowerQuery) {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 120 {
				trimmed = trimmed[:117] + "..."
			}
			return trimmed
		}
	}
	return ""
}
