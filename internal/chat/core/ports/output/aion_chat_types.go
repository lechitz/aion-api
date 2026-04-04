package output

import "github.com/lechitz/aion-api/internal/chat/core/domain"

// ConversationMessage represents a single message in the conversation history.
type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SendMessageRequest represents the request sent to the Aion-Chat service.
type SendMessageRequest struct {
	UserID              uint64                   `json:"user_id"`
	Message             string                   `json:"message"`
	ConversationHistory []ConversationMessage    `json:"conversation_history,omitempty"`
	Context             map[string]interface{}   `json:"context,omitempty"`
	Runtime             *domain.RuntimeSelection `json:"runtime,omitempty"`
}

// SendMessageResponse represents the response from the Aion-Chat service.
type SendMessageResponse struct {
	Response      string                   `json:"response"`
	UI            map[string]interface{}   `json:"ui,omitempty"`
	FunctionCalls []FunctionCall           `json:"function_calls,omitempty"`
	TokensUsed    int                      `json:"tokens_used,omitempty"`
	Sources       []map[string]interface{} `json:"sources,omitempty"`
}

// FunctionCall represents a function that was called by the LLM.
type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}
