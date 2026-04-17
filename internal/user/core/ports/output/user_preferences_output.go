package output

import (
	"context"

	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// UserPreferencesRepository defines persistence operations for user preferences.
type UserPreferencesRepository interface {
	GetPreferencesByUserID(ctx context.Context, userID uint64) (domain.UserPreferences, error)
	UpsertPreferences(ctx context.Context, prefs domain.UserPreferences) (domain.UserPreferences, error)
}
