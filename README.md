![Toolnook](docs/assets/toolnook-logo.svg)

# Toolnook Frontend (AssetScope)

**The issuer dashboard for AssetScope.** Upload normalized on-chain, ledger, and reserve snapshots and review the comparison report in one browser workflow.

[Read the beginner-friendly guide](https://toolnook.github.io/frontend/) · [Browse all Toolnook projects](https://github.com/toolnook)

## Run the dashboard

You need Go 1.23 or newer. The dashboard is served locally and calls the Backend gateway at the same origin:

```bash
go run ./cmd/assetscope --addr 127.0.0.1:8090
```

For a complete local setup, including the Contracts API and Backend gateway, follow the [integrated setup guide](https://github.com/toolnook/backend#run-the-integrated-assetscope-app-locally). Open `http://127.0.0.1:8080` through the gateway.

The form accepts three CSV snapshots with columns `asset_code,asset_issuer,amount_in_issued_asset_units,as_of,source_ref`. Amounts must already be valued in the issued asset's units. The dashboard sends those files to the Contracts API through the Backend gateway.

## Legacy DocForge utility

The original Soroban documentation generator remains available in `cmd/docforge` while AssetScope becomes the primary product. It scans annotated source and generates static contract documentation:

```bash
go run ./cmd/docforge scan -dir /path/to/your-contract
go run ./cmd/docforge build -dir /path/to/your-contract -o ./generated-docs -name MyContract
go run ./cmd/docforge serve -d ./generated-docs
```

## Explore and contribute

- [Project guide and annotation examples](https://toolnook.github.io/frontend/)
- [Source code](https://github.com/toolnook/frontend)
- [Report a problem](https://github.com/toolnook/frontend/issues)
- [Contribution guide](CONTRIBUTING.md)

Run the test suite with `go test ./...`. This project is open source under the [MIT License](LICENSE).
