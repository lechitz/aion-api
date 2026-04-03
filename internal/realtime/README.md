# Realtime Bounded Context

**Path:** `internal/realtime`

## Purpose

`internal/realtime` owns authenticated server-sent event fan-out for projection-ready record updates.
It bridges derived backend events into per-user live streams without turning the HTTP layer into a stateful business boundary.

## Current Surface

| Surface | Responsibility |
| --- | --- |
| `core/ports/input.Service.Publish` | fan out one realtime event to subscribers of the same user |
| `core/ports/input.Service.Subscribe` | open one per-user stream and return a cleanup function |
| HTTP `GET /realtime{cfg.Realtime.StreamPath}` | authenticated SSE stream for the current user |
| `adapter/secondary/kafka` | read projection-ready events from Kafka and publish them into the in-memory service |

## Current Shape

| Area | Responsibility |
| --- | --- |
| `core/usecase` | in-memory per-user publish/subscribe service |
| `adapter/primary/http/handler` | SSE transport, auth context extraction, framing, and disconnect handling |
| `adapter/secondary/kafka` | projection-event reader that feeds the service after derived rows are ready |

## Boundary Rules

- this context owns delivery mechanics for live projection updates, not the business rules that produce source events
- Kafka envelope semantics stay in upstream event contracts; this boundary only consumes projection-ready inputs
- HTTP handlers must stay transport-only and should not invent filtering, aggregation, or authorization semantics beyond authenticated user scope

## Validate

```bash
go test ./internal/realtime/...
go test ./internal/platform/server/http/...
make verify
```

## Performance Readiness

The meaningful performance question here is delivery continuity, not raw request throughput.

Current practical check:

```bash
make realtime-record-smoke
make load-test-realtime-record-created
```

Watch for:

- end-to-end delivery after projection readiness
- disconnect or reconnect stability
- dropped-event logs caused by bounded subscriber buffers
- latency correlation with projection or outbox health rather than SSE framing alone
- whether the committed realtime SSE scenario stays inside the local latency thresholds for the full async path, not just HTTP stream setup

## Risks And Compatibility Notes

- subscriber state is intentionally in-memory and process-local
- backpressure is handled by bounded subscriber buffers; slow consumers can miss events instead of stalling the whole stream
- realtime truth depends on the projection path being healthy; if projection materialization drifts, this surface degrades before it fails in transport terms

## Related Docs

- [`../record/README.md`](../record/README.md)
- [`../eventoutbox/README.md`](../eventoutbox/README.md)
- [`../platform/server/http/README.md`](../platform/server/http/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
