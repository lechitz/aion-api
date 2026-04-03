# Loki Configuration

**Path:** `infrastructure/observability/loki`

## Purpose

This folder stores the Loki config used by the local observability stack.

## Current File

| File | Purpose |
| --- | --- |
| config file in this folder | log storage and query backend configuration |

## Validate

- start the stack
- send logs through Fluent Bit
- confirm queries return recent application logs

## Boundary Rules

- keep storage and query backend behavior here, not in dashboard or app docs

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
