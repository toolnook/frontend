package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/toolnook/frontend/internal/types"
)

// RustParser parses .rs files for Soroban smart contract annotations.
//
// Soroban contracts on Stellar are written in Rust using the soroban-sdk.
// This parser looks for:
//   - #[contractimpl] attribute on impl blocks (contract entry points)
//   - #[contracttype] attribute on enums/structs (contract types)
//   - #[contractevent] attribute (emitted events)
//   - @doc, @param, @returns, @example annotations in preceding comments
type RustParser struct{}

func (p *RustParser) Language() string     { return "rust" }
func (p *RustParser) Extensions() []string  { return []string{".rs"} }

// ParseFile reads Rust/Soroban source and extracts documented symbols.
func (p *RustParser) ParseFile(path string, content []byte) (types.Module, error) {
	lines := strings.Split(string(content), "\n")
	moduleName := filepath.Base(path)

	module := types.Module{
		Name:     moduleName,
		Path:     path,
		Language: "rust",
		Imports:  extractRustImports(lines),
	}

	// Patterns for Soroban contract constructs
	contractImplRe := regexp.MustCompile(`#\[(?:contractimpl|contract_impl)\]`)
	fnRe := regexp.MustCompile(`^\s*(?:pub\s+)?fn\s+(\w+)\s*\(`)
	structRe := regexp.MustCompile(`^\s*#\[contracttype\]`)
	structNameRe := regexp.MustCompile(`^\s*(?:pub\s+)?(?:struct|enum)\s+(\w+)`)
	eventRe := regexp.MustCompile(`#\[(?:contractevent|contract_event)\]`)

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check for contractimpl attribute
		if contractImplRe.MatchString(line) {
			// Find the impl block and all fn definitions within it
			i++
			// Skip to the impl line
			for i < len(lines) && !strings.Contains(lines[i], "impl") {
				i++
			}
			i++ // skip impl line

			// Parse functions inside the impl block
			braceDepth := 1
			for i < len(lines) && braceDepth > 0 {
				trimmed := strings.TrimSpace(lines[i])
				braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")

				if m := fnRe.FindStringSubmatch(lines[i]); m != nil {
					doc := p.parseSymbol(lines, i, m[1], types.KindFunction, strings.TrimSpace(lines[i]), moduleName)
					doc.AddTag("contract")
					module.Documents = append(module.Documents, doc)
				}
				i++
			}
			continue
		}

		// Check for contracttype (struct or enum)
		if structRe.MatchString(line) {
			// Look for the struct/enum definition on the next line(s)
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if m := structNameRe.FindStringSubmatch(lines[j]); m != nil {
					doc := p.parseSymbol(lines, j, m[1], types.KindType, strings.TrimSpace(lines[j]), moduleName)
					doc.AddTag("soroban-type")
					module.Documents = append(module.Documents, doc)
					i = j
					break
				}
			}
			i++
			continue
		}

		// Check for contractevent
		if eventRe.MatchString(line) {
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if m := structNameRe.FindStringSubmatch(lines[j]); m != nil {
					doc := p.parseSymbol(lines, j, m[1], types.KindType, strings.TrimSpace(lines[j]), moduleName)
					doc.AddTag("event")
					module.Documents = append(module.Documents, doc)
					i = j
					break
				}
			}
			i++
			continue
		}

		// Regular pub fn (not in contractimpl) — still document if annotated
		if m := fnRe.FindStringSubmatch(line); m != nil && strings.Contains(line, "pub ") {
			doc := p.parseSymbol(lines, i, m[1], types.KindFunction, strings.TrimSpace(line), moduleName)
			if doc.IsDocumented {
				module.Documents = append(module.Documents, doc)
			}
		}

		i++
	}

	return module, nil
}

// parseSymbol looks upward from the declaration for comment annotations.
func (p *RustParser) parseSymbol(lines []string, lineNum int, name string, kind types.SymbolKind, signature string, moduleName string) types.Document {
	doc := types.Document{
		Name:       name,
		Kind:       kind,
		Module:     moduleName,
		Signature:  signature,
		Line:       lineNum + 1,
		Visibility: types.VisibilityPublic,
	}

	// Private if not pub
	if !strings.Contains(signature, "pub ") {
		doc.Visibility = types.VisibilityPrivate
	}

	// Look for preceding // comments
	commentLines := p.extractCommentBlock(lines, lineNum)
	if len(commentLines) > 0 {
		p.parseAnnotations(commentLines, &doc)
	}

	return doc
}

