package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

// JSONExporter streams records as newline-delimited JSON (NDJSON).
type JSONExporter struct {
	enc *json.Encoder
}

// NewJSONExporter creates a JSONExporter.
func NewJSONExporter(w io.Writer) *JSONExporter {
	return &JSONExporter{enc: json.NewEncoder(w)}
}

// Export streams a record as JSON line.
func (e *JSONExporter) Export(r *model.Record) error {
	return e.enc.Encode(r)
}

// CSVExporter formats records as CSV rows.
type CSVExporter struct {
	w   *csv.Writer
	hdr bool
}

// NewCSVExporter creates a CSVExporter.
func NewCSVExporter(w io.Writer) *CSVExporter {
	return &CSVExporter{w: csv.NewWriter(w)}
}

// Export writes a CSV row.
func (e *CSVExporter) Export(r *model.Record) error {
	if !e.hdr {
		if err := e.w.Write([]string{"timestamp", "level", "latency_ms", "status", "method", "path", "message"}); err != nil {
			return err
		}
		e.hdr = true
	}

	ts := ""
	if !r.Timestamp.IsZero() {
		ts = r.Timestamp.Format(time.RFC3339)
	}

	row := []string{
		ts,
		r.Level.String(),
		fmt.Sprintf("%.2f", float64(r.Latency.Milliseconds())),
		fmt.Sprintf("%d", r.StatusCode),
		r.Method,
		r.Path,
		r.Message,
	}

	return e.w.Write(row)
}

// Flush flushes buffered CSV output.
func (e *CSVExporter) Flush() {
	e.w.Flush()
}
