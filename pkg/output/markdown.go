package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

// MarkdownReportGenerator generates GitHub-flavored markdown incident/log reports.
type MarkdownReportGenerator struct {
	w io.Writer
}

// NewMarkdownReportGenerator creates a MarkdownReportGenerator.
func NewMarkdownReportGenerator(w io.Writer) *MarkdownReportGenerator {
	return &MarkdownReportGenerator{w: w}
}

// Generate writes a formatted markdown report.
func (g *MarkdownReportGenerator) Generate(
	totalLogs int64,
	duration time.Duration,
	lat stats.LatencySummary,
	errClusters []*stats.ErrorCluster,
	statusMap map[string]int64,
	slowEndpoints []*stats.EndpointStats,
) error {
	var sb strings.Builder

	sb.WriteString("# LogSift Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("*Generated at: %s*\n\n", time.Now().Format(time.RFC1123)))

	sb.WriteString("## 1. Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Records Processed:** %d\n", totalLogs))
	if duration > 0 {
		rate := float64(totalLogs) / duration.Seconds()
		sb.WriteString(fmt.Sprintf("- **Processing Time:** %v (%.1f records/sec)\n", duration, rate))
	}
	sb.WriteString("\n")

	sb.WriteString("## 2. Latency Percentiles\n\n")
	sb.WriteString("| Metric | Duration |\n")
	sb.WriteString("|:-------|:---------|\n")
	sb.WriteString(fmt.Sprintf("| Min | %v |\n", lat.Min))
	sb.WriteString(fmt.Sprintf("| Mean | %v |\n", lat.Mean))
	sb.WriteString(fmt.Sprintf("| p50 (Median) | %v |\n", lat.P50))
	sb.WriteString(fmt.Sprintf("| p75 | %v |\n", lat.P75))
	sb.WriteString(fmt.Sprintf("| p90 | %v |\n", lat.P90))
	sb.WriteString(fmt.Sprintf("| p95 | %v |\n", lat.P95))
	sb.WriteString(fmt.Sprintf("| p99 | %v |\n", lat.P99))
	sb.WriteString(fmt.Sprintf("| Max | %v |\n", lat.Max))
	sb.WriteString("\n")

	sb.WriteString("## 3. Top Error Clusters\n\n")
	if len(errClusters) == 0 {
		sb.WriteString("No errors recorded.\n\n")
	} else {
		sb.WriteString("| Count | Error Fingerprint | Sample Message |\n")
		sb.WriteString("|:------|:------------------|:---------------|\n")
		for _, c := range errClusters {
			sb.WriteString(fmt.Sprintf("| %d | `%s` | `%s` |\n", c.Count, c.Fingerprint, c.Sample))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 4. HTTP Status Distribution\n\n")
	sb.WriteString("| Category | Count |\n")
	sb.WriteString("|:---------|:------|\n")
	for cat, count := range statusMap {
		sb.WriteString(fmt.Sprintf("| %s | %d |\n", cat, count))
	}
	sb.WriteString("\n")

	if len(slowEndpoints) > 0 {
		sb.WriteString("## 5. Slowest Endpoints\n\n")
		sb.WriteString("| Endpoint | Requests | Error Count | P95 Latency |\n")
		sb.WriteString("|:---------|:---------|:------------|:------------|\n")
		for _, ep := range slowEndpoints {
			sb.WriteString(fmt.Sprintf("| `%s` | %d | %d | %v |\n", ep.Path, ep.Count, ep.Errors, ep.P95Latency))
		}
		sb.WriteString("\n")
	}

	_, err := g.w.Write([]byte(sb.String()))
	return err
}
