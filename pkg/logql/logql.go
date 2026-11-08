package logql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"github.com/sujanv/LogSift/pkg/model"
)

type MatchOp string

const (
	OpEqual     MatchOp = "="
	OpNotEqual  MatchOp = "!="
	OpRegex     MatchOp = "=~"
	OpNotRegex  MatchOp = "!~"
	OpContains  MatchOp = "|="
	OpExcludes  MatchOp = "!="
	OpRegMatch  MatchOp = "|~"
	OpRegExcl   MatchOp = "!~"
)

type LabelMatcher struct {
	Key   string
	Op    MatchOp
	Value string
	Re    *regexp.Regexp
}

type PipelineStage struct {
	Op    MatchOp
	Value string
	Re    *regexp.Regexp
}

type Query struct {
	Matchers   []LabelMatcher
	Stages     []PipelineStage
	MetricFunc string
	Window     time.Duration
}

func Parse(query string) (*Query, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, errors.New("empty LogQL query")
	}

	result := &Query{
		Matchers: make([]LabelMatcher, 0),
		Stages:   make([]PipelineStage, 0),
	}

	return parseInnerQuery(q, result)
}

func parseInnerQuery(q string, result *Query) (*Query, error) {
	q = strings.TrimSpace(q)
	rangeIdx := strings.Index(q, "[")
	if rangeIdx != -1 {
		endRange := strings.Index(q, "]")
		if endRange > rangeIdx {
			durStr := q[rangeIdx+1 : endRange]
			dur, err := time.ParseDuration(durStr)
			if err == nil {
				result.Window = dur
			}
			q = q[:rangeIdx] + q[endRange+1:]
		}
	}

	if strings.HasPrefix(q, "{") {
		closeIdx := strings.Index(q, "}")
		if closeIdx == -1 {
			return nil, errors.New("unclosed stream selector {")
		}
		selector := q[1:closeIdx]
		pipeline := q[closeIdx+1:]

		if err := parseLabels(selector, result); err != nil {
			return nil, err
		}
		if err := parsePipeline(pipeline, result); err != nil {
			return nil, err
		}
		return result, nil
	}

	if err := parsePipeline(q, result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseLabels(selector string, result *Query) error {
	pairs := strings.Split(selector, ",")
	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var m LabelMatcher
		if strings.Contains(p, "=~") {
			parts := strings.SplitN(p, "=~", 2)
			m = LabelMatcher{Key: strings.TrimSpace(parts[0]), Op: OpRegex, Value: unquote(parts[1])}
			re, err := regexp.Compile(m.Value)
			if err != nil {
				return fmt.Errorf("invalid regex in label matcher: %w", err)
			}
			m.Re = re
		} else if strings.Contains(p, "!~") {
			parts := strings.SplitN(p, "!~", 2)
			m = LabelMatcher{Key: strings.TrimSpace(parts[0]), Op: OpNotRegex, Value: unquote(parts[1])}
			re, err := regexp.Compile(m.Value)
			if err != nil {
				return fmt.Errorf("invalid regex in label matcher: %w", err)
			}
			m.Re = re
		} else if strings.Contains(p, "!=") {
			parts := strings.SplitN(p, "!=", 2)
			m = LabelMatcher{Key: strings.TrimSpace(parts[0]), Op: OpNotEqual, Value: unquote(parts[1])}
		} else if strings.Contains(p, "=") {
			parts := strings.SplitN(p, "=", 2)
			m = LabelMatcher{Key: strings.TrimSpace(parts[0]), Op: OpEqual, Value: unquote(parts[1])}
		} else {
			return fmt.Errorf("unrecognized label matcher: %s", p)
		}
		result.Matchers = append(result.Matchers, m)
	}
	return nil
}

func parsePipeline(p string, result *Query) error {
	tokens := strings.Fields(p)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		switch t {
		case "|=", "!=", "|~", "!~":
			if i+1 >= len(tokens) {
				return fmt.Errorf("expected pattern string after operator %s", t)
			}
			val := unquote(tokens[i+1])
			i++
			stage := PipelineStage{Op: MatchOp(t), Value: val}
			if t == "|~" || t == "!~" {
				re, err := regexp.Compile(val)
				if err != nil {
					return fmt.Errorf("invalid regex pattern: %w", err)
				}
				stage.Re = re
			}
			result.Stages = append(result.Stages, stage)
		}
	}
	return nil
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == 34 && s[len(s)-1] == 34) || (s[0] == 39 && s[len(s)-1] == 39) {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func (q *Query) Matches(rec *model.Record) bool {
	if rec == nil {
		return false
	}

	for _, m := range q.Matchers {
		var val string
		switch strings.ToLower(m.Key) {
		case "level":
			val = strings.ToLower(rec.Level.String())
		case "path":
			val = rec.Path
		case "method":
			val = rec.Method
		case "ip":
			val = rec.IP
		case "user_id":
			val = rec.UserID
		case "trace_id":
			val = rec.TraceID
		default:
			if rec.Fields != nil {
				if fv, ok := rec.Fields[m.Key]; ok {
					val = fmt.Sprintf("%v", fv)
				}
			}
		}

		match := false
		switch m.Op {
		case OpEqual:
			match = strings.EqualFold(val, m.Value)
		case OpNotEqual:
			match = !strings.EqualFold(val, m.Value)
		case OpRegex:
			match = m.Re != nil && m.Re.MatchString(val)
		case OpNotRegex:
			match = m.Re != nil && !m.Re.MatchString(val)
		}
		if !match {
			return false
		}
	}

	msg := rec.Message
	if msg == "" && len(rec.Raw) > 0 {
		msg = string(rec.Raw)
	}

	for _, stage := range q.Stages {
		switch stage.Op {
		case OpContains:
			if !strings.Contains(msg, stage.Value) {
				return false
			}
		case OpExcludes:
			if strings.Contains(msg, stage.Value) {
				return false
			}
		case OpRegMatch:
			if stage.Re != nil && !stage.Re.MatchString(msg) {
				return false
			}
		case OpRegExcl:
			if stage.Re != nil && stage.Re.MatchString(msg) {
				return false
			}
		}
	}

	return true
}

func (q *Query) FilterRecords(records []*model.Record) []*model.Record {
	var filtered []*model.Record
	for _, r := range records {
		if q.Matches(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
