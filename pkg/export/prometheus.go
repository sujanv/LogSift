package export

import (
	"fmt"
	"strings"

	"github.com/sujanv/LogSift/pkg/stats"
)

func FormatPrometheusMetrics(totalLogs, totalErrors int64, lat stats.LatencySummary, statusMap map[string]int64) string {
	var sb strings.Builder

	sb.WriteString("# HELP logsift_records_total Total log records processed.\n")
	sb.WriteString("# TYPE logsift_records_total counter\n")
	sb.WriteString(fmt.Sprintf("logsift_records_total %d\n\n", totalLogs))

	sb.WriteString("# HELP logsift_errors_total Total error-level records.\n")
	sb.WriteString("# TYPE logsift_errors_total counter\n")
	sb.WriteString(fmt.Sprintf("logsift_errors_total %d\n\n", totalErrors))

	sb.WriteString("# HELP logsift_latency_seconds Latency quantiles.\n")
	sb.WriteString("# TYPE logsift_latency_seconds summary\n")
	sb.WriteString(fmt.Sprintf("logsift_latency_seconds{quantile=\"0.5\"} %.4f\n", lat.P50.Seconds()))
	sb.WriteString(fmt.Sprintf("logsift_latency_seconds{quantile=\"0.9\"} %.4f\n", lat.P90.Seconds()))
	sb.WriteString(fmt.Sprintf("logsift_latency_seconds{quantile=\"0.95\"} %.4f\n", lat.P95.Seconds()))
	sb.WriteString(fmt.Sprintf("logsift_latency_seconds{quantile=\"0.99\"} %.4f\n", lat.P99.Seconds()))
	sb.WriteString(fmt.Sprintf("logsift_latency_seconds_count %d\n\n", lat.Count))

	if len(statusMap) > 0 {
		sb.WriteString("# HELP logsift_http_status_total HTTP status codes observed.\n")
		sb.WriteString("# TYPE logsift_http_status_total counter\n")
		for code, count := range statusMap {
			sb.WriteString(fmt.Sprintf("logsift_http_status_total{code=\"%s\"} %d\n", code, count))
		}
	}

	return sb.String()
}
