package filter

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

// TextMatcher searches for substring matches in log message or raw line.
type TextMatcher struct {
	Query         string
	CaseSensitive bool
	queryBytes    []byte
}

// NewTextMatcher creates a TextMatcher.
func NewTextMatcher(query string, caseSensitive bool) *TextMatcher {
	return &TextMatcher{
		Query:         query,
		CaseSensitive: caseSensitive,
		queryBytes:    []byte(query),
	}
}

// Match returns true if query is found in record message or raw bytes.
func (tm *TextMatcher) Match(record *model.Record) bool {
	if len(tm.Query) == 0 {
		return true
	}

	if tm.CaseSensitive {
		if strings.Contains(record.Message, tm.Query) {
			return true
		}
		if len(record.Raw) > 0 && bytes.Contains(record.Raw, tm.queryBytes) {
			return true
		}
	} else {
		if strings.Contains(strings.ToLower(record.Message), strings.ToLower(tm.Query)) {
			return true
		}
		if len(record.Raw) > 0 && bytes.Contains(bytes.ToLower(record.Raw), bytes.ToLower(tm.queryBytes)) {
			return true
		}
	}

	return false
}

// RegexMatcher matches log records against a compiled regular expression.
type RegexMatcher struct {
	re *regexp.Regexp
}

// NewRegexMatcher compiles and creates a RegexMatcher.
func NewRegexMatcher(pattern string) (*RegexMatcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexMatcher{re: re}, nil
}

// Match returns true if pattern matches message or raw line.
func (rm *RegexMatcher) Match(record *model.Record) bool {
	if rm.re == nil {
		return true
	}
	if rm.re.MatchString(record.Message) {
		return true
	}
	if len(record.Raw) > 0 && rm.re.Match(record.Raw) {
		return true
	}
	return false
}
