//nolint:testpackage // tests cover unexported tracer builder helpers.
package tracer

import (
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBuildOTLPExporter(t *testing.T) {
	t.Run("invalid timeout warns and exporter is still created", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		loggerMock := mocks.NewMockContextLogger(ctrl)
		loggerMock.EXPECT().Warnw(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		cfg := &config.Config{
			Observability: config.ObservabilityConfig{
				OtelExporterOTLPEndpoint: "localhost:4318",
				OtelExporterTimeout:      "not-a-duration",
				OtelExporterCompression:  CompressionGzip,
				OtelExporterHeaders:      "x-api-key=123,env=test",
				OtelExporterInsecure:     true,
			},
		}

		exporter, err := buildOTLPExporter(cfg, loggerMock)
		require.NoError(t, err)
		require.NotNil(t, exporter)
	})

	t.Run("invalid endpoint format warns and still builds exporter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		loggerMock := mocks.NewMockContextLogger(ctrl)
		loggerMock.EXPECT().Warnw(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		cfg := &config.Config{
			Observability: config.ObservabilityConfig{
				OtelExporterOTLPEndpoint: "http://[::1",
				OtelExporterInsecure:     true,
			},
		}

		exporter, err := buildOTLPExporter(cfg, loggerMock)
		require.NoError(t, err)
		require.NotNil(t, exporter)
	})
}

func TestInitTracerReturnsCleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	loggerMock := mocks.NewMockContextLogger(ctrl)
	loggerMock.EXPECT().Errorw(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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
		},
	}

	cleanup := InitTracer(cfg, loggerMock)
	require.NotNil(t, cleanup)
	cleanup()
}

func TestInitTracerDisabledViaKillSwitch(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	loggerMock := mocks.NewMockContextLogger(ctrl)
	loggerMock.EXPECT().Infow(gomock.Any()).Times(1)

	cfg := &config.Config{
		Observability: config.ObservabilityConfig{
			OtelExporterEnabled: false,
		},
	}

	cleanup := InitTracer(cfg, loggerMock)
	require.NotNil(t, cleanup)
	cleanup()
}
