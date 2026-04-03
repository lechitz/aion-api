# Seed Caller Tool

**Path:** `hack/tools/seed-caller`

## Purpose

This tool exercises API-driven seed flows used during local data bootstrap and diagnostics.

## Runtime Flow

- load the expected runtime config
- call the target seed path or helper flow
- report failures clearly for local investigation

## Quick Run

```bash
make seed-api-caller
go run ./hack/tools/seed-caller
```

## Boundary Rules

- this tool complements SQL seed assets; it does not replace them
- it should reuse canonical app behavior instead of forking seed semantics

## Validate

- run the Make target against a working local stack
- confirm failures are visible and actionable
