package jobs

import (
	"fmt"
	"strings"
)

// SplitCommandLine parses a command string into individual tokens/arguments,
// respecting single quotes, double quotes, and backslash escape sequences.
func SplitCommandLine(commandLine string) ([]string, error) {
	trimmed := strings.TrimSpace(commandLine)
	if trimmed == "" {
		return nil, nil
	}

	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	hasToken := false

	for i := 0; i < len(commandLine); i++ {
		b := commandLine[i]

		if escaped {
			current.WriteByte(b)
			hasToken = true
			escaped = false
			continue
		}

		if b == '\\' && !inSingleQuote {
			escaped = true
			continue
		}

		if inSingleQuote {
			if b == '\'' {
				inSingleQuote = false
			} else {
				current.WriteByte(b)
				hasToken = true
			}
			continue
		}

		if inDoubleQuote {
			if b == '"' {
				inDoubleQuote = false
			} else {
				current.WriteByte(b)
				hasToken = true
			}
			continue
		}

		switch b {
		case '\'':
			inSingleQuote = true
			hasToken = true
		case '"':
			inDoubleQuote = true
			hasToken = true
		case ' ', '\t', '\n', '\r':
			if hasToken {
				tokens = append(tokens, current.String())
				current.Reset()
				hasToken = false
			}
		default:
			current.WriteByte(b)
			hasToken = true
		}
	}

	if inSingleQuote || inDoubleQuote {
		return nil, fmt.Errorf("unclosed quote in command string")
	}

	if escaped {
		current.WriteByte('\\')
		hasToken = true
	}

	if hasToken {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}
