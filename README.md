# Toolnook Frontend — Soroban Contract Documentation

A documentation generator for Soroban smart contracts on Stellar. Parses Rust source files annotated with `@doc`, `@param`, `@returns`, and `@example` tags and generates a static, searchable documentation site with dependency graphs and coverage indicators.

Part of the [Toolnook](https://github.com/toolnook) Stellar developer tooling platform:

| Repo | Role |
|------|------|
| [toolnook/backend](https://github.com/toolnook/backend) | Horizon gateway with caching, rate limiting, circuit breaker |
| **toolnook/frontend** | Soroban contract documentation generator |
| [toolnook/contracts](https://github.com/toolnook/contracts) | XLM payment reconciliation and double-entry ledger |

## Why this helps Stellar

Soroban smart contracts on Stellar are written in Rust using the `soroban-sdk`. There's no standard tool for generating user-friendly documentation from contract source. Toolnook Frontend:

- **Parses** Soroban `#[contractimpl]`, `#[contracttype]`, and `#[contractevent]` attributes
- **Extracts** contract method signatures, parameter types, return types, and events
- **Generates** a static HTML site with search, styled for reading contract APIs
- **Visualizes** dependencies between contract modules as an interactive SVG graph
- **Reports** documentation coverage — flags undocumented public contract methods

## Quick start

```bash
git clone https://github.com/toolnook/frontend.git
cd frontend
go build -o frontend ./cmd/frontend

# Scan Soroban contract source
./frontend scan ./contracts

# Generate docs site
./frontend build -dir ./contracts -o ./docs

# Serve locally
./frontend serve -d ./docs -port 8090
```

## Annotation format

```rust
use soroban_sdk::{contractimpl, Address, Env};

/// @doc Transfers tokens between accounts
/// @param from: Address The sender's account
/// @param to: Address The receiver's account
/// @param amount: i128 The amount to transfer in base units
/// @returns Result<(), Error> Ok on success, Error on failure
/// @example
/// token.transfer(env, &from, &to, 1_000_000)
#[contractimpl]
impl TokenContract {
    /// @doc Returns the total supply
    /// @returns i128 The total token supply
    pub fn total_supply(env: Env) -> i128 {
        0
    }

    /// @doc Transfers tokens
    /// @param from: Address The sender
    /// @param to: Address The receiver
    /// @param amount: i128 The amount
    pub fn transfer(env: Env, from: Address, to: Address, amount: i128) -> Result<(), Error> {
        Ok(())
    }
}

/// @doc Error types for the token contract
#[contracttype]
pub enum Error {
    InsufficientBalance = 1,
    Unauthorized = 2,
}

/// @doc Emitted on successful transfer
#[contractevent]
pub struct TransferEvent {
    pub from: Address,
    pub to: Address,
    pub amount: i128,
}
```

## What it recognizes

| Soroban attribute | Parsed as | Tags |
|-------------------|-----------|------|
| `#[contractimpl]` | Contract entry points (functions inside impl block) | `contract` |
| `#[contracttype]` | Custom types (structs, enums) | `soroban-type` |
| `#[contractevent]` | Event definitions | `event` |

Also parses TypeScript and Python for non-contract code documentation.

## CLI commands

```bash
frontend scan ./contracts           # Parse and output summary
frontend build -dir ./contracts -o ./docs   # Generate static site
frontend serve -d ./docs             # Serve locally
```

## Running tests

```bash
go test ./... -v
```

## License

MIT — see [LICENSE](LICENSE).
