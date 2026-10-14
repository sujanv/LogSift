package schema

type SchemaSummary struct {
	TotalLogs  int64            `json:"total_logs"`
	FieldCount int              `json:"field_count"`
	Fields     map[string]int64 `json:"fields"`
}

func (s *InferredSchema) Summary() SchemaSummary {
	s.mu.Lock()
	defer s.mu.Unlock()

	fm := make(map[string]int64, len(s.Fields))
	for k, v := range s.Fields {
		fm[k] = v.Count
	}
	return SchemaSummary{
		TotalLogs:  s.TotalLogs,
		FieldCount: len(s.Fields),
		Fields:     fm,
	}
}
