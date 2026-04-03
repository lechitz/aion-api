# Platform Configuration

**Path:** `internal/platform/config`

## Purpose

This package is the canonical source for typed runtime configuration, env defaults, and cross-section validation used to bootstrap `aion-api`.

## Canonical Source

- field definitions, env tags, and defaults live in `environments.go`
- cross-field validation rules live in `config.go`

When docs conflict with those files, the code wins.

## Current Config Sections

| Section | What it controls |
| --- | --- |
| `General` | app name, environment, version |
| `Observability` | OTLP endpoint, service identity, exporter headers, compression, and timeouts |
| `ServerHTTP` | HTTP host, port, context, API root, Swagger or docs mounts, health paths, and timeouts |
| `ServerGraphql` | GraphQL host, path, and transport limits |
| `DB` | PostgreSQL connectivity and pool or retry settings |
| `Cache` | Redis address, DB isolation by bounded context, pool, and timeout |
| `Kafka` | broker list and canonical topic names |
| `Outbox` | batch size, publish interval, and enabled flag |
| `Realtime` | SSE path, consumer-group prefix, heartbeat, and subscriber buffer |
| `Cookie` | auth cookie domain, path, same-site, secure, and max-age |
| `AionChat` | external `aion-chat` base URL, service key, and timeout |
| `AvatarStorage` | S3-compatible avatar storage configuration |
| `Application` | shutdown timeout and request context timeout |

## Validation Coverage

`Config.Validate()` currently enforces:

- HTTP and GraphQL path, timeout, and header constraints
- cache and DB minimums or required fields
- observability endpoint and compression rules
- Kafka topic and broker requirements
- outbox and realtime runtime minimums
- application shutdown and runtime constraints

## Boundary Rules

- add new env keys by extending the typed structs first, then validation if needed
- do not duplicate long env key tables in distant READMEs; link back here instead
- any contract-visible change to paths, topics, cookie behavior, or external endpoints must update the nearest consumer or operator docs in the same PR

## Validate

```bash
go test ./internal/platform/config/...
make verify
```

## Risks And Compatibility Notes

- config drift is dangerous because the service may still boot with defaults that no longer reflect intended runtime behavior
- path, topic, and timeout changes are compatibility-sensitive for operators and sibling services

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
