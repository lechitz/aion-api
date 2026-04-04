// Package usecase contains the business logic for the Chat context.
package usecase

import (
	"context"
	"encoding/json"

	"github.com/lechitz/aion-api/internal/chat/core/domain"
	outputport "github.com/lechitz/aion-api/internal/chat/core/ports/output"
)

// fetchConversationHistory retrieves recent chat history from cache (Redis).
// It fetches the last N messages and converts them to conversation format
// with proper role ordering (oldest first, as required by LLMs).
// Non-blocking: returns empty slice on cache error.
func (s *ChatService) fetchConversationHistory(ctx context.Context, userID uint64, limit int) []outputport.ConversationMessage {
	conversationHistory := []outputport.ConversationMessage{}

	history, err := s.chatHistoryCache.GetLatest(ctx, userID, limit)
	if err != nil {
		// Don't fail on cache retrieval error, just log and continue without history
		s.logger.WarnwCtx(ctx, "Failed to retrieve conversation history from cache (non-blocking)",
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return conversationHistory
	}

	if len(history) == 0 {
		return conversationHistory
	}

	// Convert history to conversation messages format (reverse order - oldest first)
	// History comes DESC (newest first), we need ASC (oldest first) for LLM context
	for i := len(history) - 1; i >= 0; i-- {
		h := history[i]
		// Add user message
		conversationHistory = append(conversationHistory, outputport.ConversationMessage{
			Role:    "user",
			Content: h.Message,
		})
		// Add assistant response
		conversationHistory = append(conversationHistory, outputport.ConversationMessage{
			Role:    "assistant",
			Content: h.Response,
		})
	}

	s.logger.InfowCtx(ctx, "Including conversation history from cache",
		LogKeyUserID, userID,
		"history_messages", len(conversationHistory),
	)

	return conversationHistory
}

// buildChatRequest creates an InternalChatRequest from user input and conversation history.
func buildChatRequest(
	userID uint64,
	message string,
	history []outputport.ConversationMessage,
	context map[string]interface{},
	runtime *domain.RuntimeSelection,
) *outputport.SendMessageRequest {
	payloadContext := map[string]interface{}{
		ContextKeyTimezone:     DefaultTimezone, // Backward-compatible key
		ContextKeyUserTimezone: DefaultTimezone, // Canonical key for aion-chat
	}
	for key, value := range context {
		payloadContext[key] = value
	}
	return &outputport.SendMessageRequest{
		UserID:              userID,
		Message:             message,
		ConversationHistory: history,
		Context:             payloadContext,
		Runtime:             runtime,
	}
}

// saveChatInteraction persists the chat exchange (message + response) to database and cache.
// Converts FunctionCalls to a map with JSON-encoded arguments for storage.
// Non-blocking: designed to be called in a goroutine, logs errors without returning them.
func (s *ChatService) saveChatInteraction(ctx context.Context, userID uint64, message, response string, tokensUsed int, functionCalls []outputport.FunctionCall) {
	functionCallsMap := make(map[string]string)
	for _, call := range functionCalls {
		// Convert Args map to JSON string for storage
		argsJSON := ""
		if len(call.Args) > 0 {
			if jsonBytes, err := json.Marshal(call.Args); err == nil {
				argsJSON = string(jsonBytes)
			}
		}
		functionCallsMap[call.Name] = argsJSON
	}

	if err := s.SaveChatHistory(ctx, userID, message, response, tokensUsed, functionCallsMap); err != nil {
		s.logger.ErrorwCtx(ctx, "Failed to save chat history (non-blocking)",
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
	}
}
