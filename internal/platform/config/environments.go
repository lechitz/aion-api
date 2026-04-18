// Package config provides centralized configuration structures and loading logic.
package config

import "time"

// GeneralConfig holds a general application configuration.
type GeneralConfig struct {
	Name    string `envconfig:"APP_NAME"`
	Env     string `envconfig:"APP_ENV"`
	Version string `envconfig:"APP_VERSION"`
}

// Secret holds encryption/authentication secrets.
type Secret struct {
	Key string `envconfig:"SECRET_KEY"`
}

// ObservabilityConfig holds all observability-related configuration.
//
// OtelExporterEnabled is the kill switch. When false, the tracer/metric
// providers become no-ops and no OTLP exporter is created — this lets the
// public distribution runtime ship without a collector, and keeps the
// configured endpoint optional in that case.
type ObservabilityConfig struct {
	OtelExporterOTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"aion-dev-otel-collector:4318"`
	OtelServiceName          string `envconfig:"OTEL_SERVICE_NAME"           default:"aion-api"`
	OtelServiceVersion       string `envconfig:"OTEL_SERVICE_VERSION"        default:"0.0.1"`
	OtelExporterHeaders      string `envconfig:"OTEL_EXPORTER_HEADERS"       default:""`
	OtelExporterTimeout      string `envconfig:"OTEL_EXPORTER_TIMEOUT"       default:"5s"`
	OtelExporterCompression  string `envconfig:"OTEL_EXPORTER_COMPRESSION"   default:"none"`
	OtelExporterEnabled      bool   `envconfig:"OTEL_EXPORTER_ENABLED"       default:"true"`
	OtelExporterInsecure     bool   `envconfig:"OTEL_EXPORTER_INSECURE"      default:"true"`
}

// KafkaConfig holds Kafka broker and topic settings for canonical event publication.
type KafkaConfig struct {
	Brokers                     string `envconfig:"KAFKA_BROKERS"                        default:"aion-dev-kafka:9092"`
	RecordEventsTopic           string `envconfig:"KAFKA_TOPIC_RECORD_EVENTS"            default:"aion.record.events.v1"`
	RecordProjectionEventsTopic string `envconfig:"KAFKA_TOPIC_RECORD_PROJECTION_EVENTS" default:"aion.record_projection.events.v1"`
}

// OutboxConfig holds runtime controls for the outbox publisher loop.
type OutboxConfig struct {
	PublishEnabled  bool          `envconfig:"OUTBOX_PUBLISH_ENABLED"  default:"true"`
	PublishInterval time.Duration `envconfig:"OUTBOX_PUBLISH_INTERVAL" default:"2s"`
	BatchSize       int           `envconfig:"OUTBOX_BATCH_SIZE"       default:"50"`
}

// RealtimeConfig holds runtime controls for SSE and projection event fanout.
type RealtimeConfig struct {
	StreamPath          string        `envconfig:"REALTIME_STREAM_PATH"           default:"/events/stream"`
	ConsumerGroupPrefix string        `envconfig:"REALTIME_CONSUMER_GROUP_PREFIX" default:"aion-api-realtime"`
	HeartbeatInterval   time.Duration `envconfig:"REALTIME_HEARTBEAT_INTERVAL"    default:"15s"`
	SubscriberBuffer    int           `envconfig:"REALTIME_SUBSCRIBER_BUFFER"     default:"32"`
	Enabled             bool          `envconfig:"REALTIME_ENABLED"               default:"true"`
}

// CacheConfig holds Redis cache configuration.
// Each bounded context uses a separate Redis database for isolation.
type CacheConfig struct {
	Addr     string `envconfig:"CACHE_ADDR"     default:"aion-dev-redis:6379"`
	Password string `envconfig:"CACHE_PASSWORD"`

	AuthDB     int `envconfig:"CACHE_AUTH_DB"     default:"0"`
	CategoryDB int `envconfig:"CACHE_CATEGORY_DB" default:"1"`
	TagDB      int `envconfig:"CACHE_TAG_DB"      default:"2"`
	RecordDB   int `envconfig:"CACHE_RECORD_DB"   default:"3"`
	UserDB     int `envconfig:"CACHE_USER_DB"     default:"4"`
	ChatDB     int `envconfig:"CACHE_CHAT_DB"     default:"5"`

	PoolSize       int           `envconfig:"CACHE_POOL_SIZE"       default:"10"`
	ConnectTimeout time.Duration `envconfig:"CACHE_CONNECT_TIMEOUT" default:"5s"`
}

// CookieConfig holds cookie configuration.
type CookieConfig struct {
	Domain   string `envconfig:"COOKIE_DOMAIN"   default:"localhost"`
	Path     string `envconfig:"COOKIE_PATH"     default:"/"`
	SameSite string `envconfig:"COOKIE_SAMESITE" default:"Lax"`
	Secure   bool   `envconfig:"COOKIE_SECURE"   default:"false"`
	MaxAge   int    `envconfig:"COOKIE_MAX_AGE"  default:"0"`
}

// ServerGraphql holds GraphQL server configuration.
type ServerGraphql struct {
	Host string `envconfig:"GRAPHQL_HOST" default:"0.0.0.0"`
	Name string `envconfig:"GRAPHQL_NAME" default:"GraphQL"`
	Path string `envconfig:"GRAPHQL_PATH" default:"/graphql"`

	ReadTimeout  time.Duration `envconfig:"GRAPHQL_READ_TIMEOUT"  default:"5s"`
	WriteTimeout time.Duration `envconfig:"GRAPHQL_WRITE_TIMEOUT" default:"5s"`

	ReadHeaderTimeout time.Duration `envconfig:"GRAPHQL_READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeout       time.Duration `envconfig:"GRAPHQL_IDLE_TIMEOUT"        default:"60s"`
	MaxHeaderBytes    int           `envconfig:"GRAPHQL_MAX_HEADER_BYTES"    default:"1048576"`
}

