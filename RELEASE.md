# LogSift v1.7.0 Release Notes

LogSift v1.7.0 introduces FinOps cloud cost intelligence, zero-noise debug flight recording, and differential root-cause incident diagnosis:

- **Multi-Cloud Cost Estimator**: Calculate real-time log ingestion and retention costs across AWS CloudWatch, Datadog, Splunk, and GCP with automated cost-saving recommendations.
- **Flight Recorder Debug-on-Fault**: Ring buffer keeping debug logs in memory and automatically flushing preceding context only when an error occurs, eliminating healthy debug storage costs.
- **Deterministic HMAC Pseudonymization**: GDPR-compliant PII tokenization enabling user correlation without storing raw personal data.
- **Differential Incident Diagnosis**: Relative risk frequency comparison between baseline and incident windows to rank root-cause culprit signatures.
- **OpenTelemetry OTLP Receiver**: Native parsing of OpenTelemetry Collector `v1/logs` push payloads.
- **Regional Latency Heatmap**: Multi-region route performance tracking with edge degradation alerts.
