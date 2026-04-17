# User Bounded Context

**Path:** `internal/user`

## Purpose

`internal/user` owns account lifecycle, public registration, profile and avatar management, password changes, and soft deletion.

## Current HTTP Surface

### Public routes

- `POST /user/create`
- `POST /user/avatar/upload`
- `POST /registration/start`
- `PUT /registration/{registration_id}/profile`
- `PUT /registration/{registration_id}/avatar`
- `POST /registration/{registration_id}/complete`

### Authenticated routes

- `GET /user/all`
- `GET /user/me`
- `GET /user/preferences`
- `GET /user/{user_id}`
- `PUT /user/`
- `PUT /user/preferences`
- `DELETE /user/avatar`
- `PUT /user/password`
- `DELETE /user/`

## Runtime Contract

- password updates refresh the auth cookie or token on success
- appearance preferences are cached locally by the dashboard and can be rehydrated from `/user/preferences`, including preset selection and custom CSS overrides
- cache layers must never store password hashes or raw passwords
- registration is a multi-step public flow separate from the authenticated profile-update surface
- avatar upload and removal is owned here even when backed by external object storage

## Boundary Rules

- identity, password, and registration rules stay in core usecases
- transport adapters own request decoding, auth context extraction, and cookie refresh wiring
- auth and session semantics collaborate with `internal/auth`, but profile and account ownership stay here

## Validate

```bash
go test ./internal/user/...
go test ./internal/platform/server/http/utils/cookies/...
make verify
```

## Risks And Compatibility Notes

- registration and avatar flows are user-facing contracts and should not drift silently in transport shape
- if cookie refresh behavior or account-deletion semantics change, keep user and auth docs aligned
- the leaf handler package still contains a legacy marker file for discoverability; the canonical boundary doc remains this README

## Related Docs

- [`../auth/README.md`](../auth/README.md)
- [`../platform/server/http/utils/cookies/README.md`](../platform/server/http/utils/cookies/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
