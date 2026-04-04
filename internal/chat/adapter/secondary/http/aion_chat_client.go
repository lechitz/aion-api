// Package http provides the HTTP client adapter for communicating with Aion-Chat service.
//
//revive:disable:var-naming // package name deliberately matches HTTP adapter naming
package http

//revive:enable:var-naming

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lechitz/aion-api/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// SendMessage sends a chat message to the Aion-Chat service.
func (c *AionChatClient) SendMessage(ctx context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSendMessage)
	defer span.End()

	url := fmt.Sprintf("%s%s", c.baseURL, PathProcess)

	span.SetAttributes(
		attribute.String(AttrHTTPURL, url),
		attribute.String(AttrHTTPMethod, http.MethodPost),
		attribute.String(AttrUserID, strconv.FormatUint(req.UserID, 10)),
	)

	httpReq, err := c.buildSendMessageRequest(ctx, url, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	c.logger.InfowCtx(ctx, MsgCallingAionChatService, commonkeys.URL, url, AttrUserID, req.UserID)
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrHTTPRequestFailed)
		c.logger.ErrorwCtx(ctx, ErrHTTPRequestFailed, commonkeys.Error, err.Error(), "url", url)
		return nil, fmt.Errorf("%s: %w", ErrAionChatRequestFailed, err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	span.SetAttributes(attribute.Int(AttrHTTPStatusCode, httpResp.StatusCode))

	body, err := readResponseBody(httpResp)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedReadResponse)
		c.logger.ErrorwCtx(ctx, ErrFailedReadResponse, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedReadResponse, err)
	}

	if err := c.validateSendMessageResponse(ctx, span, httpResp.StatusCode, body); err != nil {
		return nil, err
	}

	resp, err := decodeSendMessageResponse(body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedUnmarshal)
		c.logger.ErrorwCtx(ctx, ErrFailedUnmarshal, commonkeys.Error, err.Error(), "body", string(body))
		return nil, fmt.Errorf("%s: %w", ErrFailedUnmarshal, err)
	}

	span.SetAttributes(
		attribute.Int(AttrTokensUsed, resp.TokensUsed),
		attribute.Int(AttrResponseLength, len(resp.Response)),
	)
	span.SetStatus(codes.Ok, StatusMessageSent)

	c.logger.InfowCtx(ctx, MsgAionChatResponseReceived, AttrUserID, req.UserID, AttrResponseLengthShort, len(resp.Response))

	return resp, nil
}

func (c *AionChatClient) buildSendMessageRequest(ctx context.Context, url string, req *dto.InternalChatRequest) (*http.Request, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.ErrorwCtx(ctx, ErrFailedMarshal, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedMarshal, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.ErrorwCtx(ctx, ErrFailedCreateRequest, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedCreateRequest, err)
	}
	httpReq.Header.Set(HeaderContentType, ContentTypeJSON)
	httpReq.Header.Set(HeaderAccept, ContentTypeJSON)
	return httpReq, nil
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

func (c *AionChatClient) validateSendMessageResponse(ctx context.Context, span oteltrace.Span, statusCode int, body []byte) error {
	if statusCode == http.StatusOK {
		return nil
	}
	if statusCode == StatusCodeClientClosedRequest {
		span.SetStatus(codes.Ok, StatusRequestCancelled)
		c.logger.WarnwCtx(ctx, MsgAionChatRequestCancelled, AttrStatusCode, statusCode, AttrBody, string(body))
		return fmt.Errorf("%s: %w", ErrAionChatRequestFailed, context.Canceled)
	}
	if statusCode == http.StatusBadRequest {
		span.SetStatus(codes.Error, ErrAionChatNonOK)
		c.logger.ErrorwCtx(ctx, ErrAionChatNonOK, AttrStatusCode, statusCode, AttrBody, string(body))
		return sharederrors.NewValidationError("runtime", extractAionChatErrorDetail(body))
	}

	span.SetStatus(codes.Error, ErrAionChatNonOK)
	c.logger.ErrorwCtx(ctx, ErrAionChatNonOK, AttrStatusCode, statusCode, AttrBody, string(body))
	return fmt.Errorf("%s: status %d: %s", ErrAionChatNonOK, statusCode, string(body))
}

func extractAionChatErrorDetail(body []byte) string {
	var payload struct {
		Detail any `json:"detail"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		switch detail := payload.Detail.(type) {
		case string:
			if detail != "" {
				return detail
			}
		case map[string]any:
			if errValue, ok := detail["error"].(string); ok && errValue != "" {
				return errValue
			}
		}
	}
	if len(body) == 0 {
		return ErrAionChatNonOK
	}
	return string(body)
}

func decodeSendMessageResponse(body []byte) (*dto.InternalChatResponse, error) {
	var resp dto.InternalChatResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
