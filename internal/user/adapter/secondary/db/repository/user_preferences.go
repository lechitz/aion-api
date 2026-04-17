package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

// GetPreferencesByUserID retrieves user preferences by user ID.
// Returns domain defaults if no record exists yet.
func (up UserRepository) GetPreferencesByUserID(ctx context.Context, userID uint64) (domain.UserPreferences, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.preferences.get_by_user_id")
	defer span.End()
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	var prefDB model.UserPreferencesDB
	err := up.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&prefDB).Error()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "preferences_not_found_returning_defaults")
			return domain.DefaultUserPreferences(userID), nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		up.logger.ErrorwCtx(ctx, "failed to get user preferences", commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	span.SetStatus(codes.Ok, "preferences_retrieved")
	return mapper.UserPreferencesFromDB(prefDB), nil
}

// UpsertPreferences inserts or updates user preferences.
func (up UserRepository) UpsertPreferences(ctx context.Context, prefs domain.UserPreferences) (domain.UserPreferences, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.preferences.upsert")
	defer span.End()
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(prefs.UserID, 10)))

	prefDB := mapper.UserPreferencesToDB(prefs)
	now := time.Now().UTC()
	prefDB.UpdatedAt = now

	err := up.db.WithContext(ctx).
		Raw(`INSERT INTO aion_api.user_preferences (user_id, theme_preset, theme_mode, compact_mode, reduced_motion, custom_overrides, created_at, updated_at)
		     VALUES (?, ?, ?, ?, ?, ?::jsonb, ?, ?)
		     ON CONFLICT (user_id) DO UPDATE SET
		       theme_preset = EXCLUDED.theme_preset,
		       theme_mode = EXCLUDED.theme_mode,
		       compact_mode = EXCLUDED.compact_mode,
		       reduced_motion = EXCLUDED.reduced_motion,
		       custom_overrides = EXCLUDED.custom_overrides,
		       updated_at = EXCLUDED.updated_at
		     RETURNING user_id, theme_preset, theme_mode, compact_mode, reduced_motion, custom_overrides, created_at, updated_at`,
			prefs.UserID, prefDB.ThemePreset, prefDB.ThemeMode, prefDB.CompactMode, prefDB.ReducedMotion, prefDB.CustomOverrides, now, now).
		Scan(&prefDB).Error()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		up.logger.ErrorwCtx(ctx, "failed to upsert user preferences", commonkeys.UserID, prefs.UserID, commonkeys.Error, err.Error())
		return domain.UserPreferences{}, err
	}

	span.SetStatus(codes.Ok, "preferences_upserted")
	up.logger.InfowCtx(ctx, "user preferences upserted", commonkeys.UserID, prefs.UserID)
	return mapper.UserPreferencesFromDB(prefDB), nil
}
