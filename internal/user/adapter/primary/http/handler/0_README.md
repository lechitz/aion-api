# Legacy Marker For User HTTP Handler Docs

**Path:** `internal/user/adapter/primary/http/handler`

## Status

This file is intentionally minimal.
The canonical boundary documentation for the user context lives in [`internal/user/README.md`](../../../../README.md).

## Why This File Still Exists

- this leaf package already had a local doc artifact and keeping a marker here improves discoverability during code search
- route ownership, runtime rules, validation, and compatibility notes must now be maintained only in the parent boundary README

## Rules

- do not expand this file into a second source of truth
- when handler behavior changes, update [`internal/user/README.md`](../../../../README.md) instead
