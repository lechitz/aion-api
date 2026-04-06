package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authdomain "github.com/lechitz/aion-api/internal/auth/core/domain"
	authinput "github.com/lechitz/aion-api/internal/auth/core/ports/input"
	"github.com/lechitz/aion-api/internal/chat/adapter/primary/http/dto"
	handler "github.com/lechitz/aion-api/internal/chat/adapter/primary/http/handler"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	chatinput "github.com/lechitz/aion-api/internal/chat/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

type mockChatService struct {
	processFn func(ctx context.Context, userID uint64, message string, requestContext map[string]interface{}, runtime *dto.ChatRuntimeSelection) (*domain.ChatResult, error)
}

func (m mockChatService) ProcessMessage(
	ctx context.Context,
	userID uint64,
	message string,
	requestContext map[string]interface{},
	runtime *dto.ChatRuntimeSelection,
) (*domain.ChatResult, error) {
	if m.processFn != nil {
		return m.processFn(ctx, userID, message, requestContext, runtime)
	}
	return &domain.ChatResult{}, nil
}

func (mockChatService) SaveChatHistory(context.Context, uint64, string, string, int, map[string]string) error {
	return nil
}

func (mockChatService) GetChatHistory(context.Context, uint64, int, int) ([]domain.ChatHistory, error) {
	return nil, nil
}

func (mockChatService) GetLatestChatHistory(context.Context, uint64, int) ([]domain.ChatHistory, error) {
	return nil, nil
}

func (mockChatService) GetChatContext(context.Context, uint64) (*domain.ChatContext, error) {
	return &domain.ChatContext{}, nil
}

type mockLogger struct{}

func (mockLogger) Infof(string, ...any)                      {}
func (mockLogger) Errorf(string, ...any)                     {}
func (mockLogger) Debugf(string, ...any)                     {}
func (mockLogger) Warnf(string, ...any)                      {}
func (mockLogger) Infow(string, ...any)                      {}
func (mockLogger) Errorw(string, ...any)                     {}
func (mockLogger) Debugw(string, ...any)                     {}
func (mockLogger) Warnw(string, ...any)                      {}
func (mockLogger) InfowCtx(context.Context, string, ...any)  {}
func (mockLogger) ErrorwCtx(context.Context, string, ...any) {}
func (mockLogger) WarnwCtx(context.Context, string, ...any)  {}
func (mockLogger) DebugwCtx(context.Context, string, ...any) {}

type capturedLogEntry struct {
	message string
	keyvals []any
}

type capturingLogger struct {
	entries []capturedLogEntry
}

func (l *capturingLogger) Infof(string, ...any)                      {}
func (l *capturingLogger) Errorf(string, ...any)                     {}
func (l *capturingLogger) Debugf(string, ...any)                     {}
func (l *capturingLogger) Warnf(string, ...any)                      {}
func (l *capturingLogger) Infow(string, ...any)                      {}
func (l *capturingLogger) Errorw(string, ...any)                     {}
func (l *capturingLogger) Debugw(string, ...any)                     {}
func (l *capturingLogger) Warnw(string, ...any)                      {}
func (l *capturingLogger) ErrorwCtx(context.Context, string, ...any) {}
func (l *capturingLogger) WarnwCtx(context.Context, string, ...any)  {}
func (l *capturingLogger) DebugwCtx(context.Context, string, ...any) {}
func (l *capturingLogger) InfowCtx(_ context.Context, msg string, keyvals ...any) {
	l.entries = append(l.entries, capturedLogEntry{message: msg, keyvals: keyvals})
}

type mockAuthService struct{}

func (mockAuthService) Validate(context.Context, string) (uint64, map[string]any, error) {
	return 0, nil, nil
}

func (mockAuthService) Login(context.Context, string, string) (authdomain.AuthenticatedUser, string, string, error) {
	return authdomain.AuthenticatedUser{}, "", "", nil
}
func (mockAuthService) Logout(context.Context, uint64) error { return nil }
func (mockAuthService) RefreshTokenRenewal(context.Context, string) (string, string, error) {
	return "", "", nil
}

type mockRouter struct {
	groupWithCalls int
	posts          []string
}

