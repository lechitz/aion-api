# Internal Application Layer

**Path:** `internal`

## Purpose

`internal` contains the application code that must stay inside this module.
The repo is organized around bounded contexts plus a small set of cross-cutting runtime layers.

## Current Areas

| Area | Role |
| --- | --- |
| `admin`, `audit`, `auth`, `category`, `chat`, `eventoutbox`, `realtime`, `record`, `tag`, `user` | bounded contexts with core ports, usecases, and local adapters |
| `adapter/` | shared adapter infrastructure reused across contexts |
| `platform/` | config, Fx wiring, runtime services, ports, and server composition |
| `shared/` | stable cross-cutting constants and key namespaces |

## Boundary Rules

- context business rules belong inside the owning bounded context
- shared transport or infra helpers belong in `adapter/` or `platform/` only when they are reused across contexts
- `shared/` stays intentionally small; context-specific constants should remain local when possible
- add a local `README.md` only for meaningful boundaries, not for purely structural folders

## Validate

```bash
go test ./internal/...
make verify
```

## Documentation Notes

- treat the nearest boundary README as the canonical explanation for ownership, dependencies, risks, and validation
- if a context changes transport surface, dependency direction, or runtime contract, update that local README in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
