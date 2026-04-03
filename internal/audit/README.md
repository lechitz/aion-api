# Audit Context

**Path:** `internal/audit`

## Purpose

`internal/audit` owns immutable action-event persistence and the internal diagnostics query surface for those events.

## Current Surface

| Surface | Current contract |
| --- | --- |
| `core/ports/input.Service.WriteEvent` | validate and persist one immutable audit event |
| `core/ports/input.Service.ListEvents` | return filtered audit events for diagnostics |
| HTTP `GET /audit/events` | authenticated diagnostics endpoint; self-scope by default, `user_id` cross-user queries only for admin callers |
| Storage | `aion_api.audit_action_events` |

## Current Producers

- `chat` emits UI-action audit events through `WriteEvent`
- other contexts may publish audit events, but they must treat persistence as non-blocking from a business-outcome perspective

## Boundary Rules

- `audit` owns the immutable event log, not the business workflows that generated it
- the read surface is HTTP-only in the current runtime
- consumers should not depend on audit writes for transactional guarantees

## Validate

```bash
go test ./internal/audit/...
make verify
```

## Risks And Compatibility Notes

- payloads must remain allow-listed and redacted
- append-only behavior is part of the safety model and should not be weakened by convenience updates
- if diagnostics filtering changes, keep admin/self-scope rules explicit in tests and docs

## Related Docs

- [`../chat/README.md`](../chat/README.md)
- [`../platform/server/http/README.md`](../platform/server/http/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
