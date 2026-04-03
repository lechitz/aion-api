# Primary GraphQL Adapter

**Path:** `internal/adapter/primary/graphql`

## Purpose

This package is the central GraphQL transport entrypoint for the application.
It hosts schema composition, gqlgen-generated artifacts, directive wiring, and thin resolvers that delegate to context controllers.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Server setup | build the GraphQL HTTP handler and configure transports or middleware |
| Schema composition | maintain root schema and module-based schema extensions |
| Directive wiring | register cross-cutting directives such as `@auth` |
| Resolver bridge | convert GraphQL inputs and context into controller calls |
| Codegen integration | define gqlgen generation targets and configuration |

## Structure

| Path | Role |
| --- | --- |
| `schema/root.graphqls` | root query or mutation and shared directives or scalars |
| `schema/modules/*.graphqls` | domain extensions for `category`, `tags`, `record`, `chat`, and `user` |
| `directives/auth.go` | authorization directive implementation |
| `resolver.go` | dependency wiring from services to context GraphQL controllers |
| `*.resolvers.go` | thin field resolvers generated or preserved by gqlgen |
| `server.go` | HTTP router, recovery middleware, auth middleware, and transports |
| `generated.go` | gqlgen generated execution engine; do not edit manually |
| `model/models_gen.go` | gqlgen generated GraphQL transport models |
| `gqlgen.yml` | gqlgen configuration |

## GraphQL Runtime Flow

1. HTTP request enters the GraphQL handler.
2. Recovery middleware wraps request execution.
3. Auth middleware populates request context when auth service is configured.
4. gqlgen executes the operation and runs directives.
5. Resolver reads context values and delegates to the context controller.
6. Controller maps GraphQL DTOs to usecase calls and returns response data or errors.

## Directive Model

| Directive | Behavior |
| --- | --- |
| `@auth(roles: String)` | requires `ctxkeys.UserID`; enforces role when provided; bypasses role checks for trusted `ctxkeys.ServiceAccount=true` calls |

## Configured Transports

| Transport | Status |
| --- | --- |
| `GET` | enabled |
| `POST` | enabled |
| `OPTIONS` | enabled |
| `MultipartForm` | enabled |
| `Websocket` | enabled (`KeepAlivePingInterval: 10s`) |

## Boundary Rules

- resolvers must remain thin and should not host domain business logic
- context-specific mapping or orchestration belongs to `internal/<ctx>/adapter/primary/graphql/controller`
- generated files should be managed by codegen only
- module-based schema composition keeps domain contracts isolated while exposing a single endpoint

## Validate

```bash
go test ./internal/adapter/primary/graphql/...
make graphql
make graphql.validate
```

## Risks And Compatibility Notes

- schema, resolvers, and shared contract docs must move together; drift here is a public contract regression
- trusted service-account bypass in `@auth` is a security-sensitive rule and should stay explicit

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
