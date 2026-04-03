# Admin Bounded Context

**Path:** `internal/admin`

## Purpose

`internal/admin` owns privileged user-governance operations and role hierarchy enforcement.

## Current HTTP Surface

All current admin routes are authenticated and mounted under `/admin/users/{user_id}`:

| Route | Responsibility |
| --- | --- |
| `PUT /roles` | replace role set through the legacy role-update endpoint |
| `PUT /promote-admin` | add admin privilege with hierarchy validation |
| `PUT /demote-admin` | remove admin privilege with hierarchy validation |
| `PUT /block` | apply blocked role or state |
| `PUT /unblock` | remove blocked role or state |

## Runtime Contract

- transport protection comes from auth middleware
- authorization semantics and role hierarchy validation stay inside the admin core
- role storage remains backend-owned through `aion_api.roles` and `aion_api.user_roles`
- the admin repository also acts as the source-of-truth `RolesReader` consumed by auth-related flows

## Boundary Rules

- there is no current GraphQL admin surface
- admin adapters stay thin; privileged transition rules stay in core
- this context governs roles and user-state transitions, not generic user profile editing

## Validate

```bash
go test ./internal/admin/...
make verify
```

## Risks And Compatibility Notes

- role semantics are security-sensitive and shared with auth/session validation flows
- if route shape or role meaning changes, update both this README and the related auth docs in the same PR

## Related Docs

- [`../auth/README.md`](../auth/README.md)
- [`../user/README.md`](../user/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
