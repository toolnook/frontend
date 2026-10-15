package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/forgekit/docforge/internal/types"
)

// PythonParser parses .py files for annotated documentation.
type PythonParser struct{}

func (p *PythonParser) Language() string    { return "python" }
func (p *PythonParser) Extensions() []string { return []string{".py"} }

// ParseFile reads Python source code and extracts documented symbols.
func (p *PythonParser) ParseFile(path string, content []byte) (types.Module, error) {
	lines := strings.Split(string(content), "\n")
	moduleName := filepath.Base(path)

	module := types.Module{
		Name:     moduleName,
		Path:     path,
		Language: "python",
		Imports:  extractPythonImports(lines),
	}

	// Regex patterns for Python function/class definitions
	funcRe := regexp.MustCompile(`^\s*def\s+(\w+)\s*\(`)
	asyncFuncRe := regexp.MustCompile(`^\s*async\s+def\s+(\w+)\s*\(`)
	classRe := regexp.MustCompile(`^\s*class\s+(\w+)`)

	for i, line := range lines {
		// Try async function
		if m := asyncFuncRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindFunction, line, moduleName)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try function
		if m := funcRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindFunction, line, moduleName)
			module.Documents = append(module.Documents, doc)
			continue
		}

		// Try class
		if m := classRe.FindStringSubmatch(line); m != nil {
			doc := p.parseSymbol(lines, i, m[1], types.KindClass, line, moduleName)
			module.Documents = append(module.Documents, doc)
			continue
		}
	}

	return module, nil
}

// parseSymbol looks upward from the declaration line for Python comments
// (either """docstrings""" or # @doc annotations) and extracts metadata.
func (p *PythonParser) parseSymbol(lines []string, lineNum int, name string, kind types.SymbolKind, signature string, moduleName string) types.Document {
	doc := types.Document{
		Name:         name,
		Kind:         kind,
		Module:       moduleName,
		Signature:    strings.TrimSpace(signature),
		Line:         lineNum + 1,
		Visibility:   types.VisibilityPublic,
		IsDocumented: false,
	}

	// Check if the symbol is private (prefixed with _)
	if strings.HasPrefix(name, "_") {
		doc.Visibility = types.VisibilityPrivate
	}

	// Look for preceding # comment lines with annotations
	commentLines := p.extractCommentBlock(lines, lineNum)
	if len(commentLines) > 0 {
		p.parseAnnotations(commentLines, &doc)
	}

	// Also check for docstring after the def/class
	if !doc.IsDocumented {
		docstringLines := p.extractDocstring(lines, lineNum)
		if len(docstringLines) > 0 {
			p.parseAnnotations(docstringLines, &doc)
		}
	}

	return doc
}

// extractCommentBlock looks upward from the declaration for # comment
// lines containing @doc, @param, @returns, @example annotations.
func (p *PythonParser) extractCommentBlock(lines []string, lineNum int) []string {
	var commentLines []string

	for i := lineNum - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])

		if trimmed == "" {
			continue
		}

		// Python single-line comment
		if strings.HasPrefix(trimmed, "#") {
			inner := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			commentLines = append([]string{inner}, commentLines...)
			continue
		}

		break
	}

	return commentLines
}

// extractDocstring looks for a triple-quoted string immediately after
// the function/class definition line.
func (p *PythonParser) extractDocstring(lines []string, lineNum int) []string {
	var docLines []string

	// Look for the first non-empty line after the def/class
	startIdx := -1
	for i := lineNum + 1; i < len(lines) && i < lineNum+10; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, `"""`) || strings.HasPrefix(trimmed, `'''`) {
			startIdx = i
		}
		break
	}

	if startIdx == -1 {
		return docLines
	}

	quoteChar := `"""`
	if strings.HasPrefix(strings.TrimSpace(lines[startIdx]), `'''`) {
		quoteChar = `'''`
	}

	// Check if the docstring is on a single line
	trimmed := strings.TrimSpace(lines[startIdx])
	if strings.Count(trimmed, quoteChar) >= 2 {
		inner := strings.TrimSuffix(strings.TrimPrefix(trimmed, quoteChar), quoteChar)
		if inner != "" {
			docLines = append(docLines, strings.TrimSpace(inner))
		}
		return docLines
	}

	// Multi-line docstring
	for i := startIdx + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, quoteChar) {
			break
		}
		if line != "" {
			docLines = append(docLines, line)
		}
	}

	return docLines
}

// parseAnnotations extracts @doc, @param, @returns, and @example from
// comment or docstring lines.
func (p *PythonParser) parseAnnotations(commentLines []string, doc *types.Document) {
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
			param := parsePythonParam(strings.TrimSpace(strings.TrimPrefix(line, "@param")))
			if param.Name != "" {
				doc.Parameters = append(doc.Parameters, param)
			}
		} else if strings.HasPrefix(line, "@returns") {
			doc.IsDocumented = true
			ret := strings.TrimSpace(strings.TrimPrefix(line, "@returns"))
			// Format: "(type) description" or just description
			if strings.HasPrefix(ret, "(") {
				end := strings.Index(ret, ")")
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

// parsePythonParam parses a Python @param annotation.
// Format: "name (type) description" or "name description"
func parsePythonParam(s string) types.Parameter {
	param := types.Parameter{}

	// Try: name (type) description
	parenRe := regexp.MustCompile(`^(\w+)\s+\(([^)]+)\)\s*(.*)$`)
	if m := parenRe.FindStringSubmatch(s); m != nil {
		param.Name = m[1]
		param.Type = m[2]
		param.Description = strings.TrimSpace(m[3])
		return param
	}

	// Try: name description (no type)
	parts := strings.SplitN(s, " ", 2)
	if len(parts) > 0 {
		param.Name = parts[0]
	}
	if len(parts) > 1 {
		param.Description = strings.TrimSpace(parts[1])
	}

	return param
}

// extractPythonImports finds import statements in Python source.
func extractPythonImports(lines []string) []string {
	importRe := regexp.MustCompile(`^\s*(?:from\s+(\S+)\s+import|import\s+(\S+))`)
	var imports []string
	seen := make(map[string]bool)

	for _, line := range lines {
		if m := importRe.FindStringSubmatch(line); m != nil {
			imp := m[1]
			if imp == "" {
				imp = m[2]
			}
			imp = strings.TrimSpace(imp)
			if imp != "" && !seen[imp] {
				imports = append(imports, imp)
				seen[imp] = true
			}
		}
	}

	return imports
}
