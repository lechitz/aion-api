// Package handler implements HTTP handlers for chat endpoints.
package handler

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerChatHandler is the tracer name for chat handler.
// Format: aion-api.<domain>.<layer> .
const TracerChatHandler = "aion-api.chat.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanChatHandler is the span name for handling chat requests.
	SpanChatHandler = "chat.handler.handle"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventDecodeRequest indicates the decoding of the chat request.
	EventDecodeRequest = "chat.handler.decode_request"
	// EventValidateRequest indicates the validation of the chat request.
	EventValidateRequest = "chat.handler.validate_request"
	// EventCallService indicates the call to the chat service.
	EventCallService = "chat.handler.call_service"
	// EventChatSuccess indicates successful chat processing.
	EventChatSuccess = "chat.handler.success"
	// EventChatError indicates an error occurred during chat processing.
	EventChatError = "chat.handler.error"
	// EventChatCancelled indicates the request was cancelled by the client.
	EventChatCancelled = "chat.handler.cancelled"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusChatSuccess indicates successful chat processing.
	StatusChatSuccess = "chat processed successfully"
	// StatusChatCancelled indicates successful cancellation acknowledgement.
	StatusChatCancelled = "chat request cancelled by client"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	// ErrChat indicates a failure in processing the chat message.
	ErrChat = "failed to process chat message"
	// ErrRequiredMessage indicates that the message field is required.
	ErrRequiredMessage = "message is required"
	// ErrMessageTooShort indicates that the message is too short.
	ErrMessageTooShort = "message must be at least 1 character"
	// ErrMessageTooLong indicates that the message exceeds the maximum length.
	ErrMessageTooLong = "message must not exceed 2000 characters"
	// ErrUserIDNotFound indicates that the user ID was not found in the context.
	ErrUserIDNotFound = "user ID not found in context"
)

// Log messages.
const (
	MsgChatSuccess      = "Chat processed successfully"
	MsgChatRequestStart = "Processing chat request"
	// MsgChatCancelled indicates chat request cancellation by the client.
	MsgChatCancelled = "Chat request cancelled by client"
	// MsgChatRequestIncludesUIAction indicates request includes ui_action metadata.
	MsgChatRequestIncludesUIAction = "HTTP chat request includes UI action"
)

// Response messages.
const (
	// MsgChatCancelledResponse is the success payload message returned on cancel.
	MsgChatCancelledResponse = "Chat request cancelled"
)

// HTTP status codes.
const (
	// HTTPStatusClientClosedRequest is the non-standard code used for cancelled requests.
	HTTPStatusClientClosedRequest = 499
)

// Error text markers used to detect client-side cancellation.
const (
	ErrorTextContextCanceled   = "context canceled"
	ErrorTextRequestCanceled   = "request canceled"
	ErrorTextOperationCanceled = "operation was canceled"
)

// Validation constraints.
const (
	// MinMessageLength is the minimum length for a chat message.
	MinMessageLength = 1
	// MaxMessageLength is the maximum length for a chat message.
	MaxMessageLength = 2000
	// MaxAudioSize is the maximum size for audio files (10MB).
	MaxAudioSize = 10 << 20
)

// =============================================================================
// VOICE CHAT - Magic Strings and Constants
// =============================================================================

// Span names for voice chat.
const (
	// SpanChatVoice is the span name for voice chat requests.
	SpanChatVoice = "chat.handler.voice"
)

// Event names for voice chat.
const (
	// EventParseMultipartForm indicates parsing of multipart form data.
	EventParseMultipartForm = "parse_multipart_form"
	// EventForwardToAionChat indicates forwarding request to aion-chat service.
	EventForwardToAionChat = "forward_to_aion_chat"
)

