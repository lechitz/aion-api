// Package config provides configuration loading and validation for the application.
package config

import (
	"errors"
	"fmt"
)

// Config holds all configuration sections required to bootstrap the application.
type Config struct {
	Realtime      RealtimeConfig
	Kafka         KafkaConfig
	General       GeneralConfig
	Secret        Secret
	AvatarStorage AvatarStorageConfig
	Observability ObservabilityConfig
	AionChat      AionChatConfig
	Cookie        CookieConfig
	ServerHTTP    ServerHTTP
	DB            DBConfig
	ServerGraphql ServerGraphql
	Cache         CacheConfig
	Outbox        OutboxConfig
	Application   Application
}

// Validate checks if the configuration is valid, returning the first validation error encountered.
func (c *Config) Validate() error {
	if err := c.validateHTTP(); err != nil {
		return err
	}
	if err := c.validateGraphQL(); err != nil {
		return err
	}
	if err := c.validateCache(); err != nil {
		return err
	}
	if err := c.validateDB(); err != nil {
		return err
	}
	if err := c.validateObservability(); err != nil {
		return err
	}
	if err := c.validateKafka(); err != nil {
		return err
	}
	if err := c.validateApp(); err != nil {
		return err
	}
	return nil
}

func (c *Config) validateKafka() error {
	if c.Kafka.Brokers == "" {
		return errors.New(ErrKafkaBrokersEmpty)
	}
	if c.Kafka.RecordEventsTopic == "" {
		return errors.New(ErrKafkaRecordEventsTopicEmpty)
	}
	if c.Kafka.RecordProjectionEventsTopic == "" {
		return errors.New(ErrKafkaRecordProjectionEventsTopicEmpty)
	}
	if c.Outbox.PublishInterval < MinOutboxPublishInterval {
		return fmt.Errorf(ErrOutboxPublishIntervalMin, MinOutboxPublishInterval)
	}
	if c.Outbox.BatchSize < MinOutboxBatchSize {
		return fmt.Errorf(ErrOutboxBatchSizeMin, MinOutboxBatchSize)
	}
	if c.Realtime.Enabled {
		if err := validateHTTPPath(
			c.Realtime.StreamPath,
			false,
			ErrRealtimeStreamPathEmpty,
			ErrRealtimeStreamPathMustStart,
			ErrRealtimeStreamPathTooShort,
			ErrRealtimeStreamPathMustNotEndSlash,
		); err != nil {
			return err
		}
		if c.Realtime.HeartbeatInterval < MinRealtimeHeartbeatInterval {
			return fmt.Errorf(ErrRealtimeHeartbeatIntervalMin, MinRealtimeHeartbeatInterval)
		}
		if c.Realtime.SubscriberBuffer < MinRealtimeSubscriberBuffer {
			return fmt.Errorf(ErrRealtimeSubscriberBufferMin, MinRealtimeSubscriberBuffer)
		}
		if c.Realtime.ConsumerGroupPrefix == "" {
			return errors.New(ErrRealtimeConsumerGroupPrefixEmpty)
		}
	}
	return nil
}

func (c *Config) validateHTTP() error {
	if c.ServerHTTP.Host == "" {
		return errors.New(ErrHTTPHostRequired)
	}
	if c.ServerHTTP.Port == "" {
		return errors.New(ErrHTTPPortRequired)
	}

	// Context: allow root "/" but disallow trailing slash otherwise.
	if err := validateHTTPPath(
		c.ServerHTTP.Context,
		true,
		ErrHTTPContextPathEmpty,
		ErrHTTPContextMustStart,
		"",
		ErrHTTPContextMustNotEndWithSlash,
	); err != nil {
		return err
	}

	// API root (versioned base, e.g., "/api/v1"): must start with '/', cannot be just "/", cannot end with "/".
	if err := validateHTTPPath(
		c.ServerHTTP.APIRoot,
		false,
		ErrHTTPAPIRootEmpty,
		ErrHTTPAPIRootMustStart,
		ErrHTTPAPIRootTooShort,
		ErrHTTPAPIRootMustNotEndSlash,
	); err != nil {
		return err
	}

	// Swagger mount: must start with '/', cannot be just "/", cannot end with "/".
	if err := validateHTTPPath(
		c.ServerHTTP.SwaggerMountPath,
		false,
		ErrHTTPSwaggerMountPathEmpty,
		ErrHTTPSwaggerMountMustStart,
		ErrHTTPSwaggerMountTooShort,
		ErrHTTPSwaggerMountMustNotEndSlash,
	); err != nil {
		return err
	}

	// Docs alias: same rules as swagger mount.
	if err := validateHTTPPath(
		c.ServerHTTP.DocsAliasPath,
		false,
		ErrHTTPDocsAliasPathEmpty,
		ErrHTTPDocsAliasMustStart,
		ErrHTTPDocsAliasTooShort,
		ErrHTTPDocsAliasMustNotEndSlash,
	); err != nil {
		return err
	}

	// Health route: same rules as swagger mount.
	if err := validateHTTPPath(
		c.ServerHTTP.HealthRoute,
		false,
		ErrHTTPHealthRouteEmpty,
		ErrHTTPHealthRouteMustStart,
		ErrHTTPHealthRouteTooShort,
		ErrHTTPHealthRouteMustNotEndSlash,
	); err != nil {
		return err
	}

	// Timeouts and header limits
	if c.ServerHTTP.ReadTimeout < MinHTTPTimeout {
		return fmt.Errorf(ErrHTTPReadTimeoutMin, MinHTTPTimeout)
	}
	if c.ServerHTTP.WriteTimeout != 0 && c.ServerHTTP.WriteTimeout < MinHTTPTimeout {
		return fmt.Errorf(ErrHTTPWriteTimeoutMin, MinHTTPTimeout)
	}
	if c.ServerHTTP.ReadHeaderTimeout <= 0 {
		return errors.New(ErrHTTPReadHeaderTimeoutMin)
	}
	if c.ServerHTTP.IdleTimeout <= 0 {
		return errors.New(ErrHTTPIdleTimeoutMin)
	}
	if c.ServerHTTP.MaxHeaderBytes <= 0 {
		return errors.New(ErrHTTPMaxHeaderBytesMin)
	}

	return nil
}

