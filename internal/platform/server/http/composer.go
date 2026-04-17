// Package http provides the HTTP server composition for all adapters and routes.
//
//revive:disable:var-naming // keep using http to clearly denote the server layer
package http

//revive:enable:var-naming

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql"
	httpSwagger "github.com/swaggo/http-swagger"

	adminhandler "github.com/lechitz/aion-api/internal/admin/adapter/primary/http/handler"
	audithandler "github.com/lechitz/aion-api/internal/audit/adapter/primary/http/handler"
	authhandler "github.com/lechitz/aion-api/internal/auth/adapter/primary/http/handler"
	chathandler "github.com/lechitz/aion-api/internal/chat/adapter/primary/http/handler"
	realtimehandler "github.com/lechitz/aion-api/internal/realtime/adapter/primary/http/handler"
	userhandler "github.com/lechitz/aion-api/internal/user/adapter/primary/http/handler"

	"github.com/lechitz/aion-api/internal/platform/app"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	genericHandler "github.com/lechitz/aion-api/internal/platform/server/http/generic/handler"
	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/cors"
	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/recovery"
	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/requestid"
	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/servicetoken"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
	"github.com/lechitz/aion-api/internal/platform/server/http/router/chi"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type httpRouteConfig struct {
	apiContext   string
	swaggerMount string
	docsAlias    string
	routeHealth  string
}

// ComposeHandler assembles the HTTP handler with platform middlewares, domain routes, Swagger UI and GraphQL.
func ComposeHandler(cfg *config.Config, deps *app.Dependencies, log logger.ContextLogger) (http.Handler, error) {
	// Define the main router for the HTTP server
	router := chi.New()

	genericHandlers := genericHandler.New(log, cfg.General)

	// Global middlewares
	router.Use(
		requestid.New(),
		recovery.New(genericHandlers),
		cors.New(),
	)

	// Default handlers
	router.SetNotFound(http.HandlerFunc(genericHandlers.NotFoundHandler))
	router.SetMethodNotAllowed(http.HandlerFunc(genericHandlers.MethodNotAllowedHandler))
	router.SetError(genericHandlers.ErrorHandler)

	routes := resolveHTTPRouteConfig(cfg)
	router.Group(routes.apiContext, func(api ports.Router) {
		mountSwaggerAndDocs(api, routes)
		api.Group(cfg.ServerHTTP.APIRoot, func(v1 ports.Router) {
			registerDomainRoutes(v1, cfg, deps, log)
		})
	})

	// OpenTelemetry HTTP wrapper: instrument the main router but expose health route uninstrumented
	instrumented := otelhttp.NewHandler(router, fmt.Sprintf(OTelHTTPHandlerNameFormat, cfg.Observability.OtelServiceName))

	mux := http.NewServeMux()

	corsMiddleware := cors.New()
	requestIDMiddleware := requestid.New()
	healthHandler := requestIDMiddleware(http.HandlerFunc(genericHandlers.HealthCheck))

	mountHealthRoutes(mux, corsMiddleware(healthHandler), routes, cfg.ServerHTTP.APIRoot)

	mux.Handle("/", instrumented)

	return mux, nil
}

func resolveHTTPRouteConfig(cfg *config.Config) httpRouteConfig {
	apiContext := cfg.ServerHTTP.Context
	if apiContext == "" {
		apiContext = "/"
	}

	swaggerMount := cfg.ServerHTTP.SwaggerMountPath
	if swaggerMount == "" {
		swaggerMount = DefaultSwaggerMountPath
	}

	docsAlias := cfg.ServerHTTP.DocsAliasPath
	if docsAlias == "" {
		docsAlias = DefaultDocsAliasPath
	}

	routeHealth := cfg.ServerHTTP.HealthRoute
	if routeHealth == "" {
		routeHealth = DefaultRouteHealth
	}

	return httpRouteConfig{
		apiContext:   apiContext,
		swaggerMount: swaggerMount,
		docsAlias:    docsAlias,
		routeHealth:  routeHealth,
	}
}

func mountSwaggerAndDocs(api ports.Router, routes httpRouteConfig) {
	swaggerDocURL := path.Clean(routes.apiContext + "/" +
		strings.TrimPrefix(routes.swaggerMount, "/") + "/" + DefaultSwaggerDocFile)

	api.Mount(routes.swaggerMount, httpSwagger.Handler(
		httpSwagger.URL(swaggerDocURL),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	api.GET(routes.docsAlias, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path.Join(routes.apiContext, routes.swaggerMount, DefaultSwaggerIndexFile), http.StatusTemporaryRedirect)
	}))
}

func registerDomainRoutes(v1 ports.Router, cfg *config.Config, deps *app.Dependencies, log logger.ContextLogger) {
	if deps.AuthService != nil {
		ah := authhandler.New(deps.AuthService, cfg, log)
		authhandler.RegisterHTTP(v1, ah)
	}

	if deps.UserService != nil {
		uh := userhandler.New(deps.UserService, cfg, log)
		var ph *userhandler.PreferencesHandler
		if deps.UserPreferencesService != nil {
			ph = userhandler.NewPreferencesHandler(deps.UserPreferencesService, cfg, log)
		}
		userhandler.RegisterHTTP(v1, uh, ph, deps.AuthService, log)
	}

	if deps.AdminService != nil {
		ah := adminhandler.New(deps.AdminService, cfg, log)
		adminhandler.RegisterHTTP(v1, ah, deps.AuthService, log)
	}

	if deps.ChatService != nil {
		ch := chathandler.New(deps.ChatService, cfg, log)
		chathandler.RegisterHTTP(v1, ch, deps.AuthService, log)
	}

	if deps.AuditService != nil {
		ah := audithandler.New(deps.AuditService, cfg, log)
		audithandler.RegisterHTTP(v1, ah, deps.AuthService, log)
	}

	if deps.RealtimeService != nil {
		rh := realtimehandler.New(deps.RealtimeService, cfg, log)
		realtimehandler.RegisterHTTP(v1, rh, deps.AuthService, log)
	}

	gqlHandler, err := graphql.NewGraphqlHandler(
		deps.AuthService,
		deps.CategoryService,
		deps.TagService,
		deps.RecordService,
		deps.ChatService,
		deps.UserService,
		log,
		cfg,
	)
	if err != nil {
		log.Errorw(LogErrComposeGraphQL, commonkeys.Error, err)
		return
	}

	wrappedGQL := servicetoken.New(cfg, log)(gqlHandler)
	v1.Mount(cfg.ServerGraphql.Path, wrappedGQL)
}

func mountHealthRoutes(mux *http.ServeMux, healthHandler http.Handler, routes httpRouteConfig, apiRoot string) {
	pathClean := path.Clean(routes.apiContext + "/" + strings.TrimPrefix(routes.routeHealth, "/"))
	mux.Handle(pathClean, healthHandler)
	mux.Handle(pathClean+"/", healthHandler)

	altHealth := path.Clean(routes.apiContext + "/" + strings.TrimPrefix(apiRoot, "/") + "/" + strings.TrimPrefix(routes.routeHealth, "/"))
	mux.Handle(altHealth, healthHandler)
	mux.Handle(altHealth+"/", healthHandler)
}
