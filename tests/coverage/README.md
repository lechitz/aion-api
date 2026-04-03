# Coverage Artifacts

**Path:** `tests/coverage`

## Purpose

Location for generated test coverage outputs and related reports.

## Typical Workflow

```bash
make test-cover
make test-cover-detail
```

## Boundary Rules

- treat files in this folder as generated artifacts
- avoid manual edits to generated reports
- keep report generation reproducible via Make targets

## Risks And Compatibility Notes

- stale reports create false confidence when they are treated as current evidence

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
