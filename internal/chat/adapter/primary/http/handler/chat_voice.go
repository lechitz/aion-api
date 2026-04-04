package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ChatVoice processes an audio message from the user and returns the AI response.
//
// @Summary      Send voice chat message
// @Description  Sends an audio file to be transcribed and processed by the AI assistant. Requires authentication.
// @Tags         ChatText
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer token"
// @Param        audio          formData  file    true   "Audio file (webm, wav, mp3, max 10MB, max 60s)"
// @Param        language       formData  string  false  "Language code (pt, en, es) or auto-detect if empty"
// @Param        provider       formData  string  false  "Optional runtime provider override (for example: openai, ollama)"
// @Param        model          formData  string  false  "Optional runtime model override when provider is present"
// @Success      200            {object}  map[string]interface{}  "Voice chat response with transcription and AI response"
// @Failure      400            {string}  string                  "Invalid audio file or validation error"
// @Failure      401            {string}  string                  "Unauthorized - missing or invalid token"
// @Failure      413            {string}  string                  "Audio file too large"
// @Failure      500            {string}  string                  "Internal server error"
// @Failure      503            {string}  string                  "Service unavailable - AI service is down"
// @Router       /chat/audio [post]
// @Security     BearerAuth.
func (h *Handler) ChatVoice(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerChatHandler).Start(r.Context(), SpanChatVoice)
	defer span.End()

	userID, ok := h.extractUserID(ctx, w, span)
	if !ok {
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, r.RemoteAddr),
	)

	file, header, language, runtimeProvider, runtimeModel, ok := h.parseVoiceRequest(ctx, w, r, span)
	if !ok {
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			h.Logger.WarnwCtx(ctx, LogFailedCloseAudioFile, commonkeys.Error, closeErr)
		}
	}()

	span.SetAttributes(
		attribute.String(AttrAudioFilename, header.Filename),
		attribute.Int64(AttrAudioSizeBytes, header.Size),
		attribute.String(AttrAudioContentType, header.Header.Get(HeaderContentType)),
	)
	if language != "" {
		span.SetAttributes(attribute.String(AttrAudioLanguage, language))
	}

	span.AddEvent(EventForwardToAionChat)
	buf, contentType, ok := h.buildMultipartRequest(ctx, w, span, file, header, userID, language, runtimeProvider, runtimeModel)
	if !ok {
		return
	}

	responseBody, statusCode, ok := h.forwardToAionChat(ctx, w, span, buf, contentType)
	if !ok {
		return
	}

	if statusCode != http.StatusOK {
		h.writeErrorResponse(ctx, w, span, statusCode, responseBody)
		return
	}

	h.writeSuccessResponse(ctx, w, span, userID, header.Size, statusCode, responseBody)
}

func (h *Handler) extractUserID(ctx context.Context, w http.ResponseWriter, span trace.Span) (uint64, bool) {
	userIDValue := ctx.Value(ctxkeys.UserID)
	if userIDValue == nil {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrUserIDNotFound), h.Logger)
		return 0, false
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		span.SetStatus(codes.Error, ErrInvalidUserIDType)
		h.Logger.ErrorwCtx(ctx, LogInvalidUserIDType, LogKeyValue, userIDValue)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrInvalidUserID), h.Logger)
		return 0, false
	}

	return userID, true
}

