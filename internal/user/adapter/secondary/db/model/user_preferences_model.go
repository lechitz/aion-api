package model

import "time"

const (
	// TableUserPreferences is the name of the database table for user preferences.
	TableUserPreferences = "aion_api.user_preferences"
)

// UserPreferencesDB represents the database model for user preferences.
type UserPreferencesDB struct {
	ThemePreset     string    `gorm:"column:theme_preset;default:midnight"`
	ThemeMode       string    `gorm:"column:theme_mode;default:system"`
	CustomOverrides string    `gorm:"column:custom_overrides;type:jsonb;default:'{}'"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	UserID          uint64    `gorm:"primaryKey;column:user_id"`
	CompactMode     bool      `gorm:"column:compact_mode;default:false"`
	ReducedMotion   bool      `gorm:"column:reduced_motion;default:false"`
}

// TableName specifies the custom database table name for the UserPreferencesDB model.
func (UserPreferencesDB) TableName() string {
	return TableUserPreferences
}
