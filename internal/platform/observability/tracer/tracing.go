// Package tracer provides a wrapper for configuring and managing the OpenTelemetry tracer.
package tracer

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/observability"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	// ErrFailedToInitializeOTLPExporter is logged when the OTLP trace exporter cannot be created.
	ErrFailedToInitializeOTLPExporter = "failed to initialize OTLP trace exporter"

	// WarnTracingDisabled is logged when tracing must be disabled because OTLP exporter bootstrap failed.
	WarnTracingDisabled = "tracing disabled because OTLP exporter initialization failed"

	// ErrFailedToShutdownTracerProvider is logged when the tracer provider fails to shut down.
	ErrFailedToShutdownTracerProvider = "failed to shutdown tracer provider"

	// WarnInvalidOTLPExporterTimeout is logged when the timeout string cannot be parsed and the default is used.
	WarnInvalidOTLPExporterTimeout = "invalid OTLP exporter timeout, using default"

	// CompressionGzip is the string value that enables gzip compression for the OTLP HTTP exporter.
	CompressionGzip = "gzip"
)

// InitTracer initializes the OpenTelemetry tracer using the provided configuration,
// installs it as the global tracer provider, and returns a cleanup function to shut it down gracefully.
func InitTracer(cfg *config.Config, logger logger.ContextLogger) func() {
	if !cfg.Observability.OtelExporterEnabled {
		if logger != nil {
			logger.Infow("tracing disabled via OTEL_EXPORTER_ENABLED=false")
		}
		return func() {}
	}

	exporter, err := buildOTLPExporter(cfg, logger)
	if err != nil {
		logger.Errorw(ErrFailedToInitializeOTLPExporter, commonkeys.Error, err)
		logger.Warnw(WarnTracingDisabled, commonkeys.Error, err)
		return func() {}
	}

	resources := buildResource(cfg)

	// configure sampler from environment variables if provided
	sampler := configureSamplerFromEnv()

	providerOpts := []trace.TracerProviderOption{
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	}
	if sampler != nil {
		providerOpts = append(providerOpts, trace.WithSampler(sampler))
	}

	traceProvider := trace.NewTracerProvider(providerOpts...)

	otel.SetTracerProvider(traceProvider)

	// Set global propagators (W3C TraceContext + Baggage) so spans are propagated across services
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			logger.Errorw(ErrFailedToShutdownTracerProvider, commonkeys.Error, err)
		}
	}
}

// buildOTLPExporter creates and configures the OTLP HTTP trace exporter from cfg.
func buildOTLPExporter(cfg *config.Config, logger logger.ContextLogger) (trace.SpanExporter, error) {
	// Normalize endpoint: accept host:port or http(s)://host:port
	endpointVal := cfg.Observability.OtelExporterOTLPEndpoint
	normalized, err := observability.NormalizeEndpoint(endpointVal)
	if err != nil {
		if logger != nil {
			logger.Warnw("invalid OTEL_EXPORTER_OTLP_ENDPOINT, using raw value", commonkeys.Error, err)
		}
		normalized = endpointVal
	}
	// If normalized includes a scheme, extract host:port for exporter.WithEndpoint
	endpointForExporter := normalized
	if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") {
		if u, err := url.Parse(normalized); err == nil {
			if u.Host != "" {
				endpointForExporter = u.Host
			}
		}
	}

	opts := []otlptrace.Option{
		otlptrace.WithEndpoint(endpointForExporter),
	}

	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlptrace.WithInsecure())
	}

	if cfg.Observability.OtelExporterTimeout != "" {
		if timeout, err := time.ParseDuration(cfg.Observability.OtelExporterTimeout); err == nil {
			opts = append(opts, otlptrace.WithTimeout(timeout))
		} else if logger != nil {
			logger.Warnw(WarnInvalidOTLPExporterTimeout, commonkeys.Error, err)
		}
	}

	if cfg.Observability.OtelExporterCompression == CompressionGzip {
		opts = append(opts, otlptrace.WithCompression(otlptrace.GzipCompression))
	}

	headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders)
	if len(headers) > 0 {
		opts = append(opts, otlptrace.WithHeaders(headers))
	}

	exporter, err := otlptrace.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

// buildResource constructs service resource attributes used by the tracer provider.
func buildResource(cfg *config.Config) *resource.Resource {
	// Resolve host name and a stable instance id for this process
	hostname, _ := os.Hostname()
	instanceID := uuid.NewString()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.Observability.OtelServiceName),
		semconv.ServiceVersionKey.String(cfg.Observability.OtelServiceVersion),
		attribute.String("deployment.environment", cfg.General.Env),
		attribute.String("host.name", hostname),
		attribute.String("service.instance.id", instanceID),
	)
}

// configureSamplerFromEnv reads OTEL_TRACES_SAMPLER* env vars and returns a configured sampler or nil.
func configureSamplerFromEnv() trace.Sampler {
	var sampler trace.Sampler
	samplerName := strings.TrimSpace(strings.ToLower(os.Getenv("OTEL_TRACES_SAMPLER")))
	samplerArg := strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG"))
	if samplerName == "" {
		return nil
	}

	switch samplerName {
	case "always_on", "always":
		sampler = trace.AlwaysSample()
	case "always_off", "never":
		sampler = trace.NeverSample()
	case "traceidratio", "traceidratiobased", "ratio":
		if f, err := strconv.ParseFloat(samplerArg, 64); err == nil {
			sampler = trace.ParentBased(trace.TraceIDRatioBased(f))
		}
	case "parentbased", "parentbased_traceidratio":
		if f, err := strconv.ParseFloat(samplerArg, 64); err == nil {
			sampler = trace.ParentBased(trace.TraceIDRatioBased(f))
		}
	default:
		// unknown sampler — ignore and fall back to default
	}

	return sampler
}
