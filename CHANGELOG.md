# Changelog

All notable changes to DocForge are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] — 2026-11-03

### Added
- Rust/Soroban parser for Stellar smart contract source files
- Parses `#[contractimpl]`, `#[contracttype]`, and `#[contractevent]` attributes
- Extracts contract method signatures, parameters, return types, and events
- Annotation parsing from `@doc`, `@param`, `@returns`, `@example` comment tags
- TypeScript parser with JSDoc block comment support
- Python parser with comment and docstring support
- Pluggable language parser registry with file extension detection
- Static HTML/CSS/JS site generator with dark theme
- Client-side full-text search with JSON index
- Dependency graph builder with interactive SVG visualization
- Documentation coverage calculator with missing-symbols report
- CLI commands: scan, build, serve, version
- Local HTTP server for development preview
- Full test suite including Soroban contract parsing tests
- CI workflow for GitHub Actions
- Issue templates and PR template

[0.1.0]: https://github.com/toolnook/frontend/releases/tag/v0.1.0
