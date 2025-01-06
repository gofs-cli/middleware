// TODO: document that this can panic on init
package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var DefaultTracing = &defaultTracing

var defaultTracing Tracing

type Tracing struct {
	exporter trace.SpanExporter
	Provider *trace.TracerProvider
}

func New(exporter trace.SpanExporter, serviceName string) *Tracing {
	if exporter == nil {
		return nil
	}

	tracer := Tracing{exporter: exporter}
	batcher := trace.NewBatchSpanProcessor(exporter)

	tracer.Provider = trace.NewTracerProvider(
		trace.WithSpanProcessor(batcher),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(tracer.Provider)
	return &tracer
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewMiddleware(r.Method+" "+r.URL.Path)(h).ServeHTTP(w, r)
	})
}
