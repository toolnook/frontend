---
name: soroban-docforge
description: Use when extending, debugging, or documenting DocForge's source parsing and Soroban contract documentation generation.
---

# Soroban DocForge

Help maintain the Go tool that turns annotated Soroban source into static API documentation.

- Read the parser registry, language parser, generator templates, and corresponding tests before changing behavior.
- Preserve support for existing source languages while extending Soroban Rust parsing for contract implementations, types, and events.
- Keep generated documentation deterministic and safe to publish as static files; do not imply that generated docs verify or deploy a live contract.
- Use real Soroban/Rust terminology, and define contract function, type, and event in beginner-facing documentation.
- Add focused parser or generator tests for behavior changes and run `go test ./...` and `go vet ./...`.
- Do not include credentials or add contributor/author entries that were not requested.
