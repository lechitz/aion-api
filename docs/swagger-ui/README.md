# Swagger UI Static Bundle

**Path:** `docs/swagger-ui`

## Purpose

This folder stores the static Swagger UI bundle published by the docs site.

## Current Relationship To Runtime

- the runtime mounts the generated OpenAPI contract
- this folder stores the static assets that render that contract in the published documentation surface
- the REST contract authority still lives under `contracts/openapi`

## Files

| File Pattern | Purpose |
| --- | --- |
| static bundle assets | Swagger UI rendering support |
| config or index assets | site-level wiring for the published REST explorer |

## Validate

```bash
make swag
make docs-verify
```

## Boundary Rules

- do not edit this folder to redefine REST semantics
- if the static UI depends on a new asset or config change, keep it aligned with the generated OpenAPI artifacts

## Risks And Compatibility Notes

- static bundle drift can make the published explorer stale even when the OpenAPI contract itself is correct

---

<!-- doc-nav:start -->
## Navigation
- [Back to docs index](../index.md)
<!-- doc-nav:end -->