func (m *mockRouter) Use(...ports.Middleware)          {}
func (m *mockRouter) Group(string, func(ports.Router)) {}
func (m *mockRouter) GroupWith(_ ports.Middleware, fn func(ports.Router)) {
	m.groupWithCalls++
	fn(m)
}
func (m *mockRouter) Mount(string, http.Handler)                               {}
func (m *mockRouter) Handle(string, string, http.Handler)                      {}
func (m *mockRouter) GET(string, http.Handler)                                 {}
func (m *mockRouter) POST(path string, _ http.Handler)                         { m.posts = append(m.posts, path) }
func (m *mockRouter) PUT(string, http.Handler)                                 {}
func (m *mockRouter) DELETE(string, http.Handler)                              {}
func (m *mockRouter) SetNotFound(http.Handler)                                 {}
func (m *mockRouter) SetMethodNotAllowed(http.Handler)                         {}
func (m *mockRouter) SetError(func(http.ResponseWriter, *http.Request, error)) {}
func (m *mockRouter) ServeHTTP(http.ResponseWriter, *http.Request)             {}

func TestNewAndRegisterHTTP(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{}, mockLogger{})
	require.NotNil(t, h)

	router := &mockRouter{}
	handler.RegisterHTTP(router, h, nil, mockLogger{})
	require.Equal(t, 0, router.groupWithCalls)
	require.Empty(t, router.posts)

	handler.RegisterHTTP(router, h, mockAuthService{}, mockLogger{})
	require.Equal(t, 1, router.groupWithCalls)
	require.ElementsMatch(t, []string{"/chat/text", "/chat/cancel", "/chat/audio"}, router.posts)
}

