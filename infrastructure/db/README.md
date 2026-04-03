# Database Infrastructure

**Path:** `infrastructure/db`

## Purpose

This folder owns versioned schema changes and deterministic local data bootstrap assets.

## Current Areas

| Area | Responsibility |
| --- | --- |
| `migrations/` | schema evolution and rollback-safe SQL |
| `seed/` | disposable local or QA data sets and helper seed assets |

## Boundary Rules

- migrations own schema shape; seeds must not quietly compensate for missing schema work
- SQL here is operational code and should stay reviewable and deterministic

## Validate

```bash
make migrate-up
make seed-all
```

## Risks And Compatibility Notes

- drift between migrations and seeds creates false confidence in local setups

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
