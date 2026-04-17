package input

import (
	"context"

	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// SavePreferencesCommand carries the fields to persist for user preferences.
type SavePreferencesCommand struct {
	ThemePreset     string
	ThemeMode       string
	CustomOverrides map[string]string
	CompactMode     *bool
	ReducedMotion   *bool
}

// UserPreferencesReader defines read access to user preferences.
type UserPreferencesReader interface {
	GetPreferences(ctx context.Context, userID uint64) (domain.UserPreferences, error)
}

// UserPreferencesWriter defines write access to user preferences.
type UserPreferencesWriter interface {
	SavePreferences(ctx context.Context, userID uint64, cmd SavePreferencesCommand) (domain.UserPreferences, error)
}

// UserPreferencesService aggregates read and write operations for user preferences.
type UserPreferencesService interface {
	UserPreferencesReader
	UserPreferencesWriter
}