func (h *Handler) parseVoiceRequest(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	span trace.Span,
) (multipart.File, *multipart.FileHeader, string, string, string, bool) {
	span.AddEvent(EventParseMultipartForm)
	r.Body = http.MaxBytesReader(w, r.Body, MaxAudioSize)
	if err := r.ParseMultipartForm(MaxAudioSize); err != nil {
		if isMultipartTooLargeError(err) {
			span.SetStatus(codes.Error, ErrAudioFileTooLarge)
			h.Logger.ErrorwCtx(ctx, LogAudioFileTooLarge, commonkeys.Error, err, LogKeyMax, MaxAudioSize)
			httpresponse.WriteValidationErrorSpan(ctx, w, span,
				sharederrors.NewValidationError(FormFieldAudio,
					fmt.Sprintf("Audio file too large (max: %d bytes)", MaxAudioSize)), h.Logger)
			return nil, nil, "", "", "", false
		}

		span.SetStatus(codes.Error, ErrFailedParseMultipartForm)
		h.Logger.ErrorwCtx(ctx, LogFailedParseMultipartForm, commonkeys.Error, err)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio, ErrInvalidMultipartForm), h.Logger)
		return nil, nil, "", "", "", false
	}

	file, header, err := r.FormFile(FormFieldAudio)
	if err != nil {
		span.SetStatus(codes.Error, ErrMissingAudioFile)
		h.Logger.ErrorwCtx(ctx, LogMissingAudioFile, commonkeys.Error, err)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio, ErrAudioFileRequired), h.Logger)
		return nil, nil, "", "", "", false
	}

	if header.Size > MaxAudioSize {
		span.SetStatus(codes.Error, ErrAudioFileTooLarge)
		h.Logger.ErrorwCtx(ctx, LogAudioFileTooLarge, LogKeySize, header.Size, LogKeyMax, MaxAudioSize)
		httpresponse.WriteValidationErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio,
				fmt.Sprintf("Audio file too large: %d bytes (max: %d)", header.Size, MaxAudioSize)), h.Logger)
		return nil, nil, "", "", "", false
	}

	language := ""
	if r.MultipartForm != nil {
		languageValues := r.MultipartForm.Value[FormFieldLanguage]
		if len(languageValues) > 0 {
			language = languageValues[0]
		}
	}

	runtimeProvider, runtimeModel, err := extractRuntimeOverrideFromMultipart(r.MultipartForm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		h.Logger.ErrorwCtx(ctx, err.Error())
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldMessage, err.Error()), h.Logger)
		return nil, nil, "", "", "", false
	}

	return file, header, language, runtimeProvider, runtimeModel, true
}

func isMultipartTooLargeError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, multipart.ErrMessageTooLarge) {
		return true
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "request body too large") ||
		strings.Contains(message, "multipart: message too large") ||
		strings.Contains(message, "message too large")
}

