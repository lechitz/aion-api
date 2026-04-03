# Database Migrations

**Path:** `infrastructure/db/migrations`

## Purpose

This folder stores versioned SQL migrations that define schema evolution for `aion-api`.

## Current Shape

- one migration pair per change set
- forward and rollback SQL kept explicit
- migration order encoded in file naming

## Validate

```bash
make migrate-up
make migrate-down
make verify
```

## Boundary Rules

- migrations own schema history and must stay reversible when the change shape allows it
- seeds and application code should not be used to compensate for broken or incomplete migrations

## Risks And Compatibility Notes

- irreversible or partially reversible migrations require stronger review because they outlive the code that introduced them

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
