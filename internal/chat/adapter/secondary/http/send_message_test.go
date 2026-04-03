package http_test

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"strings"
	"testing"

	chatdto "github.com/lechitz/aion-api/internal/chat/adapter/primary/http/dto"
	chathttp "github.com/lechitz/aion-api/internal/chat/adapter/secondary/http"
	"github.com/stretchr/testify/require"
)

type mockChatHTTPClient struct {
	doFn func(req *stdhttp.Request) (*stdhttp.Response, error)
}

func (m mockChatHTTPClient) Do(req *stdhttp.Request) (*stdhttp.Response, error) {
	return m.doFn(req)
}

func (m mockChatHTTPClient) Get(context.Context, string) (*stdhttp.Response, error) {
	return nil, errors.New("not implemented")
}

func (m mockChatHTTPClient) Post(context.Context, string, string, interface{}) (*stdhttp.Response, error) {
	return nil, errors.New("not implemented")
}

type mockChatHTTPLogger struct{}

func (mockChatHTTPLogger) Infof(string, ...any)                      {}
func (mockChatHTTPLogger) Errorf(string, ...any)                     {}
func (mockChatHTTPLogger) Debugf(string, ...any)                     {}
func (mockChatHTTPLogger) Warnf(string, ...any)                      {}
func (mockChatHTTPLogger) Infow(string, ...any)                      {}
func (mockChatHTTPLogger) Errorw(string, ...any)                     {}
func (mockChatHTTPLogger) Debugw(string, ...any)                     {}
func (mockChatHTTPLogger) Warnw(string, ...any)                      {}
func (mockChatHTTPLogger) InfowCtx(context.Context, string, ...any)  {}
func (mockChatHTTPLogger) ErrorwCtx(context.Context, string, ...any) {}
func (mockChatHTTPLogger) WarnwCtx(context.Context, string, ...any)  {}
func (mockChatHTTPLogger) DebugwCtx(context.Context, string, ...any) {}

type errReadCloser struct{}

func (errReadCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (errReadCloser) Close() error             { return nil }

func newRequest() *chatdto.InternalChatRequest {
	return &chatdto.InternalChatRequest{
		UserID:  1,
		Message: "hello",
	}
}

func TestNewClient_PanicsWhenHTTPClientNil(t *testing.T) {
	require.Panics(t, func() {
		_ = chathttp.New(nil, "http://aion-dev-chat:8000", mockChatHTTPLogger{})
	})
}

func TestSendMessage_Success(t *testing.T) {
	client := chathttp.New(mockChatHTTPClient{
		doFn: func(req *stdhttp.Request) (*stdhttp.Response, error) {
			require.Equal(t, stdhttp.MethodPost, req.Method)
			require.Contains(t, req.URL.String(), "/internal/process")
			return &stdhttp.Response{
				StatusCode: stdhttp.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"response":"ok","tokens_used":10}`)),
			}, nil
		},
	}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})

	resp, err := client.SendMessage(t.Context(), newRequest())
	require.NoError(t, err)
	require.Equal(t, "ok", resp.Response)
	require.Equal(t, 10, resp.TokensUsed)
}

func TestSendMessage_Errors(t *testing.T) {
	t.Run("marshal error", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return nil, errors.New("should not call")
			},
		}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})
		req := newRequest()
		req.Context = map[string]interface{}{"bad": func() {}}

		_, err := client.SendMessage(t.Context(), req)
		require.Error(t, err)
	})

	t.Run("request build error", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return nil, errors.New("should not call")
			},
		}, "://bad", mockChatHTTPLogger{})

		_, err := client.SendMessage(t.Context(), newRequest())
		require.Error(t, err)
	})

	t.Run("http do error", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return nil, errors.New("http down")
			},
		}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})

		_, err := client.SendMessage(t.Context(), newRequest())
		require.Error(t, err)
	})

	t.Run("read body error", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{StatusCode: stdhttp.StatusOK, Body: errReadCloser{}}, nil
			},
		}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})

		_, err := client.SendMessage(t.Context(), newRequest())
		require.Error(t, err)
	})

	t.Run("non 200", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{
					StatusCode: stdhttp.StatusBadGateway,
					Body:       io.NopCloser(strings.NewReader("upstream failed")),
				}, nil
			},
		}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})

		_, err := client.SendMessage(t.Context(), newRequest())
		require.Error(t, err)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		client := chathttp.New(mockChatHTTPClient{
			doFn: func(_ *stdhttp.Request) (*stdhttp.Response, error) {
				return &stdhttp.Response{
					StatusCode: stdhttp.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{invalid-json")),
				}, nil
			},
		}, "http://aion-dev-chat:8000", mockChatHTTPLogger{})

		_, err := client.SendMessage(t.Context(), newRequest())
		require.Error(t, err)
	})
}
