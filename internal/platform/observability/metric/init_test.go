//nolint:testpackage // tests need access to unexported Init wiring behavior.
package metric

import (
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestInitOtelMetricsReturnsCleanup(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Env: "test",
		},
		Observability: config.ObservabilityConfig{
			OtelExporterEnabled:      true,
			OtelExporterOTLPEndpoint: "localhost:4318",
			OtelServiceName:          "aion-test",
			OtelServiceVersion:       "v1",
			OtelExporterInsecure:     true,
			OtelExporterTimeout:      "1s",
		},
	}

	cleanup := InitOtelMetrics(cfg, noopLogger{})
	require.NotNil(t, cleanup)
	cleanup()
}

func TestInitOtelMetricsDisabledViaKillSwitch(t *testing.T) {
	cfg := &config.Config{
		Observability: config.ObservabilityConfig{
			OtelExporterEnabled: false,
		},
	}

	cleanup := InitOtelMetrics(cfg, noopLogger{})
	require.NotNil(t, cleanup)
	cleanup()
}
