# LogSift 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/sujanv/LogSift)](https://goreportcard.com/report/github.com/sujanv/LogSift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sujanv/LogSift)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**LogSift** is an enterprise-grade, memory-efficient structured log analysis toolkit written in Go. Designed to process multi-gigabyte log files and high-throughput streams in bounded memory ($O(1)$ RAM), LogSift combines streaming parsers, reservoir percentiles, PII data masking, automated error clustering, log drift & regression diffing, distributed trace waterfalls, and an interactive Terminal UI (TUI).

---

## ⚡ Highlights & Key Features

- **🚀 Constant-Memory Streaming I/O**: Process multi-gigabyte files (< 1MB heap allocation on 100,000+ records) using custom buffered readers (`bufio.Reader.ReadLine()`), recycled `sync.Pool` byte buffers, and reservoir sampling.
- **📦 Transparent Gzip Decompression**: Automatically detects gzip magic bytes (`0x1f, 0x8b`) and streams `.gz` archives without extracting them to disk.
- **🧠 Auto-Detecting Multi-Format Parser**:
  - **Structured JSON**: Zerolog, Zap, Bunyan, Pino, Logrus, Spring Boot.
  - **Logfmt**: Key-value pairs (`level=info msg="user login" duration=45ms`).
  - **Apache / Nginx Access Logs**: Common (CLF) and Combined formats.
  - **Syslog**: RFC 3164 and RFC 5424 formats.
- **🔒 PII & Data Sanitization Engine (`pkg/mask`)**:
  - Scans and redacts Credit Cards (Luhn algorithm verified), Social Security Numbers (SSN), Emails, Bearer Tokens, Passwords, and JWTs.
  - Configurable masking strategies: `redact` (`[REDACTED_...]`), `mask` (`jo***@domain.com`), and `hash` (SHA-256 prefix).
- **🎯 Expressive Search & Filter Engine**:
  - Field comparison operators: `=`, `!=`, `>`, `>=`, `<`, `<=`, `contains`, `in`.
  - Boolean query parser: `level=ERROR AND (status>=500 OR latency>250ms) AND NOT path=/healthz`.
  - Human relative duration filters: `--since 15m`, `--since 2h`, `--since 7d`.
- **📊 Real-time Statistical Analytics**:
  - Bounded-space reservoir sampling for latency percentiles ($p50, p75, p90, p95, p99, p99.9$).
  - Error Clustering & Fingerprinting: Automatically tokenizes dynamic parameters (UUIDs, IP addresses, timestamps, numbers, hashes) to cluster recurring root-cause error templates.
  - HTTP status breakdown (2xx, 3xx, 4xx, 5xx) and slowest endpoint discovery.
  - Time-bucketed volume histogram with Unicode sparklines (` ▂▃▄▅▆▇█`).
- **🔍 Log Drift & Regression Diff (`pkg/diff`)**:
  - Compare pre-release vs. post-release log sessions.
  - Detects error rate deltas, $p95$/$p99$ latency regressions, new error fingerprints introduced, and resolved errors.
- **🧵 Distributed Trace Correlation (`pkg/correlate`)**:
  - Correlate records by `trace_id` across log files and render an ASCII waterfall timeline.
- **🖥️ Interactive Terminal UI (`pkg/tui`)**:
  - Bubble Tea & Lipgloss interactive log viewer with live search (`/`), level filters (`1`-`5`), bookmarking (`m`), and JSON detail inspector (`Enter`).
- **🌐 Offline HTML Dashboard & Embedded Server**:
  - Generate self-contained HTML reports with embedded SVG charts (zero CDN dependencies).
  - Embedded REST API and live Server-Sent Events (SSE) log stream (`logsift serve --port 8080`).
- **💾 Embedded Pure-Go SQLite Indexing (`pkg/storage`)**:
  - Bulk ingests logs into SQLite at **65,000+ rows/sec** in WAL mode for ad-hoc SQL analytics.
- **⚡ Synthetic Log Generator (`cmd/loggen`)**:
  - Generates realistic multi-gigabyte log datasets at **200,000+ logs/sec**.

---

## 🏗️ Architecture

```
                 +--------------------------------+
                 | Log Stream (File / Gzip / Stdin)|
                 +----------------+---------------+
                                  |
                           [ StreamReader ]
                      (Buffered + Chunk Recycler)
                                  |
                                  v
                        [ AutoDetect / Parser ]
                  (JSON / Logfmt / CLF / Syslog)
                                  |
                    +-------------+-------------+
                    |                           |
                    v                           v
          [ Filter / Masker ]          [ Streaming Aggregator ]
       (Fields, Regex, Timestamps,   (Percentiles, Error Fingerprints,
         Luhn PII Sanitization)          Status Code Breakdown)
                    |                           |
       +------------+------------+              |
       |            |            |              |
       v            v            v              v
  [Terminal]     [JSON]        [CSV]     [HTML / Markdown Report]
       |
       +--------------------+-------------------+
       |                    |                   |
       v                    v                   v
  [Bubble Tea TUI]   [SQLite Index]    [Embedded REST / SSE]
```

