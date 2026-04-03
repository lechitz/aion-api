# Chat Bounded Context

**Path:** `internal/chat`

## Purpose

`internal/chat` owns authenticated chat interaction flows, conversation history and context retrieval, and the integration boundary with the external `aion-chat` service.

## Current Transport Surface

| Surface | Current contract |
| --- | --- |
| HTTP `POST /chat/text` | authenticated message send; returns assistant response, UI payload, sources, and optional usage |
| HTTP `POST /chat/cancel` | authenticated cancel request; proxies cancellation to `AION_CHAT_URL/internal/cancel` |
| HTTP `POST /chat/audio` | authenticated voice or chat ingestion path |
| GraphQL read surface | context controllers expose chat history and aggregated chat context |

There is no shared GraphQL mutation contract for chat in the current backend-owned public surface.

## Runtime Flow

1. Primary adapters authenticate the request and extract `userID`.
2. `ProcessMessage` loads recent cached chat messages for context.
3. The secondary HTTP adapter forwards the request to `aion-chat`.
4. The usecase maps response text, UI payload, sources, and token usage into the local domain result.
5. UI-action audit persistence is attempted best-effort through the `audit` bounded context.
6. Chat history is saved asynchronously so request cancellation does not drop persistence work.

## Boundary Rules

- provider-specific HTTP semantics stay in the secondary adapter
- audit persistence failures must not fail the main chat response path
- UI-action metadata belongs to request context and transport contracts; business ownership of audit storage remains in `internal/audit`

## Validate

```bash
go test ./internal/chat/...
go test ./internal/adapter/primary/graphql/...
make verify
```

## Risks And Compatibility Notes

- `aion-chat` integration behavior is a cross-repo contract and should stay explicit when payload or timeout semantics change
- async history persistence intentionally favors response latency over strict coupling to request cancellation
- if chat read surfaces change, keep shared GraphQL query docs aligned

## Related Docs

- [`../audit/README.md`](../audit/README.md)
- [`../adapter/primary/graphql/README.md`](../adapter/primary/graphql/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
