package pseudonym

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"github.com/sujanv/LogSift/pkg/model"
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

func (p *Pseudonymizer) RotateSalt(newSalt string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.salt = []byte(newSalt)
}

func (p *Pseudonymizer) TokenizeRecord(rec *model.Record) {
	if rec == nil {
		return
	}

	if rec.UserID != "" {
		rec.UserID = p.Tokenize(rec.UserID)
	}
	if rec.IP != "" {
		rec.IP = p.Tokenize(rec.IP)
	}

	if rec.Fields != nil {
		for _, key := range []string{"email", "user", "username", "ssn", "client_ip"} {
			if val, ok := rec.Fields[key].(string); ok && val != "" {
				rec.Fields[key] = p.Tokenize(val)
			}
		}
	}
}
