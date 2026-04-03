# Observability Quickstart

Use this guide to apply and validate the local observability stack (Prometheus, Grafana, Jaeger, Loki).

## Prerequisites

- Docker and Docker Compose running
- Local stack bootstrapped (`make dev`)

## Apply Observability Setup

```bash
./infrastructure/observability/scripts/setup-improvements.sh
```

The script restarts required services and performs baseline checks.
It validates the currently provisioned dashboards and datasources before checking endpoints.

## Validate Services

| Service | URL | Expected |
| --- | --- | --- |
| Grafana | `http://localhost:3000` | Login page available (`aion/aion`) |
| Prometheus | `http://localhost:9090` | Targets up |
| Jaeger | `http://localhost:16686` | Service list available |

## Generate Traffic

```bash
for i in {1..50}; do curl -s http://localhost:5001/aion/api/v1/health > /dev/null; done
```

Then confirm traces and metrics appear in dashboards.

## Grafana Verification

1. Open **aion-api - RED Metrics Dashboard (Professional)**.
2. Confirm latency, error-rate, and throughput panels are populated.
3. Use trace links from exemplars to jump into Jaeger.

## Record A Baseline

For performance-ready documentation, capture at least:

| Metric | Why it matters |
| --- | --- |
| request volume or smoke source | makes the observation reproducible |
| p95 latency | catches slow-tail regressions |
| error rate | separates slow from broken |
| top slow endpoints | identifies where to zoom in |
| trace exemplar or trace id | supports root-cause follow-up |

If the flow under review is projection or realtime sensitive, combine dashboard data with:

```bash
make outbox-diagnose
make record-projection-smoke
make realtime-record-smoke
```

## Loki Verification

In Grafana Explore, choose Loki datasource and run:

```logql
{container_name="/aion-dev-api"} | json
```

Filter by fields like `trace_id` and `request_id` to correlate logs with traces.

## Troubleshooting

- Empty dashboards: verify Prometheus scrape targets are healthy.
- No traces in Jaeger: check OTLP endpoint env vars in API container.
- Missing logs in Loki: verify Fluent Bit container status and labels.

## Next Step

Read [Platform Runtime](platform.md) for runtime-level observability wiring and conventions.
For performance-specific documentation, use [Performance Readiness](performance-readiness.md).
