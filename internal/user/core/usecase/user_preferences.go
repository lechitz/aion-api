package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"github.com/lechitz/aion-api/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// PreferencesService implements input.UserPreferencesService.
type PreferencesService struct {
	repo   output.UserPreferencesRepository
	logger logger.ContextLogger
}

// NewPreferencesService creates a new PreferencesService.
func NewPreferencesService(repo output.UserPreferencesRepository, logger logger.ContextLogger) *PreferencesService {
	return &PreferencesService{
		repo:   repo,
		logger: logger,
	}
}

// GetPreferences returns user preferences, defaulting if none exist.
func (s *PreferencesService) GetPreferences(ctx context.Context, userID uint64) (domain.UserPreferences, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanGetPreferences)
	defer span.End()
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	s.logger.InfowCtx(ctx, "getting user preferences", commonkeys.UserID, userID)

	prefs, err := s.repo.GetPreferencesByUserID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorwCtx(ctx, "failed to get user preferences", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	span.SetStatus(codes.Ok, "preferences_retrieved")
	return prefs, nil
}

// SavePreferences validates and persists user preferences.
func (s *PreferencesService) SavePreferences(ctx context.Context, userID uint64, cmd input.SavePreferencesCommand) (domain.UserPreferences, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanSavePreferences)
	defer span.End()
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	s.logger.InfowCtx(ctx, "saving user preferences", commonkeys.UserID, userID)

	if cmd.ThemePreset != "" && !domain.AllowedThemePresets[cmd.ThemePreset] {
		err := sharederrors.NewValidationError("theme_preset", "invalid theme preset")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.UserPreferences{}, err
	}
	if cmd.ThemeMode != "" && !domain.AllowedThemeModes[cmd.ThemeMode] {
		err := sharederrors.NewValidationError("theme_mode", "invalid theme mode")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	// Start from current preferences to support partial updates
	current, err := s.repo.GetPreferencesByUserID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorwCtx(ctx, "failed to get current preferences for update", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	if cmd.ThemePreset != "" {
		current.ThemePreset = cmd.ThemePreset
	}
	if cmd.ThemeMode != "" {
		current.ThemeMode = cmd.ThemeMode
	}
	if cmd.CompactMode != nil {
		current.CompactMode = *cmd.CompactMode
	}
	if cmd.ReducedMotion != nil {
		current.ReducedMotion = *cmd.ReducedMotion
	}
	if cmd.CustomOverrides != nil {
		current.CustomOverrides = cmd.CustomOverrides
	}
	current.UserID = userID

	saved, err := s.repo.UpsertPreferences(ctx, current)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorwCtx(ctx, "failed to save user preferences", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	span.SetStatus(codes.Ok, "preferences_saved")
	s.logger.InfowCtx(ctx, "user preferences saved", commonkeys.UserID, userID)
	return saved, nil
}
