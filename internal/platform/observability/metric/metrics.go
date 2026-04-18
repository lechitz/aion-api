// Package metric provides utilities for initializing and managing OpenTelemetry metrics.
package metric

import (
	"context"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/observability"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	// ErrFailedToInitializeOTLPMetricsExporter is logged when the OTLP metrics exporter cannot be created.
	ErrFailedToInitializeOTLPMetricsExporter = "failed to initialize OTLP metric exporter"

	// WarnMetricsDisabled is logged when metrics must be disabled because OTLP exporter bootstrap failed.
	WarnMetricsDisabled = "metrics disabled because OTLP metric exporter initialization failed"

	// ErrInvalidOTELExporterTimeout is logged when the timeout string cannot be parsed as a valid duration.
	ErrInvalidOTELExporterTimeout = "invalid OTLP exporter timeout"

	// ErrToShutDownOTELMetrics is logged when shutting down the metrics provider fails.
	ErrToShutDownOTELMetrics = "failed to shutdown OTLP metrics provider"
)

// InitOtelMetrics sets up the OpenTelemetry MeterProvider using the given configuration,
// installs it as the global provider, and returns a cleanup function to gracefully shut it down.
func InitOtelMetrics(cfg *config.Config, logger logger.ContextLogger) func() {
	if !cfg.Observability.OtelExporterEnabled {
		if logger != nil {
			logger.Infow("metrics disabled via OTEL_EXPORTER_ENABLED=false")
		}
		return func() {}
	}

	opts := buildMetricOptions(cfg, logger)

	exporter, err := otlpmetrichttp.New(context.Background(), opts...)
	if err != nil {
		if logger != nil {
			logger.Errorw(ErrFailedToInitializeOTLPMetricsExporter, commonkeys.Error, err)
			logger.Warnw(WarnMetricsDisabled, commonkeys.Error, err)
		}
		return func() {}
	}

	// Build resource with common attributes
	hostname, _ := os.Hostname()
	instanceID := uuid.NewString()

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Observability.OtelServiceName),
			semconv.ServiceVersionKey.String(cfg.Observability.OtelServiceVersion),
			attribute.String("deployment.environment", cfg.General.Env),
			attribute.String("host.name", hostname),
			attribute.String("service.instance.id", instanceID),
		)),
	)

	otel.SetMeterProvider(provider)

	return func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			if logger != nil {
				logger.Errorw(ErrToShutDownOTELMetrics, commonkeys.Error, err)
			}
		}
	}
}

// buildMetricOptions constructs OTLP metric exporter options from config.
func buildMetricOptions(cfg *config.Config, logger logger.ContextLogger) []otlpmetrichttp.Option {
	endpointForExporter := computeEndpointForExporter(cfg.Observability.OtelExporterOTLPEndpoint, logger)

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpointForExporter),
	}

	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if cfg.Observability.OtelExporterTimeout != "" {
		if timeout, err := time.ParseDuration(cfg.Observability.OtelExporterTimeout); err == nil {
			opts = append(opts, otlpmetrichttp.WithTimeout(timeout))
		} else if logger != nil {
			logger.Warnw(ErrInvalidOTELExporterTimeout, commonkeys.Error, err)
		}
	}

	if cfg.Observability.OtelExporterCompression == "gzip" {
		opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	}

	headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders)
	if len(headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(headers))
	}

	return opts
}

// computeEndpointForExporter normalizes the configured endpoint and returns a host:port
// value suitable for otlpmetrichttp.WithEndpoint. It accepts either host:port or full URL.
func computeEndpointForExporter(raw string, logger logger.ContextLogger) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}

	normalized, err := observability.NormalizeEndpoint(raw)
	if err != nil {
		if logger != nil {
			logger.Warnw("invalid OTEL_EXPORTER_OTLP_ENDPOINT, using raw value", commonkeys.Error, err)
		}
		return raw
	}

	// If it includes a scheme, extract host:port
	if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") {
		if u, err := url.Parse(normalized); err == nil {
			if u.Host != "" {
				return u.Host
			}
		}
	}
	return normalized
}
