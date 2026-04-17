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

// AllowedThemePresets defines valid preset names.
var AllowedThemePresets = map[string]bool{
	"default":   true,
	"midnight":  true,
	"warm-gray": true,
	"ocean":     true,
	"frost":     true,
	"ember":     true,
	"charcoal":  true,
	"custom":    true,
}

// AllowedThemeModes defines valid theme mode values.
var AllowedThemeModes = map[string]bool{
	"system": true,
	"light":  true,
	"dark":   true,
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
