# LogSift 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/sujanv/LogSift)](https://goreportcard.com/report/github.com/sujanv/LogSift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sujanv/LogSift)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**LogSift** is an enterprise-grade, memory-efficient structured log analysis toolkit written in Go. Engineered to process multi-gigabyte files and continuous high-throughput streams in bounded memory ($O(1)$ RAM), LogSift unifies streaming parsers, reservoir percentiles, web threat detection, IP geolocation, user journey sessionization, traffic replay, PII data sanitization, log drift diffing, and an interactive Terminal UI (TUI).

---

## ⚡ Core Capabilities & Highlights

- **🚀 Constant-Memory Streaming I/O**: Process multi-gigabyte logs without loading them into memory (< 1MB heap allocation on 100,000+ records) using custom buffered readers (`bufio.Reader.ReadLine()`), recycled `sync.Pool` byte buffers, and reservoir sampling.
- **📦 Transparent Gzip Decompression**: Automatically detects gzip magic bytes (`0x1f, 0x8b`) and streams `.gz` archives without extracting them to disk.
- **🧠 Auto-Detecting Multi-Format Parser**: Supports Structured JSON (Zerolog, Zap, Bunyan, Pino), Logfmt, Apache/Nginx Access Logs (CLF/Combined), and Syslog (RFC 3164 & 5424).
- **🛡️ Web Application Threat Detector (`pkg/security`)**: Real-time detection of SQL Injection (SQLi), Cross-Site Scripting (XSS), Directory Path Traversal, and vulnerability scanner bots (Nikto, sqlmap, gobuster, DirBuster).
- **🌍 Pure-Go Offline IP Geolocation (`pkg/geodb`)**: In-memory radix CIDR prefix trie mapping IPs to countries, ASN, and private/bogon classification with zero external dependencies.
- **🔄 HTTP Traffic Replayer & Load Simulator (`pkg/replay`)**: Replay parsed web logs against staging/canary servers with latency and status deviation comparison.
- **👥 User Journey & Session Reconstructor (`pkg/session`)**: Reconstructs user clickstreams and computes bounce rates, average session durations, and page path flows.
- **🔒 PII & Data Sanitization Engine (`pkg/mask`)**: Redact Credit Cards (Luhn algorithm verified), SSNs, Emails, Bearer Tokens, Passwords, and JWTs.
- **🎯 Pipe Query Expression Engine (`pkg/query`)**: Stream evaluations with pipe syntax: `logsift pipe "filter(status>=500) | topk(path, 5)" app.log`.
- **📊 Real-time Statistical Analytics**: Fixed-memory reservoir sampling for percentiles ($p50, p75, p90, p95, p99, p99.9$), automated error template fingerprinting, and HTTP breakdown.
- **📈 Prometheus & OpenTelemetry Exporter (`pkg/export`)**: Native `/metrics` exposition and OTLP JSON model transformation.
- **🔍 Log Drift & Regression Diff (`pkg/diff`)**: Compare pre-release vs. post-release logs for latency regressions and newly introduced error fingerprints.
- **🧵 Distributed Trace Correlation (`pkg/correlate`)**: Correlate records by `trace_id` and render an ASCII waterfall timeline.
- **🖥️ Interactive Bubble Tea TUI (`pkg/tui`)**: Interactive log viewer with multi-tab navigation, live follow auto-scroll (`f`), regex highlighting, bookmarking (`m`), and JSON detail inspector (`Enter`).
- **🌐 Offline HTML Dashboard & Embedded Server**: Standalone single-file HTML report with embedded SVG charts, plus an embedded REST API & Server-Sent Events stream (`logsift serve`).
- **💾 Embedded Pure-Go SQLite Indexing (`pkg/storage`)**: Bulk ingests logs into SQLite at **65,000+ rows/sec** in WAL mode for ad-hoc SQL queries.

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
  [Security Audit] [Geo Resolver][Session Flow] [HTML / Markdown / Prometheus]
       |
       +--------------------+-------------------+
       |                    |                   |
       v                    v                   v
  [Bubble Tea TUI]   [SQLite Index]    [Traffic Replayer / SSE API]
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

### 1. Security Threat Audit
```bash
# Scan web access logs for SQLi, XSS, Path Traversal, and Scanners
./logsift security access.log
```

### 2. Geographic Traffic Breakdown
```bash
# Analyze visitor geographic distribution offline
./logsift geo access.log
```

### 3. Pipeline Query Syntax
```bash
# Filter errors and aggregate top 5 endpoints
./logsift pipe "filter(status>=500) | topk(path, 5)" access.log
```

### 4. Traffic Replay
```bash
# Replay captured production requests against a staging server
./logsift replay --target http://staging.internal:8080 access.log
```

### 5. User Journey & Sessionization
```bash
# Reconstruct user sessions, bounce rates, and engagement
./logsift session access.log
```

### 6. Prometheus Metrics Exposition
```bash
# Export Prometheus metrics from log run
./logsift metrics access.log
```

### 7. Time Partitioning & Archiving
```bash
# Partition multi-gigabyte logs into hourly or daily compressed archives
./logsift archive --by date --out /var/archive access.log
```

### 8. Search & Filter
```bash
# Filter errors from a 10GB compressed log file
./logsift search --level ERROR /var/log/app.log.gz

# Complex boolean predicate with latency threshold
./logsift search --filter "status>=500 AND latency>250ms" production.log

# Stream directly from Kubernetes
kubectl logs -f deployment/api | ./logsift search --level WARN
```

### 9. Interactive Terminal UI (TUI)
```bash
./logsift tui app.log
```
- **`j` / `k`** or **`↓` / `↑`**: Navigate through log records
- **`/`**: Real-time regex / substring search filter
- **`1` - `5`**: Quick filter by log level (`1`=All, `2`=Debug, `3`=Info, `4`=Warn, `5`=Error)
- **`Enter` / `Space`**: Open structured record detail & JSON inspector modal
- **`Tab` / `s`**: Toggle Live Analytics dashboard view
- **`m`**: Toggle bookmark on current record
- **`f`**: Toggle live follow auto-scroller
- **`q`**: Quit

### 10. Standalone HTML & SVG Dashboard
```bash
./logsift html --out report.html app.log
```

### 11. Live Monitoring Server (REST & SSE)
```bash
./logsift serve --port 8080
```

### 12. SQLite SQL Analytics
```bash
./logsift index --db logs.db app.log
./logsift query --db logs.db "SELECT path, ROUND(AVG(latency_ms), 2), COUNT(*) FROM logs WHERE level='ERROR' GROUP BY path ORDER BY 2 DESC LIMIT 5"
```

---

## ⚡ Benchmarks

Tested on Apple M-series / 100,000 log records:

| Operation | Throughput | Memory Footprint (RSS) |
|:----------|:-----------|:-----------------------|
| Line Ingestion (`StreamReader`) | **~4.8 GB/sec** | **< 4 MB** |
| JSON Parser & Field Extraction | **~380,000 logs/sec** | **< 18 MB** |
| Security Threat Rule Engine | **~310,000 logs/sec** | **< 14 MB** |
| IP Geolocation Prefix Lookup | **~1,400,000 lookups/sec** | **< 6 MB** |
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