---

## 📦 Installation

```bash
git clone https://github.com/sujanv/LogSift.git
cd LogSift
make build
```

Binaries will be generated in the root directory:
- `logsift`: Main analysis and query CLI tool.
- `loggen`: Synthetic log generator tool.

Or install directly via `go install`:
```bash
go install github.com/sujanv/LogSift/cmd/logsift@latest
```

---

## 💻 CLI Reference & Recipes

### 1. Search & Filter
```bash
# Filter errors from a 10GB compressed log file
./logsift search --level ERROR /var/log/app.log.gz

# Complex boolean predicate with latency threshold
./logsift search --filter "status>=500 AND latency>250ms" production.log

# Stream directly from Kubernetes or systemd journal
kubectl logs -f deployment/api | ./logsift search --level WARN
```

### 2. Summarize Latency & Errors
```bash
# Terminal summary table
./logsift summarize access.log

# Export post-mortem Markdown report
./logsift summarize --report incident-report.md access.log
```

### 3. PII Masking & Sanitization
```bash
# Redact emails, credit cards, and bearer tokens
./logsift mask --strategy redact app.log > sanitized.log

# Mask using asterisks (e.g. jo***@domain.com)
./logsift mask --strategy mask app.log

# Mask using cryptographic SHA-256 prefixes
./logsift mask --strategy hash app.log
```

### 4. Pre/Post Deployment Regression Diff
```bash
# Compare pre-release (v1.0.0) and post-release (v1.1.0) logs
./logsift diff v1.0.0.log v1.1.0.log
```
*Output highlights:*
- Overall error rate change (e.g. `+1.85% (Regression)`)
- $p50$, $p95$, $p99$ latency deltas
- Newly introduced error fingerprints vs. resolved error templates

### 5. Distributed Trace Waterfall
```bash
# Assemble and display trace timeline across microservices
./logsift correlate --trace 8f7e2a9b-4321 cluster.log
```

### 6. Interactive Terminal UI (TUI)
```bash
./logsift tui app.log
```
- **`j` / `k`** or **`↓` / `↑`**: Navigate through log records
- **`/`**: Real-time regex / substring search filter
- **`1` - `5`**: Quick filter by log level (`1`=All, `2`=Debug, `3`=Info, `4`=Warn, `5`=Error)
- **`Enter` / `Space`**: Open structured record detail & JSON inspector modal
- **`Tab` / `s`**: Toggle Live Analytics dashboard view
- **`m`**: Toggle bookmark on current record
- **`q`**: Quit

### 7. Standalone HTML & SVG Dashboard
```bash
# Generate single-file offline HTML dashboard with embedded SVG charts
./logsift html --out report.html app.log
```

### 8. Live HTTP Monitoring Server
```bash
# Start embedded server on port 8080
./logsift serve --port 8080
```
- `GET /health`: Health status.
- `GET /api/v1/metrics`: Latency summary JSON ($p50, p75, p90, p95, p99, p99.9$).
- `GET /api/v1/errors`: Top recurring error clusters JSON.

### 9. SQL Analytics via Embedded SQLite
```bash
# Ingest 100,000 logs into SQLite
./logsift index --db logs.db app.log

# Query slowest endpoints with SQL
./logsift query --db logs.db "SELECT path, ROUND(AVG(latency_ms), 2), COUNT(*) FROM logs WHERE level='ERROR' GROUP BY path ORDER BY 2 DESC LIMIT 5"
```

### 10. Generate Benchmark Synthetic Logs
```bash
# Generate 1,000,000 synthetic JSON log entries
./loggen -n 1000000 -o benchmark.log
```

---

## ⚡ Benchmarks

Tested on Apple M-series / 100,000 log records:

| Operation | Throughput | Memory Footprint (RSS) |
|:----------|:-----------|:-----------------------|
| Line Ingestion (`StreamReader`) | **~4.8 GB/sec** | **< 4 MB** |
| JSON Parser & Field Extraction | **~380,000 logs/sec** | **< 18 MB** |
| PII Masking & Regex Sanitization | **~240,000 logs/sec** | **< 12 MB** |
| Streaming Reservoir Sampling ($p99$) | **~1,200,000 logs/sec** | **< 8 MB** |
| Full Search & Filter Pipeline | **~310,000 logs/sec** | **< 22 MB** |
| SQLite Batch Insertion | **~65,000 logs/sec** | **< 45 MB** |

---

## 🧪 Testing

```bash
# Run unit tests across all packages
make test

# Run memory-bound assertion & end-to-end suite
go test -v ./test/...

# Run performance benchmarks
make bench
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
