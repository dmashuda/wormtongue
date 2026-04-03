package examples

// Example represents a single code style/pattern example in the library.
type Example struct {
	Name     string // filename without extension, e.g. "worker-pool"
	Language string // top-level directory, e.g. "go"
	Category string // subdirectory, e.g. "concurrency"
	Path     string // relative path: "go/concurrency/worker-pool"
	FullPath string // absolute filesystem path to the .md file
}

// Filter specifies optional criteria for listing examples.
type Filter struct {
	Language string
	Category string
}

// SearchResult pairs a matched example with the first matching line.
type SearchResult struct {
	Example   Example
	MatchLine string // trimmed line that matched the query
}
