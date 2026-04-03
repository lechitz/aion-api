//nolint:testpackage // tests need access to unexported option builders for coverage.
package metric

import (
	"context"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/stretchr/testify/require"
)

type noopLogger struct{}

func (noopLogger) Infof(string, ...any)                      {}
func (noopLogger) Errorf(string, ...any)                     {}
func (noopLogger) Debugf(string, ...any)                     {}
func (noopLogger) Warnf(string, ...any)                      {}
func (noopLogger) Infow(string, ...any)                      {}
func (noopLogger) Errorw(string, ...any)                     {}
func (noopLogger) Debugw(string, ...any)                     {}
func (noopLogger) Warnw(string, ...any)                      {}
func (noopLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopLogger) DebugwCtx(context.Context, string, ...any) {}

func TestComputeEndpointForExporter(t *testing.T) {
	require.Empty(t, computeEndpointForExporter("", nil))
	require.Equal(t, "aion-dev-otel-collector:4318", computeEndpointForExporter("http://aion-dev-otel-collector:4318", nil))
	require.Equal(t, "aion-dev-otel-collector:4318", computeEndpointForExporter("aion-dev-otel-collector:4318", nil))
	require.Equal(t, "http://[::1", computeEndpointForExporter("http://[::1", noopLogger{}))
}

func TestBuildMetricOptions(t *testing.T) {
	cfg := &config.Config{
		Observability: config.ObservabilityConfig{
			OtelExporterOTLPEndpoint: "aion-dev-otel-collector:4318",
			OtelExporterInsecure:     true,
			OtelExporterTimeout:      "5s",
			OtelExporterCompression:  "gzip",
			OtelExporterHeaders:      "x-api-key=abc",
		},
	}
	opts := buildMetricOptions(cfg, noopLogger{})
	require.Len(t, opts, 5)

	cfg.Observability.OtelExporterTimeout = "invalid-duration"
	opts = buildMetricOptions(cfg, noopLogger{})
	require.Len(t, opts, 4)
}
