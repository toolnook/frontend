package parser

import "errors"

// ErrUnsupportedLanguage is returned when no parser is registered for
// a file extension.
var ErrUnsupportedLanguage = errors.New("unsupported language: no parser registered for this file type")