// ServerHTTP holds HTTP server configuration.
type ServerHTTP struct {
	Host string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Name string `envconfig:"HTTP_NAME" default:"HTTP"`
	Port string `envconfig:"HTTP_PORT" default:"5001"    required:"true"`

	Context          string `envconfig:"HTTP_CONTEXT"            default:"/aion"`
	APIRoot          string `envconfig:"HTTP_API_ROOT"           default:"/api/v1"`
	SwaggerMountPath string `envconfig:"HTTP_SWAGGER_MOUNT_PATH" default:"/swagger"`
	DocsAliasPath    string `envconfig:"HTTP_DOCS_ALIAS_PATH"    default:"/docs"`
	HealthRoute      string `envconfig:"HTTP_HEALTH_ROUTE"       default:"/health"`

	ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT"        default:"10s"`
	WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"       default:"10s"`
	ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeout       time.Duration `envconfig:"HTTP_IDLE_TIMEOUT"        default:"60s"`
	MaxHeaderBytes    int           `envconfig:"HTTP_MAX_HEADER_BYTES"    default:"1048576"` // 1<<20
}

// DBConfig holds database connection configuration.
type DBConfig struct {
	Type            string        `envconfig:"DB_TYPE"                   default:"postgres"`
	Name            string        `envconfig:"DB_NAME"                                       required:"true"`
	Host            string        `envconfig:"DB_HOST"                   default:"localhost"`
	Port            string        `envconfig:"DB_PORT"                   default:"5432"`
	User            string        `envconfig:"DB_USER"                                       required:"true"`
	Password        string        `envconfig:"DB_PASSWORD"                                   required:"true"`
	SSLMode         string        `envconfig:"DB_SSL_MODE"               default:"disable"`
	TimeZone        string        `envconfig:"TIME_ZONE"                 default:"UTC"`
	MaxOpenConns    int           `envconfig:"DB_MAX_CONNECTIONS"        default:"10"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNECTIONS"   default:"5"`
	MaxRetries      int           `envconfig:"DB_CONNECT_MAX_RETRIES"    default:"3"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME"      default:"30m"`
	RetryInterval   time.Duration `envconfig:"DB_CONNECT_RETRY_INTERVAL" default:"3s"`
}

// Application holds general application-related configuration.
type Application struct {
	Timeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`

	ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
}

// AionChatConfig holds configuration for the Aion-Chat service (Python AI service).
type AionChatConfig struct {
	BaseURL    string        `envconfig:"AION_CHAT_URL"         default:"http://aion-dev-chat:8000"`
	ServiceKey string        `envconfig:"AION_CHAT_SERVICE_KEY" default:""`
	Timeout    time.Duration `envconfig:"AION_CHAT_TIMEOUT"     default:"30s"`
}

// AvatarStorageConfig holds S3-compatible storage config for avatar uploads.
type AvatarStorageConfig struct {
	Provider      string `envconfig:"AVATAR_STORAGE_PROVIDER"     default:"s3"`
	S3Endpoint    string `envconfig:"AVATAR_S3_ENDPOINT"          default:"http://aion-dev-localstack:4566"`
	S3Region      string `envconfig:"AVATAR_S3_REGION"            default:"us-east-1"`
	S3Bucket      string `envconfig:"AVATAR_S3_BUCKET"            default:"aion-assets"`
	S3Prefix      string `envconfig:"AVATAR_S3_PREFIX"            default:"avatars"`
	PublicBaseURL string `envconfig:"AVATAR_PUBLIC_BASE_URL"      default:"http://localhost:4566/aion-assets"`
	AccessKeyID   string `envconfig:"AVATAR_S3_ACCESS_KEY_ID"     default:"test"`
	SecretKey     string `envconfig:"AVATAR_S3_SECRET_ACCESS_KEY" default:"test"`
	MaxUploadMB   int    `envconfig:"AVATAR_MAX_UPLOAD_MB"        default:"20"`
}
