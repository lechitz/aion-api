# Auth Bounded Context

**Path:** `internal/auth`

## Purpose

`internal/auth` owns login, session validation, refresh-token rotation, logout, and role or session cache integration.
It is the backend authority for authenticated HTTP session state used by REST handlers and GraphQL auth enforcement.

## Current HTTP Surface

| Route | Access | Current behavior |
| --- | --- | --- |
| `POST /auth/login` | public | validate credentials, return access token plus user snapshot, set auth and refresh cookies |
| `POST /auth/refresh` | public | rotate access and refresh tokens from the refresh cookie |
| `GET /auth/session` | authenticated | validate bearer token or auth cookie and return session snapshot |
| `POST /auth/logout` | authenticated | revoke current server-side session state and clear cookies |

## Internal Shape

| Area | Responsibility |
| --- | --- |
| `core/usecase` | `Login`, `Validate`, `RefreshTokenRenewal`, `Logout`, and role-cache orchestration |
| `core/ports/output` | auth provider, auth store, roles reader, and role cache contracts |
| `adapter/secondary/cache` | Redis-backed token, session, and role-cache operations |
| `adapter/primary/http/handler` | REST transport mapping for `/auth/*` routes |
| `adapter/primary/http/middleware` | protected-route middleware that injects authenticated user context |

## Boundary Rules

- cookie transport rules are owned by `internal/platform/server/http/utils/cookies`, not by core usecases
- GraphQL `@auth` enforcement lives in the central GraphQL adapter, but validation semantics ultimately come from this context
- security-sensitive integrations stay behind output ports; transport adapters only decode, map, and emit cookies or responses

## Validate

```bash
go test ./internal/auth/...
go test ./internal/platform/server/http/utils/cookies/...
make verify
```

## Performance Readiness

Current local load validation for this context is the versioned login scenario:

```bash
make load-test-auth-login
```

Use it to watch:

- login transport latency through the real HTTP path
- auth error rate under concurrent local traffic
- regressions caused by session or cache changes even when functional tests still pass

## Risks And Compatibility Notes

- cookie names, max-age behavior, and refresh semantics are consumer-visible session contracts
- role-cache drift can create authorization surprises even when login still works
- if auth semantics change, keep REST handlers, GraphQL directive behavior, and cookie docs aligned

## Related Docs

- [`../platform/server/http/utils/cookies/README.md`](../platform/server/http/utils/cookies/README.md)
- [`../platform/server/http/middleware/servicetoken/README.md`](../platform/server/http/middleware/servicetoken/README.md)
- [`../adapter/primary/graphql/README.md`](../adapter/primary/graphql/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
