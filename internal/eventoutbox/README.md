# Event Outbox

**Path:** `internal/eventoutbox`

## Purpose

`internal/eventoutbox` owns the durable outbox that persists canonical backend events before publication to Kafka.
It is the relay boundary between transactional writes inside `aion-api` and the wider event backbone.

## Current Surface

| Interface | Responsibility |
| --- | --- |
| `core/ports/input.Service.Enqueue` | validate and persist one outbox event |
| `core/ports/input.PublisherService.PublishPending` | publish one batch of pending events and mark or reschedule rows |
| `core/ports/output.EventRepository` | save, list pending, mark published, reschedule, and expose stats |
| `adapter/secondary/kafka` | publish normalized outbox envelopes to Kafka |

## Runtime Contract

- durable rows are stored in `aion_api.event_outbox`
- newly enqueued events use the backend-owned canonical envelope and version defaults
- the publisher loop reads pending rows in batches, publishes externally, then either marks rows as published or reschedules them with backoff and last-error metadata
- aggregate stats are available through repository support code for operator diagnostics

## Boundary Rules

- producer contexts own business semantics and decide when an event should be enqueued
- `eventoutbox` owns durability and publication mechanics, not business behavior
- consumers, projections, realtime fanout, and downstream retries are outside this bounded context
- this package is not directly exposed through REST or GraphQL

## Validate

```bash
go test ./internal/eventoutbox/...
make verify
```

## Performance Readiness

The meaningful performance questions here are backlog growth, oldest-pending age, and publish or reschedule health.

Current practical checks:

```bash
make outbox-diagnose
make record-projection-smoke
```

Track:

- pending count
- oldest pending age
- failed or rescheduled sample rows
- whether downstream projection consumers recover after publish

## Risks And Compatibility Notes

- envelope versioning and topic semantics are compatibility-sensitive for downstream consumers
- reschedule and backoff behavior must stay visible for operator diagnostics
- if publication cadence or retry logic changes, keep this README aligned with `cmd/outbox-publisher`

## Related Docs

- [`../platform/config/README.md`](../platform/config/README.md)
- [`../../cmd/outbox-publisher/README.md`](../../cmd/outbox-publisher/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
