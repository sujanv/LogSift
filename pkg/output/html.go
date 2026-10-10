package output

import (
	"fmt"
	"io"
	"time"

	"github.com/sujanv/LogSift/pkg/stats"
)

// HTMLReportGenerator generates an offline standalone HTML dashboard.
type HTMLReportGenerator struct {
	w io.Writer
}

// NewHTMLReportGenerator creates an HTMLReportGenerator.
func NewHTMLReportGenerator(w io.Writer) *HTMLReportGenerator {
	return &HTMLReportGenerator{w: w}
}

// Generate renders the HTML dashboard with dark mode styling and metrics.
func (g *HTMLReportGenerator) Generate(
	totalLogs int64,
	duration time.Duration,
	lat stats.LatencySummary,
	errClusters []*stats.ErrorCluster,
	statusMap map[string]int64,
) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>LogSift Dashboard</title>
<style>
  body { background: #0f172a; color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 0; padding: 2rem; }
  .container { max-width: 1100px; margin: 0 auto; }
  h1 { color: #818cf8; margin-bottom: 0.5rem; }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; margin: 1.5rem 0; }
  .card { background: #1e293b; padding: 1.5rem; border-radius: 8px; border: 1px solid #334155; }
  .metric { font-size: 1.8rem; font-weight: bold; color: #38bdf8; }
  .label { font-size: 0.9rem; color: #94a3b8; }
  table { width: 100%%; border-collapse: collapse; margin-top: 1rem; background: #1e293b; border-radius: 8px; overflow: hidden; }
  th, td { padding: 0.75rem 1rem; text-align: left; border-bottom: 1px solid #334155; }
  th { background: #0f172a; color: #94a3b8; }
  .code { font-family: monospace; color: #f43f5e; }
</style>
</head>
<body>
<div class="container">
  <h1>LogSift Analytics Report</h1>
  <p style="color:#64748b;">Generated at %s</p>

  <div class="grid">
    <div class="card"><div class="label">Total Records</div><div class="metric">%d</div></div>
    <div class="card"><div class="label">Median (p50) Latency</div><div class="metric">%v</div></div>
    <div class="card"><div class="label">p95 Latency</div><div class="metric">%v</div></div>
    <div class="card"><div class="label">p99 Latency</div><div class="metric">%v</div></div>
  </div>

  <h2>Top Recurring Error Clusters</h2>
  <table>
    <thead><tr><th>Count</th><th>Fingerprint</th><th>Sample</th></tr></thead>
    <tbody>
`, time.Now().Format(time.RFC1123), totalLogs, lat.P50, lat.P95, lat.P99)

	for _, c := range errClusters {
		html += fmt.Sprintf("<tr><td>%d</td><td class='code'>%s</td><td>%s</td></tr>\n", c.Count, c.Fingerprint, c.Sample)
	}

	html += `    </tbody>
  </table>
</div>
</body>
</html>`

	_, err := g.w.Write([]byte(html))
	return err
}
