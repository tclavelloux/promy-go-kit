package fxotlphttp

import (
	"context"

	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SetNewExporter replaces the exporter constructor and returns a restore function.
func SetNewExporter(fn func(ctx context.Context) (otelsdktrace.SpanExporter, error)) func() {
	previous := newExporter
	newExporter = fn

	return func() { newExporter = previous }
}