// Error messages for voice chat.
const (
	// ErrInvalidUserIDType indicates that the user ID has an invalid type.
	ErrInvalidUserIDType = "invalid user ID type"
	// ErrInvalidUserID indicates an invalid user ID.
	ErrInvalidUserID = "invalid user ID"
	// ErrFailedParseMultipartForm indicates failure to parse multipart form.
	ErrFailedParseMultipartForm = "failed to parse multipart form"
	// ErrInvalidMultipartForm is the user-facing error for invalid multipart form.
	ErrInvalidMultipartForm = "Invalid multipart form or file too large"
	// ErrMissingAudioFile indicates missing audio file in the request.
	ErrMissingAudioFile = "missing audio file"
	// ErrAudioFileRequired is the user-facing error for missing audio file.
	ErrAudioFileRequired = "Audio file is required"
	// ErrAudioFileTooLarge indicates audio file exceeds maximum size.
	ErrAudioFileTooLarge = "audio file too large"
	// ErrFailedCreateFormFile indicates failure to create form file.
	ErrFailedCreateFormFile = "failed to create form file"
	// ErrFailedCopyAudioData indicates failure to copy audio data.
	ErrFailedCopyAudioData = "failed to copy audio data"
	// ErrFailedCreateRequest indicates failure to create HTTP request.
	ErrFailedCreateRequest = "failed to create request"
	// ErrFailedCallAionChat indicates failure to call aion-chat service.
	ErrFailedCallAionChat = "failed to call aion-chat"
	// ErrFailedReadResponse indicates failure to read response.
	ErrFailedReadResponse = "failed to read response"
	// ErrAionChatReturnedError indicates aion-chat service returned an error.
	ErrAionChatReturnedError = "aion-chat returned error"
	// ErrFailedWriteUserIDField indicates failure to write user_id field.
	ErrFailedWriteUserIDField = "failed to write user_id field"
	// ErrFailedWriteLanguageField indicates failure to write language field.
	ErrFailedWriteLanguageField = "failed to write language field"
	// ErrFailedCloseMultipartWriter indicates failure to close multipart writer.
	ErrFailedCloseMultipartWriter = "failed to close multipart writer"
)

// User-facing messages for voice chat.
const (
	// MsgFailedProcessAudio is the user-facing message for audio processing failure.
	MsgFailedProcessAudio = "Failed to process audio"
	// MsgInternalServerError is the user-facing message for internal server errors.
	MsgInternalServerError = "Internal server error"
	// MsgAIServiceUnavailable is the user-facing message when AI service is unavailable.
	MsgAIServiceUnavailable = "AI service temporarily unavailable"
)

// Status descriptions for voice chat.
const (
	// StatusVoiceChatSuccess indicates successful voice chat processing.
	StatusVoiceChatSuccess = "voice chat processed successfully"
)

// Log messages for voice chat.
const (
	// LogInvalidUserIDType indicates invalid user ID type in context.
	LogInvalidUserIDType = "Invalid user ID type in context"
	// LogFailedParseMultipartForm indicates failure to parse multipart form.
	LogFailedParseMultipartForm = "Failed to parse multipart form"
	// LogMissingAudioFile indicates missing audio file.
	LogMissingAudioFile = "Missing audio file"
	// LogAudioFileTooLarge indicates audio file is too large.
	LogAudioFileTooLarge = "Audio file too large"
	// LogFailedCreateFormFile indicates failure to create form file.
	LogFailedCreateFormFile = "Failed to create form file"
	// LogFailedCopyAudioData indicates failure to copy audio data.
	LogFailedCopyAudioData = "Failed to copy audio data"
	// LogFailedCreateRequest indicates failure to create request to aion-chat.
	LogFailedCreateRequest = "Failed to create request to aion-chat"
	// LogFailedCallAionChat indicates failure to call aion-chat service.
	LogFailedCallAionChat = "Failed to call aion-chat service"
	// LogFailedReadResponse indicates failure to read aion-chat response.
	LogFailedReadResponse = "Failed to read aion-chat response"
	// LogAionChatError indicates aion-chat returned an error.
	LogAionChatError = "Aion-chat returned error"
	// LogVoiceChatSuccess indicates successful voice chat processing.
	LogVoiceChatSuccess = "Voice chat processed successfully"
	// LogFailedCloseAudioFile indicates failure to close audio file.
	LogFailedCloseAudioFile = "failed to close audio file"
	// LogFailedCloseResponseBody indicates failure to close response body.
	LogFailedCloseResponseBody = "failed to close response body"
	// LogFailedWriteUserIDField indicates failure to write user_id field.
	LogFailedWriteUserIDField = "failed to write user_id field"
	// LogFailedWriteLanguageField indicates failure to write language field.
	LogFailedWriteLanguageField = "failed to write language field"
	// LogFailedCloseMultipartWriter indicates failure to close multipart writer.
	LogFailedCloseMultipartWriter = "failed to close multipart writer"
	// LogFailedWriteErrorResponse indicates failure to write error response.
	LogFailedWriteErrorResponse = "failed to write error response"
	// LogFailedWriteSuccessResponse indicates failure to write success response.
	LogFailedWriteSuccessResponse = "failed to write success response"
)

