// Package dto contains data transfer objects for chat HTTP requests/responses.
package dto

import (
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	outputport "github.com/lechitz/aion-api/internal/chat/core/ports/output"
)

// ChatRequest represents the incoming chat message from the client.
type ChatRequest struct {
	Message string                 `json:"message"           validate:"required,min=1,max=2000" example:"Quanto de água eu bebi hoje?"`
	Context map[string]interface{} `json:"context,omitempty"`
	Runtime *ChatRuntimeSelection  `json:"runtime,omitempty"`
}

// ChatRuntimeSelection represents the requested LLM runtime selection.
type ChatRuntimeSelection = domain.RuntimeSelection

// ChatResponse represents the response returned to the client.
type ChatResponse struct {
	Response string                   `json:"response"          example:"Você bebeu 2.5 litros de água hoje..."`
	UI       map[string]interface{}   `json:"ui,omitempty"`
	Sources  []map[string]interface{} `json:"sources,omitempty"`
	Usage    *TokenUsage              `json:"usage,omitempty"`
}

// ChatCancelRequest represents a cancel request for active chat processing.
type ChatCancelRequest struct{}

// ChatCancelResponse represents cancel response status.
type ChatCancelResponse struct {
	Cancelled bool   `json:"cancelled"`
	Message   string `json:"message"`
}

// TokenUsage represents LLM token consumption statistics.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"     example:"50"`
	CompletionTokens int `json:"completion_tokens,omitempty" example:"100"`
	TotalTokens      int `json:"total_tokens"                example:"150"`
}

// ConversationMessage represents a single message in the conversation history.
type ConversationMessage = outputport.ConversationMessage

// InternalChatRequest represents the request sent to the Aion-Chat service (Python).
type InternalChatRequest = outputport.SendMessageRequest

// InternalChatResponse represents the response from the Aion-Chat service.
type InternalChatResponse = outputport.SendMessageResponse

// FunctionCall represents a function that was called by the LLM.
type FunctionCall = outputport.FunctionCall
