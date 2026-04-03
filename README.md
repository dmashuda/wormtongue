# wormtongue

A CLI and MCP server that provides curated code style and pattern examples for LLMs.

Maintain a library of well-written code examples organized by language and category. Query them via CLI commands or let LLMs access them directly through the Model Context Protocol.

## Install

```bash
go install github.com/dmashuda/wormtongue@latest
```

Or download a binary from [Releases](https://github.com/dmashuda/wormtongue/releases).

## Usage

### List examples

```bash
wormtongue list
wormtongue list --language go
wormtongue list --language csharp --category patterns
```

### Show an example

```bash
wormtongue show go/concurrency/worker-pool
```

### Search examples

```bash
wormtongue search "async"
wormtongue search "error handling" --limit 5
```

### Add an example

```bash
# With inline content
wormtongue add go testing table-tests --content "# Table Tests\n\nContent here."

# From a file via stdin
cat example.md | wormtongue add go testing table-tests

# Overwrite an existing example
wormtongue add go testing table-tests --content "# Updated" --force
```

### Manage external sources

```bash
wormtongue source add team-standards /path/to/team/examples
wormtongue source list
wormtongue source remove team-standards
```

### Start MCP server

```bash
wormtongue serve
```

## MCP Configuration

Add wormtongue to your MCP client configuration:

### Claude Desktop / Claude Code

```json
{
  "mcpServers": {
    "wormtongue": {
      "command": "wormtongue",
      "args": ["serve"]
    }
  }
}
```

### MCP Tools

The MCP server exposes four tools:

- **list_examples** — List available examples, optionally filtered by `language` and `category`
- **get_example** — Retrieve the full content of an example by path
- **search_examples** — Search examples by keyword with optional `limit`
- **add_example** — Add a new example with `language`, `category`, `name`, and `content` (optional `force` to overwrite)

## Example Format

Examples are markdown files organized as `examples/<language>/<category>/<name>.md`:

```
examples/
  go/
    concurrency/
      worker-pool.md
    error-handling/
      sentinel-errors.md
    patterns/
      functional-options.md
  csharp/
    async/
      async-await-patterns.md
    patterns/
      builder-pattern.md
```

Names must be lowercase alphanumeric with hyphens (e.g. `worker-pool`, `async-await-patterns`). New examples are written to the first configured source directory.

Each file follows this structure:

```markdown
# Example Title

Brief description.

## When to Use
- Use case bullets

## Example
\```go
// code here
\```

## Key Points
- Important takeaways
```

## Configuration

Config file: `~/.config/wormtongue/config.yaml`

```yaml
sources:
  - name: team-standards
    path: /path/to/external/examples
```

The built-in `examples/` directory is always included. Set `WORMTONGUE_EXAMPLES` to override its location.

## Development

```bash
make build       # Build binary to bin/
make test        # Run tests
make vet         # Run go vet
make run ARGS="list"  # Run with local examples
```
