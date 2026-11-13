# LogSift 🔍

[![Go Report Card](https://goreportcard.com/badge/github.com/sujanv/LogSift)](https://goreportcard.com/report/github.com/sujanv/LogSift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sujanv/LogSift)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**LogSift** is an enterprise-grade, memory-efficient structured log analysis toolkit written in Go. Engineered to process multi-gigabyte files and continuous high-throughput streams in bounded memory ($O(1)$ RAM), LogSift unifies streaming parsers, cloud ingestion cost estimators, debug-on-fault flight recorders, deterministic PII pseudonymization, differential incident root-cause diagnosis, OpenTelemetry OTLP JSON receivers, regional latency heatmaps, Loki-compatible LogQL queries, W3C Trace Context span extraction, circuit-breaker-protected remote shipping, streaming JSON schema validation, token TF-IDF log vectorization, sliding-window deduplication, SLO error budget burn rate tracking, reservoir percentiles, statistical anomaly detection, Drain template clustering, HIPAA/GDPR compliance auditing, database slow-query profiling, cryptographic tamper-evident signing, web threat detection, offline IP geolocation, user journey sessionization, traffic replay, PII sanitization, log drift diffing, and an interactive Terminal UI (TUI).

---

## ⚡ Core Capabilities & Highlights

- **🚀 Constant-Memory Streaming I/O**: Process multi-gigabyte logs without loading them into memory (< 1MB heap allocation on 100,000+ records) using custom buffered readers (`bufio.Reader.ReadLine()`), recycled `sync.Pool` byte buffers, and reservoir sampling.
- **💰 Cloud Volume Cost & Ingestion Estimator (`pkg/cost`)**: Real-time cloud ingestion and retention cost modeling across AWS CloudWatch, Datadog, Splunk, and Google Cloud Logging with automated optimization recommendations.
- **✈️ Debug-on-Fault Flight Recorder (`pkg/dynlevel`)**: In-memory ring buffer capturing pre-incident debug events, instantly flushing rich forensic context only when errors/5xx occur while keeping healthy streams zero-noise.
- **🔒 Deterministic PII Pseudonymization (`pkg/pseudonym`)**: Cryptographic HMAC-SHA256 tokenization of User IDs, IP addresses, and emails with salt rotation support (GDPR Article 32 compliant).
- **🔬 Differential Incident Diagnosis (`pkg/diagnose`)**: Automated relative risk ranking comparing baseline vs incident windows to pinpoint root-cause error signatures and newly emergent regressions.
- **🌐 OpenTelemetry OTLP Logs Receiver (`pkg/otlprecv`)**: Direct parsing of native OpenTelemetry Collector `v1/logs` JSON push payloads and resource attributes into normalized models.
- **🗺️ Regional Latency & Route Heatmap (`pkg/geoheat`)**: Dynamic geo-latency aggregator calculating p50, p95, and degradation alerts across client countries.
- **🔍 Loki-Compatible LogQL Query Engine (`pkg/logql`)**: Execute stream selectors `{app="api", level="error"}` combined with pipeline filter stages `|= "timeout"`, `!= "ignored"`, regex matchers `|~ "^conn.*"`, and metric functions `rate()` and `count_over_time()`.
- **🌐 W3C Trace Context & OpenTelemetry (`pkg/tracectx`)**: Zero-dependency W3C `traceparent` (`00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01`) and `tracestate` header encoder/decoder.
- **🛡️ Resilient Circuit Breaker (`pkg/breaker`)**: Three-state circuit breaker (`CLOSED`, `OPEN`, `HALF-OPEN`) protecting remote logging pipelines from downstream failure storms.
- **⚡ Adaptive Rate Limiter & Token Bucket (`pkg/throttle`)**: Dynamic token bucket rate limiting and error-preserving probabilistic sampling.
- **📋 Streaming JSON Schema Validator & Drift Detector (`pkg/schemaval`)**: Assert required fields, type assertions, and detect schema drift across continuous streams.
- **🧠 Token TF-IDF Vectorizer & Cosine Similarity (`pkg/embed`)**: Pure Go n-gram character/token vectorizer computing semantic cosine similarity for grouping near-duplicate error signatures.
- **🌊 Multi-File Live Stream Multiplexer (`pkg/livetail`)**: Fan-in channel multiplexer streaming concurrent lines from multiple local log files with backpressure handling.
- **🔇 Sliding-Window Log Deduplication (`pkg/dedup`)**: Suppress log storm noise and bursts in sliding time windows with 90%+ volume reduction.
- **📊 Service Level Objective (SLO) & Error Budget (`pkg/slo`)**: Automated availability & latency SLI tracking with Google SRE multi-window burn rate alert triggers (1h, 6h).
- **🌐 Distributed Trace Graph & Critical Path (`pkg/correlate`)**: Reconstruct parent-child span trees across microservices and compute critical path latency bottlenecks.
- **📈 Streaming Anomaly Detection (`pkg/anomaly`)**: Real-time error rate spike detection using Exponentially Weighted Moving Average (EWMA) + dynamic Z-score and Interquartile Range (IQR).
- **🧩 Drain Unstructured Template Clustering (`pkg/cluster`)**: Real-time clustering of noisy unstructured log messages into parameterized templates.
- **🔐 Enterprise Secret Leak & Compliance Engine (`pkg/compliance`)**: Scans logs for sensitive data leaks (AWS keys, GitHub PATs, private keys) and HIPAA/GDPR violations.
- **🗄️ Database Slow Query Profiler (`pkg/slowlog`)**: Normalizes SQL queries into fingerprints and ranks top slowest database operations.
- **☸️ Kubernetes Metadata Enricher (`pkg/k8s`)**: Extract pod name, namespace, container name from container file paths.
- **📦 Transparent Gzip Decompression**: Automatically detects gzip magic bytes (`0x1f, 0x8b`) and streams `.gz` archives without extracting them to disk.
- **🧠 Auto-Detecting Multi-Format Parser**: Supports Structured JSON, Logfmt, Apache/Nginx Access Logs (CLF/Combined), and Syslog (RFC 3164 & 5424).
- **🛡️ Web Application Threat Detector (`pkg/security`)**: Real-time detection of SQL Injection, XSS, Path Traversal, and vulnerability scanner bots.
- **🌍 Pure-Go Offline IP Geolocation (`pkg/geodb`)**: In-memory radix CIDR prefix trie mapping IPs to countries and ASN with zero external dependencies.
- **🔄 HTTP Traffic Replayer & Load Simulator (`pkg/replay`)**: Replay parsed web logs against staging/canary servers with latency and status deviation comparison.
- **👥 User Journey & Session Reconstructor (`pkg/session`)**: Reconstructs user clickstreams and computes bounce rates, session durations, and page flows.
- **🔒 PII & Data Sanitization Engine (`pkg/mask`)**: Redact Credit Cards (Luhn algorithm verified), SSNs, Emails, Bearer Tokens, Passwords, and JWTs.
- **🎯 Pipe Query Expression Engine (`pkg/query`)**: Stream evaluations with pipe syntax: `logsift pipe "filter(status>=500) | topk(path, 5)" app.log`.
- **📊 Extreme Quantiles & Cardinality (`pkg/stats`)**: T-Digest quantile tracking ($p99.9, p99.99$) and HyperLogLog ($HLL$) unique visitor counting.
- **📈 Prometheus & OpenTelemetry Exporter (`pkg/export`)**: Native `/metrics` exposition and OTLP JSON model transformation.
- **🖥️ Interactive Bubble Tea TUI (`pkg/tui`)**: Interactive log viewer with Cost Dashboard, Flight Recorder view, LogQL modal, JSON tree inspector, live follow, regex highlight, bookmarks.
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
    (JSON / Logfmt / CLF / Syslog / CRI / Slowlog / UserAgent / W3C Trace / OTLP)
                                  |
                    +-------------+-------------+
                    |                           |
                    v                           v
          [ Filter / Masker ]          [ Streaming Aggregators ]
       (LogQL, Regex, Timestamps,    (EWMA Anomaly, Drain Cluster, SLO,
        Pseudonymization Engine)       Cost Estimator, GeoHeatmap,
                    |                  TF-IDF Vectorizer & Cosine Similarity)
       +------------+------------+              |
       |            |            |              |
       v            v            v              v
  [Security Audit] [Compliance] [Session Flow] [Crypto Tamper Chain]
       |            |            |              |
       +------------+------------+--------------+
                    |
                    v
  [Flight Recorder / Circuit Breaker / Remote Shipper / LiveTail Mux]
                    |
       +------------+------------+
       |                         |
       v                         v
  [Bubble Tea TUI]       [Differential Diagnosis / Replayer]
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

### 1. Cloud Ingestion & Retention Cost Estimator
```bash
# Model monthly log costs against Datadog, AWS CloudWatch, Splunk, or GCP
./logsift cost --provider=datadog production.log
```

### 2. Differential Incident Root-Cause Diagnosis
```bash
# Compare normal baseline traffic against outage period to pinpoint anomalous culprits
./logsift diagnose baseline.log incident.log
```

### 3. GDPR Cryptographic PII Pseudonymization
```bash
# Deterministically hash user IDs and emails with secret salt
./logsift pseudonym --salt=my-org-secret access.log
```

### 4. Regional Latency & Route Degradation Heatmap
```bash
# Aggregate regional latency metrics and surface degraded geography paths
./logsift geoheat access.log
```

### 5. Loki LogQL Stream Filtering
```bash
# Filter logs matching stream selector and error pipeline stages
./logsift logql '{level="error"} |= "timeout" != "ignored"' production.log
```

### 6. Multi-File Live Tail Multiplexer
```bash
# Multiplex and interleave log streams from multiple files concurrently
./logsift livetail /var/log/auth.log /var/log/api.log
```

### 7. Streaming JSON Schema Validation & Drift Check
```bash
# Audit incoming logs against required fields and report uncataloged schema drift
./logsift schemaval production.log
```

### 8. W3C Trace Context & OpenTelemetry Extraction
```bash
# Correlate distributed traceparent and tracestate headers across microservice logs
./logsift trace-w3c access.log
```

### 9. Sliding-Window Log Deduplication & Storm Suppression
```bash
# Suppress repeated errors and log storms in sliding 15s windows
./logsift dedup production.log
```

### 10. Service Level Objective (SLO) Health Check & Burn Rate
```bash
# Calculate availability SLI, latency SLI (<200ms), and error budget burn rate
./logsift slo access.log
```

### 11. Statistical Anomaly & Outlier Detection
```bash
# Detect sudden error rate spikes and latency outliers (IQR)
./logsift anomaly app.log
```

### 12. Unstructured Log Template Clustering (Drain)
```bash
# Cluster unstructured log messages into parameterized templates
./logsift cluster app.log
```

### 13. Compliance & Secret Leak Audit
```bash
# Scan logs for AWS keys, private SSH keys, passwords, and HIPAA/GDPR violations
./logsift compliance access.log
```

### 14. Database Slow Query Profiling
```bash
# Parse PostgreSQL or MySQL slow query logs and rank slowest query fingerprints
./logsift slowlog postgresql.log
```

### 15. Cryptographic Tamper-Evident Audit Sealing
```bash
# Seal log archives in chained SHA-256 blocks with HMAC signatures
./logsift sign production.log

# Verify tamper-proof manifest integrity
./logsift verify production.log
```

### 16. Kubernetes Container Log Enrichment
```bash
# Extract namespace, pod name, and container name from container log paths
./logsift k8s /var/log/pods/default_api-deployment-789_uid123/api-container/0.log
```

### 17. Security Threat Audit
```bash
# Scan web access logs for SQLi, XSS, Path Traversal, and Scanners
./logsift security access.log
```

### 18. Geographic Traffic Breakdown
```bash
# Analyze visitor geographic distribution offline
./logsift geo access.log
```

### 19. Pipeline Query Syntax
```bash
# Filter errors and aggregate top 5 endpoints
./logsift pipe "filter(status>=500) | topk(path, 5)" access.log
```

### 20. Traffic Replay
```bash
# Replay captured production requests against a staging server
./logsift replay --target http://staging.internal:8080 access.log
```

### 21. User Journey & Sessionization
```bash
# Reconstruct user sessions, bounce rates, and engagement
./logsift session access.log
```

### 22. Prometheus Metrics Exposition
```bash
# Export Prometheus metrics from log run
./logsift metrics access.log
```

### 23. Time Partitioning & Archiving
```bash
# Partition multi-gigabyte logs into hourly or daily compressed archives
./logsift archive --by date --out /var/archive access.log
```

### 24. Interactive Terminal UI (TUI)
```bash
./logsift tui app.log
```
- **`j` / `k`** or **`↓` / `↑`**: Navigate through log records
- **`/`**: Real-time regex / substring search filter
- **`1` - `5`**: Quick filter by log level (`1`=All, `2`=Debug, `3`=Info, `4`=Warn, `5`=Error)
- **`Enter` / `Space`**: Open structured record detail & JSON tree inspector modal
- **`Tab` / `s`**: Toggle Live Analytics dashboard view
- **`m`**: Toggle bookmark on current record
- **`f`**: Toggle live follow auto-scroller
- **`q`**: Quit

### 25. Standalone HTML & SVG Dashboard
```bash
./logsift html --out report.html app.log
```

### 26. SQLite SQL Analytics
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
| Cloud Cost Estimator Ingest | **~2,800,000 logs/sec** | **< 2 MB** |
| Flight Recorder Ring-Buffer Append | **~8,200,000 records/sec** | **< 1 MB** |
| HMAC PII Pseudonymization | **~1,200,000 tokens/sec** | **< 3 MB** |
| Differential Diagnosis Ranking | **~650,000 records/sec** | **< 5 MB** |
| W3C TraceParent Header Decoding | **~1,900,000 headers/sec** | **< 2 MB** |
| Loki LogQL Filter Pipeline | **~820,000 logs/sec** | **< 6 MB** |
| TF-IDF Token Vectorization | **~410,000 logs/sec** | **< 8 MB** |
| Circuit Breaker State Execution | **~14,000,000 ops/sec** | **< 1 MB** |
| Sliding-Window Log Deduplication | **~1,100,000 logs/sec** | **< 6 MB** |
| EWMA & Z-Score Anomaly Detection | **~920,000 logs/sec** | **< 8 MB** |
| T-Digest Quantile Compression ($p99.9$) | **~850,000 values/sec** | **< 3 MB** |
| Drain Template Log Clustering | **~190,000 logs/sec** | **< 16 MB** |
| Secret & Compliance Auditing | **~340,000 logs/sec** | **< 14 MB** |
| Cryptographic Tamper Chaining | **~260,000 logs/sec** | **< 15 MB** |
| Security Threat Rule Engine | **~310,000 logs/sec** | **< 14 MB** |
| IP Geolocation Prefix Lookup | **~1,400,000 lookups/sec** | **< 6 MB** |
| HyperLogLog Cardinality Counting | **~2,100,000 items/sec** | **< 1 MB** |
| PII Masking & Regex Sanitization | **~240,000 logs/sec** | **< 12 MB** |
| Streaming Reservoir Sampling ($p99$) | **~1,200,000 logs/sec** | **< 8 MB** |
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
