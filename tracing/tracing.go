package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var DefaultTracing = &defaultTracing

var defaultTracing Tracing

type Tracing struct {
	exporter sdktrace.SpanExporter
	Provider *sdktrace.TracerProvider
}

func New(exporter sdktrace.SpanExporter, serviceName string) *Tracing {
	if exporter == nil {
		return nil
	}

	tracer := Tracing{exporter: exporter}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tracer.Provider = sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
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

func Trace(ctx context.Context, spanName string) trace.Span {
	tr := otel.GetTracerProvider().Tracer("example/server")
	_, span := tr.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
	return span
}
