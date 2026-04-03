# Observability Infrastructure

**Path:** `infrastructure/observability`

## Purpose

This folder owns the repo-local telemetry configs mounted by the Docker stack.
It covers traces, metrics, logs, and Grafana provisioning for local and prod-like execution paths.

## Current Signal Flow

1. services emit OTLP telemetry to the collector
2. `otel/` exports traces and metrics onward
3. `prometheus/` scrapes the collector
4. `fluentbit/` tails Docker logs and forwards them to `loki/`
5. `grafana/` provisions datasources and dashboards over Prometheus, Loki, and trace backends

## Current Areas

| Area | Responsibility |
| --- | --- |
| `otel/` | telemetry ingestion and routing |
| `prometheus/` | metrics scraping |
| `loki/` | log storage and query backend |
| `fluentbit/` | log collection and forwarding |
| `grafana/` | datasource and dashboard provisioning |
| `scripts/` | helper automation for validation, not the primary source of truth |

## Boundary Rules

- compose profiles and config files are canonical; helper scripts are secondary
- keep telemetry wiring deterministic and versioned in-repo
- if a dashboard or query depends on a label or exporter, update the relevant config in the same change

## Validate

```bash
make dev
./infrastructure/observability/scripts/setup-improvements.sh
```

## Performance Signals

This boundary is where transport and pipeline performance become visible.

Current high-signal checks:

- RED dashboard p50, p75, and p95 latency
- error-rate panels
- top slow endpoints
- trace exemplars for slow requests
- log correlation by `trace_id` and `request_id`

For the full operator-facing workflow, see [`../../docs/performance-readiness.md`](../../docs/performance-readiness.md).

## Risks And Compatibility Notes

- telemetry label changes can silently break dashboards and alert queries even when the stack still starts

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
