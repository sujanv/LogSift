package pseudonym

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
)

type Pseudonymizer struct {
	mu   sync.RWMutex
	salt []byte
}

func NewPseudonymizer(salt string) *Pseudonymizer {
	if salt == "" {
		salt = "logsift-default-privacy-salt-2026"
	}
	return &Pseudonymizer{
		salt: []byte(salt),
	}
}

func (p *Pseudonymizer) Tokenize(raw string) string {
	if raw == "" {
		return ""
	}
	p.mu.RLock()
	defer p.mu.RUnlock()

	mac := hmac.New(sha256.New, p.salt)
	mac.Write([]byte(strings.ToLower(strings.TrimSpace(raw))))
	hashed := mac.Sum(nil)
	return "pseudo_" + hex.EncodeToString(hashed[:8])
}
