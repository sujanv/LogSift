# LogSift v1.5.0 Release Notes

LogSift v1.5.0 brings major enhancements across enterprise observability, security, and streaming data processing:

- **Sliding-Window Deduplication**: Suppress log storms and repeated error bursts.
- **SLO & Error Budget Tracking**: Automated availability and latency SLIs with Google SRE burn rate triggers.
- **Distributed Trace Graphs**: Reconstruct parent-child span graphs and identify critical path bottlenecks.
- **T-Digest Extreme Quantiles**: Extreme accuracy at tail latencies (p99.9, p99.99).
- **Chunk Bloom Filters**: Prune multi-gigabyte log partitions in sub-milliseconds.
- **Resilient Disk Spooler**: Zero log loss with encrypted local store-and-forward.
