package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"github.com/lechitz/aion-api/internal/user/core/usecase"
	"github.com/stretchr/testify/require"
)

type preferencesRepoStub struct {
	current domain.UserPreferences
	saved   domain.UserPreferences
}

func (r *preferencesRepoStub) GetPreferencesByUserID(context.Context, uint64) (domain.UserPreferences, error) {
	return r.current, nil
}

func (r *preferencesRepoStub) UpsertPreferences(_ context.Context, prefs domain.UserPreferences) (domain.UserPreferences, error) {
	r.saved = prefs
	r.saved.UpdatedAt = time.Now().UTC()
	return r.saved, nil
}

type preferencesLoggerStub struct{}

func (preferencesLoggerStub) Infof(string, ...any)                      {}
func (preferencesLoggerStub) Errorf(string, ...any)                     {}
func (preferencesLoggerStub) Debugf(string, ...any)                     {}
func (preferencesLoggerStub) Warnf(string, ...any)                      {}
func (preferencesLoggerStub) Infow(string, ...any)                      {}
func (preferencesLoggerStub) Errorw(string, ...any)                     {}
func (preferencesLoggerStub) Debugw(string, ...any)                     {}
func (preferencesLoggerStub) Warnw(string, ...any)                      {}
func (preferencesLoggerStub) InfowCtx(context.Context, string, ...any)  {}
func (preferencesLoggerStub) ErrorwCtx(context.Context, string, ...any) {}
func (preferencesLoggerStub) WarnwCtx(context.Context, string, ...any)  {}
func (preferencesLoggerStub) DebugwCtx(context.Context, string, ...any) {}

func TestPreferencesService_SavePreferencesPreservesOmittedBooleans(t *testing.T) {
	current := domain.DefaultUserPreferences(7)
	current.ThemePreset = "midnight"
	current.ThemeMode = "dark"
	current.CompactMode = true
	current.ReducedMotion = true
	current.CustomOverrides = map[string]string{"--surface": "#111111"}

	repo := &preferencesRepoStub{current: current}
	svc := usecase.NewPreferencesService(repo, preferencesLoggerStub{})

	saved, err := svc.SavePreferences(t.Context(), 7, input.SavePreferencesCommand{
		ThemePreset: "ocean",
	})

	require.NoError(t, err)
	require.Equal(t, uint64(7), saved.UserID)
	require.Equal(t, "ocean", saved.ThemePreset)
	require.Equal(t, "dark", saved.ThemeMode)
	require.True(t, saved.CompactMode)
	require.True(t, saved.ReducedMotion)
	require.Equal(t, map[string]string{"--surface": "#111111"}, saved.CustomOverrides)
	require.Equal(t, saved, repo.saved)
}

func TestPreferencesService_SavePreferencesAppliesExplicitBooleans(t *testing.T) {
	current := domain.DefaultUserPreferences(9)
	current.CompactMode = true
	current.ReducedMotion = true

	compactMode := false
	reducedMotion := false
	repo := &preferencesRepoStub{current: current}
	svc := usecase.NewPreferencesService(repo, preferencesLoggerStub{})

	saved, err := svc.SavePreferences(t.Context(), 9, input.SavePreferencesCommand{
		CompactMode:   &compactMode,
		ReducedMotion: &reducedMotion,
	})

	require.NoError(t, err)
	require.False(t, saved.CompactMode)
	require.False(t, saved.ReducedMotion)
}

func TestPreferencesService_SavePreferencesRejectsInvalidPreset(t *testing.T) {
	repo := &preferencesRepoStub{current: domain.DefaultUserPreferences(11)}
	svc := usecase.NewPreferencesService(repo, preferencesLoggerStub{})

	_, err := svc.SavePreferences(t.Context(), 11, input.SavePreferencesCommand{
		ThemePreset: "unknown",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid theme preset")
	require.Empty(t, repo.saved)
}
