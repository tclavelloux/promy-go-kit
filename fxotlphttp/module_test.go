package fxotlphttp_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/tclavelloux/promy-go-kit/fxotlphttp"
)

const warnNoEndpoint = "otlp-http trace processor without OTEL_EXPORTER_OTLP_ENDPOINT, spans are not exported"

type receivedRequest struct {
	path          string
	contentType   string
	authorization string
	bodySize      int
}

type receiver struct {
	mu       sync.Mutex
	requests []receivedRequest
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	r.mu.Lock()
	r.requests = append(r.requests, receivedRequest{
		path:          req.URL.Path,
		contentType:   req.Header.Get("Content-Type"),
		authorization: req.Header.Get("Authorization"),
		bodySize:      len(body),
	})
	r.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (r *receiver) received() []receivedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]receivedRequest(nil), r.requests...)
}

type deps struct {
	factory  trace.TracerProviderFactory
	provider oteltrace.TracerProvider
	logs     logtest.TestLogBuffer
}

func startApp(t *testing.T, processorType string) deps {
	t.Helper()

	t.Setenv("APP_CONFIG_PATH", "testdata")
	t.Setenv("APP_ENV", "dev")
	t.Setenv("TRACE_PROCESSOR_TYPE", processorType)

	var d deps

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxotlphttp.Module,
		fx.Populate(&d.factory, &d.provider, &d.logs),
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return d
}

func TestModule_ExportsSpansOverHTTPWithHeaders(t *testing.T) {
	recv := &receiver{}
	server := httptest.NewServer(recv)
	t.Cleanup(server.Close)

	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", server.URL)
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "Authorization=Basic dGVzdDp0b2tlbg==")

	d := startApp(t, "otlp-http")

	assert.IsType(t, &fxotlphttp.Factory{}, d.factory)

	_, span := d.provider.Tracer("fxotlphttp-test").Start(context.Background(), "Test.Span")
	span.End()

	sdkProvider, ok := d.provider.(*otelsdktrace.TracerProvider)
	require.True(t, ok)
	require.NoError(t, sdkProvider.ForceFlush(context.Background()))

	requests := recv.received()
	require.Len(t, requests, 1)
	assert.Equal(t, "/v1/traces", requests[0].path)
	assert.Equal(t, "application/x-protobuf", requests[0].contentType)
	assert.Equal(t, "Basic dGVzdDp0b2tlbg==", requests[0].authorization)
	assert.Positive(t, requests[0].bodySize)
}

func TestModule_TracesEndpointAloneEnablesExport(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "http://127.0.0.1:4318/v1/traces")

	d := startApp(t, "OTLP-HTTP")

	assert.IsType(t, &fxotlphttp.Factory{}, d.factory)
}

func TestModule_NoEndpointWarnsAndAddsNothing(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")

	d := startApp(t, "otlp-http")

	assert.IsType(t, &trace.DefaultTracerProviderFactory{}, d.factory)
	logtest.AssertHasLogRecord(t, d.logs, map[string]any{
		"level":   "warn",
		"module":  "fxotlphttp",
		"message": warnNoEndpoint,
	})
}

func TestModule_OtherProcessorTypesUntouched(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:4318")

	for _, processorType := range []string{"noop", "stdout", "otlp-grpc", ""} {
		t.Run(processorType, func(t *testing.T) {
			d := startApp(t, processorType)

			assert.IsType(t, &trace.DefaultTracerProviderFactory{}, d.factory)
			logtest.AssertHasNotLogRecord(t, d.logs, map[string]any{"module": "fxotlphttp"})
		})
	}
}

func TestFactory_ExporterErrorFallsBackWithoutFailing(t *testing.T) {
	restore := fxotlphttp.SetNewExporter(func(context.Context) (otelsdktrace.SpanExporter, error) {
		return nil, errors.New("boom")
	})
	t.Cleanup(restore)

	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:4318")

	d := startApp(t, "otlp-http")

	assert.NotNil(t, d.provider)
	logtest.AssertHasLogRecord(t, d.logs, map[string]any{
		"level":   "warn",
		"module":  "fxotlphttp",
		"message": "cannot create otlp-http trace exporter, spans are not exported",
	})
}
