package filter

import (
	"fmt"
	"strings"

	"github.com/sujanv/LogSift/pkg/model"
)

// Predicate is an evaluation node in an expression tree.
type Predicate interface {
	Match(record *model.Record) bool
}

// AndPredicate evaluates logical AND of multiple predicates.
type AndPredicate struct {
	Children []Predicate
}

func (a *AndPredicate) Match(record *model.Record) bool {
	for _, child := range a.Children {
		if !child.Match(record) {
			return false
		}
	}
	return true
}

// OrPredicate evaluates logical OR of multiple predicates.
type OrPredicate struct {
	Children []Predicate
}

func (o *OrPredicate) Match(record *model.Record) bool {
	if len(o.Children) == 0 {
		return true
	}
	for _, child := range o.Children {
		if child.Match(record) {
			return true
		}
	}
	return false
}

// NotPredicate evaluates logical NOT.
type NotPredicate struct {
	Child Predicate
}

func (n *NotPredicate) Match(record *model.Record) bool {
	return !n.Child.Match(record)
}

// ParseExpression parses query strings into an executable Predicate tree.
// Examples:
// "level=ERROR AND (status>=500 OR latency>500ms)"
func ParseExpression(query string) (Predicate, error) {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return &AndPredicate{}, nil
	}

	tokens := tokenize(query)
	p := &exprParser{tokens: tokens, pos: 0}
	return p.parseOr()
}

type exprParser struct {
	tokens []string
	pos    int
}

func (p *exprParser) peek() string {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return ""
}

func (p *exprParser) next() string {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *exprParser) parseOr() (Predicate, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	children := []Predicate{left}
	for strings.EqualFold(p.peek(), "OR") {
		p.next() // consume OR
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}

	if len(children) == 1 {
		return children[0], nil
	}
	return &OrPredicate{Children: children}, nil
}

func (p *exprParser) parseAnd() (Predicate, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	children := []Predicate{left}
	for strings.EqualFold(p.peek(), "AND") {
		p.next() // consume AND
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}

	if len(children) == 1 {
		return children[0], nil
	}
	return &AndPredicate{Children: children}, nil
}

func (p *exprParser) parseUnary() (Predicate, error) {
	if strings.EqualFold(p.peek(), "NOT") {
		p.next() // consume NOT
		child, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &NotPredicate{Child: child}, nil
	}
	return p.parsePrimary()
}

func (p *exprParser) parsePrimary() (Predicate, error) {
	tok := p.peek()
	if tok == "(" {
		p.next() // consume (
		expr, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.peek() != ")" {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		p.next() // consume )
		return expr, nil
	}

	// Term: e.g. "level=ERROR", "status>=500", "path=/api"
	term := p.next()
	if len(term) == 0 {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	return parseConditionTerm(term)
}

func parseConditionTerm(term string) (Predicate, error) {
	operators := []Op{OpGte, OpLte, OpNeq, OpEq, OpGt, OpLt, OpContains, OpIn}
	for _, op := range operators {
		opStr := string(op)
		idx := strings.Index(term, opStr)
		if idx != -1 {
			field := strings.TrimSpace(term[:idx])
			val := strings.Trim(strings.TrimSpace(term[idx+len(opStr):]), "\"'")
			return NewFieldMatcher(field, op, val), nil
		}
	}

	// Bare term: substring search
	return NewTextMatcher(term, false), nil
}

func tokenize(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			if c == quoteChar {
				inQuote = false
			}
			cur.WriteByte(c)
			continue
		}

		if c == '"' || c == '\'' {
			inQuote = true
			quoteChar = c
			cur.WriteByte(c)
			continue
		}

		if c == '(' || c == ')' {
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			tokens = append(tokens, string(c))
			continue
		}

		if c == ' ' || c == '\t' {
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			continue
		}

		cur.WriteByte(c)
	}

	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}

	return tokens
}
