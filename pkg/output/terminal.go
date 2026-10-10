package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/model"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorGray    = "\033[90m"
	ColorBold    = "\033[1m"
)

// TerminalFormatter renders colored terminal logs.
type TerminalFormatter struct {
	w       io.Writer
	noColor bool
}

// NewTerminalFormatter creates a TerminalFormatter.
func NewTerminalFormatter(w io.Writer, noColor bool) *TerminalFormatter {
	if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}
	return &TerminalFormatter{w: w, noColor: noColor}
}

// FormatRecord prints a single record nicely formatted.
func (tf *TerminalFormatter) FormatRecord(r *model.Record) {
	tsStr := ""
	if !r.Timestamp.IsZero() {
		tsStr = r.Timestamp.Format("2006-01-02 15:04:05.000")
	}

	lvlBadge := tf.colorizeLevel(r.Level)
	latencyStr := tf.colorizeLatency(r.Latency)

	msg := r.Message
	if msg == "" && len(r.Raw) > 0 {
		msg = string(r.Raw)
	}

	extra := ""
	if r.StatusCode > 0 {
		extra += fmt.Sprintf(" status=%d", r.StatusCode)
	}
	if r.Path != "" {
		extra += fmt.Sprintf(" path=%s", r.Path)
	}
	if r.TraceID != "" {
		extra += fmt.Sprintf(" trace=%s", r.TraceID)
	}

	fmt.Fprintf(tf.w, "%s %s %s %s%s\n", tsStr, lvlBadge, latencyStr, msg, extra)
}

func (tf *TerminalFormatter) colorizeLevel(l model.LogLevel) string {
	s := fmt.Sprintf("[%-5s]", l.String())
	if tf.noColor {
		return s
	}

	switch l {
	case model.LevelFatal, model.LevelError:
		return fmt.Sprintf("%s%s%s%s", ColorBold, ColorRed, s, ColorReset)
	case model.LevelWarn:
		return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
	case model.LevelInfo:
		return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
	case model.LevelDebug:
		return fmt.Sprintf("%s%s%s", ColorCyan, s, ColorReset)
	default:
		return fmt.Sprintf("%s%s%s", ColorGray, s, ColorReset)
	}
}

func (tf *TerminalFormatter) colorizeLatency(d time.Duration) string {
	if d <= 0 {
		return "      "
	}
	s := fmt.Sprintf("%6s", d.Round(time.Millisecond).String())
	if tf.noColor {
		return s
	}

	switch {
	case d >= 500*time.Millisecond:
		return fmt.Sprintf("%s%s%s", ColorRed, s, ColorReset)
	case d >= 200*time.Millisecond:
		return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
	default:
		return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
	}
}

// PrintTable prints tabular summary rows.
func (tf *TerminalFormatter) PrintTable(headers []string, rows [][]string) {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print headers
	for i, h := range headers {
		fmt.Fprintf(tf.w, "%-*s  ", colWidths[i], h)
	}
	fmt.Fprintln(tf.w)

	// Print divider
	for _, w := range colWidths {
		fmt.Fprintf(tf.w, "%s  ", strings.Repeat("-", w))
	}
	fmt.Fprintln(tf.w)

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				fmt.Fprintf(tf.w, "%-*s  ", colWidths[i], cell)
			}
		}
		fmt.Fprintln(tf.w)
	}
}
