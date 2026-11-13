package pseudonym

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestPseudonymizer(t *testing.T) {
	p := NewPseudonymizer("enterprise-salt-123")

	t1 := p.Tokenize("alice@example.com")
	t2 := p.Tokenize("alice@example.com")
	t3 := p.Tokenize("bob@example.com")

	if t1 != t2 {
		t.Errorf("pseudonym must be deterministic for identical inputs: %s != %s", t1, t2)
	}
	if t1 == t3 {
		t.Errorf("different inputs must not produce identical tokens")
	}

	rec := &model.Record{
		UserID: "usr_99812",
		IP:     "192.168.1.55",
		Fields: map[string]interface{}{
			"email": "alice@example.com",
		},
	}
	p.TokenizeRecord(rec)

	if rec.UserID == "usr_99812" || rec.IP == "192.168.1.55" {
		t.Errorf("PII fields were not tokenized: %+v", rec)
	}
	if rec.Fields["email"] != t1 {
		t.Errorf("expected email to match pseudonym token t1")
	}
}
