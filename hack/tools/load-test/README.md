# Versioned Load Test Tool

**Path:** `hack/tools/load-test`

## Purpose

This tool provides committed, repeatable local load scenarios for `aion-api`.
It is the current versioned layer between smoke checks and future microbenchmarks.

## Current Scenarios

| Scenario | Purpose |
| --- | --- |
| `auth-login` | exercise the public auth login path with the seeded test user |
| `record-projections-latest` | exercise the authenticated GraphQL read path for derived record projections |
| `dashboard-snapshot` | exercise the authenticated GraphQL dashboard aggregation path |
| `realtime-record-created` | exercise the authenticated SSE stream until a created record event is delivered through the full async pipeline |

Thresholds for these scenarios live in `thresholds.json`.

Current committed thresholds:

- `auth-login`: error rate `0%`, p50 `<= 100ms`, p95 `<= 150ms`
- `record-projections-latest`: error rate `0%`, p50 `<= 25ms`, p95 `<= 75ms`
- `dashboard-snapshot`: error rate `0%`, p50 `<= 15ms`, p95 `<= 50ms`
- `realtime-record-created`: error rate `0%`, p50 `<= 12500ms`, p95 `<= 16500ms`

## Usage

```bash
go run ./hack/tools/load-test --scenario auth-login
go run ./hack/tools/load-test --scenario record-projections-latest
go run ./hack/tools/load-test --scenario dashboard-snapshot
go run ./hack/tools/load-test --scenario realtime-record-created
make load-test
make load-test-baseline
```

When invoked as `--scenario realtime-record-created` without extra flags, the tool automatically switches to the committed local baseline profile for that async path:

- `requests=20`
- `concurrency=4`
- `warmup=2`
- `timeout=30s`

## Boundary Rules

- keep thresholds committed and reviewable; do not hide them in shell history
- use seeded, reproducible credentials and shared query documents where possible
- treat these scenarios as local readiness checks, not production SLAs

## Validate

```bash
go test ./hack/tools/load-test
make load-test-baseline
```

## Risks And Compatibility Notes

- results reflect the local multi-repo stack and current observability wiring
- thresholds should move only after a measured run and a documented reason
- avoid running doc rebuilds or other hot-reload-triggering commands in parallel with this tool; they can distort local latency and connection stability
- the realtime scenario measures outbox, Kafka, projection, and SSE delivery together; do not compare its numbers directly to synchronous HTTP or GraphQL reads
- the current realtime thresholds intentionally include small headroom above recent clean local runs because this path is dominated by async pipeline timing rather than transport-only latency
