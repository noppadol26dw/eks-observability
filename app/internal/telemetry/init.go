package telemetry

import (
	"context"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/detectors/aws/eks"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const defaultEndpoint = "localhost:4317"
const defaultServiceName = "eks-observability-app"

func otelEndpoint() string {
	s := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if s == "" {
		return defaultEndpoint
	}
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	return s
}

func otelInsecure() bool {
	return os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true"
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	detector := eks.NewResourceDetector()
	detected, err := detector.Detect(ctx)
	if err != nil {
		detected = resource.Empty()
	}
	defaultRes, err := resource.Merge(
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(defaultServiceName),
		),
		detected,
	)
	if err != nil {
		return nil, err
	}
	return defaultRes, nil
}

func Init(ctx context.Context) error {
	endpoint := otelEndpoint()
	insecure := otelInsecure()

	res, err := newResource(ctx)
	if err != nil {
		return err
	}

	// Trace
	traceOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
	}
	if insecure {
		traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
	}
	traceExporter, err := otlptracegrpc.New(ctx, traceOpts...)
	if err != nil {
		return err
	}

	idg := xray.NewIDGenerator()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	// Metrics
	metricOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}
	if insecure {
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	}
	metricExporter, err := otlpmetricgrpc.New(ctx, metricOpts...)
	if err != nil {
		log.Printf("metric exporter init failed: %v", err)
	} else {
		mp := metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		)
		otel.SetMeterProvider(mp)
	}

	// Logs
	logOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
	}
	if insecure {
		logOpts = append(logOpts, otlploggrpc.WithInsecure())
	}
	logExporter, err := otlploggrpc.New(ctx, logOpts...)
	if err != nil {
		log.Printf("log exporter init failed: %v", err)
	} else {
		processor := sdklog.NewBatchProcessor(logExporter)
		lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(processor), sdklog.WithResource(res))
		global.SetLoggerProvider(lp)
	}

	return nil
}
