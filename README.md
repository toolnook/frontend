# DocForge

A documentation-as-code generator that turns annotated source code into a static, searchable API reference site. Think of it as a lightweight, dependency-minimal alternative to Swagger/OpenAPI for teams that don't want a full framework.

## What it does

DocForge scans your source files for custom annotation comments and generates a complete static documentation website:

- **Parses annotations** from `@doc`, `@param`, `@returns`, and `@example` comment tags
- **Supports TypeScript and Python** (extensible to other languages)
- **Generates static HTML/CSS/JS** with built-in search
- **Visualizes dependencies** as an interactive SVG graph
- **Reports documentation coverage** — shows which public functions lack docs
- **Serves locally** for development preview

## Why it exists

Most API documentation tools (Swagger, TypeDoc, Sphinx) are powerful but heavy — they require specific frameworks, build tooling, or configuration formats. DocForge uses a simple comment convention that works in any language and produces a standalone static site with zero runtime dependencies.

## Quick start

```bash
# Clone and build
git clone https://github.com/toolnook/docforge.git
cd docforge
go build -o docforge ./cmd/docforge

# Scan source files
./docforge scan ./src

# Generate documentation site
./docforge build -dir ./src -o ./docs

# Serve locally
./docforge serve -d ./docs -port 8090
```

## Annotation format

### TypeScript

```typescript
/**
 * @doc Adds two numbers together
 * @param {number} a - The first number
 * @param {number} b - The second number
 * @returns {number} The sum of a and b
 * @example
 * add(1, 2) // returns 3
 */
export function add(a: number, b: number): number {
  return a + b;
}
```

### Python

```python
# @doc Adds two numbers
# @param a (int) The first number
# @param b (int) The second number
# @returns (int) The sum of a and b
# @example
# add(1, 2)  # returns 3
def add(a, b):
    return a + b
```

## CLI commands

### `docforge scan <dir>`
Parse source files and output a summary of documented symbols.

```bash
./docforge scan ./src -format json
```

### `docforge build`
Generate the static documentation site.

```bash
./docforge build -dir ./src -o ./docs -version 1.0.0
```

Options:
- `-dir` — Source directory to scan (default: `.`)
- `-o` — Output directory (default: `./docs`)
- `-version` — Version tag for the docs (default: `0.1.0`)
- `-name` — Project name (defaults to directory name)

### `docforge serve`
Serve generated docs locally.

```bash
./docforge serve -d ./docs -port 8090
```

## Architecture overview

```
┌─────────────────────────────────────────────────────┐
│                    CLI (main.go)                      │
│       scan / build / serve command dispatch           │
├─────────────┬──────────────┬────────────────────────┤
│   Parser    │  Generator   │  Server                 │
│ (parser/)   │ (generator/) │ (server/)              │
│             │              │                         │
│ TypeScript  │ HTML pages    │ Static file             │
│ Python      │ CSS + JS      │ serving                 │
│ Registry    │ Search index  │                         │
├─────────────┤ Graph data    ├────────────────────────┤
│   Search    │ Coverage      │  Graph                  │
│ (search/)   │ (coverage/)   │ (graph/)                │
│             │               │                         │
│ Index build │ Coverage %    │ Dependency              │
│ Query       │ Missing list  │ graph + SVG             │
├─────────────┴───────────────┴────────────────────────┤
│              Types (types.go)                         │
│  Document, Module, Project, CoverageReport           │
└─────────────────────────────────────────────────────┘
```

| Package | Responsibility |
|---------|---------------|
| `internal/types` | Core data structures (Document, Module, Project) |
| `internal/parser` | Source file parsing with annotation extraction |
| `internal/generator` | Static HTML/CSS/JS site generation |
| `internal/search` | Search index building and querying |
| `internal/graph` | Dependency graph construction and visualization |
| `internal/coverage` | Documentation coverage calculation |
| `internal/server` | Local HTTP server for development |

## Unique features

### Dependency graph

DocForge analyzes import statements across modules and generates an interactive SVG visualization showing which modules depend on which. The graph is rendered client-side with size proportional to dependency degree.

### Coverage indicator

The generated site shows a progress bar indicating what percentage of public symbols have documentation. Undocumented symbols are highlighted with a "missing docs" badge, making it easy to identify gaps.

### Annotation convention

The annotation format is deliberately simple — just `@tag value` in comments. No special syntax, no custom file format, no preprocessing step. It works with any editor and any language that supports comments.

## Running tests

```bash
go test ./... -v
```

## Roadmap

- [ ] Support for Go, Rust, and Java annotations
- [ ] Markdown page support (intermixed with API docs)
- [ ] Custom theme support
- [ ] Versioned documentation with diff view
- [ ] Export to OpenAPI/JSON Schema
- [ ] Plugin system for custom parsers

## License

MIT — see [LICENSE](LICENSE).
