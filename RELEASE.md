# LogSift v1.6.0 Release Notes

LogSift v1.6.0 introduces cloud-native observability, distributed tracing, and resilient log processing features:

- **Loki LogQL Engine**: Filter log streams and compute stream rates with Loki-compatible queries.
- **W3C Trace Context & OpenTelemetry**: Native extraction of W3C `traceparent` and `tracestate` headers.
- **Resilient Circuit Breaker**: Closed/Open/Half-Open state machine safeguarding remote log shippers.
- **Adaptive Rate Limiting & Sampling**: Token bucket and error-preserving probabilistic sampling.
- **Streaming JSON Schema Validation & Drift**: Enforce log contracts and automatically surface schema changes.
- **TF-IDF Semantic Vectorization**: Cosine similarity clustering for grouping near-duplicate error logs.
- **Live Stream Multiplexing**: Interleave concurrent logs from multiple sources with backpressure.