// Form field names.
const (
	// FormFieldAudio is the form field name for audio file.
	FormFieldAudio = "audio"
	// FormFieldUserID is the form field name for user ID.
	FormFieldUserID = "user_id"
	// FormFieldLanguage is the form field name for language.
	FormFieldLanguage = "language"
	// FormFieldProvider is the form field name for runtime provider override.
	FormFieldProvider = "provider"
	// FormFieldModel is the form field name for runtime model override.
	FormFieldModel = "model"
	// FormFieldMessage is the form field name for chat message.
	FormFieldMessage = "message"
)

// Attribute names for tracing.
const (
	// AttrAudioFilename is the attribute name for audio filename.
	AttrAudioFilename = "audio.filename"
	// AttrAudioSizeBytes is the attribute name for audio size in bytes.
	AttrAudioSizeBytes = "audio.size_bytes"
	// AttrAudioContentType is the attribute name for audio content type.
	AttrAudioContentType = "audio.content_type"
	// AttrAudioLanguage is the attribute name for audio language.
	AttrAudioLanguage = "audio.language"
	// AttrAionChatStatusCode is the attribute name for aion-chat status code.
	AttrAionChatStatusCode = "aion_chat.status_code"
	// AttrMessageLength is the attribute name for message length.
	AttrMessageLength = "message_length"
	// AttrTokensUsed is the attribute name for tokens used.
	AttrTokensUsed = "tokens_used"
	// AttrResponseLength is the attribute name for response length.
	AttrResponseLength = "response_length"
	// AttrSourcesCount is the attribute name for sources count.
	AttrSourcesCount = "sources_count"
)

// HTTP methods.
const (
	// HTTPMethodPost is the POST HTTP method.
	HTTPMethodPost = "POST"
)

// HTTP headers.
const (
	// HeaderContentType is the Content-Type header name.
	HeaderContentType = "Content-Type"
	// ContentTypeJSON is the application/json content type.
	ContentTypeJSON = "application/json"
)

// Service paths.
const (
	// PathProcessAudio is the path for processing audio in aion-chat service.
	PathProcessAudio = "/internal/process-audio"
)

// Log keys for metadata.
const (
	// LogKeyValue is the log key for generic value.
	LogKeyValue = "value"
	// LogKeySize is the log key for size.
	LogKeySize = "size"
	// LogKeyMax is the log key for maximum value.
	LogKeyMax = "max"
	// LogKeyStatusCode is the log key for HTTP status code.
	LogKeyStatusCode = "status_code"
	// LogKeyResponse is the log key for response body.
	LogKeyResponse = "response"
	// LogKeyAudioSize is the log key for audio file size.
	LogKeyAudioSize = "audio_size"
	// LogKeyUIActionType is the log key for ui action type.
	LogKeyUIActionType = "ui_action_type"
	// LogKeyDraftID is the log key for draft identifier.
	LogKeyDraftID = "draft_id"
	// LogKeyConsentRequired is the log key for consent requirement flag.
	LogKeyConsentRequired = "consent_required"
	// LogKeyConsentConfirmed is the log key for consent confirmation flag.
	LogKeyConsentConfirmed = "consent_confirmed"
	// LogKeyConsentPolicyVersion is the log key for consent policy version.
	LogKeyConsentPolicyVersion = "consent_policy_version"
	// LogKeyQuickAddContractVersion is the log key for quick_add contract version.
	LogKeyQuickAddContractVersion = "quick_add_contract_version"
	// LogKeyQuickAddEntity is the log key for quick_add entity.
	LogKeyQuickAddEntity = "quick_add_entity"
	// LogKeyQuickAddOperation is the log key for quick_add operation.
	LogKeyQuickAddOperation = "quick_add_operation"
	// LogKeyQuickAddIdempotencyKey is the log key for quick_add idempotency key.
	LogKeyQuickAddIdempotencyKey = "quick_add_idempotency_key"
)

// Context payload keys for chat request metadata.
const (
	ContextKeyUIAction     = "ui_action"
	ContextKeyUIActionType = "type"
	ContextKeyDraftID      = "draft_id"
	ContextKeyConsent      = "consent"
	ContextKeyQuickAdd     = "quick_add"
)
