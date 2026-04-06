// Package dto contains DTOs for user HTTP responses.
package dto

import "time"

// UserMeResponse represents the payload returned by GET /user/me.
type UserMeResponse struct {
	ID                  uint64    `json:"id"                   example:"42"`
	Name                string    `json:"name"                 example:"Alice Doe"`
	Username            string    `json:"username"             example:"alice"`
	Email               string    `json:"email"                example:"alice@example.com"`
	CreatedAt           time.Time `json:"created_at"           example:"2024-01-02T15:04:05Z"`
	Locale              *string   `json:"locale,omitempty"     example:"en-US"`
	Timezone            *string   `json:"timezone,omitempty"   example:"America/Sao_Paulo"`
	Location            *string   `json:"location,omitempty"   example:"São Paulo, BR"`
	Bio                 *string   `json:"bio,omitempty"        example:"Backend engineer passionate about observability."`
	AvatarURL           *string   `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`
	OnboardingCompleted bool      `json:"onboarding_completed" example:"false"`
}