// extractCommentBlock looks upward for /// or // comment lines.
func (p *RustParser) extractCommentBlock(lines []string, lineNum int) []string {
	var commentLines []string

	for i := lineNum - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])

		if trimmed == "" {
			continue
		}

		// Rust doc comments (///) or regular comments (//)
		if strings.HasPrefix(trimmed, "///") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "///"))
			commentLines = append([]string{inner}, commentLines...)
			continue
		}
		if strings.HasPrefix(trimmed, "//") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
			commentLines = append([]string{inner}, commentLines...)
			continue
		}

		break
	}

	return commentLines
}

// parseAnnotations extracts @doc, @param, @returns, @example from comments.
func (p *RustParser) parseAnnotations(commentLines []string, doc *types.Document) {
	if len(commentLines) == 0 {
		return
	}

	var descParts []string
	inExample := false
	var currentExample *types.Example

	for _, line := range commentLines {
		if strings.HasPrefix(line, "@doc") {
			desc := strings.TrimSpace(strings.TrimPrefix(line, "@doc"))
			if desc != "" {
				descParts = append(descParts, desc)
			}
			doc.IsDocumented = true
		} else if strings.HasPrefix(line, "@param") {
			doc.IsDocumented = true
			param := parseRustParam(strings.TrimSpace(strings.TrimPrefix(line, "@param")))
			if param.Name != "" {
				doc.Parameters = append(doc.Parameters, param)
			}
		} else if strings.HasPrefix(line, "@returns") {
			doc.IsDocumented = true
			ret := strings.TrimSpace(strings.TrimPrefix(line, "@returns"))
			// Rust types can contain spaces inside generics (e.g., "Result<(), Error>")
			// If the type contains "<", find the matching ">" and split after it
			if gtIdx := strings.Index(ret, ">"); gtIdx >= 0 && strings.Contains(ret[:gtIdx], "<") {
				doc.Returns = ret[:gtIdx+1]
				doc.ReturnsDesc = strings.TrimSpace(ret[gtIdx+1:])
			} else {
				// Simple type: "Type description"
				parts := strings.SplitN(ret, " ", 2)
				if len(parts) > 0 {
					doc.Returns = parts[0]
				}
				if len(parts) > 1 {
					doc.ReturnsDesc = parts[1]
				}
			}
		} else if strings.HasPrefix(line, "@example") {
			doc.IsDocumented = true
			desc := strings.TrimSpace(strings.TrimPrefix(line, "@example"))
			currentExample = &types.Example{Description: desc}
			inExample = true
			doc.Examples = append(doc.Examples, *currentExample)
			currentExample = &doc.Examples[len(doc.Examples)-1]
		} else if inExample {
			if currentExample != nil {
				currentExample.Code += line + "\n"
			}
		} else if !strings.HasPrefix(line, "@") {
			if line != "" {
				descParts = append(descParts, line)
				doc.IsDocumented = true
			}
		}

		if inExample && strings.HasPrefix(line, "@") && !strings.HasPrefix(line, "@example") {
			inExample = false
		}
	}

	doc.Description = strings.Join(descParts, " ")
}

// parseRustParam parses a @param annotation for Rust.
// Format: "name: Type description" or "name Type description"
func parseRustParam(s string) types.Parameter {
	param := types.Parameter{}

	// Try: name: Type description
	if idx := strings.Index(s, ":"); idx > 0 {
		param.Name = strings.TrimSpace(s[:idx])
		rest := strings.TrimSpace(s[idx+1:])
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) > 0 {
			param.Type = parts[0]
		}
		if len(parts) > 1 {
			param.Description = parts[1]
		}
		return param
	}

	// Try: name Type description
	parts := strings.SplitN(s, " ", 3)
	if len(parts) > 0 {
		param.Name = parts[0]
	}
	if len(parts) > 1 {
		param.Type = parts[1]
	}
	if len(parts) > 2 {
		param.Description = parts[2]
	}

	return param
}

// extractRustImports finds use statements in Rust source.
func extractRustImports(lines []string) []string {
	useRe := regexp.MustCompile(`^\s*use\s+([^;]+);`)
	var imports []string
	seen := make(map[string]bool)

	for _, line := range lines {
		if m := useRe.FindStringSubmatch(line); m != nil {
			imp := strings.TrimSpace(m[1])
			if !seen[imp] {
				imports = append(imports, imp)
				seen[imp] = true
			}
		}
	}

	return imports
}
