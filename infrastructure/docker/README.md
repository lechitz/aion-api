# Docker Infrastructure

**Path:** `infrastructure/docker`

## Purpose

This folder owns the container image, entrypoint, and Docker-facing runtime assets used to run `aion-api`.

## Current Layout

| Path | Responsibility |
| --- | --- |
| `Dockerfile` | multi-stage image that builds `aion-api` and `aion-api-outbox-publisher` |
| `scripts/entrypoint.sh` | default container entrypoint; starts `aion-api` |
| profile-specific compose assets in this area | runtime wiring for local and prod-like execution paths |

## Operational Flows

```bash
make build-dev
make dev
make rebuild-dev
make prod-up
```

The main local stack also runs the outbox publisher, observability stack, and sibling services from `aion-chat`, `aion-web`, `aion-ingest`, and `aion-streams`.

## Boundary Rules

- an isolated clone of `aion-api` can build the image, but the full `make dev` flow assumes the complete `/Aion` workspace
- keep machine-specific values out of the root `Dockerfile`
- if container behavior differs from runtime code, the entrypoint script, compose profile, and Make targets are the canonical sources
- this folder owns image and compose wiring only; app configuration remains in `internal/platform/config`

## Validate

```bash
make build-dev
make dev
```

## Risks And Compatibility Notes

- image, entrypoint, and compose wiring must stay aligned with `cmd/api`, `cmd/outbox-publisher`, and `internal/platform/config`
- the Docker surface is workspace-aware; document local assumptions here rather than leaking them into unrelated runtime READMEs

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
