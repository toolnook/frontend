![Toolnook](docs/assets/toolnook-logo.svg)

# Toolnook Frontend (DocForge)

**Readable documentation for Soroban smart contracts.** DocForge scans annotated source code and turns it into a searchable website, so people can understand a contract's functions, types, and events without reading every source file first.

[Read the beginner-friendly guide](https://toolnook.github.io/frontend/) · [Browse all Toolnook projects](https://github.com/toolnook)

## What does it do?

Soroban is Stellar's smart-contract platform. A Soroban contract may expose functions, define custom types, and emit events. DocForge reads Rust source comments and Soroban attributes, then creates a static HTML reference site that you can share or host.

DocForge describes source code; it does not deploy contracts or verify what a live contract does.

## Try it locally

You need Go 1.23 or newer. From this repository, scan a directory containing your Soroban Rust source:

```bash
git clone https://github.com/toolnook/frontend.git
cd frontend
go run ./cmd/docforge scan -dir /path/to/your-contract
```

Generate and preview a site:

```bash
go run ./cmd/docforge build -dir /path/to/your-contract -o ./generated-docs -name MyContract
go run ./cmd/docforge serve -d ./generated-docs
```

Open the local address printed by the server (usually `http://localhost:8090`). Add tags such as `@doc`, `@param`, `@returns`, and `@example` to source comments to give readers useful explanations.

## Concepts in plain language

| Term | Meaning |
| --- | --- |
| Soroban | Stellar's platform for smart contracts. |
| Contract function | An operation a user or another contract can call. |
| Type | The shape of data used by a contract, such as a record or enum. |
| Event | A message a contract emits to describe something that happened. |
| Static site | Plain HTML files that can be hosted without running an application server. |

## Explore and contribute

- [Project guide and annotation examples](https://toolnook.github.io/frontend/)
- [Source code](https://github.com/toolnook/frontend)
- [Report a problem](https://github.com/toolnook/frontend/issues)
- [Contribution guide](CONTRIBUTING.md)

Run the test suite with `go test ./...`. This project is open source under the [MIT License](LICENSE).
