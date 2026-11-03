# Toolnook Frontend

The documentation generator for the Toolnook platform. Turns annotated source code into a static, searchable API reference site with dependency graphs and coverage indicators.

Part of the [Toolnook](https://github.com/toolnook) open-source developer infrastructure platform:

| Repo | Role |
|------|------|
| [toolnook/backend](https://github.com/toolnook/backend) | API gateway, reverse proxy, rate limiting, circuit breaker |
| **toolnook/frontend** | Documentation-as-code generator with dependency graphs |
| [toolnook/contracts](https://github.com/toolnook/contracts) | Double-entry accounting CLI with project profitability |

## What it does

Toolnook Frontend scans your source files for custom annotation comments and generates a complete static documentation website:

- **Parses** `@doc`, `@param`, `@returns`, `@example` annotations
- **Supports** TypeScript and Python (extensible to more languages)
- **Generates** static HTML/CSS/JS with built-in search
- **Visualizes** module dependencies as an interactive SVG graph
- **Reports** documentation coverage — flags undocumented public symbols

## Quick start

```bash
git clone https://github.com/toolnook/frontend.git
cd frontend
go build -o frontend ./cmd/frontend

# Scan source files
./frontend scan ./src

# Generate docs site
./frontend build -dir ./src -o ./docs

# Serve locally
./frontend serve -d ./docs -port 8090
```

## Annotation format

### TypeScript

```typescript
/**
 * @doc Adds two numbers
 * @param {number} a - The first number
 * @param {number} b - The second number
 * @returns {number} The sum
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
# @returns (int) The sum
def add(a, b):
    return a + b
```

## CLI commands

```bash
frontend scan ./src           # Parse and output summary
frontend build -dir ./src -o ./docs   # Generate static site
frontend serve -d ./docs       # Serve locally
```

## Architecture

```
Source files → Parser → Documents → Generator → Static HTML site
                ↓                        ↓
           Coverage %              Search index (JSON)
                ↓                        ↓
          Coverage badge            Dependency graph (SVG)
```

| Package | Responsibility |
|---------|---------------|
| `internal/types` | Core data structures |
| `internal/parser` | TS + Python annotation parsing |
| `internal/generator` | HTML/CSS/JS site generation |
| `internal/search` | Search index builder |
| `internal/graph` | Dependency graph + SVG |
| `internal/coverage` | Documentation coverage |
| `internal/server` | Local dev server |

## Running tests

```bash
go test ./... -v
```

## Roadmap

- [ ] Go, Rust, Java parser support
- [ ] Markdown page support
- [ ] Custom themes
- [ ] Export to OpenAPI/JSON Schema

## License

MIT — see [LICENSE](LICENSE).
