package mapper

import (
	"encoding/json"

	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// UserPreferencesFromDB converts a model.UserPreferencesDB into a domain.UserPreferences.
func UserPreferencesFromDB(p model.UserPreferencesDB) domain.UserPreferences {
	overrides := make(map[string]string)
	if p.CustomOverrides != "" && p.CustomOverrides != "{}" {
		_ = json.Unmarshal([]byte(p.CustomOverrides), &overrides)
	}

	return domain.UserPreferences{
		UserID:          p.UserID,
		ThemePreset:     p.ThemePreset,
		ThemeMode:       p.ThemeMode,
		CompactMode:     p.CompactMode,
		ReducedMotion:   p.ReducedMotion,
		CustomOverrides: overrides,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// UserPreferencesToDB converts a domain.UserPreferences into a model.UserPreferencesDB.
func UserPreferencesToDB(p domain.UserPreferences) model.UserPreferencesDB {
	overridesJSON := "{}"
	if len(p.CustomOverrides) > 0 {
		if b, err := json.Marshal(p.CustomOverrides); err == nil {
			overridesJSON = string(b)
		}
	}

	return model.UserPreferencesDB{
		UserID:          p.UserID,
		ThemePreset:     p.ThemePreset,
		ThemeMode:       p.ThemeMode,
		CompactMode:     p.CompactMode,
		ReducedMotion:   p.ReducedMotion,
		CustomOverrides: overridesJSON,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}
