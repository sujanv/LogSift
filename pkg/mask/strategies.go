package mask

type Strategy interface {
	Apply(val, category string) string
}

type RedactStrategy struct{}

func (RedactStrategy) Apply(val, category string) string {
	return "[REDACTED_" + category + "]"
}
