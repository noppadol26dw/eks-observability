package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "github.com/eks-observability/app/handler"
const tracerName = "github.com/eks-observability/app/handler"

var (
	requestCount metric.Int64Counter
	requestOnce  sync.Once
)

func getRequestCount() metric.Int64Counter {
	requestOnce.Do(func() {
		meter := otel.Meter(meterName)
		requestCount, _ = meter.Int64Counter(
			"http_requests_total",
			metric.WithDescription("Total HTTP requests"),
		)
	})
	return requestCount
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer(tracerName)

	ctx, span := tracer.Start(ctx, "Hello")
	defer span.End()

	if c := getRequestCount(); c != nil {
		c.Add(ctx, 1, metric.WithAttributes(
			attribute.String("handler", "hello"),
			attribute.String("method", r.Method),
		))
	}

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", "/hello"),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "hello",
		"time":   time.Now().Format(time.RFC3339),
	})
}
