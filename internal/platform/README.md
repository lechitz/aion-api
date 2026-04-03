# Platform Layer

**Path:** `internal/platform`

## Purpose

`internal/platform` provides the domain-agnostic runtime foundation used by all bounded contexts.

## Current Areas

| Area | Role |
| --- | --- |
| `app/` | shared dependency bundle exposed to primary adapters |
| `config/` | typed env loading and validation |
| `fxapp/` | Fx module composition and lifecycle wiring |
| `httpclient/` | shared OTEL-instrumented outbound HTTP client |
| `observability/` | tracer and metrics bootstrap |
| `ports/` | reusable output-port contracts |
| `server/` | server composition and HTTP runtime wiring |

## Boundary Rules

- keep platform code free from product or domain semantics
- centralize runtime wiring here instead of scattering it across contexts
- if a helper is specific to one context, keep it with that context instead of moving it into `platform`

## Validate

```bash
go test ./internal/platform/...
make verify
```

## Risks And Compatibility Notes

- platform changes tend to have cross-repo and cross-context blast radius even when the code diff looks small
- config, server, and observability changes should stay aligned with the matching infrastructure READMEs

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
