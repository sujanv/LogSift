package mask

import (
	"regexp"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

var (
	emailRegex  = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	ssnRegex    = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
	bearerRegex = regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9._~+/-]+=*`)
	jwtRegex    = regexp.MustCompile(`eyJ[a-zA-Z0-9_-]{10,}\.eyJ[a-zA-Z0-9_-]{10,}\.[a-zA-Z0-9_-]+`)
	cardRegex   = regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
)

// Masker sanitizes sensitive information in log streams.
type Masker struct {
	strategy Strategy
}

// NewMasker creates a Masker with the specified strategy.
func NewMasker(s Strategy) *Masker {
	if s == nil {
		s = RedactStrategy{}
	}
	return &Masker{strategy: s}
}

// MaskText scans a string and masks emails, SSNs, credit cards, bearer tokens, and JWTs.
func (m *Masker) MaskText(s string) string {
	s = emailRegex.ReplaceAllStringFunc(s, func(val string) string {
		return m.strategy.Apply(val, "EMAIL")
	})
	s = ssnRegex.ReplaceAllStringFunc(s, func(val string) string {
		return m.strategy.Apply(val, "SSN")
	})
	s = jwtRegex.ReplaceAllStringFunc(s, func(val string) string {
		return m.strategy.Apply(val, "JWT")
	})
	s = bearerRegex.ReplaceAllStringFunc(s, func(val string) string {
		return "Bearer " + m.strategy.Apply(val[7:], "TOKEN")
	})
	s = cardRegex.ReplaceAllStringFunc(s, func(val string) string {
		clean := strings.ReplaceAll(strings.ReplaceAll(val, "-", ""), " ", "")
		if isValidLuhn(clean) {
			return m.strategy.Apply(clean, "CARD")
		}
		return val
	})
	return s
}

// MaskRecord sanitizes message, raw bytes, and sensitive fields map.
func (m *Masker) MaskRecord(r *model.Record) {
	if r.Message != "" {
		r.Message = m.MaskText(r.Message)
	}
	if len(r.Raw) > 0 {
		r.Raw = []byte(m.MaskText(string(r.Raw)))
	}

	sensitiveKeys := []string{"password", "passwd", "token", "secret", "auth", "authorization", "credit_card", "ssn"}
	if r.Fields != nil {
		for k, v := range r.Fields {
			lowerK := strings.ToLower(k)
			for _, sk := range sensitiveKeys {
				if strings.Contains(lowerK, sk) {
					r.Fields[k] = "[PROTECTED]"
					break
				}
			}
			if strVal, ok := v.(string); ok {
				r.Fields[k] = m.MaskText(strVal)
			}
		}
	}
}

func isValidLuhn(number string) bool {
	if len(number) < 13 || len(number) > 19 {
		return false
	}
	sum := 0
	alternate := false
	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
		if n < 0 || n > 9 {
			return false
		}
		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}
	return sum%10 == 0
}
