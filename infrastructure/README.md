# Infrastructure Layer

**Path:** `infrastructure`

## Purpose

This layer owns repo-local infrastructure assets used to run, provision, and observe `aion-api`.

## Current Areas

| Area | Responsibility |
| --- | --- |
| `db/` | schema migrations and local seed datasets |
| `docker/` | image build assets and container runtime wiring |
| `observability/` | OTEL, Prometheus, Loki, Fluent Bit, and Grafana configs |

## Boundary Rules

- keep business logic out of `infrastructure`
- treat versioned SQL, Docker assets, and telemetry configs as code
- cross-repo orchestration still depends on sibling repos in the workspace; this folder owns only the `aion-api` side of that wiring

## Validate

```bash
make dev
make verify
```

## Risks And Compatibility Notes

- runtime assets here often become the first source of operator truth during incidents, so stale docs create real operational drift

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
