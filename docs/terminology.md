# Observability Terminology

Short reference for terms used in this project. See [OpenTelemetry Glossary](https://opentelemetry.io/docs/concepts/glossary/) and [Observability Primer](https://opentelemetry.io/docs/concepts/observability-primer/) for more.

| Term | Meaning |
|------|--------|
| **Observability** | Ability to understand a system from the outside and answer "what happened / why" without knowing internals. |
| **Telemetry** | Data emitted by the system or app: logs, metrics, traces. |
| **Logs** | Timestamped records (text or JSON) for specific events; good for detail and errors. |
| **Metrics** | Numeric aggregates over time (e.g. error rate, latency, req/s); good for alerts and dashboards. |
| **Traces** | Path of a request across services; each step is a **Span**. |
| **Span** | Single unit of work (e.g. one HTTP request, one DB call); has name, timing, attributes. |
| **OTLP** | OpenTelemetry Protocol — standard way to send telemetry to a collector or backend. |
| **Signals** | OTel term for outputs that describe system activity: traces, metrics, logs, and baggage. |
| **Receiver** | Collects telemetry from sources. Can be pull-based (scrape) or push-based (listen). Examples: OTLP, file log. |
| **Processor** | Transforms or filters telemetry (e.g. batch, add attributes, sampling). |
| **Connector** | Consumes one signal type and produces another (e.g. span-to-metrics). |
| **Extension** | Adds capabilities to the collector (e.g. storage for file offsets, health checks). |
| **Pipeline** | Flow of data: receivers → processors → exporters. Defined under `service.pipelines` in config. |
| **Baggage** | Context passed between signals to correlate logs, traces, and metrics. |
| **Exporter** | Component that sends data (trace/metric/log) to a backend (e.g. OTLP exporter). |
| **TracerProvider / MeterProvider / LoggerProvider** | Global setup for creating traces, metrics, and logs in the app. |
| **Propagator** | How trace context (trace ID, etc.) is passed across services (e.g. X-Ray propagator). |
| **Resource** | Environment info where the app runs (service name, EKS cluster, pod); **Resource detector** (e.g. EKS detector) discovers it. |
| **SLI / SLO** | SLI = measure of service behavior; SLO = target tied to business (e.g. 99% latency &lt; 200ms). |
| **Cardinality** | Number of unique values for an attribute; too high can make metrics/logs expensive and slow. |

For the full OTEL + ADOT architecture, see [architecture](architecture.md).

## Links

- [Collector Configuration](https://opentelemetry.io/docs/collector/configuration/) — receivers, processors, exporters, pipelines
- [Collector Contrib (receivers, exporters)](https://github.com/open-telemetry/opentelemetry-collector-contrib) — file log, Kafka, many backends
- [OpenTelemetry Glossary](https://opentelemetry.io/docs/concepts/glossary/)
- [OpenTelemetry Observability Primer](https://opentelemetry.io/docs/concepts/observability-primer/)
- [ADOT Go SDK Manual Instrumentation](https://aws-otel.github.io/docs/getting-started/go-sdk/manual-instr/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
- [ADOT EKS Add-on](https://aws-otel.github.io/docs/getting-started/adot-eks-add-on/)
- [ADOT Collector](https://aws-otel.github.io/docs/getting-started/collector/)
- [Container Insights with ADOT](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-EKS-otel.html)
