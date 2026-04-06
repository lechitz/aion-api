// Package input defines the inbound ports (use cases) for the chat module.
package input

import (
	"context"

	"github.com/lechitz/aion-api/internal/chat/core/domain"
)

// ChatService defines the interface for chat operations (use cases).
type ChatService interface {
	// ProcessMessage sends a message to the AI and returns the response.
	ProcessMessage(
		ctx context.Context,
		userID uint64,
		message string,
		requestContext map[string]interface{},
		runtime *domain.RuntimeSelection,
	) (*domain.ChatResult, error)

	// SaveChatHistory persists a chat interaction to the database.
	SaveChatHistory(ctx context.Context, userID uint64, message, response string, tokensUsed int, functionCalls map[string]string) error

	// GetChatHistory retrieves chat history for a user with pagination.
	GetChatHistory(ctx context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error)

	// GetLatestChatHistory retrieves the N most recent chat entries for a user.
	GetLatestChatHistory(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error)

	// GetChatContext retrieves aggregated context for AI including recent activity.
	GetChatContext(ctx context.Context, userID uint64) (*domain.ChatContext, error)
}
