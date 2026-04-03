# Developer Utilities (`hack`)

**Path:** `hack`

## Purpose

This folder contains development-only tools and scripts.
It follows the common `hack/` convention for non-production operational helpers.

## Subfolders

| Folder | Responsibility |
| --- | --- |
| `tools/` | Go CLIs for seed and utility workflows |
| `dev/` | shell scripts for diagnostics and local troubleshooting |

## Usage Examples

```bash
make seed-api-caller
make seed-helper
go run ./hack/tools/graph-projection-export --user-id 999
make record-projection-smoke
make realtime-record-smoke
make record-projection-page-smoke
make ingest-event-smoke
make outbox-diagnose
make load-test-baseline
make event-backbone-gate-preflight
make event-backbone-gate
bash hack/dev/test-chat.sh
```

## Boundary Rules

- keep this folder out of production image and runtime paths
- use it for reproducible local workflows and debugging support
- domain logic must remain in `internal/`, not in `hack/` scripts or tools

## Risks And Compatibility Notes

- utilities here are useful for validation, but they are not public runtime contracts
- if a tool becomes required for normal operator flow, promote it into a better-owned boundary instead of hiding it under `hack/`

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
