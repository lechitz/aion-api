// Package usecase implements the chat use cases (business logic).
package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/chat/core/domain"
	outputport "github.com/lechitz/aion-api/internal/chat/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ProcessMessage processes a chat message by forwarding it to the Aion-Chat service.
func (s *ChatService) ProcessMessage(ctx context.Context, userID uint64, message string, requestContext map[string]interface{}, runtime *domain.RuntimeSelection) (*domain.ChatResult, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanProcessMessage)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttrMessageLength, len(message)),
	)

	s.logger.InfowCtx(ctx, LogProcessingChatMessage, LogKeyUserID, userID, LogKeyMessageLength, len(message))
	if uiActionType, draftID := extractUIActionMetadata(requestContext); uiActionType != "" {
		s.logger.InfowCtx(
			ctx,
			LogChatRequestIncludesUIAction,
			LogKeyUserID, userID,
			LogKeyUIActionType, uiActionType,
			LogKeyDraftID, draftID,
		)
	}

	// Fetch recent conversation history from cache (6 messages = 3 exchanges for context)
	const historyLimit = 6
	conversationHistory := s.fetchConversationHistory(ctx, userID, historyLimit)

	// Build request with conversation context
	req := buildChatRequest(userID, message, conversationHistory, requestContext, runtime)

	// Call external Aion-Chat service
	resp, err := s.aionChatClient.SendMessage(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, StatusFailedToCallAionChat)
		s.logger.ErrorwCtx(ctx, LogFailedToCallAionChat,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return nil, err
	}

	// Build result
	result := &domain.ChatResult{
		Response:      resp.Response,
		UI:            resp.UI,
		Sources:       convertSources(resp.Sources),
		TokensUsed:    resp.TokensUsed,
		FunctionCalls: extractFunctionNames(resp.FunctionCalls),
	}

	span.SetAttributes(
		attribute.Int(AttrTokensUsed, resp.TokensUsed),
		attribute.Int(AttrFunctionCallsCount, len(resp.FunctionCalls)),
	)
	span.SetStatus(codes.Ok, StatusMessageProcessedSuccessfully)

	s.logger.InfowCtx(ctx, LogChatMessageProcessedSuccessfully,
		LogKeyUserID, userID,
		LogKeyTokensUsed, resp.TokensUsed,
		LogKeyResponseLength, len(resp.Response),
	)

	// Best-effort sync audit persistence for UI actions.
	s.persistAuditActionEvent(ctx, userID, requestContext, result)

	// Save chat history asynchronously while preserving request values without inheriting cancellation.
	go s.saveChatInteraction(context.WithoutCancel(ctx), userID, message, resp.Response, resp.TokensUsed, resp.FunctionCalls)

	return result, nil
}

func extractUIActionMetadata(requestContext map[string]interface{}) (string, string) {
	if requestContext == nil {
		return "", ""
	}
	uiAction, ok := requestContext[ContextKeyUIAction].(map[string]interface{})
	if !ok || uiAction == nil {
		return "", ""
	}
	actionType, _ := uiAction[ContextKeyUIActionType].(string)
	draftID, _ := uiAction[ContextKeyDraftID].(string)
	return actionType, draftID
}

// convertSources converts the sources from the internal response to domain format.
func convertSources(sources []map[string]interface{}) []interface{} {
	if sources == nil {
		return nil
	}
	result := make([]interface{}, len(sources))
	for i, source := range sources {
		result[i] = source
	}
	return result
}

// extractFunctionNames extracts function names from function calls.
func extractFunctionNames(calls []outputport.FunctionCall) []string {
	if calls == nil {
		return nil
	}
	names := make([]string, len(calls))
	for i, call := range calls {
		names[i] = call.Name
	}
	return names
}
