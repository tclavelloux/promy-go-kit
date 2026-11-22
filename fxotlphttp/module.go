// Package fxotlphttp exports Yokai traces over OTLP/HTTP.
//
// Yokai's trace module (github.com/ankorstore/yokai/fxtrace) only exports over plaintext
// gRPC with no auth headers, and maps any unknown modules.trace.processor.type to noop.
// Railway native Traces and Grafana Cloud both take OTLP/HTTP, so this package fills the
// gap (decision: https://github.com/tclavelloux/promy-steering/issues/16).
//
// Add [Module] next to Yokai's modules and set modules.trace.processor.type to otlp-http.
// fxtrace then registers its noop processor, and [Module] appends an OTLP/HTTP batch span
// processor. The exporter reads only the standard OTEL_EXPORTER_OTLP_* variables
// (endpoint, headers, protocol, timeout), so the sink changes by configuration, never code.
//
// Scope: promy-go-kit holds Yokai-aware shared Go code. promy-event-bus stays framework-free.
package fxotlphttp

import (
	"context"
	"os"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

// ProcessorType is the modules.trace.processor.type value that enables the OTLP/HTTP exporter.
const ProcessorType = "otlp-http"

// Module decorates Yokai's [trace.TracerProviderFactory] with an OTLP/HTTP span processor.
//
// It is a root-level option, not an fx.Module: fx scopes a decorator to the module declaring
// it, and the factory is consumed inside Yokai's trace module.
var Module = fx.Options(
	fx.Decorate(DecorateTracerProviderFactory),
)

// DecorateParams are the dependencies of [DecorateTracerProviderFactory].
type DecorateParams struct {
	fx.In
	Factory trace.TracerProviderFactory
	Config  *config.Config
	Logger  *log.Logger
}

// DecorateTracerProviderFactory wraps the factory when modules.trace.processor.type is
// otlp-http and an OTLP endpoint is set. Otherwise it returns the factory unchanged.
// It never fails: tracing must not stop a service from starting.
func DecorateTracerProviderFactory(p DecorateParams) trace.TracerProviderFactory {
	if !strings.EqualFold(p.Config.GetString("modules.trace.processor.type"), ProcessorType) {
		return p.Factory
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") == "" && os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		p.Logger.Warn().Str("module", "fxotlphttp").Msg("otlp-http trace processor without OTEL_EXPORTER_OTLP_ENDPOINT, spans are not exported")

		return p.Factory
	}

	return &Factory{base: p.Factory, logger: p.Logger}
}

// Factory is a [trace.TracerProviderFactory] that appends an OTLP/HTTP batch span processor.
type Factory struct {
	base   trace.TracerProviderFactory
	logger *log.Logger
}

var _ trace.TracerProviderFactory = (*Factory)(nil)

// newExporter is a seam for tests; it reads the OTEL_EXPORTER_OTLP_* variables.
var newExporter = func(ctx context.Context) (otelsdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx)
}

// Create builds the tracer provider with the OTLP/HTTP span processor appended to options.
// If the exporter cannot be built, it logs a warning and builds the provider without it.
func (f *Factory) Create(options ...trace.TracerProviderOption) (*otelsdktrace.TracerProvider, error) {
	exporter, err := newExporter(context.Background())
	if err != nil {
		f.logger.Warn().Err(err).Str("module", "fxotlphttp").Msg("cannot create otlp-http trace exporter, spans are not exported")

		return f.base.Create(options...)
	}

	return f.base.Create(append(options, trace.WithSpanProcessor(otelsdktrace.NewBatchSpanProcessor(exporter)))...)
}
