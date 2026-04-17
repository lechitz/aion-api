package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/aion-api/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/aion-api/internal/auth/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
)

// RegisterHTTP registers the user-related HTTP handlers with the provided router.
func RegisterHTTP(
	r ports.Router,
	h *Handler,
	ph *PreferencesHandler,
	authService authinput.AuthService,
	lg logger.ContextLogger,
) {
	r.Group("/user", func(ur ports.Router) {
		// Public
		ur.POST("/create", http.HandlerFunc(h.Create))
		ur.POST("/avatar/upload", http.HandlerFunc(h.UploadAvatar))

		// Private
		if authService != nil {
			mw := authMiddleware.New(authService, lg)
			ur.GroupWith(mw.Auth, func(pr ports.Router) {
				pr.GET("/all", http.HandlerFunc(h.ListAll))
				pr.GET("/me", http.HandlerFunc(h.GetMe))
				pr.GET("/{user_id}", http.HandlerFunc(h.GetUserByID))
				pr.PUT("/", http.HandlerFunc(h.UpdateUser))
				pr.DELETE("/avatar", http.HandlerFunc(h.DeleteAvatar))
				pr.PUT("/password", http.HandlerFunc(h.UpdateUserPassword))
				pr.DELETE("/", http.HandlerFunc(h.SoftDeleteUser))
				if ph != nil {
					pr.GET("/preferences", http.HandlerFunc(ph.GetPreferences))
					pr.PUT("/preferences", http.HandlerFunc(ph.SavePreferences))
				}
			})
		}
	})

	r.Group("/registration", func(rr ports.Router) {
		rr.POST("/start", http.HandlerFunc(h.StartRegistration))
		rr.PUT("/{registration_id}/profile", http.HandlerFunc(h.UpdateRegistrationProfile))
		rr.PUT("/{registration_id}/avatar", http.HandlerFunc(h.UpdateRegistrationAvatar))
		rr.POST("/{registration_id}/complete", http.HandlerFunc(h.CompleteRegistration))
	})
}
