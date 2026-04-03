# API Collections

**Path:** `docs/collections`

## Purpose

This folder publishes consumer-friendly request collections for manual QA and local integration checks.
The canonical REST contract still lives in `contracts/openapi/swagger.yaml`.

## Current Artifact

| Path | Purpose |
| --- | --- |
| `postman/aion-api.postman_collection.json` | Postman collection covering auth, user, admin, chat, GraphQL, and health flows |

## Validate

- import the collection into Postman
- set `{{baseURL}}` to the target API origin
- use collection variables or cookies for auth
- after changing REST or GraphQL flows, verify that the checked-in collection still exercises the intended happy paths

## Boundary Rules

- update this artifact when consumer-facing manual QA flows change materially
- do not treat the collection as a source of truth over OpenAPI or runtime behavior
- keep secrets, personal profiles, and local-only values outside the checked-in JSON

---

<!-- doc-nav:start -->
## Navigation
- [Back to docs index](../index.md)
<!-- doc-nav:end -->
