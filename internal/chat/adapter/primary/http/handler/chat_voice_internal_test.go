package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"go.opentelemetry.io/otel"
)

type noopVoiceLogger struct{}

func (noopVoiceLogger) Infof(string, ...any)                      {}
func (noopVoiceLogger) Errorf(string, ...any)                     {}
func (noopVoiceLogger) Debugf(string, ...any)                     {}
func (noopVoiceLogger) Warnf(string, ...any)                      {}
func (noopVoiceLogger) Infow(string, ...any)                      {}
func (noopVoiceLogger) Errorw(string, ...any)                     {}
func (noopVoiceLogger) Debugw(string, ...any)                     {}
func (noopVoiceLogger) Warnw(string, ...any)                      {}
func (noopVoiceLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopVoiceLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopVoiceLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopVoiceLogger) DebugwCtx(context.Context, string, ...any) {}

type memMultipartFile struct {
	*bytes.Reader
}

func (m *memMultipartFile) Close() error { return nil }

type errMultipartFile struct{}

func (e *errMultipartFile) Read([]byte) (int, error)          { return 0, errors.New("read fail") }
func (e *errMultipartFile) ReadAt([]byte, int64) (int, error) { return 0, errors.New("readat fail") }
func (e *errMultipartFile) Seek(int64, int) (int64, error)    { return 0, errors.New("seek fail") }
func (e *errMultipartFile) Close() error                      { return nil }

func TestVoiceExtractUserID(t *testing.T) {
	h := &Handler{Logger: noopVoiceLogger{}, Config: &config.Config{}}
	_, span := otel.Tracer("test").Start(t.Context(), "extract")
	defer span.End()

	t.Run("missing user id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		_, ok := h.extractUserID(t.Context(), rec, span)
		if ok {
			t.Fatal("expected false when user id is missing")
		}
	})

	t.Run("invalid user id type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, "10")
		rec := httptest.NewRecorder()
		_, ok := h.extractUserID(ctx, rec, span)
		if ok {
			t.Fatal("expected false when user id type is invalid")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ctxkeys.UserID, uint64(10))
		rec := httptest.NewRecorder()
		userID, ok := h.extractUserID(ctx, rec, span)
		if !ok || userID != 10 {
			t.Fatalf("unexpected result: ok=%v userID=%d", ok, userID)
		}
	})
}

func TestVoiceBuildMultipartRequest(t *testing.T) {
	h := &Handler{Logger: noopVoiceLogger{}, Config: &config.Config{}}
	_, span := otel.Tracer("test").Start(t.Context(), "build")
	defer span.End()

	header := &multipart.FileHeader{Filename: "audio.wav"}

	t.Run("success with language", func(t *testing.T) {
		file := &memMultipartFile{Reader: bytes.NewReader([]byte("audio-bytes"))}
		rec := httptest.NewRecorder()
		buf, contentType, ok := h.buildMultipartRequest(t.Context(), rec, span, file, header, 5, "pt", "openai", "gpt-4.1-mini")
		if !ok {
			t.Fatal("expected success")
		}
		if buf == nil || contentType == "" {
			t.Fatal("expected non-empty multipart payload")
		}
		if !bytes.Contains(buf.Bytes(), []byte("audio-bytes")) {
			t.Fatal("expected multipart payload to contain audio bytes")
		}
		if !bytes.Contains(buf.Bytes(), []byte(`name="provider"`)) {
			t.Fatal("expected multipart payload to contain provider field")
		}
		if !bytes.Contains(buf.Bytes(), []byte(`gpt-4.1-mini`)) {
			t.Fatal("expected multipart payload to contain model field")
		}
	})

	t.Run("copy error", func(t *testing.T) {
		file := &errMultipartFile{}
		rec := httptest.NewRecorder()
		_, _, ok := h.buildMultipartRequest(t.Context(), rec, span, file, header, 5, "", "", "")
		if ok {
			t.Fatal("expected failure when audio copy fails")
		}
	})
}

func TestExtractRuntimeOverrideFromMultipart(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		provider, model, err := extractRuntimeOverrideFromMultipart(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if provider != "" || model != "" {
			t.Fatalf("unexpected runtime override: provider=%q model=%q", provider, model)
		}
	})

	t.Run("provider required when model present", func(t *testing.T) {
		form := &multipart.Form{Value: map[string][]string{
			FormFieldModel: {"gpt-4.1-mini"},
		}}
		_, _, err := extractRuntimeOverrideFromMultipart(form)
		if err == nil || err.Error() != "runtime.provider is required when runtime is present" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("model required when provider present", func(t *testing.T) {
		form := &multipart.Form{Value: map[string][]string{
			FormFieldProvider: {"openai"},
		}}
		_, _, err := extractRuntimeOverrideFromMultipart(form)
		if err == nil || err.Error() != "runtime.model is required when runtime is present" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		form := &multipart.Form{Value: map[string][]string{
			FormFieldProvider: {" openai "},
			FormFieldModel:    {" gpt-4.1-mini "},
		}}
		provider, model, err := extractRuntimeOverrideFromMultipart(form)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if provider != "openai" || model != "gpt-4.1-mini" {
			t.Fatalf("unexpected runtime override: provider=%q model=%q", provider, model)
		}
	})
}

func TestVoiceWriteResponses(t *testing.T) {
	h := &Handler{Logger: noopVoiceLogger{}, Config: &config.Config{}}
	_, span := otel.Tracer("test").Start(t.Context(), "write")
	defer span.End()

	t.Run("error response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.writeErrorResponse(t.Context(), rec, span, 503, []byte(`{"error":"unavailable"}`))
		if rec.Code != 503 {
			t.Fatalf("expected 503, got %d", rec.Code)
		}
		if rec.Body.String() != `{"error":"unavailable"}` {
			t.Fatalf("unexpected body: %s", rec.Body.String())
		}
	})

	t.Run("success response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.writeSuccessResponse(t.Context(), rec, span, 7, 12, 200, []byte(`{"message":"ok"}`))
		if rec.Code != 200 {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if rec.Body.String() != `{"message":"ok"}` {
			t.Fatalf("unexpected body: %s", rec.Body.String())
		}
	})
}

var _ io.Reader = (*memMultipartFile)(nil)
