# Seed Helper Tool

**Path:** `hack/tools/seed-helper`

## Purpose

This tool supports local seed and setup workflows for `aion-api`.

## Main Commands

- use the matching Make target when available
- use direct `go run` only when validating the tool itself or debugging its behavior

## Boundary Rules

- this tool supports seed workflows; it does not own canonical seed data or schema shape
- SQL seeds stay under `infrastructure/db/seed`

## Validate

```bash
make seed-helper
go test ./hack/tools/seed-helper/...
```
