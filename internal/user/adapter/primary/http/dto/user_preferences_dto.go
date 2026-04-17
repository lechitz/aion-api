package dto

import (
	"time"

	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
)

// UserPreferencesResponse represents the payload returned by GET /user/preferences.
type UserPreferencesResponse struct {
	ThemePreset     string            `json:"theme_preset"     example:"default"`
	ThemeMode       string            `json:"theme_mode"       example:"system"`
	CompactMode     bool              `json:"compact_mode"     example:"false"`
	ReducedMotion   bool              `json:"reduced_motion"   example:"false"`
	CustomOverrides map[string]string `json:"custom_overrides"`
	CreatedAt       time.Time         `json:"created_at"       example:"2024-01-02T15:04:05Z"`
	UpdatedAt       time.Time         `json:"updated_at"       example:"2024-01-02T15:04:05Z"`
}

// UserPreferencesResponseFromDomain maps a domain.UserPreferences to a UserPreferencesResponse.
func UserPreferencesResponseFromDomain(p domain.UserPreferences) UserPreferencesResponse {
	overrides := p.CustomOverrides
	if overrides == nil {
		overrides = map[string]string{}
	}
	return UserPreferencesResponse{
		ThemePreset:     p.ThemePreset,
		ThemeMode:       p.ThemeMode,
		CompactMode:     p.CompactMode,
		ReducedMotion:   p.ReducedMotion,
		CustomOverrides: overrides,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// SaveUserPreferencesRequest represents the payload for PUT /user/preferences.
type SaveUserPreferencesRequest struct {
	ThemePreset     string            `json:"theme_preset"`
	ThemeMode       string            `json:"theme_mode"`
	CompactMode     *bool             `json:"compact_mode"`
	ReducedMotion   *bool             `json:"reduced_motion"`
	CustomOverrides map[string]string `json:"custom_overrides"`
}

// ToCommand converts the request DTO to an input port command.
func (r SaveUserPreferencesRequest) ToCommand() input.SavePreferencesCommand {
	return input.SavePreferencesCommand{
		ThemePreset:     r.ThemePreset,
		ThemeMode:       r.ThemeMode,
		CompactMode:     r.CompactMode,
		ReducedMotion:   r.ReducedMotion,
		CustomOverrides: r.CustomOverrides,
	}
}
