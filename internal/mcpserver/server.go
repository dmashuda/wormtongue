package mcpserver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dmashuda/wormtongue/internal/examples"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Run starts the MCP server over stdio, registering tools that query the example store.
func Run(ctx context.Context, store *examples.ExampleStore) error {
	s := server.NewMCPServer(
		"wormtongue",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerTools(s, store)

	transport := server.NewStdioServer(s)
	return transport.Listen(ctx, os.Stdin, os.Stdout)
}

func registerTools(s *server.MCPServer, store *examples.ExampleStore) {
	// list_examples
	s.AddTool(
		mcp.NewTool("list_examples",
			mcp.WithDescription("List available code style/pattern examples, optionally filtered by language and category"),
			mcp.WithString("language", mcp.Description("Filter by programming language (e.g. go, csharp)")),
			mcp.WithString("category", mcp.Description("Filter by category (e.g. concurrency, patterns)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			language := req.GetString("language", "")
			category := req.GetString("category", "")

			results := store.List(examples.Filter{
				Language: language,
				Category: category,
			})

			if len(results) == 0 {
				return mcp.NewToolResultText("No examples found."), nil
			}

			var sb strings.Builder
			for _, ex := range results {
				fmt.Fprintf(&sb, "%s/%s/%s\n", ex.Language, ex.Category, ex.Name)
			}
			return mcp.NewToolResultText(sb.String()), nil
		},
	)

	// get_example
	s.AddTool(
		mcp.NewTool("get_example",
			mcp.WithDescription("Retrieve the full content of a specific code example by its path"),
			mcp.WithString("path", mcp.Required(), mcp.Description("Example path (e.g. go/concurrency/worker-pool)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path, err := req.RequireString("path")
			if err != nil {
				return mcp.NewToolResultError("path is required"), nil
			}

			_, content, err := store.Get(path)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(content), nil
		},
	)

	// add_example
	s.AddTool(
		mcp.NewTool("add_example",
			mcp.WithDescription("Add a new code style/pattern example to the library. Use this to capture coding styles and patterns from user feedback or PR reviews."),
			mcp.WithString("language", mcp.Required(), mcp.Description("Programming language (e.g. go, csharp)")),
			mcp.WithString("category", mcp.Required(), mcp.Description("Category (e.g. concurrency, patterns)")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Example name in kebab-case (e.g. worker-pool)")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Full markdown content of the example")),
			mcp.WithBoolean("force", mcp.Description("Overwrite existing example if true (default false)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			language, err := req.RequireString("language")
			if err != nil {
				return mcp.NewToolResultError("language is required"), nil
			}
			category, err := req.RequireString("category")
			if err != nil {
				return mcp.NewToolResultError("category is required"), nil
			}
			name, err := req.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError("name is required"), nil
			}
			content, err := req.RequireString("content")
			if err != nil {
				return mcp.NewToolResultError("content is required"), nil
			}
			force := req.GetBool("force", false)

			ex, addErr := store.Add(language, category, name, content, examples.AddOptions{
				Force: force,
			})
			if addErr != nil {
				return mcp.NewToolResultError(addErr.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Added example: %s", ex.Path)), nil
		},
	)

	// list_languages
	s.AddTool(
		mcp.NewTool("list_languages",
			mcp.WithDescription("List all available programming languages in the example library"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			langs := store.Languages()
			if len(langs) == 0 {
				return mcp.NewToolResultText("No languages found."), nil
			}

			return mcp.NewToolResultText(strings.Join(langs, "\n")), nil
		},
	)

	// search_examples
	s.AddTool(
		mcp.NewTool("search_examples",
			mcp.WithDescription("Search code examples by keyword across names, categories, and content"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of results (default 10)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := req.RequireString("query")
			if err != nil {
				return mcp.NewToolResultError("query is required"), nil
			}

			limit := req.GetInt("limit", 10)

			results := store.Search(query, limit)
			if len(results) == 0 {
				return mcp.NewToolResultText("No matching examples found."), nil
			}

			var sb strings.Builder
			for _, r := range results {
				fmt.Fprintf(&sb, "%s/%s/%s\n", r.Example.Language, r.Example.Category, r.Example.Name)
				if r.MatchLine != "" {
					fmt.Fprintf(&sb, "  ...%s...\n", r.MatchLine)
				}
			}
			return mcp.NewToolResultText(sb.String()), nil
		},
	)
}
