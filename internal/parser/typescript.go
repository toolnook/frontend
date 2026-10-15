package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/forgekit/docforge/internal/types"
)

// TypeScriptParser parses .ts and .tsx files for annotated documentation.
type TypeScriptParser struct{}

func (p *TypeScriptParser) Language() string     { return "typescript" }
func (p *TypeScriptParser) Extensions() []string  { return []string{".ts", ".tsx"} }

// ParseFile reads TypeScript source code and extracts documented symbols.
func (p *TypeScriptParser) ParseFile(path string, content []byte) (types.Module, error) {
	lines := strings.Split(string(content), "\n")
	moduleName := filepath.Base(path)

	module := types.Module{
		Name:     moduleName,
		Path:     path,
		Language: "typescript",
		Imports:  extractTSImports(lines),
	}

	// Regex patterns for finding function/class/type declarations
	funcRe := regexp.MustCompile(`^\s*(?:export\s+)?(?:async\s+)?function\s+(\w+)`)
	constFuncRe := regexp.MustCompile(`^\s*(?:export\s+)?const\s+(\w+)\s*=\s*(?:async\s*)?\(`)
	classRe := regexp.MustCompile(`^\s*(?:export\s+)?class\s+(\w+)`)
	typeRe := regexp.MustCompile(`^\s*(?:export\s+)?type\s+(\w+)`)
	interfaceRe := regexp.MustCompile(`^\s*(?:export\s+)?interface\s+(\w+)`)

	for i, line := range lines {
		// Try function declaration
		if m := funcRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindFunction, line)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try const arrow function
		if m := constFuncRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindFunction, line)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try class declaration
		if m := classRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindClass, line)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try type declaration
		if m := typeRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindType, line)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try interface declaration
		if m := interfaceRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindType, line)
			module.Documents = append(module.Documents, doc)
			continue
		}
	}

	return module, nil
}

// parseSymbol looks upward from the declaration line for a JSDoc comment
// block and extracts @doc, @param, @returns, and @example annotations.
func (p *TypeScriptParser) parseSymbol(lines []string, lineNum int, name string, kind types.SymbolKind, signature string) types.Document {
	doc := types.Document{
		Name:        name,
		Kind:        kind,
		Module:      filepath.Base(lines[0]), // placeholder, set by caller
		Signature:   strings.TrimSpace(signature),
		Line:        lineNum + 1,
		Visibility:  types.VisibilityPublic, // default; will override if private
		IsDocumented: false,
	}

	// Check if the symbol is private (prefixed with _)
	if strings.HasPrefix(name, "_") {
		doc.Visibility = types.VisibilityPrivate
	}

	// Check if it's an export (public) or not
	if !strings.Contains(signature, "export") && !strings.Contains(signature, "Export") {
		// Non-exported declarations in TS are module-scoped, but we treat
		// them as public for documentation purposes unless prefixed with _
	}

	// Look for preceding comment block
	commentLines := p.extractCommentBlock(lines, lineNum)
	if len(commentLines) > 0 {
		p.parseAnnotations(commentLines, &doc)
	}

	return doc
}

// extractCommentBlock looks upward from the declaration line for a
// JSDoc comment block (/** ... */) and returns the inner lines.
func (p *TypeScriptParser) extractCommentBlock(lines []string, lineNum int) []string {
	var commentLines []string

	// Walk backwards from the line before the declaration
	for i := lineNum - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])

		// Empty lines are skipped
		if trimmed == "" {
			continue
		}

		// End of comment block (walking backward, this is the opening)
		if strings.Contains(trimmed, "/**") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "/**"))
			inner = strings.TrimSuffix(inner, "*/")
			if inner != "" && inner != "*" {
				commentLines = append([]string{inner}, commentLines...)
			}
			break
		}

		// Inside a comment block
		if strings.HasPrefix(trimmed, "*") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "*"))
			if strings.HasSuffix(trimmed, "*/") {
				inner = strings.TrimSuffix(inner, "*/")
				inner = strings.TrimSpace(inner)
			}
			if inner != "" {
				commentLines = append([]string{inner}, commentLines...)
			}
			continue
		}

		// Line comment (//)
		if strings.HasPrefix(trimmed, "//") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
			commentLines = append([]string{inner}, commentLines...)
			continue
		}

		// Non-comment line means the block ended
		break
	}

	return commentLines
}

