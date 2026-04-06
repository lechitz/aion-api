// Package handler constants contains constants used throughout the generic handler.
package handler

import (
	httperrors "github.com/lechitz/aion-api/internal/platform/server/http/errors"
)

// Tracer names for OpenTelemetry generic handler operations.
const (
	TracerGenericHandler     = "aion-api.generic.handler" // Main tracer for generic handler
	TracerHealthCheckHandler = "generic.health_check"     // Span name for health check
	TracerErrorHandler       = "generic.error_handler"    // Span name for internal error handler
	TracerRecoveryHandler    = "generic.recovery_handler" // Span name for recovery from panic
)

// EventHealthCheck Event names for key points within generic handler spans.
const (
	EventHealthCheck = "health_check_event"
)

// StatusHealthCheckOK Status names for semantic span states.
const (
	StatusHealthCheckOK = "health_check_ok" // Semantic status for a healthy state
)

// HealthStatusHealthy Other string constants for health status.
const (
	HealthStatusHealthy = "healthy"
)

// Error variables for generic handler (to be used as an error interface).
var (
	// ErrMethodNotAllowed indicates that the HTTP method used is not allowed for the requested resource.
	ErrMethodNotAllowed = httperrors.ErrMethodNotAllowed
	// ErrResourceNotFound indicates that the requested resource could not be found.
	ErrResourceNotFound = httperrors.ErrResourceNotFound
	// ErrInternalServer indicates a generic internal server error.
	ErrInternalServer = httperrors.ErrInternalServer
)

// Standardized messages for logs, responses, and traces.
const (
	MsgServiceIsHealthy     = "service is healthy"
	MsgMethodNotAllowed     = "the requested method is not allowed"
	MsgResourceNotFound     = "resource not found"
	MsgRecoveredFromPanic   = "application recovered from panic"
	MsgInternalServerError  = "internal server error"
	MsgRecoveryHandlerFired = "recovery handler triggered"
)

// StacktraceFormat is the printf format for recording recovered panic and stacktrace in tracing/logs.
const StacktraceFormat = "%v\n%s"

// RecoveredFormat is the printf format for recording recovered panic and stacktrace in tracing/logs.
const RecoveredFormat = "recovered from panic: %s\n%s"
