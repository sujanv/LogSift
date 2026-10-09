package parser

import "github.com/sujanv/LogSift/pkg/model"

// Parser defines the common interface for log parsers.
type Parser interface {
	// Parse extracts structured information from a raw log line into record.
	Parse(line []byte, record *model.Record) error
	// Name returns the identifier of the parser format.
	Name() string
}
