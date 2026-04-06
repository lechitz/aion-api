# Makefile Modules

**Path:** `makefiles`

## Purpose

Modular Makefile fragments included by the root `Makefile`.
They centralize commands for build, environment, migrations, codegen, tests, and quality workflows.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| command modularization | split make targets by concern or domain |
| developer workflow | expose reproducible local commands |
| CI parity | keep local verify steps aligned with pipeline checks |

## Typical Command Areas

- Docker and runtime
- migrations and seeds
- codegen such as GraphQL and mocks
- tests and coverage
- performance-readiness smokes and diagnostics
- versioned local load scenarios
- lint and verify

## Performance-Ready Commands

These commands are the current repo-supported performance-readiness path:

```bash
make outbox-diagnose
make record-projection-smoke
make realtime-record-smoke
make record-projection-page-smoke
make load-test-auth-login
make load-test-record-projections
make load-test-dashboard-snapshot
make load-test-realtime-record-created
make event-backbone-gate
```

`make event-backbone-gate` is the strongest integrated readiness signal because it combines the backend smokes with the dashboard records E2E.
`make load-test-baseline` is the current versioned local load layer; it enforces committed thresholds for auth login, derived record GraphQL reads, dashboard snapshot aggregation, and realtime SSE delivery.
These commands are not microbenchmarks. They validate system behavior at the boundaries where latency and pipeline health matter most today.
The realtime load target is intentionally preconfigured with a smaller request profile than the synchronous scenarios because it measures the full async chain from record creation through outbox, Kafka, projection, and SSE delivery.

## Docker Build Targets (`docker.mk`)

Key targets that affect local disk usage and build performance:

| Target | What it does |
| --- | --- |
| `build-dev` | builds all service images via Docker Compose, then runs `docker image prune -f` to remove dangling images automatically |
| `rebuild-dev` | forces a clean rebuild of all images, then runs `docker image prune -f` before resuming the stack |

**Why it matters:** before the auto-prune step, repeated `make build-dev` or `make rebuild-dev` calls would accumulate tens of gigabytes of dangling images over iterative development sessions. The automatic cleanup keeps disk usage stable without requiring manual `docker system prune` runs.

If `docker image prune` is too aggressive for your local setup, scope the cleanup manually with `docker image prune --filter "until=24h"`.

## Boundary Rules

- keep the root `Makefile` thin; logic belongs in modules
- keep target naming predictable and discoverable
- keep environment variable requirements documented near the target that uses them

## Validate

```bash
make verify
```

## Risks And Compatibility Notes

- target drift creates confusion between local workflows and CI behavior even when the commands still exist

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
