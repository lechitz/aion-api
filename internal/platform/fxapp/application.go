// Package fxapp wires the application using Uber Fx modules.
package fxapp

import (
	"github.com/lechitz/aion-api/internal/adapter/secondary/hasher"
	"github.com/lechitz/aion-api/internal/adapter/secondary/token"
	adminRepo "github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/repository"
	admin "github.com/lechitz/aion-api/internal/admin/core/usecase"
	auditRepo "github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/repository"
	audit "github.com/lechitz/aion-api/internal/audit/core/usecase"
	authCache "github.com/lechitz/aion-api/internal/auth/adapter/secondary/cache"
	auth "github.com/lechitz/aion-api/internal/auth/core/usecase"
	categoryCache "github.com/lechitz/aion-api/internal/category/adapter/secondary/cache"
	categoryRepo "github.com/lechitz/aion-api/internal/category/adapter/secondary/db/repository"
	category "github.com/lechitz/aion-api/internal/category/core/usecase"
	chatCache "github.com/lechitz/aion-api/internal/chat/adapter/secondary/cache"
	chatHistoryRepo "github.com/lechitz/aion-api/internal/chat/adapter/secondary/db/repository"
	chatClient "github.com/lechitz/aion-api/internal/chat/adapter/secondary/http"
	chat "github.com/lechitz/aion-api/internal/chat/core/usecase"
	eventOutboxRepo "github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/repository"
	eventOutbox "github.com/lechitz/aion-api/internal/eventoutbox/core/usecase"
	"github.com/lechitz/aion-api/internal/platform/app"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/cache"
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/httpclient"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	realtime "github.com/lechitz/aion-api/internal/realtime/core/usecase"
	recordCache "github.com/lechitz/aion-api/internal/record/adapter/secondary/cache"
	recordRepo "github.com/lechitz/aion-api/internal/record/adapter/secondary/db/repository"
	record "github.com/lechitz/aion-api/internal/record/core/usecase"
	tagCache "github.com/lechitz/aion-api/internal/tag/adapter/secondary/cache"
	tagRepo "github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/repository"
	tag "github.com/lechitz/aion-api/internal/tag/core/usecase"
	userCache "github.com/lechitz/aion-api/internal/user/adapter/secondary/cache"
	userRepo "github.com/lechitz/aion-api/internal/user/adapter/secondary/db/repository"
	userAvatarStorage "github.com/lechitz/aion-api/internal/user/adapter/secondary/storage/s3"
	user "github.com/lechitz/aion-api/internal/user/core/usecase"
	"go.uber.org/fx"
)

// ApplicationModule wires the application layer (use cases, repositories, adapters) and exposes Dependencies for HTTP composition.
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var ApplicationModule = fx.Options(fx.Provide(ProvideAppDependencies))

// AppDependencies is a type alias for app.Dependencies for backwards compatibility.
type AppDependencies = app.Dependencies

// appDepsParams groups all dependencies needed for application wiring.
type appDepsParams struct {
	fx.In

	Cfg           *config.Config
	DB            db.DB
	AuthCache     cache.Cache `name:"authCache"`
	CategoryCache cache.Cache `name:"categoryCache"`
	TagCache      cache.Cache `name:"tagCache"`
	RecordCache   cache.Cache `name:"recordCache"`
	UserCache     cache.Cache `name:"userCache"`
	ChatCache     cache.Cache `name:"chatCache"`
	HTTPClient    httpclient.HTTPClient
	Log           logger.ContextLogger
}

// ProvideAppDependencies composes repositories and use cases using shared infrastructure.
// Receives db.DB interface (not *gorm.DB) from InfraModule, following Dependency Inversion Principle.
// Each bounded context uses its own Redis database for cache isolation.
func ProvideAppDependencies(deps appDepsParams) *AppDependencies {
	hasherProvider := hasher.New()
	tokenProvider := token.NewProvider(deps.Cfg.Secret.Key)

	adminRepository := adminRepo.New(deps.DB, deps.Log)

	userRepository := userRepo.New(deps.DB, deps.Log, adminRepository)
	avatarStorage, avatarErr := userAvatarStorage.NewAvatarStorage(deps.Cfg.AvatarStorage, deps.Log)
	if avatarErr != nil {
		deps.Log.Errorw("failed to initialize avatar storage", "error", avatarErr)
		avatarStorage = nil
	}

	categoryRepository := categoryRepo.New(deps.DB, deps.Log)
	tagRepository := tagRepo.New(deps.DB, deps.Log)
	recordRepository := recordRepo.New(deps.DB, deps.Log)
	chatHistoryRepository := chatHistoryRepo.New(deps.DB, deps.Log)
	auditActionEventRepository := auditRepo.NewAuditActionEventRepository(deps.DB, deps.Log)
	eventOutboxRepository := eventOutboxRepo.NewEventRepository(deps.DB, deps.Log)

	authCacheStore := authCache.NewStore(deps.AuthCache, deps.Log)
	userCacheStore := userCache.NewStore(deps.UserCache, deps.Log)
	categoryCacheStore := categoryCache.NewStore(deps.CategoryCache, deps.Log)
	tagCacheStore := tagCache.NewStore(deps.TagCache, deps.Log)
	recordCacheStore := recordCache.NewStore(deps.RecordCache, deps.Log)
	chatHistoryCacheStore := chatCache.NewStore(deps.ChatCache, deps.Log)
	chatHTTPClient := chatClient.New(deps.HTTPClient, deps.Cfg.AionChat.BaseURL, deps.Log)
	auditService := audit.NewService(auditActionEventRepository, deps.Log)
	outboxService := eventOutbox.NewService(eventOutboxRepository, deps.Log)
	realtimeService := realtime.NewService(deps.Log, deps.Cfg.Realtime.SubscriberBuffer)

	authService := auth.NewService(adminRepository, authCacheStore, userRepository, userCacheStore, authCacheStore, tokenProvider, hasherProvider, deps.Log)
	userService := user.NewService(userRepository, userRepository, userCacheStore, avatarStorage, authCacheStore, tokenProvider, hasherProvider, deps.Log)
	userPreferencesService := user.NewPreferencesService(userRepository, deps.Log)
	adminService := admin.NewService(adminRepository, authCacheStore, authCacheStore, deps.Log)
	categoryService := category.NewService(categoryRepository, categoryCacheStore, deps.Log)
	tagService := tag.NewService(tagRepository, tagCacheStore, deps.Log)
	recordService := record.NewService(recordRepository, recordCacheStore, tagRepository, deps.Log).
		WithOutbox(outboxService).
		WithTransactionManager(deps.DB).
		WithProjectionReader(recordRepository)
	chatService := chat.NewService(chatHTTPClient, chatHistoryRepository, chatHistoryCacheStore, auditService, deps.Log)

	return &AppDependencies{
		AuthService:            authService,
		UserService:            userService,
		UserPreferencesService: userPreferencesService,
		AdminService:           adminService,
		CategoryService:        categoryService,
		TagService:             tagService,
		RecordService:          recordService,
		ChatService:            chatService,
		AuditService:           auditService,
		OutboxService:          outboxService,
		RealtimeService:        realtimeService,
		Logger:                 deps.Log,
	}
}
