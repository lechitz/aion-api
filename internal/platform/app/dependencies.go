// Package app defines the core application contracts.
package app

import (
	inputAdmin "github.com/lechitz/aion-api/internal/admin/core/ports/input"
	inputAudit "github.com/lechitz/aion-api/internal/audit/core/ports/input"
	inputAuth "github.com/lechitz/aion-api/internal/auth/core/ports/input"
	inputCategory "github.com/lechitz/aion-api/internal/category/core/ports/input"
	inputChat "github.com/lechitz/aion-api/internal/chat/core/ports/input"
	inputEventOutbox "github.com/lechitz/aion-api/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	inputRealtime "github.com/lechitz/aion-api/internal/realtime/core/ports/input"
	inputRecord "github.com/lechitz/aion-api/internal/record/core/ports/input"
	inputTag "github.com/lechitz/aion-api/internal/tag/core/ports/input"
	inputUser "github.com/lechitz/aion-api/internal/user/core/ports/input"
)

// Dependencies exposes application services that primary adapters (HTTP/GraphQL) consume.
// This is the contract between the application layer and presentation layer.
type Dependencies struct {
	AuthService            inputAuth.AuthService
	UserService            inputUser.UserService
	UserPreferencesService inputUser.UserPreferencesService
	AdminService           inputAdmin.AdminService
	CategoryService        inputCategory.CategoryService
	TagService             inputTag.TagService
	RecordService          inputRecord.RecordService
	ChatService            inputChat.ChatService
	AuditService           inputAudit.Service
	OutboxService          inputEventOutbox.Service
	RealtimeService        inputRealtime.Service
	Logger                 logger.ContextLogger
}