func TestChatText_Success(t *testing.T) {
	h := handler.New(mockChatService{
		processFn: func(_ context.Context, userID uint64, message string, requestContext map[string]interface{}, runtime *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
			require.Equal(t, uint64(7), userID)
			require.Equal(t, "hello", message)
			require.Equal(t, "v", requestContext["k"])
			require.Nil(t, runtime)
			return &domain.ChatResult{
				Response:   "ok",
				UI:         map[string]interface{}{"type": "simple"},
				Sources:    []interface{}{map[string]interface{}{"id": 1}, "ignored"},
				TokensUsed: 12,
			}, nil
		},
	}, &config.Config{}, mockLogger{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello","context":{"k":"v"}}`))
	req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ChatText(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Chat processed successfully")
	require.Contains(t, rec.Body.String(), "\"total_tokens\":12")
	require.Contains(t, rec.Body.String(), "\"sources\":[{")
}

func TestChatText_Success_WithRuntimeSelection(t *testing.T) {
	h := handler.New(mockChatService{
		processFn: func(_ context.Context, userID uint64, message string, requestContext map[string]interface{}, runtime *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
			require.Equal(t, uint64(7), userID)
			require.Equal(t, "hello", message)
			require.NotNil(t, runtime)
			require.Equal(t, "openai", runtime.Provider)
			require.Equal(t, "gpt-5.4-mini", runtime.Model)
			require.Equal(t, "v", requestContext["k"])
			return &domain.ChatResult{Response: "ok"}, nil
		},
	}, &config.Config{}, mockLogger{})

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/chat/text",
		strings.NewReader(`{"message":"hello","context":{"k":"v"},"runtime":{"provider":"openai","model":"gpt-5.4-mini"}}`),
	)
	req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ChatText(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Chat processed successfully")
}

func TestChatText_AuthErrors(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{}, mockLogger{})

	t.Run("missing user id", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid user id type", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, "7"))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestChatText_RequestErrors(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{}, mockLogger{})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"   "}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error runtime provider", func(t *testing.T) {
		req := httptest.NewRequestWithContext(
			t.Context(),
			http.MethodPost,
			"/chat/text",
			strings.NewReader(`{"message":"hello","runtime":{"provider":" ","model":"gpt-5.4-mini"}}`),
		)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error runtime model", func(t *testing.T) {
		req := httptest.NewRequestWithContext(
			t.Context(),
			http.MethodPost,
			"/chat/text",
			strings.NewReader(`{"message":"hello","runtime":{"provider":"openai","model":" "}}`),
		)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestChatText_ServiceErrors(t *testing.T) {
	h := handler.New(mockChatService{processFn: func(context.Context, uint64, string, map[string]interface{}, *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
		return nil, errors.New("boom")
	}}, &config.Config{}, mockLogger{})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()
		h.ChatText(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("service validation error stays explicit", func(t *testing.T) {
		validationHandler := handler.New(mockChatService{
			processFn: func(context.Context, uint64, string, map[string]interface{}, *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
				return nil, sharederrors.NewValidationError("runtime", "Invalid runtime provider 'invalid-provider'")
			},
		}, &config.Config{}, mockLogger{})

		req := httptest.NewRequestWithContext(
			t.Context(),
			http.MethodPost,
			"/chat/text",
			strings.NewReader(`{"message":"hello","runtime":{"provider":"invalid-provider","model":"x"}}`),
		)
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()

		validationHandler.ChatText(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "validation error on runtime")
		require.Contains(t, rec.Body.String(), "Invalid runtime provider")
	})

	t.Run("service canceled", func(t *testing.T) {
		cancelHandler := handler.New(mockChatService{
			processFn: func(context.Context, uint64, string, map[string]interface{}, *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
				return nil, context.Canceled
			},
		}, &config.Config{}, mockLogger{})

		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/text", strings.NewReader(`{"message":"hello"}`))
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
		rec := httptest.NewRecorder()

		cancelHandler.ChatText(rec, req)
		require.Equal(t, 499, rec.Code)
		require.Contains(t, rec.Body.String(), "Chat request cancelled")
	})
}

func TestChatText_LogsUIActionMetadataWithConsent(t *testing.T) {
	logger := &capturingLogger{}
	h := handler.New(mockChatService{
		processFn: func(_ context.Context, userID uint64, message string, _ map[string]interface{}, runtime *dto.ChatRuntimeSelection) (*domain.ChatResult, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, "confirmar", message)
			require.Nil(t, runtime)
			return &domain.ChatResult{Response: "ok"}, nil
		},
	}, &config.Config{}, logger)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/chat/text",
		strings.NewReader(
			`{"message":"confirmar","context":{"ui_action":{"type":"draft_accept","draft_id":"draft-xyz","consent":{"required":true,"confirmed":true,"policy_version":"consent-v1"},"quick_add":{"contract_version":" quick-add-v1 ","entity":"category","operation":"create","idempotency_key":"qa-1"}}}}`,
		),
	)
	req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(9)))
	rec := httptest.NewRecorder()

	h.ChatText(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var uiActionLog *capturedLogEntry
	for idx := range logger.entries {
		entry := &logger.entries[idx]
		if entry.message == handler.MsgChatRequestIncludesUIAction {
			uiActionLog = entry
			break
		}
	}
	require.NotNil(t, uiActionLog, "expected ui_action metadata log entry")

	logMap := make(map[string]any)
	for i := 0; i+1 < len(uiActionLog.keyvals); i += 2 {
		key, ok := uiActionLog.keyvals[i].(string)
		if !ok {
			continue
		}
		logMap[key] = uiActionLog.keyvals[i+1]
	}

	require.Equal(t, "draft_accept", logMap[handler.LogKeyUIActionType])
	require.Equal(t, "draft-xyz", logMap[handler.LogKeyDraftID])
	require.Equal(t, true, logMap[handler.LogKeyConsentRequired])
	require.Equal(t, true, logMap[handler.LogKeyConsentConfirmed])
	require.Equal(t, "consent-v1", logMap[handler.LogKeyConsentPolicyVersion])
	require.Equal(t, "quick-add-v1", logMap[handler.LogKeyQuickAddContractVersion])
	require.Equal(t, "category", logMap[handler.LogKeyQuickAddEntity])
	require.Equal(t, "create", logMap[handler.LogKeyQuickAddOperation])
	require.Equal(t, "qa-1", logMap[handler.LogKeyQuickAddIdempotencyKey])
}

var (
	_ chatinput.ChatService = (*mockChatService)(nil)
	_ authinput.AuthService = (*mockAuthService)(nil)
)
