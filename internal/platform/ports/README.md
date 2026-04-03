# Platform Output Ports

**Path:** `internal/platform/ports`

## Purpose

This folder defines shared output-port contracts reused across the application.

## Current Output Ports

| Package | Role |
| --- | --- |
| `output/cache` | cache operations and lifecycle |
| `output/db` | database abstraction and transaction support |
| `output/hasher` | password hashing |
| `output/httpclient` | outbound HTTP requests with instrumentation |
| `output/keygen` | key generation helpers |
| `output/logger` | structured logger contract |

## Boundary Rules

- keep interfaces small and technology-agnostic
- ports exist to decouple core logic from vendor implementations
- adding a new port here means the contract is shared across multiple boundaries, not just one context

## Validate

```bash
go test ./internal/platform/ports/...
make verify
```

## Risks And Compatibility Notes

- port changes can cascade across multiple contexts, mocks, and secondary adapters
- if a contract is only useful to one bounded context, keep it local instead of promoting it here

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
