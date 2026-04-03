# Development Scripts

**Path:** `hack/dev`

## Purpose

This folder stores shell helpers for local diagnostics and dev-only troubleshooting.

## Script Inventory

- repo-local scripts in this folder are expected to support investigation, not production execution

## Boundary Rules

- keep scripts explicit, readable, and safe to rerun
- do not move business logic or canonical runtime orchestration here

## Validate

- run the target script against the intended local stack
- confirm it fails clearly when prerequisites are missing

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
