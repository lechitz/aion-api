# Test Setup Helpers

**Path:** `tests/setup`

## Purpose

Shared test builders and helpers for unit testing service and usecase layers.
These helpers reduce boilerplate around gomock setup and default fixtures.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| suite builders | build ready-to-use test suites per domain service |
| mock wiring | create and expose required mocks consistently |
| shared fixtures | provide default entities and helpers for tests |

## Usage Pattern

```go
suite := setup.UserServiceTest(t)
defer suite.Ctrl.Finish()
```

## Boundary Rules

- keep helpers deterministic and focused
- favor explicit fixture builders over hidden global state
- keep suite APIs stable across refactors

## Validate

```bash
go test ./tests/setup/...
go test ./internal/...
```

## Risks And Compatibility Notes

- helper drift can hide missing mock expectations or overfit tests to implementation details

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
