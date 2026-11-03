package parser

import (
	"testing"

	"github.com/toolnook/docforge/internal/types"
)

func TestRegistryDetectLanguage(t *testing.T) {
	r := NewRegistry()
	if r.DetectLanguage("foo.ts") != "typescript" {
		t.Error("expected typescript for .ts")
	}
	if r.DetectLanguage("foo.py") != "python" {
		t.Error("expected python for .py")
	}
	if r.DetectLanguage("foo.rb") != "" {
		t.Error("expected empty for unsupported extension")
	}
}

func TestRegistrySupportedExtensions(t *testing.T) {
	r := NewRegistry()
	exts := r.SupportedExtensions()
	if len(exts) < 2 {
		t.Errorf("expected at least 2 extensions, got %d", len(exts))
	}
}

func TestParseTypeScriptFunction(t *testing.T) {
	p := &TypeScriptParser{}
	module, err := p.ParseFile("test.ts", []byte(`/**
 * @doc Adds two numbers together
 * @param {number} a - The first number
 * @param {number} b - The second number
 * @returns {number} The sum
 * @example
 * add(1, 2) // returns 3
 */
export function add(a: number, b: number): number {
  return a + b;
}

export function subtract(a: number, b: number): number {
  return a - b;
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(module.Documents))
	}

	addDoc := module.Documents[0]
	if addDoc.Name != "add" {
		t.Errorf("expected name 'add', got '%s'", addDoc.Name)
	}
	if !addDoc.IsDocumented {
		t.Error("expected add to be documented")
	}
	if addDoc.Description == "" {
		t.Error("expected non-empty description")
	}
	if len(addDoc.Parameters) != 2 {
		t.Errorf("expected 2 params, got %d", len(addDoc.Parameters))
	}
	if addDoc.Parameters[0].Name != "a" {
		t.Errorf("expected first param name 'a', got '%s'", addDoc.Parameters[0].Name)
	}
	if addDoc.Returns != "number" {
		t.Errorf("expected returns 'number', got '%s'", addDoc.Returns)
	}
	if len(addDoc.Examples) != 1 {
		t.Errorf("expected 1 example, got %d", len(addDoc.Examples))
	}

	// subtract should be undocumented
	subDoc := module.Documents[1]
	if subDoc.IsDocumented {
		t.Error("expected subtract to be undocumented")
	}
}

func TestParseTypeScriptClass(t *testing.T) {
	p := &TypeScriptParser{}
	module, err := p.ParseFile("test.ts", []byte(`/**
 * @doc A calculator class
 */
export class Calculator {
  add(a, b) { return a + b; }
}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(module.Documents))
	}
	if module.Documents[0].Kind != types.KindClass {
		t.Errorf("expected kind class, got %s", module.Documents[0].Kind)
	}
}

func TestParsePythonFunction(t *testing.T) {
	p := &PythonParser{}
	module, err := p.ParseFile("test.py", []byte(`# @doc Adds two numbers
# @param a (int) The first number
# @param b (int) The second number
# @returns (int) The sum
# @example
# add(1, 2)  # returns 3
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(module.Documents))
	}

	addDoc := module.Documents[0]
	if addDoc.Name != "add" {
		t.Errorf("expected name 'add', got '%s'", addDoc.Name)
	}
	if !addDoc.IsDocumented {
		t.Error("expected add to be documented")
	}
	if len(addDoc.Parameters) != 2 {
		t.Errorf("expected 2 params, got %d", len(addDoc.Parameters))
	}
	if addDoc.Returns != "int" {
		t.Errorf("expected returns 'int', got '%s'", addDoc.Returns)
	}

	subDoc := module.Documents[1]
	if subDoc.IsDocumented {
		t.Error("expected subtract to be undocumented")
	}
}

func TestParsePythonClass(t *testing.T) {
	p := &PythonParser{}
	module, err := p.ParseFile("test.py", []byte(`# @doc A calculator class
class Calculator:
    pass
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(module.Documents))
	}
	if module.Documents[0].Kind != types.KindClass {
		t.Errorf("expected kind class, got %s", module.Documents[0].Kind)
	}
	if !module.Documents[0].IsDocumented {
		t.Error("expected Calculator to be documented")
	}
}

func TestParseTypeScriptImports(t *testing.T) {
	p := &TypeScriptParser{}
	module, err := p.ParseFile("test.ts", []byte(`import { foo } from './utils';
import { bar } from '../lib/bar';
export function main() {}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Imports) != 2 {
		t.Errorf("expected 2 imports, got %d", len(module.Imports))
	}
	if module.Imports[0] != "./utils" {
		t.Errorf("expected first import './utils', got '%s'", module.Imports[0])
	}
}

func TestParsePythonImports(t *testing.T) {
	p := &PythonParser{}
	module, err := p.ParseFile("test.py", []byte(`from utils import helper
import os
import sys

def main():
    pass
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
}
	if len(module.Imports) != 3 {
		t.Errorf("expected 3 imports, got %d", len(module.Imports))
	}
}

func TestParseUnsupportedFileReturnsError(t *testing.T) {
	r := NewRegistry()
	_, err := r.ParseFile("test.unknown")
	if err != ErrUnsupportedLanguage {
		t.Errorf("expected ErrUnsupportedLanguage, got %v", err)
	}
}

func TestPrivateSymbolDetection(t *testing.T) {
	p := &TypeScriptParser{}
	module, err := p.ParseFile("test.ts", []byte(`function _private() {}
export function public() {}
`))
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(module.Documents) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(module.Documents))
	}
	if module.Documents[0].Visibility != types.VisibilityPrivate {
		t.Error("expected _private to be private")
	}
	if module.Documents[1].Visibility != types.VisibilityPublic {
		t.Error("expected public to be public")
	}
}
