package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/internal/user/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// PreferencesHandler handles user preferences HTTP operations.
type PreferencesHandler struct {
	Service input.UserPreferencesService
	Logger  logger.ContextLogger
	Config  *config.Config
}

// NewPreferencesHandler returns a PreferencesHandler with dependencies injected.
func NewPreferencesHandler(service input.UserPreferencesService, cfg *config.Config, logger logger.ContextLogger) *PreferencesHandler {
	return &PreferencesHandler{
		Service: service,
		Config:  cfg,
		Logger:  logger,
	}
}

// GetPreferences returns the authenticated user's preferences.
func (h *PreferencesHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer(TracerUserHandler)
	ctx, span := tr.Start(r.Context(), "user.preferences.get")
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam)
		httpresponse.WriteAuthError(w, sharederrors.ErrMissingUserIDParam, h.Logger)
		return
	}

	prefs, err := h.Service.GetPreferences(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, err.Error())
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "get_preferences", h.Logger)
		return
	}

	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))
	span.SetStatus(codes.Ok, "preferences_retrieved")

	response := dto.UserPreferencesResponseFromDomain(prefs)
	httpresponse.WriteSuccess(w, http.StatusOK, response, "preferences_retrieved")
}

// SavePreferences updates the authenticated user's preferences.
func (h *PreferencesHandler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer(TracerUserHandler)
	ctx, span := tr.Start(r.Context(), "user.preferences.save")
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam)
		httpresponse.WriteAuthError(w, sharederrors.ErrMissingUserIDParam, h.Logger)
		return
	}

	var req dto.SaveUserPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	cmd := req.ToCommand()

	saved, err := h.Service.SavePreferences(ctx, userID, cmd)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "save_preferences", h.Logger)
		return
	}

	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))
	span.SetStatus(codes.Ok, "preferences_saved")

	response := dto.UserPreferencesResponseFromDomain(saved)
	httpresponse.WriteSuccess(w, http.StatusOK, response, "preferences_saved")
}
