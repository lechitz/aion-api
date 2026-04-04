// Package output defines the outbound ports (interfaces to external services).
package output

import (
	"context"
)

// AionChatClient defines the interface for communicating with the Aion-Chat service (Python).
type AionChatClient interface {
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
}
