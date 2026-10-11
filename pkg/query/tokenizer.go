package query

import (
	"strings"
)

// PipelineStage represents an operation in a pipeline pipe string.
type PipelineStage struct {
	Command string
	Args    []string
}

// ParsePipeline parses strings like "filter(status >= 500) | topk(path, 5) | stats(p99)"
func ParsePipeline(query string) ([]PipelineStage, error) {
	rawStages := strings.Split(query, "|")
	var stages []PipelineStage

	for _, s := range rawStages {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			continue
		}

		openParen := strings.IndexByte(trimmed, '(')
		closeParen := strings.LastIndexByte(trimmed, ')')

		if openParen != -1 && closeParen != -1 && closeParen > openParen {
			cmd := strings.TrimSpace(trimmed[:openParen])
			inner := strings.TrimSpace(trimmed[openParen+1 : closeParen])
			args := splitArgs(inner)
			stages = append(stages, PipelineStage{
				Command: cmd,
				Args:    args,
			})
		} else {
			parts := strings.Fields(trimmed)
			if len(parts) > 0 {
				stages = append(stages, PipelineStage{
					Command: parts[0],
					Args:    parts[1:],
				})
			}
		}
	}

	return stages, nil
}

func splitArgs(s string) []string {
	if len(s) == 0 {
		return nil
	}
	raw := strings.Split(s, ",")
	args := make([]string, len(raw))
	for i, a := range raw {
		args[i] = strings.Trim(strings.TrimSpace(a), "\"'")
	}
	return args
}
