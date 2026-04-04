// Package usecase_test contains tests for chat use cases.
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	auditdomain "github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProcessMessage_Success(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(456)
	message := TestMessageHello

	// Mock: Cache returns empty history (no error)
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return([]domain.ChatHistory{}, nil)

	// Mock: AionChatClient returns successful response
	expectedResponse := &dto.InternalChatResponse{
		Response:   TestResponseWellThanks,
		TokensUsed: 25,
		Sources: []map[string]interface{}{
			{TestSourceType: TestSourceType, "id": TestSourceID},
		},
		FunctionCalls: []dto.FunctionCall{
			{Name: TestFunctionGetWeather, Args: map[string]interface{}{"city": "NYC"}},
		},
	}
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	// Mock: SaveChatHistory (async goroutine) - repository and cache
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TestResponseWellThanks, result.Response)
	require.Equal(t, 25, result.TokensUsed)
	require.Len(t, result.Sources, 1)
	require.Len(t, result.FunctionCalls, 1)
	require.Equal(t, TestFunctionGetWeather, result.FunctionCalls[0])

	// Allow async SaveChatHistory goroutine to complete
	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_Success_WithConversationHistory(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(789)
	message := TestMessageWhatAskedBefore

	// Mock: Cache returns recent history (DESC order - newest first)
	cachedHistory := []domain.ChatHistory{
		{
			ChatID:   101,
			UserID:   userID,
			Message:  TestHistoryQuestionInventor,
			Response: TestHistoryAnswerInventor,
		},
		{
			ChatID:   100,
			UserID:   userID,
			Message:  TestHistoryQuestionAI,
			Response: TestHistoryAnswerAI,
		},
	}
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return(cachedHistory, nil)

	// Mock: AionChatClient receives request with history
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error) {
			// Verify conversation history was included (4 messages: 2 user + 2 assistant)
			require.Len(t, req.ConversationHistory, 4)
			require.Equal(t, TestRoleUser, req.ConversationHistory[0].Role)
			require.Equal(t, TestHistoryQuestionAI, req.ConversationHistory[0].Content)
			require.Equal(t, TestRoleAssistant, req.ConversationHistory[1].Role)
			require.Equal(t, TestHistoryAnswerAI, req.ConversationHistory[1].Content)

			return &dto.InternalChatResponse{
				Response:   TestResponseAIInventor,
				TokensUsed: 15,
			}, nil
		})

	// Mock: SaveChatHistory (async) - allow it to succeed
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TestResponseAIInventor, result.Response)

	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_Success_WithoutHistory_CacheError(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	message := TestMessageFirst

	// Mock: Cache fails (non-blocking)
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return(nil, errors.New(TestErrorCacheUnavailable))

	// Mock: AionChatClient still processes (without history)
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error) {
			// Verify history is empty when cache fails
			require.Empty(t, req.ConversationHistory)
			return &dto.InternalChatResponse{
				Response:   TestResponseHelloHelp,
				TokensUsed: 10,
			}, nil
		})

	// Mock: SaveChatHistory (async) - allow it to succeed
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert: Should succeed despite cache error (non-blocking)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TestResponseHelloHelp, result.Response)

	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_AionChatClientError(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(111)
	message := "Test message"

	// Mock: Cache succeeds
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return([]domain.ChatHistory{}, nil)

	// Mock: AionChatClient fails
	clientError := errors.New(TestErrorServiceUnavailable)
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(nil, clientError)

	// No SaveChatHistory mocks needed (never reaches that point)

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert: Should propagate error
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, clientError, err)
}

