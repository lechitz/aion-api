# Performance Readiness

Use this guide when you need to validate or document performance-sensitive behavior in `aion-api`.

## Current State

The repo does not yet ship committed Go microbenchmarks such as `BenchmarkXxx`.

What it does have today:

- RED-style dashboards for HTTP and GraphQL behavior
- smoke checks for record projection, SSE delivery, ingest flow, and outbox health
- operator diagnostics for backlog and pending-event age
- committed local load scenarios with reviewable thresholds for auth login, derived GraphQL reads, dashboard aggregation, and realtime SSE delivery

That means current performance documentation should be based on observed boundary behavior, not invented numbers.

## Validated Local Baseline

Validated on 2026-03-28 against the local multi-repo development stack with:

```bash
./infrastructure/observability/scripts/setup-improvements.sh
make event-backbone-gate
```

Observed readiness signals from that run:

- API health endpoint returned healthy after the stack restart
- Prometheus, Grafana, Jaeger, and the OTel collector responded successfully
- Prometheus reported the `aion-dev-otel-collector` target as up
- Grafana provisioned Prometheus and Jaeger datasources and exposed the RED dashboard
- outbox diagnosis reported `pending_count=0`, `failed_count=0`, and `oldest_pending_age=n/a`
- record projection smoke, realtime SSE smoke, and projection pagination smoke all passed
- ingest smoke passed repeatedly with unique event ids after the smoke payload was made per-run unique
- the integrated dashboard records E2E passed at the end of `make event-backbone-gate`

Use that baseline as readiness evidence, not as a latency budget or throughput SLA.

## Versioned Local Load Scenarios

The repo now ships committed local load scenarios under `hack/tools/load-test/`.

Current scenarios:

| Scenario | Boundary | Command |
| --- | --- | --- |
| `auth-login` | public auth HTTP path | `make load-test-auth-login` |
| `record-projections-latest` | authenticated GraphQL derived read path | `make load-test-record-projections` |
| `dashboard-snapshot` | authenticated GraphQL dashboard aggregation path | `make load-test-dashboard-snapshot` |
| `realtime-record-created` | authenticated SSE delivery path for projection-ready record events | `make load-test-realtime-record-created` |

Combined baseline:

```bash
make load-test-baseline
```

Thresholds live in `hack/tools/load-test/thresholds.json`.
Treat them as local readiness budgets for this workspace, not as production SLOs.

Recent clean local runs on 2026-03-28:

- `auth-login`: 60 measured requests, concurrency 6, p50 about `49-50ms`, p95 about `55-128ms`, error rate `0%`
- `record-projections-latest`: 80 measured requests, concurrency 8, p50 about `6-10ms`, p95 about `15-16ms`, error rate `0%`
- `dashboard-snapshot`: 60-80 measured requests, concurrency 6-8, p50 about `3-7ms`, p95 about `6-26ms`, error rate `0%`
- `realtime-record-created`: 20 measured requests, concurrency 4, p50 about `10-12s`, p95 about `15-16s`, error rate `0%`

Run the load baseline while the dev stack is otherwise idle.
Hot-reload rebuilds or large doc-generation commands can perturb the auth scenario and produce transport noise that is not representative of the API itself.
The realtime scenario is intentionally different: it measures the full local async chain from record creation through outbox publication, Kafka consumption, projection readiness, and final SSE delivery.
Its threshold therefore carries slight headroom above the latest clean local branch measurements rather than pretending the path behaves like a synchronous request.

## What To Measure Today

| Boundary | Current signals |
| --- | --- |
| HTTP and GraphQL transport | throughput, error rate, p50/p75/p95 latency from Grafana |
| Record projection pipeline | successful projection materialization, pagination behavior, backlog age, pending count |
| Realtime SSE | delivery continuity, disconnect behavior, dropped-event logs, and end-to-end async delivery latency |
| Outbox publisher | pending count, oldest pending age, publish or reschedule health |

## Recommended Validation Flow

1. Boot the local stack.

```bash
make dev
./infrastructure/observability/scripts/setup-improvements.sh
```

2. Confirm dashboards and collectors are healthy.

Read [Observability Quickstart](observability-quickstart.md).

3. Exercise the critical boundaries.

```bash
make outbox-diagnose
make record-projection-smoke
make realtime-record-smoke
make record-projection-page-smoke
make load-test-baseline
make event-backbone-gate
```

4. Capture the metrics that matter.

At minimum, record:

- operation or boundary under test
- traffic source or smoke command
- throughput or request volume
- p95 latency when the RED dashboard provides it
- error rate
- any backlog, retry, or dropped-event signal

## When To Add Real Benchmarks

Do not document per-function numbers in a README until the repo has a committed benchmark that can reproduce them.

Good first candidates:

- record search and projection query paths
- dashboard snapshot assembly and insight generation
- outbox publish batching
- GraphQL resolver hot paths that fan into record-heavy reads

Use local benchmarks only when:

- the hot path is isolated enough to measure without full-stack noise
- the benchmark input can be kept deterministic
- the result changes an engineering decision, not just curiosity

## Documentation Rule

Only put performance notes into a boundary README when they are:

- measurable with a committed command or dashboard
- local to that boundary
- likely to influence future changes

If the note is cross-cutting, keep it here and link from the boundary README instead of duplicating a protocol everywhere.
