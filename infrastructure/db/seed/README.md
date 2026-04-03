# Database Seed Scripts

**Path:** `infrastructure/db/seed`

## Purpose

This folder owns direct SQL seed artifacts for local development, QA, and deterministic reset flows.
It is optimized for fast Postgres bootstrap, not for production data loading.

## Current Files

| Pattern | Purpose |
| --- | --- |
| `roles.sql`, `user_roles.sql`, `admin_user.sql` | role and admin bootstrap |
| `*_generate.sql` | generated categories, tags, users, and records |
| `test_*.sql` | scenario-oriented datasets for timeline and demo coverage |
| `.env.example` | helper env template for local seed tooling |

## Validate

```bash
make seed-essential
make seed-test
make seed-clean-all
make db-full
```

## Boundary Rules

- keep seed data representative but disposable
- do not let seed scripts redefine schema ownership or hide migration gaps
- API-driven seed callers and synthetic generation flows belong under `hack/tools` and complement, rather than replace, these SQL assets

## Risks And Compatibility Notes

- fake but stable data is useful; fake but misleading data is dangerous when it changes expected contracts or flows

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
