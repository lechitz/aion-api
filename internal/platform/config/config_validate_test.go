package config_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func baseConfig() config.Config {
	return config.Config{
		ServerHTTP: config.ServerHTTP{
			Host:              "0.0.0.0",
			Port:              "5001",
			Context:           "/aion",
			APIRoot:           "/api/v1",
			SwaggerMountPath:  "/swagger",
			DocsAliasPath:     "/docs",
			HealthRoute:       "/health",
			ReadTimeout:       2 * time.Second,
			WriteTimeout:      2 * time.Second,
			ReadHeaderTimeout: time.Second,
			IdleTimeout:       time.Second,
			MaxHeaderBytes:    1024,
		},
		ServerGraphql: config.ServerGraphql{
			Path: "/graphql",
		},
		Cache: config.CacheConfig{
			Addr:     "aion-dev-redis:6379",
			PoolSize: 2,
		},
		DB: config.DBConfig{
			Type:            "postgres",
			Host:            "localhost",
			Port:            "5432",
			Name:            "aion",
			User:            "aion",
			Password:        "aion",
			SSLMode:         "disable",
			TimeZone:        "UTC",
			MaxOpenConns:    1,
			MaxIdleConns:    0,
			ConnMaxLifetime: 0,
			RetryInterval:   time.Second,
			MaxRetries:      1,
		},
		Observability: config.ObservabilityConfig{
			OtelExporterOTLPEndpoint: "aion-dev-otel-collector:4318",
			OtelExporterCompression:  "none",
		},
		Kafka: config.KafkaConfig{
			Brokers:                     "aion-dev-kafka:9092",
			RecordEventsTopic:           "aion.record.events.v1",
			RecordProjectionEventsTopic: "aion.record_projection.events.v1",
		},
		Outbox: config.OutboxConfig{
			PublishEnabled:  true,
			PublishInterval: 2 * time.Second,
			BatchSize:       50,
		},
		Realtime: config.RealtimeConfig{
			Enabled:             true,
			StreamPath:          "/events/stream",
			HeartbeatInterval:   15 * time.Second,
			SubscriberBuffer:    32,
			ConsumerGroupPrefix: "aion-api-realtime",
		},
		Application: config.Application{
			ContextRequest: config.MinContextRequest,
			Timeout:        config.MinShutdownTimeout,
		},
	}
}

func TestConfigValidate_Success(t *testing.T) {
	cfg := baseConfig()
	require.NoError(t, cfg.Validate())
}

func TestConfigValidate_HTTPAndGraphQLErrors(t *testing.T) {
	cfg := baseConfig()
	cfg.ServerHTTP.Host = ""
	require.EqualError(t, cfg.Validate(), config.ErrHTTPHostRequired)

	cfg = baseConfig()
	cfg.ServerHTTP.APIRoot = "/"
	require.EqualError(t, cfg.Validate(), config.ErrHTTPAPIRootTooShort)

	cfg = baseConfig()
	cfg.ServerHTTP.SwaggerMountPath = "/swagger/"
	require.EqualError(t, cfg.Validate(), config.ErrHTTPSwaggerMountMustNotEndSlash)

	cfg = baseConfig()
	cfg.ServerGraphql.Path = "graphql"
	require.EqualError(t, cfg.Validate(), config.ErrGraphqlPathMustStart)
}

func TestConfigValidate_HTTPPathEdgeCases(t *testing.T) {
	t.Run("context root is allowed", func(t *testing.T) {
		cfg := baseConfig()
		cfg.ServerHTTP.Context = "/"
		require.NoError(t, cfg.Validate())
	})

	t.Run("context with trailing slash is invalid", func(t *testing.T) {
		cfg := baseConfig()
		cfg.ServerHTTP.Context = "/aion/"
		require.EqualError(t, cfg.Validate(), config.ErrHTTPContextMustNotEndWithSlash)
	})

	t.Run("api root without leading slash is invalid", func(t *testing.T) {
		cfg := baseConfig()
		cfg.ServerHTTP.APIRoot = "api/v1"
		require.EqualError(t, cfg.Validate(), config.ErrHTTPAPIRootMustStart)
	})

	t.Run("docs alias cannot be empty", func(t *testing.T) {
		cfg := baseConfig()
		cfg.ServerHTTP.DocsAliasPath = ""
		require.EqualError(t, cfg.Validate(), config.ErrHTTPDocsAliasPathEmpty)
	})

	t.Run("health route cannot end with slash", func(t *testing.T) {
		cfg := baseConfig()
		cfg.ServerHTTP.HealthRoute = "/health/"
		require.EqualError(t, cfg.Validate(), config.ErrHTTPHealthRouteMustNotEndSlash)
	})
}

func TestConfigValidate_CacheAndDBErrors(t *testing.T) {
	cfg := baseConfig()
	cfg.Cache.PoolSize = 0
	require.EqualError(t, cfg.Validate(), "CACHE_POOL_SIZE must be at least 1")

	cfg = baseConfig()
	cfg.DB.SSLMode = "invalid"
	require.EqualError(t, cfg.Validate(), "invalid DB_SSL_MODE value: invalid")

	cfg = baseConfig()
	cfg.DB.MaxRetries = 0
	require.EqualError(t, cfg.Validate(), "DB_CONNECT_MAX_RETRIES must be at least 1")

	cfg = baseConfig()
	cfg.Cache.Addr = ""
	require.EqualError(t, cfg.Validate(), config.ErrCacheAddrEmpty)

	cfg = baseConfig()
	cfg.DB.TimeZone = ""
	require.EqualError(t, cfg.Validate(), config.ErrDBTimeZoneEmpty)

	cfg = baseConfig()
	cfg.DB.MaxOpenConns = 0
	require.EqualError(t, cfg.Validate(), "DB_MAX_CONNECTIONS must be at least 1")
}

func TestConfigValidate_ObservabilityAndAppErrors(t *testing.T) {
	cfg := baseConfig()
	cfg.Observability.OtelExporterCompression = "brotli"
	require.EqualError(t, cfg.Validate(), "OTel Exporter compression must be either 'none' or 'gzip', got: brotli")

	cfg = baseConfig()
	cfg.Application.ContextRequest = 100 * time.Millisecond
	require.EqualError(t, cfg.Validate(), "context request timeout must be at least 500ms")

	cfg = baseConfig()
	cfg.Application.Timeout = 500 * time.Millisecond
	require.EqualError(t, cfg.Validate(), "shutdown timeout must be at least 1s second")

	cfg = baseConfig()
	cfg.Kafka.Brokers = ""
	require.EqualError(t, cfg.Validate(), config.ErrKafkaBrokersEmpty)

	cfg = baseConfig()
	cfg.Outbox.BatchSize = 0
	require.EqualError(t, cfg.Validate(), "OUTBOX_BATCH_SIZE must be at least 1")

	cfg = baseConfig()
	cfg.Kafka.RecordProjectionEventsTopic = ""
	require.EqualError(t, cfg.Validate(), config.ErrKafkaRecordProjectionEventsTopicEmpty)

	cfg = baseConfig()
	cfg.Realtime.StreamPath = "/"
	require.EqualError(t, cfg.Validate(), config.ErrRealtimeStreamPathTooShort)

	cfg = baseConfig()
	cfg.Realtime.HeartbeatInterval = 500 * time.Millisecond
	require.EqualError(t, cfg.Validate(), "REALTIME_HEARTBEAT_INTERVAL must be at least 1s")
}
