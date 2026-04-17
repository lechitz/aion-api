package handler_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	handler "github.com/lechitz/aion-api/internal/user/adapter/primary/http/handler"
	"github.com/stretchr/testify/require"
)

func TestNewAndRegisterHTTP(t *testing.T) {
	svc := &mockUserService{}
	h := handler.New(svc, &config.Config{}, mockLogger{})
	require.NotNil(t, h)

	r := &mockRouter{}
	handler.RegisterHTTP(r, h, nil, nil, mockLogger{})
	require.Equal(t, []string{"/user", "/registration"}, r.groups)
	require.Equal(t, []string{"/create", "/avatar/upload", "/start", "/{registration_id}/complete"}, r.posts)
	require.Equal(t, []string{"/{registration_id}/profile", "/{registration_id}/avatar"}, r.puts)
	require.Equal(t, 0, r.groupWithCalls)

	r = &mockRouter{}
	handler.RegisterHTTP(r, h, nil, mockAuthService{}, mockLogger{})
	require.Equal(t, []string{"/user", "/registration"}, r.groups)
	require.Equal(t, []string{"/create", "/avatar/upload", "/start", "/{registration_id}/complete"}, r.posts)
	require.Equal(t, []string{"/all", "/me", "/{user_id}"}, r.gets)
	require.Equal(t, []string{"/", "/password", "/{registration_id}/profile", "/{registration_id}/avatar"}, r.puts)
	require.Equal(t, []string{"/avatar", "/"}, r.deletes)
	require.Equal(t, 1, r.groupWithCalls)

	r = &mockRouter{}
	handler.RegisterHTTP(r, h, &handler.PreferencesHandler{}, mockAuthService{}, mockLogger{})
	require.Equal(t, []string{"/all", "/me", "/{user_id}", "/preferences"}, r.gets)
	require.Equal(t, []string{"/", "/password", "/preferences", "/{registration_id}/profile", "/{registration_id}/avatar"}, r.puts)
}