func (c *Config) validateGraphQL() error {
	if c.ServerGraphql.Path == "" {
		return errors.New(ErrGraphqlPathRequired)
	}
	if c.ServerGraphql.Path[0] != '/' {
		return errors.New(ErrGraphqlPathMustStart)
	}
	return nil
}

func (c *Config) validateCache() error {
	if c.Cache.PoolSize < MinCachePoolSize {
		return fmt.Errorf(ErrCachePoolSizeMin, MinCachePoolSize)
	}
	if c.Cache.Addr == "" {
		return errors.New(ErrCacheAddrEmpty)
	}
	return nil
}

func (c *Config) validateDB() error {
	if c.DB.Type == "" {
		return errors.New(ErrDBTypeEmpty)
	}
	if c.DB.Host == "" {
		return errors.New(ErrDBHostEmpty)
	}
	if c.DB.Port == "" {
		return errors.New(ErrDBPortEmpty)
	}
	if c.DB.Name == "" {
		return errors.New(ErrDBNameRequired)
	}
	if c.DB.User == "" {
		return errors.New(ErrDBUserRequired)
	}
	if c.DB.Password == "" {
		return errors.New(ErrDBPasswordRequired)
	}
	if c.DB.TimeZone == "" {
		return errors.New(ErrDBTimeZoneEmpty)
	}

	switch c.DB.SSLMode {
	case "disable", "require", "verify-ca", "verify-full":
		// valid
	default:
		return fmt.Errorf(ErrDBSSLModeInvalid, c.DB.SSLMode)
	}

	if c.DB.MaxOpenConns < MinDBMaxOpenConns {
		return fmt.Errorf(ErrDBMaxOpenConnsMin, MinDBMaxOpenConns)
	}
	if c.DB.MaxIdleConns < MinDBMaxIdleConns {
		return errors.New(ErrDBMaxIdleConnsNeg)
	}
	if c.DB.ConnMaxLifetime < MinDBConnMaxLifetimeMin {
		return errors.New(ErrDBConnMaxLifetimeNeg)
	}
	if c.DB.RetryInterval < MinDBRetryInterval {
		return fmt.Errorf(ErrDBRetryIntervalMin, MinDBRetryInterval)
	}
	if c.DB.MaxRetries < MinDBMaxRetries {
		return fmt.Errorf(ErrDBMaxRetriesMin, MinDBMaxRetries)
	}

	return nil
}

func (c *Config) validateObservability() error {
	if c.Observability.OtelExporterEnabled && c.Observability.OtelExporterOTLPEndpoint == "" {
		return errors.New(ErrOtelEndpointEmpty)
	}
	if c.Observability.OtelExporterCompression != "" {
		switch c.Observability.OtelExporterCompression {
		case "none", "gzip":
			// valid
		default:
			return fmt.Errorf(ErrOtelCompressionInvalid, c.Observability.OtelExporterCompression)
		}
	}
	return nil
}

func (c *Config) validateApp() error {
	if c.Application.ContextRequest < MinContextRequest {
		return fmt.Errorf(ErrAppContextReqMin, MinContextRequest)
	}
	if c.Application.Timeout < MinShutdownTimeout {
		return fmt.Errorf(ErrAppShutdownTimeoutMin, MinShutdownTimeout.String())
	}
	return nil
}

// validateHTTPPath centralizes path validation rules.
// - allowRoot=true: "/" is allowed; any longer value must not end with "/".
// - allowRoot=false: must start with "/", cannot be just "/", cannot end with "/".
func validateHTTPPath(v string, allowRoot bool, errEmpty, errMustStart, errTooShort, errMustNotEnd string) error {
	if v == "" {
		return errors.New(errEmpty)
	}
	if v[0] != '/' {
		return errors.New(errMustStart)
	}
	if allowRoot {
		// "/" is allowed; disallow trailing slash only when longer than 1.
		if len(v) > 1 && v[len(v)-1] == '/' {
			return errors.New(errMustNotEnd)
		}
		return nil
	}

	if len(v) == 1 {
		return errors.New(errTooShort)
	}

	if v[len(v)-1] == '/' {
		return errors.New(errMustNotEnd)
	}
	return nil
}
