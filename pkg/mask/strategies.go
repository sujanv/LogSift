package mask

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Strategy defines how sensitive substrings are sanitized.
type Strategy interface {
	Apply(val, category string) string
}

// RedactStrategy replaces data with [REDACTED_CATEGORY].
type RedactStrategy struct{}

func (RedactStrategy) Apply(val, category string) string {
	return "[REDACTED_" + category + "]"
}

// MaskStrategy masks middle characters with asterisks (e.g. j***@domain.com).
type MaskStrategy struct{}

func (MaskStrategy) Apply(val, category string) string {
	if len(val) <= 4 {
		return strings.Repeat("*", len(val))
	}
	return val[:2] + strings.Repeat("*", len(val)-4) + val[len(val)-2:]
}

// HashStrategy converts values into truncated cryptographic hashes.
type HashStrategy struct{}

func (HashStrategy) Apply(val, category string) string {
	h := sha256.Sum256([]byte(val))
	return fmt.Sprintf("[HASH_%s_%x]", category, h[:4])
}

// ParseStrategy returns the named strategy implementation.
func ParseStrategy(name string) Strategy {
	switch strings.ToLower(name) {
	case "mask":
		return MaskStrategy{}
	case "hash":
		return HashStrategy{}
	default:
		return RedactStrategy{}
	}
}
