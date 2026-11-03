package parser

import (
	"testing"

	"github.com/toolnook/frontend/internal/types"
)

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func TestParseSorobanContract(t *testing.T) {
	p := &RustParser{}
	module, err := p.ParseFile("contract.rs", []byte(`use soroban_sdk::contractimpl;

/// @doc Transfers tokens between accounts
/// @param from: Address The sender
/// @param to: Address The receiver
/// @param amount: i128 The amount to transfer
/// @returns Result<(), Error> Success or error
/// @example
/// token.transfer(&from, &to, 100)
#[contractimpl]
impl TokenContract {
    /// @doc Returns the contract version
    pub fn version() -> u32 {
        1
    }

    /// @doc Transfers tokens
    /// @param from: Address The sender
    /// @param to: Address The receiver
    /// @param amount: i128 The amount
    /// @returns Result<(), Error> Success or error
    pub fn transfer(env: Env, from: Address, to: Address, amount: i128) -> Result<(), Error> {
        // implementation
        Ok(())
    }

    /// @doc Gets the balance of an account
    /// @param account: Address The account to query
    /// @returns i128 The balance
    pub fn balance(env: Env, account: Address) -> i128 {
        0
    }
}

/// @doc Error types for the token contract
#[contracttype]
pub enum Error {
    InsufficientBalance = 1,
    Unauthorized = 2,
}

/// @doc Transfer event emitted on successful transfer
#[contractevent]
pub struct TransferEvent {
    pub from: Address,
    pub to: Address,
    pub amount: i128,
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	// Should find: version, transfer, balance (functions) + Error (type) + TransferEvent (type)
	if len(module.Documents) < 4 {
		t.Fatalf("expected at least 4 documents, got %d", len(module.Documents))
	}

	// Check transfer function
	var transfer *types.Document
	for _, d := range module.Documents {
		if d.Name == "transfer" {
			transfer = &d
			break
		}
	}
	if transfer == nil {
		t.Fatal("expected to find 'transfer' function")
	}
	if !transfer.IsDocumented {
		t.Error("expected transfer to be documented")
	}
	if len(transfer.Parameters) != 3 {
		t.Errorf("expected 3 params, got %d", len(transfer.Parameters))
	}
	if transfer.Parameters[0].Name != "from" {
		t.Errorf("expected first param 'from', got '%s'", transfer.Parameters[0].Name)
	}
	if transfer.Parameters[0].Type != "Address" {
		t.Errorf("expected first param type 'Address', got '%s'", transfer.Parameters[0].Type)
	}
	if transfer.Returns != "Result<(), Error>" {
		t.Errorf("expected returns 'Result<(), Error>', got '%s'", transfer.Returns)
	}
	if !contains(transfer.Tags, "contract") {
		t.Error("expected 'contract' tag")
	}
}

func TestParseRustContractType(t *testing.T) {
	p := &RustParser{}
	module, err := p.ParseFile("types.rs", []byte(`/// @doc Error types for the contract
#[contracttype]
pub enum Error {
    InsufficientBalance = 1,
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(module.Documents))
	}
	if module.Documents[0].Name != "Error" {
		t.Errorf("expected 'Error', got '%s'", module.Documents[0].Name)
	}
	if !contains(module.Documents[0].Tags, "soroban-type") {
		t.Error("expected 'soroban-type' tag")
	}
}

func TestParseRustContractEvent(t *testing.T) {
	p := &RustParser{}
	module, err := p.ParseFile("events.rs", []byte(`/// @doc Emitted on transfer
#[contractevent]
pub struct TransferEvent {
    pub from: Address,
    pub to: Address,
    pub amount: i128,
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(module.Documents))
	}
	if !contains(module.Documents[0].Tags, "event") {
		t.Error("expected 'event' tag")
	}
}

func TestParseRustImports(t *testing.T) {
	p := &RustParser{}
	module, err := p.ParseFile("contract.rs", []byte(`use soroban_sdk::{contractimpl, Address, Env};
use crate::types::Error;

pub fn dummy() {}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Imports) < 1 {
		t.Errorf("expected at least 1 import, got %d", len(module.Imports))
	}
}

func TestParseRustUndocumentedFunction(t *testing.T) {
	p := &RustParser{}
	module, err := p.ParseFile("test.rs", []byte(`#[contractimpl]
impl Contract {
    pub fn undocumented_fn() -> u32 {
        42
    }
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(module.Documents))
	}
	if module.Documents[0].IsDocumented {
		t.Error("expected function to be undocumented")
	}
}

func TestRustParserLanguage(t *testing.T) {
	p := &RustParser{}
	if p.Language() != "rust" {
		t.Errorf("expected 'rust', got '%s'", p.Language())
	}
}

func TestRegistryDetectsRust(t *testing.T) {
	r := NewRegistry()
	if r.DetectLanguage("contract.rs") != "rust" {
		t.Error("expected rust for .rs file")
	}
}
