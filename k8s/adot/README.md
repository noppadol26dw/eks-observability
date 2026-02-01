# ADOT Collector on EKS

Install ADOT via the [EKS Add-on](https://aws-otel.github.io/docs/getting-started/adot-eks-add-on/) and configure pipelines for OTLP.

1. Install the ADOT EKS add-on for your cluster.
2. Create an ADOT Collector config (ConfigMap or CRD) that:
   - Receives OTLP gRPC on port 4317.
   - Exports traces to AWS X-Ray.
   - Exports metrics to CloudWatch and/or Amazon Managed Prometheus (AMP).
   - Exports logs to CloudWatch.
3. Deploy the Collector (DaemonSet or Deployment) in the `observability` namespace.
4. Ensure the app Deployment's `OTEL_EXPORTER_OTLP_ENDPOINT` points to the Collector Service (e.g. `http://adot-collector.observability.svc.cluster.local:4317`).

See:
- [ADOT EKS Add-on](https://aws-otel.github.io/docs/getting-started/adot-eks-add-on/)
- [ADOT Collector Configuration](https://aws-otel.github.io/docs/getting-started/collector/)
- [Container Insights with ADOT](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-EKS-otel.html)
