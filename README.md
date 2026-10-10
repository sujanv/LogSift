# LogSift 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/sujanv/LogSift)](https://goreportcard.com/report/github.com/sujanv/LogSift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sujanv/LogSift)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**LogSift** is an ultra-fast, memory-efficient structured log analysis tool written in Go. It parses multi-gigabyte JSON and plain-text log files, filters by fields, clusters recurring errors, calculates latency percentiles in constant $O(1)$ space, and features an interactive Terminal User Interface (TUI).

---

## ✨ Features

- **🚀 Constant-Memory Streaming I/O**: Process multi-gigabyte log files without loading them into memory. Powered by custom buffered readers, recycled `sync.Pool` byte buffers, and reservoir sampling.
- **📦 Transparent Gzip Decompression**: Stream `.gz` logs automatically without extracting them to disk.
- **🧠 Auto-Detecting Multi-Format Parser**:
  - Structured JSON (Zap, Zerolog, Bunyan, Pino, Logrus, Spring Boot)
  - Logfmt (`key=val key="quoted val"`)
  - Apache / Nginx Combined & Common Log Format (CLF)
  - Syslog (RFC 3164 & RFC 5424)
- **🎯 Expressive Search & Filter Engine**:
  - Field operators: `status>=500`, `latency>250ms`, `path contains "/api/v1"`, `level=ERROR`
  - Boolean expressions: `level=ERROR AND (status>=500 OR latency>1s) AND NOT path=/healthz`
  - Time ranges: `--since 15m`, `--since 2h`, `--since 7d`, `--from`, `--to`
- **📊 Real-time Statistical Aggregations**:
  - Latency percentiles ($p50, p75, p90, p95, p99, p99.9$) using streaming reservoir sampling.
  - Intelligent Error Fingerprinting: Clusters recurring errors by normalizing UUIDs, IPs, timestamps, and numbers.
  - HTTP status breakdown (2xx, 3xx, 4xx, 5xx) and slowest endpoint discovery.
  - Time-bucketed volume histogram with Unicode sparklines (` ▂▃▄▅▆▇█`).
- **🖥️ Interactive Bubble Tea TUI**:
  - Live log browsing with keyboard navigation (`j`/`k`, `/` search, `1`-`5` level filter).
  - Detail inspector with formatted JSON syntax tree.
  - Live analytics dashboard tab.
- **💾 Optional Pure-Go SQLite Indexing**: Ingest logs into indexed SQLite database at 100k+ rows/sec for instant SQL queries.
- **⚡ Synthetic Log Generator**: Generate gigabytes of realistic JSON logs at 500k+ records/sec.

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
          [ Filter / Matcher ]         [ Streaming Aggregator ]
       (Fields, Regex, Timestamps)  (Percentiles, Error Fingerprints)
                    |                           |
       +------------+------------+              |
       |            |            |              |
       v            v            v              v
  [Terminal]     [JSON]        [CSV]     [Markdown Report]
       |
       v
  [Bubble Tea TUI] / [SQLite Index]
```

---

## 📦 Installation

```bash
git clone https://github.com/sujanv/LogSift.git
cd LogSift
make build
```

Or install via `go install`:

```bash
go install github.com/sujanv/LogSift/cmd/logsift@latest
```

---

## 🚀 Quick Start

### 1. Search and Stream Logs
```bash
# Filter error logs from a 10GB file
logsift search --level ERROR /var/log/app.log

# Complex Boolean expression with latency filter
logsift search --filter "status>=500 AND latency>250ms" production.log.gz

# Stream live container logs via pipe
kubectl logs -f deployment/api | logsift search --level WARN
```

### 2. Summarize Latency & Errors
```bash
# Generate instant terminal summary
logsift summarize access.log

# Export incident report to Markdown
logsift summarize --report incident-report.md access.log
```

### 3. Interactive Terminal UI
```bash
# Launch interactive TUI
logsift tui app.log
```
- **`j` / `k`**: Navigate rows
- **`/`**: Real-time search filter
- **`1`-`5`**: Quick level filter (`1`=All, `2`=Debug, `3`=Info, `4`=Warn, `5`=Error)
- **`Enter` / `Space`**: View structured record details
- **`Tab` / `s`**: Toggle Live Analytics dashboard
- **`q`**: Quit

### 4. SQL Analytics via Embedded SQLite
```bash
# Index 1M logs into SQLite
logsift index --db logs.db app.log

# Run SQL analytics
logsift query --db logs.db "SELECT path, AVG(latency_ms), COUNT(*) FROM logs WHERE level='ERROR' GROUP BY path ORDER BY 2 DESC LIMIT 5"
```

### 5. Generate Benchmark Logs
```bash
# Generate 1,000,000 synthetic JSON logs
go run ./cmd/loggen -n 1000000 -o benchmark.log
```

---

## ⚡ Benchmarks

Tested on Apple M-series / 100,000 log records:

| Operation | Throughput | Memory (RSS) |
|:----------|:-----------|:-------------|
| StreamReader Line Ingestion | **~4.8 GB/sec** | **< 4 MB** |
| JSON Parser & Field Extraction | **~380,000 logs/sec** | **< 18 MB** |
| Streaming Latency ($p99$ Reservoir) | **~1,200,000 logs/sec** | **< 8 MB** |
| Full Pipeline (Search + Filter) | **~310,000 logs/sec** | **< 22 MB** |
| SQLite Batch Insertion | **~110,000 logs/sec** | **< 45 MB** |

---

## 🧪 Testing

```bash
# Run unit tests
make test

# Run memory bound & integration tests
go test -v ./test/...

# Run benchmarks
make bench
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