// parseAnnotations extracts @doc, @param, @returns, and @example from
// comment lines and populates the Document.
func (p *TypeScriptParser) parseAnnotations(commentLines []string, doc *types.Document) {
	if len(commentLines) == 0 {
		return
	}

	var descParts []string
	var currentExample *types.Example
	inExample := false

	for _, line := range commentLines {
		if strings.HasPrefix(line, "@doc") {
			desc := strings.TrimSpace(strings.TrimPrefix(line, "@doc"))
			if desc != "" {
				descParts = append(descParts, desc)
			}
			doc.IsDocumented = true
		} else if strings.HasPrefix(line, "@param") {
			doc.IsDocumented = true
			param := parseParam(strings.TrimSpace(strings.TrimPrefix(line, "@param")))
			if param.Name != "" {
				doc.Parameters = append(doc.Parameters, param)
			}
		} else if strings.HasPrefix(line, "@returns") {
			doc.IsDocumented = true
			ret := strings.TrimSpace(strings.TrimPrefix(line, "@returns"))
			// Split type and description: "{number} The result"
			if strings.HasPrefix(ret, "{") {
				end := strings.Index(ret, "}")
				if end > 0 {
					doc.Returns = ret[1:end]
					doc.ReturnsDesc = strings.TrimSpace(ret[end+1:])
				} else {
					doc.ReturnsDesc = ret
				}
			} else {
				doc.ReturnsDesc = ret
			}
		} else if strings.HasPrefix(line, "@example") {
			doc.IsDocumented = true
			desc := strings.TrimSpace(strings.TrimPrefix(line, "@example"))
			if desc != "" {
				currentExample = &types.Example{Description: desc}
			} else {
				currentExample = &types.Example{}
			}
			inExample = true
			if currentExample != nil {
				doc.Examples = append(doc.Examples, *currentExample)
				currentExample = &doc.Examples[len(doc.Examples)-1]
			}
		} else if inExample {
			// Continuation of example code
			if currentExample != nil {
				currentExample.Code += line + "\n"
			}
		} else if !strings.HasPrefix(line, "@") {
			// Regular description text
			if line != "" && line != "*" {
				descParts = append(descParts, line)
				doc.IsDocumented = true
			}
		}

		// Reset example mode on new annotations
		if inExample && strings.HasPrefix(line, "@") && !strings.HasPrefix(line, "@example") {
			inExample = false
		}
	}

	doc.Description = strings.Join(descParts, " ")
}

// parseParam parses a @param annotation value.
// Format: "{type} name - description" or "name type description"
func parseParam(s string) types.Parameter {
	param := types.Parameter{}

	// Try {type} name - description
	if strings.HasPrefix(s, "{") {
		end := strings.Index(s, "}")
		if end > 0 {
			param.Type = s[1:end]
			rest := strings.TrimSpace(s[end+1:])
			// name - description
			parts := strings.SplitN(rest, " - ", 2)
			if len(parts) > 0 {
				param.Name = strings.TrimSpace(parts[0])
			}
			if len(parts) > 1 {
				param.Description = strings.TrimSpace(parts[1])
			}
		}
		return param
	}

	// Try name - description (no type)
	parts := strings.SplitN(s, " - ", 2)
	if len(parts) > 0 {
		param.Name = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		param.Description = strings.TrimSpace(parts[1])
	}

	return param
}

// extractTSImports finds import statements in TypeScript source.
func extractTSImports(lines []string) []string {
	importRe := regexp.MustCompile(`^\s*import\s+.*\s+from\s+['"]([^'"]+)['"]`)
	var imports []string
	seen := make(map[string]bool)

	for _, line := range lines {
		if m := importRe.FindStringSubmatch(line); m != nil {
			imp := m[1]
			if !seen[imp] {
				imports = append(imports, imp)
				seen[imp] = true
			}
		}
	}

	return imports
}
