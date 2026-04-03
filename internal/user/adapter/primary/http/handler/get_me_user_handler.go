// Package handler user implements HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/internal/user/adapter/primary/http/dto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetMe returns the authenticated user's data.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer(TracerUserHandler)
	ctx, span := tr.Start(r.Context(), "user.get_me")
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.SetStatus(codes.Error, ErrMissingUserIDParam)
		h.Logger.ErrorwCtx(ctx, ErrMissingUserIDParam)
		httpresponse.WriteAuthError(w, sharederrors.ErrMissingUserIDParam, h.Logger)
		return
	}

	user, err := h.UserService.GetByID(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, err.Error())
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "get_me", h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
		attribute.String(commonkeys.Username, user.Username),
		attribute.String(commonkeys.Email, user.Email),
		attribute.String(commonkeys.Name, user.Name),
	)
	span.SetStatus(codes.Ok, "user_me_success")

	response := dto.UserMeResponse{
		Name:                user.Name,
		Username:            user.Username,
		Email:               user.Email,
		ID:                  user.ID,
		CreatedAt:           user.CreatedAt,
		Locale:              user.Locale,
		Timezone:            user.Timezone,
		Location:            user.Location,
		Bio:                 user.Bio,
		AvatarURL:           user.AvatarURL,
		OnboardingCompleted: user.OnboardingCompleted,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, response, "user_me_success")
}
