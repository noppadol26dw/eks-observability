package telemetry

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const shutdownTimeout = 15 * time.Second

func Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("tracer provider shutdown: %v", err)
		}
	}

	if mp := otel.GetMeterProvider(); mp != nil {
		if m, ok := mp.(*metric.MeterProvider); ok {
			if err := m.Shutdown(ctx); err != nil {
				log.Printf("meter provider shutdown: %v", err)
			}
		}
	}

	if lp := global.GetLoggerProvider(); lp != nil {
		if l, ok := lp.(*sdklog.LoggerProvider); ok {
			if err := l.Shutdown(ctx); err != nil {
				log.Printf("logger provider shutdown: %v", err)
			}
		}
	}
}