func TestProcessMessage_EmptyMessage(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(222)
	message := TestMessageEmpty

	// Mock: Cache succeeds
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return([]domain.ChatHistory{}, nil)

	// Mock: AionChatClient processes empty message (validation happens at handler layer)
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(&dto.InternalChatResponse{
			Response:   TestResponseNoMessage,
			TokensUsed: 5,
		}, nil)

	// Mock: SaveChatHistory (async)
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Contains(t, result.Response, "didn't receive")

	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_WithMultipleFunctionCalls(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(333)
	message := TestMessageWaterIntake

	// Mock: Cache returns empty
	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return([]domain.ChatHistory{}, nil)

	// Mock: AionChatClient returns multiple function calls
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(&dto.InternalChatResponse{
			Response:   TestResponseWaterReminder,
			TokensUsed: 40,
			FunctionCalls: []dto.FunctionCall{
				{Name: TestFunctionGetWaterIntake, Args: map[string]interface{}{"date": "today"}},
				{Name: TestFunctionSetReminder, Args: map[string]interface{}{"time": "15:00", "message": "Drink water"}},
			},
		}, nil)

	// Mock: SaveChatHistory (async)
	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	// Execute
	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.FunctionCalls, 2)
	require.Equal(t, TestFunctionGetWaterIntake, result.FunctionCalls[0])
	require.Equal(t, TestFunctionSetReminder, result.FunctionCalls[1])

	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_ForwardsRuntimeSelection(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(334)
	message := "use this runtime"
	runtime := &domain.RuntimeSelection{
		Provider: "openai",
		Model:    "gpt-5.4-mini",
	}

	suite.HistoryCache.EXPECT().
		GetLatest(gomock.Any(), userID, 6).
		Return([]domain.ChatHistory{}, nil)

	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error) {
			require.NotNil(t, req.Runtime)
			require.Equal(t, runtime.Provider, req.Runtime.Provider)
			require.Equal(t, runtime.Model, req.Runtime.Model)
			return &dto.InternalChatResponse{
				Response: "ok",
			}, nil
		})

	suite.HistoryRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(domain.ChatHistory{}, nil).
		AnyTimes()
	suite.HistoryCache.EXPECT().
		Add(gomock.Any(), userID, gomock.Any()).
		Return(nil).
		AnyTimes()

	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, message, nil, runtime)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "ok", result.Response)

	time.Sleep(10 * time.Millisecond)
}

func TestProcessMessage_UIActionPersistsAuditEvent(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(444)
	message := "confirm this draft"
	ctx := context.WithValue(suite.Ctx, ctxkeys.TraceID, "trace-123")
	ctx = context.WithValue(ctx, ctxkeys.RequestID, "req-123")
	requestContext := map[string]interface{}{
		"ui_action": map[string]interface{}{
			"type":     "draft_accept",
			"draft_id": "draft-xyz",
			"consent": map[string]interface{}{
				"required":       true,
				"confirmed":      true,
				"policy_version": "consent-v1",
			},
			"quick_add": map[string]interface{}{
				"contract_version": "quick-add-v1",
				"entity":           "category",
				"operation":        "create",
				"idempotency_key":  "qa-1",
			},
		},
	}

	suite.HistoryCache.EXPECT().GetLatest(gomock.Any(), userID, 6).Return([]domain.ChatHistory{}, nil)
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(&dto.InternalChatResponse{
			Response: "ok",
			UI: map[string]interface{}{
				"action_result": map[string]interface{}{
					"status":       "success",
					"entity_id":    "10293",
					"message_code": "action.success",
				},
			},
		}, nil)
	suite.AuditService.EXPECT().
		WriteEvent(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, event auditdomain.AuditActionEvent) error {
			require.Equal(t, userID, event.UserID)
			require.Equal(t, "draft_accept", event.UIActionType)
			require.Equal(t, "draft-xyz", event.DraftID)
			require.Equal(t, "success", event.Status)
			require.Equal(t, "trace-123", event.TraceID)
			require.Equal(t, "req-123", event.RequestID)
			require.Equal(t, "category", event.Entity)
			require.Equal(t, "create", event.Operation)
			require.Equal(t, "action.success", event.MessageCode)
			require.False(t, event.TimestampUTC.IsZero())
			return nil
		})
	suite.HistoryRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.ChatHistory{}, nil).AnyTimes()
	suite.HistoryCache.EXPECT().Add(gomock.Any(), userID, gomock.Any()).Return(nil).AnyTimes()

	result, err := suite.ChatService.ProcessMessage(ctx, userID, message, requestContext, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestProcessMessage_AuditFailureDoesNotFailBusinessFlow(t *testing.T) {
	suite := setup.ChatServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(445)
	requestContext := map[string]interface{}{
		"ui_action": map[string]interface{}{
			"type":     "draft_cancel",
			"draft_id": "draft-abcd",
		},
	}

	suite.HistoryCache.EXPECT().GetLatest(gomock.Any(), userID, 6).Return([]domain.ChatHistory{}, nil)
	suite.AionChatClient.EXPECT().
		SendMessage(gomock.Any(), gomock.Any()).
		Return(&dto.InternalChatResponse{Response: "ok"}, nil)
	suite.AuditService.EXPECT().
		WriteEvent(gomock.Any(), gomock.Any()).
		Return(errors.New(TestErrorDatabaseTimeout))
	suite.HistoryRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.ChatHistory{}, nil).AnyTimes()
	suite.HistoryCache.EXPECT().Add(gomock.Any(), userID, gomock.Any()).Return(nil).AnyTimes()

	result, err := suite.ChatService.ProcessMessage(suite.Ctx, userID, "cancel", requestContext, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "ok", result.Response)
}
