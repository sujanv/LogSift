package mask

import (
	"strings"
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestMasker_EmailAndSSN(t *testing.T) {
	m := NewMasker(RedactStrategy{})

	input := "user alice@example.com registered with SSN 123-45-6789"
	masked := m.MaskText(input)

	if !strings.Contains(masked, "[REDACTED_EMAIL]") {
		t.Errorf("expected email to be redacted, got: %s", masked)
	}
	if !strings.Contains(masked, "[REDACTED_SSN]") {
		t.Errorf("expected SSN to be redacted, got: %s", masked)
	}
}

func TestMasker_MaskRecord(t *testing.T) {
	m := NewMasker(MaskStrategy{})

	r := &model.Record{
		Message: "Authentication error for bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.doNotLeakSignature",
		Fields: map[string]interface{}{
			"user_email": "john.doe@company.org",
			"password":   "SuperSecret99!",
		},
	}

	m.MaskRecord(r)

	if r.Fields["password"] != "[PROTECTED]" {
		t.Errorf("expected password field to be protected, got: %v", r.Fields["password"])
	}
	if !strings.Contains(r.Fields["user_email"].(string), "***") {
		t.Errorf("expected masked email, got: %v", r.Fields["user_email"])
	}
}

func BenchmarkMasker(b *testing.B) {
	m := NewMasker(RedactStrategy{})
	text := "user test@example.com with token Bearer secret-token-12345 requested /api/v1/profile"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m.MaskText(text)
	}
}
