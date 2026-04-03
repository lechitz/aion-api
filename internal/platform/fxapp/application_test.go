//nolint:testpackage // tests require unexported appDepsParams to exercise Fx wiring.
package fxapp

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/config"
	httpclientPort "github.com/lechitz/aion-api/internal/platform/ports/output/httpclient"
	"github.com/stretchr/testify/require"
)

type stubHTTPClient struct{}

func (stubHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

func (stubHTTPClient) Get(context.Context, string) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

func (stubHTTPClient) Post(context.Context, string, string, interface{}) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

type noopLoggerFx struct{}

func (noopLoggerFx) Infof(string, ...any)                      {}
func (noopLoggerFx) Errorf(string, ...any)                     {}
func (noopLoggerFx) Debugf(string, ...any)                     {}
func (noopLoggerFx) Warnf(string, ...any)                      {}
func (noopLoggerFx) Infow(string, ...any)                      {}
func (noopLoggerFx) Errorw(string, ...any)                     {}
func (noopLoggerFx) Debugw(string, ...any)                     {}
func (noopLoggerFx) Warnw(string, ...any)                      {}
func (noopLoggerFx) InfowCtx(context.Context, string, ...any)  {}
func (noopLoggerFx) ErrorwCtx(context.Context, string, ...any) {}
func (noopLoggerFx) WarnwCtx(context.Context, string, ...any)  {}
func (noopLoggerFx) DebugwCtx(context.Context, string, ...any) {}

func TestProvideAppDependencies(t *testing.T) {
	var httpClient httpclientPort.HTTPClient = stubHTTPClient{}

	deps := appDepsParams{
		Cfg: &config.Config{
			Secret: config.Secret{
				Key: "test-secret",
			},
			AionChat: config.AionChatConfig{
				BaseURL: "http://aion-dev-chat:8000",
			},
		},
		HTTPClient: httpClient,
		Log:        noopLoggerFx{},
	}

	got := ProvideAppDependencies(deps)
	require.NotNil(t, got)
	require.NotNil(t, got.AuthService)
	require.NotNil(t, got.UserService)
	require.NotNil(t, got.AdminService)
	require.NotNil(t, got.CategoryService)
	require.NotNil(t, got.TagService)
	require.NotNil(t, got.RecordService)
	require.NotNil(t, got.ChatService)
	require.NotNil(t, got.RealtimeService)
}

func TestProvideHTTPClientAndKeyGenerator(t *testing.T) {
	cfg := &config.Config{
		AionChat: config.AionChatConfig{
			Timeout: 2 * time.Second,
		},
	}

	client := ProvideHTTPClient(cfg)
	require.NotNil(t, client)

	keyGen := ProvideKeyGenerator()
	require.NotNil(t, keyGen)
	key, err := keyGen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, key)
}
