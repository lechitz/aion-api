package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/chat/adapter/primary/graphql/controller"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	chatinput "github.com/lechitz/aion-api/internal/chat/core/ports/input"
	"github.com/stretchr/testify/require"
)

type chatServiceStub struct {
	getHistoryFn func(context.Context, uint64, int, int) ([]domain.ChatHistory, error)
	getContextFn func(context.Context, uint64) (*domain.ChatContext, error)
}

func (s *chatServiceStub) ProcessMessage(context.Context, uint64, string, map[string]interface{}, *domain.RuntimeSelection) (*domain.ChatResult, error) {
	panic("unexpected ProcessMessage call")
}

func (s *chatServiceStub) SaveChatHistory(context.Context, uint64, string, string, int, map[string]string) error {
	panic("unexpected SaveChatHistory call")
}

func (s *chatServiceStub) GetChatHistory(ctx context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error) {
	if s.getHistoryFn == nil {
		panic("unexpected GetChatHistory call")
	}
	return s.getHistoryFn(ctx, userID, limit, offset)
}

func (s *chatServiceStub) GetLatestChatHistory(context.Context, uint64, int) ([]domain.ChatHistory, error) {
	panic("unexpected GetLatestChatHistory call")
}

func (s *chatServiceStub) GetChatContext(ctx context.Context, userID uint64) (*domain.ChatContext, error) {
	if s.getContextFn == nil {
		panic("unexpected GetChatContext call")
	}
	return s.getContextFn(ctx, userID)
}

type chatLoggerStub struct{}

func (chatLoggerStub) Infof(string, ...any)                      {}
func (chatLoggerStub) Errorf(string, ...any)                     {}
func (chatLoggerStub) Debugf(string, ...any)                     {}
func (chatLoggerStub) Warnf(string, ...any)                      {}
func (chatLoggerStub) Infow(string, ...any)                      {}
func (chatLoggerStub) Errorw(string, ...any)                     {}
func (chatLoggerStub) Debugw(string, ...any)                     {}
func (chatLoggerStub) Warnw(string, ...any)                      {}
func (chatLoggerStub) InfowCtx(context.Context, string, ...any)  {}
func (chatLoggerStub) ErrorwCtx(context.Context, string, ...any) {}
func (chatLoggerStub) WarnwCtx(context.Context, string, ...any)  {}
func (chatLoggerStub) DebugwCtx(context.Context, string, ...any) {}

func TestChatController_Basics(t *testing.T) {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)

	ctrl := controller.NewController(&chatServiceStub{
		getHistoryFn: func(_ context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error) {
			require.Equal(t, uint64(8), userID)
			require.Equal(t, 5, limit)
			require.Equal(t, 2, offset)
			return []domain.ChatHistory{{ChatID: 1, UserID: userID, Message: "m", Response: "r", TokensUsed: 12, CreatedAt: now, UpdatedAt: now}}, nil
		},
		getContextFn: func(_ context.Context, userID uint64) (*domain.ChatContext, error) {
			require.Equal(t, uint64(8), userID)
			return &domain.ChatContext{
				RecentChats:     []domain.ChatHistory{{ChatID: 2, UserID: userID, Message: "x", Response: "y", CreatedAt: now, UpdatedAt: now}},
				TotalRecords:    10,
				TotalCategories: 4,
				TotalTags:       3,
			}, nil
		},
	}, chatLoggerStub{})

	history, err := ctrl.GetChatHistory(t.Context(), 8, 5, 2)
	require.NoError(t, err)
	require.Len(t, history, 1)
	require.Equal(t, "1", history[0].ID)
	require.Equal(t, int32(12), history[0].TokensUsed)

	ctxData, err := ctrl.GetChatContext(t.Context(), 8)
	require.NoError(t, err)
	require.Equal(t, int32(10), ctxData.TotalRecords)
	require.Equal(t, int32(4), ctxData.TotalCategories)
	require.Equal(t, int32(3), ctxData.TotalTags)
	require.Len(t, ctxData.RecentChats, 1)
}

func TestChatController_ErrorsAndSafeInt32Overflow(t *testing.T) {
	ctrl := controller.NewController(&chatServiceStub{
		getHistoryFn: func(context.Context, uint64, int, int) ([]domain.ChatHistory, error) {
			return nil, errors.New("history failed")
		},
		getContextFn: func(context.Context, uint64) (*domain.ChatContext, error) {
			return &domain.ChatContext{TotalRecords: 1 << 40, TotalCategories: 1 << 40, TotalTags: 1 << 40}, nil
		},
	}, chatLoggerStub{})

	_, err := ctrl.GetChatHistory(t.Context(), 1, 10, 0)
	require.EqualError(t, err, "history failed")

	ctxData, err := ctrl.GetChatContext(t.Context(), 1)
	require.NoError(t, err)
	require.Equal(t, int32(0), ctxData.TotalRecords)
	require.Equal(t, int32(0), ctxData.TotalCategories)
	require.Equal(t, int32(0), ctxData.TotalTags)

	ctrlErr := controller.NewController(&chatServiceStub{
		getHistoryFn: func(context.Context, uint64, int, int) ([]domain.ChatHistory, error) { return nil, nil },
		getContextFn: func(context.Context, uint64) (*domain.ChatContext, error) { return nil, errors.New("context failed") },
	}, chatLoggerStub{})

	_, err = ctrlErr.GetChatContext(t.Context(), 1)
	require.EqualError(t, err, "context failed")
}

var _ chatinput.ChatService = (*chatServiceStub)(nil)
