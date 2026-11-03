# Contributing to DocForge

Thanks for your interest in improving DocForge! This guide covers the development workflow.

## Getting started

```bash
git clone https://github.com/crafstack/docforge.git
cd docforge
go build -o docforge ./cmd/docforge
go test ./... -v
```

## Development workflow

1. Create a branch for your change:
   ```bash
   git checkout -b fix/your-bug-or-feature
   ```

2. Make your changes. Follow the existing code style.

3. Add or update tests. All new features and bug fixes should include test coverage:
   ```bash
   go test ./... -v
   ```

4. Write clear, descriptive commit messages:
   ```
   Add Go parser for struct annotations
   Fix import resolution for relative paths
   ```

5. Push and open a pull request against `main`.

## Code style

- Follow standard Go conventions (gofmt, go vet).
- Every exported function should have a docstring comment.
- Keep packages focused — one responsibility per package.
- Prefer the standard library; avoid external dependencies.

## Testing guidelines

- Test files go alongside source files (e.g. `parser_test.go` tests `parser.go`).
- Use table-driven tests where appropriate.
- Test names should start with `Test` and describe the behaviour.
- Use `t.TempDir()` for file-based tests.

## Adding a new language parser

1. Create a new file in `internal/parser/` (e.g. `rust.go`).
2. Implement the `LanguageParser` interface:
   - `ParseFile(path string, content []byte) (types.Module, error)`
   - `Language() string`
   - `Extensions() []string`
3. Register it in `NewRegistry()`.
4. Add tests in the parser test file.

## Reporting issues

- Use the bug report template for bugs.
- Use the feature request template for new features.
- Include the DocForge version, Go version, and a minimal reproduction case.

## Pull request checklist

- [ ] Tests pass (`go test ./... -v`)
- [ ] Code passes `go vet`
- [ ] No external dependencies added (keep it stdlib-only)
- [ ] Commit messages are descriptive
- [ ] Documentation updated if behaviour changes