func (h *Handler) buildMultipartRequest(
	ctx context.Context,
	w http.ResponseWriter,
	span trace.Span,
	file multipart.File,
	header *multipart.FileHeader,
	userID uint64,
	language string,
	runtimeProvider string,
	runtimeModel string,
) (*bytes.Buffer, string, bool) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(FormFieldAudio, header.Filename)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCreateFormFile)
		h.Logger.ErrorwCtx(ctx, LogFailedCreateFormFile, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return nil, "", false
	}

	if _, err := io.Copy(part, file); err != nil {
		span.SetStatus(codes.Error, ErrFailedCopyAudioData)
		h.Logger.ErrorwCtx(ctx, LogFailedCopyAudioData, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return nil, "", false
	}

	if err := writer.WriteField(FormFieldUserID, strconv.FormatUint(userID, 10)); err != nil {
		span.SetStatus(codes.Error, ErrFailedWriteUserIDField)
		h.Logger.ErrorwCtx(ctx, LogFailedWriteUserIDField, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return nil, "", false
	}

	if language != "" {
		if err := writer.WriteField(FormFieldLanguage, language); err != nil {
			span.SetStatus(codes.Error, ErrFailedWriteLanguageField)
			h.Logger.ErrorwCtx(ctx, LogFailedWriteLanguageField, commonkeys.Error, err)
			httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
			return nil, "", false
		}
	}

	if runtimeProvider != "" {
		if err := writer.WriteField(FormFieldProvider, runtimeProvider); err != nil {
			span.SetStatus(codes.Error, ErrFailedWriteLanguageField)
			h.Logger.ErrorwCtx(ctx, "failed to write provider field", commonkeys.Error, err)
			httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
			return nil, "", false
		}
		if err := writer.WriteField(FormFieldModel, runtimeModel); err != nil {
			span.SetStatus(codes.Error, ErrFailedWriteLanguageField)
			h.Logger.ErrorwCtx(ctx, "failed to write model field", commonkeys.Error, err)
			httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
			return nil, "", false
		}
	}

	if err := writer.Close(); err != nil {
		span.SetStatus(codes.Error, ErrFailedCloseMultipartWriter)
		h.Logger.ErrorwCtx(ctx, LogFailedCloseMultipartWriter, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return nil, "", false
	}

	return &buf, writer.FormDataContentType(), true
}

func extractRuntimeOverrideFromMultipart(form *multipart.Form) (string, string, error) {
	if form == nil {
		return "", "", nil
	}

	var provider string
	if values := form.Value[FormFieldProvider]; len(values) > 0 {
		provider = strings.TrimSpace(values[0])
	}

	var model string
	if values := form.Value[FormFieldModel]; len(values) > 0 {
		model = strings.TrimSpace(values[0])
	}

	if provider == "" && model == "" {
		return "", "", nil
	}
	if provider == "" {
		return "", "", errors.New("runtime.provider is required when runtime is present")
	}
	if model == "" {
		return "", "", errors.New("runtime.model is required when runtime is present")
	}

	return provider, model, nil
}

func (h *Handler) forwardToAionChat(
	ctx context.Context,
	w http.ResponseWriter,
	span trace.Span,
	buf *bytes.Buffer,
	contentType string,
) ([]byte, int, bool) {
	aionChatURL := h.Config.AionChat.BaseURL + PathProcessAudio
	req, err := http.NewRequestWithContext(ctx, HTTPMethodPost, aionChatURL, buf)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCreateRequest)
		h.Logger.ErrorwCtx(ctx, LogFailedCreateRequest, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgInternalServerError, h.Logger)
		return nil, 0, false
	}

	req.Header.Set(HeaderContentType, contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCallAionChat)
		h.Logger.ErrorwCtx(ctx, LogFailedCallAionChat, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgAIServiceUnavailable, h.Logger)
		return nil, 0, false
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			h.Logger.WarnwCtx(ctx, LogFailedCloseResponseBody, commonkeys.Error, closeErr)
		}
	}()

	span.SetAttributes(attribute.Int(AttrAionChatStatusCode, resp.StatusCode))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedReadResponse)
		h.Logger.ErrorwCtx(ctx, LogFailedReadResponse, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgInternalServerError, h.Logger)
		return nil, 0, false
	}

	return responseBody, resp.StatusCode, true
}

func (h *Handler) writeErrorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	span trace.Span,
	statusCode int,
	responseBody []byte,
) {
	span.SetStatus(codes.Error, ErrAionChatReturnedError)
	h.Logger.ErrorwCtx(ctx, LogAionChatError, LogKeyStatusCode, statusCode, LogKeyResponse, string(responseBody))
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)
	if _, writeErr := w.Write(responseBody); writeErr != nil {
		h.Logger.ErrorwCtx(ctx, LogFailedWriteErrorResponse, commonkeys.Error, writeErr)
	}
}

func (h *Handler) writeSuccessResponse(
	ctx context.Context,
	w http.ResponseWriter,
	span trace.Span,
	userID uint64,
	audioSize int64,
	statusCode int,
	responseBody []byte,
) {
	span.SetStatus(codes.Ok, StatusVoiceChatSuccess)
	h.Logger.InfowCtx(ctx, LogVoiceChatSuccess,
		commonkeys.UserID, userID, LogKeyAudioSize, audioSize, LogKeyStatusCode, statusCode)
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseBody); err != nil {
		h.Logger.ErrorwCtx(ctx, LogFailedWriteSuccessResponse, commonkeys.Error, err)
	}
}
