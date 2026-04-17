package domain

import "time"

// UserPreferences holds the user's UI and theme preferences.
type UserPreferences struct {
	CreatedAt       time.Time         // Timestamp of when preferences were first created
	UpdatedAt       time.Time         // Timestamp of last update
	ThemePreset     string            // Theme preset name (e.g. "default", "warm-gray", "midnight")
	ThemeMode       string            // Theme mode: "system", "light", or "dark"
	CustomOverrides map[string]string // Per-variable CSS overrides (future use)
	UserID          uint64            // Owner user ID
	CompactMode     bool              // Whether compact density is enabled
	ReducedMotion   bool              // Whether animations are reduced
}

// IsAllowedThemePreset reports whether a preset name is accepted by the UI contract.
func IsAllowedThemePreset(value string) bool {
	switch value {
	case "default", "midnight", "warm-gray", "ocean", "frost", "ember", "charcoal", "custom":
		return true
	default:
		return false
	}
}

// IsAllowedThemeMode reports whether a theme mode is accepted by the UI contract.
func IsAllowedThemeMode(value string) bool {
	switch value {
	case "system", "light", "dark":
		return true
	default:
		return false
	}
}

// DefaultUserPreferences returns the default preferences for a user.
func DefaultUserPreferences(userID uint64) UserPreferences {
	return UserPreferences{
		UserID:          userID,
		ThemePreset:     "midnight",
		ThemeMode:       "system",
		CompactMode:     false,
		ReducedMotion:   false,
		CustomOverrides: map[string]string{},
	}
}
