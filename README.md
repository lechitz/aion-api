# aion-api

aion-api is a production-oriented Go backend that exposes REST and GraphQL APIs for habit and diary workflows, built with Hexagonal/Clean Architecture and strong observability.

## Why This Project

aion-api focuses on three goals:

- keep business logic isolated from transport and infrastructure
- provide stable API contracts for multiple clients
- keep operations visible and debuggable in local and production-like stacks

## Quick Links

- Documentation portal: [aion-api Docs](https://lechitz.github.io/aion-api/)
- REST explorer: [Swagger UI](https://lechitz.github.io/aion-api/swagger-ui/)
- Shared GraphQL contracts: [`contracts/graphql/README.md`](./contracts/graphql/README.md)
- OpenAPI contract: `contracts/openapi/swagger.yaml`
- GraphQL schema artifact: [`docs/graphql/schema.graphql`](./docs/graphql/schema.graphql)
- Documentation ownership map: [`.github/DOCUMENTATION_OWNERSHIP.md`](./.github/DOCUMENTATION_OWNERSHIP.md)

## Architecture At A Glance

| Layer | Purpose |
| --- | --- |
| `internal/<ctx>/core` | domain, ports, and usecases |
| `internal/<ctx>/adapter/primary` | HTTP and GraphQL input adapters |
| `internal/<ctx>/adapter/secondary` | DB, cache, Kafka, and provider adapters |
| `internal/platform` | bootstrap, config, server, observability, and shared runtime contracts |
| `internal/realtime` | authenticated SSE fan-out for projection-ready updates |
| `infrastructure` | Docker, migrations, seeds, and observability assets |

## Core Stack

- Go
- Chi
- gqlgen
- PostgreSQL + GORM
- Redis
- OpenTelemetry + Prometheus + Grafana + Loki
- Docker / Docker Compose

## Fast Local Workflow

```bash
make tools-install
make dev
make migrate-up
make seed-all
make verify
```

## Workspace Model

`aion-api` is the operational hub of the current Aion v2 local stack.

Current integrated development assumes a multi-repo workspace with sibling repositories beside this one, including:

- `aion-web`
- `aion-chat`
- `aion-ingest`
- `aion-streams`

Implications:

- `make build-dev` and `make dev` are intended for this multi-repo workspace, not for an isolated clone of `aion-api`
- the `event-backbone-gate` workflow and preflight are designed for a self-hosted runner with that workspace already available
- if you clone only `aion-api`, some integrated dev and runtime validation flows will not work until those sibling repos are also present

## Quality Gates

```bash
make test
make test-cover-detail
make docs-verify
make graphql.queries graphql.manifest graphql.validate
make verify
```

## Performance Readiness

Current state:

- the repo has strong smoke and observability support, but no committed Go microbenchmarks yet
- performance validation today is boundary-level and system-level, not synthetic per-function benchmarking

Current practical checks:

```bash
make dev
./infrastructure/observability/scripts/setup-improvements.sh
make outbox-diagnose
make record-projection-smoke
make realtime-record-smoke
make load-test-baseline
make event-backbone-gate
```

Read the protocol in [`docs/performance-readiness.md`](./docs/performance-readiness.md) before documenting any latency or throughput claim in a boundary README.
That guide now also records the latest validated local baseline and the integrated gate that produced it.

## GraphQL Contract Workflow

```bash
make graphql.queries
make graphql.manifest
make graphql.validate
make graphql.check-dirty
```

## Documentation Model

When docs overlap, use this order:

1. canonical contracts and generated artifacts under `contracts/` and `docs/graphql/`
2. the nearest boundary `README.md`
3. this root README

Practical rules:

- update the nearest boundary README in the same PR when ownership, dependencies, runtime behavior, or validation changes
- keep folder-level READMEs focused on the boundary they own; do not duplicate policy text already covered elsewhere
- use [`.github/DOCUMENTATION_OWNERSHIP.md`](./.github/DOCUMENTATION_OWNERSHIP.md) when two docs appear to overlap

## Canonical v1 Insight Surface

The v1 personal-intelligence layer is intentionally narrow and backend-owned.

Canonical GraphQL operations:

- `insightFeed`
- `analyticsSeries`

Current contract rules:

- `aion-api` is the authority for schema, resolver behavior, and shared GraphQL artifacts
- shared query documents under `contracts/graphql` must stay aligned with the live schema
- consumers such as `aion-web` and `aion-chat` may adapt presentation, but must not invent richer business semantics than the backend exposes

Current v1 scope model:

- recency windows: `WINDOW_7D`, `WINDOW_30D`, `WINDOW_90D`
- optional `date`
- optional `timezone`
- optional `categoryId`
- optional `tagIds`

Current v1 series support:

- `analyticsSeries` is intentionally narrow
- `records.count` is the canonical v1 series key

Current v1 insight semantics:

- deterministic, explainable insights
- dominant insight is the first item in `insightFeed`
- secondary insights remain ordered after the dominant item
- consumers should treat `status`, `confidence`, `summary`, `recommendedAction`, and `evidence` as backend-owned meaning

Related references:

- [`contracts/graphql/queries/README.md`](./contracts/graphql/queries/README.md)
- [`docs/graphql/README.md`](./docs/graphql/README.md)
- [`internal/record/README.md`](./internal/record/README.md)

<!-- docs-index:start -->
<details>
<summary><strong>Documentation Index</strong></summary>

Repository README map by area.

### cmd
- [`cmd/README.md`](./cmd/README.md)
- [`cmd/api/README.md`](./cmd/api/README.md)
- [`cmd/outbox-publisher/README.md`](./cmd/outbox-publisher/README.md)

### contracts
- [`contracts/graphql/README.md`](./contracts/graphql/README.md)
- [`contracts/graphql/mutations/README.md`](./contracts/graphql/mutations/README.md)
- [`contracts/graphql/queries/README.md`](./contracts/graphql/queries/README.md)
- [`contracts/openapi/README.md`](./contracts/openapi/README.md)

### docs
- [`docs/assets/README.md`](./docs/assets/README.md)
- [`docs/collections/README.md`](./docs/collections/README.md)
- [`docs/diagram/README.md`](./docs/diagram/README.md)
- [`docs/graphql/README.md`](./docs/graphql/README.md)
- [`docs/performance-readiness.md`](./docs/performance-readiness.md)
- [`docs/swagger-ui/README.md`](./docs/swagger-ui/README.md)

### hack
- [`hack/README.md`](./hack/README.md)
- [`hack/dev/README.md`](./hack/dev/README.md)
- [`hack/tools/graph-projection-export/README.md`](./hack/tools/graph-projection-export/README.md)
- [`hack/tools/load-test/README.md`](./hack/tools/load-test/README.md)
- [`hack/tools/seed-caller/README.md`](./hack/tools/seed-caller/README.md)
- [`hack/tools/seed-helper/README.md`](./hack/tools/seed-helper/README.md)

### infrastructure
- [`infrastructure/README.md`](./infrastructure/README.md)
- [`infrastructure/db/README.md`](./infrastructure/db/README.md)
- [`infrastructure/db/migrations/README.md`](./infrastructure/db/migrations/README.md)
- [`infrastructure/db/seed/README.md`](./infrastructure/db/seed/README.md)
- [`infrastructure/docker/README.md`](./infrastructure/docker/README.md)
- [`infrastructure/docker/scripts/README.md`](./infrastructure/docker/scripts/README.md)
- [`infrastructure/observability/README.md`](./infrastructure/observability/README.md)
- [`infrastructure/observability/fluentbit/README.md`](./infrastructure/observability/fluentbit/README.md)
- [`infrastructure/observability/grafana/README.md`](./infrastructure/observability/grafana/README.md)
- [`infrastructure/observability/loki/README.md`](./infrastructure/observability/loki/README.md)
- [`infrastructure/observability/otel/README.md`](./infrastructure/observability/otel/README.md)
- [`infrastructure/observability/prometheus/README.md`](./infrastructure/observability/prometheus/README.md)
- [`infrastructure/observability/scripts/README.md`](./infrastructure/observability/scripts/README.md)

### internal
- [`internal/README.md`](./internal/README.md)
- [`internal/adapter/README.md`](./internal/adapter/README.md)
- [`internal/adapter/primary/README.md`](./internal/adapter/primary/README.md)
- [`internal/adapter/primary/graphql/README.md`](./internal/adapter/primary/graphql/README.md)
- [`internal/adapter/secondary/README.md`](./internal/adapter/secondary/README.md)
- [`internal/admin/README.md`](./internal/admin/README.md)
- [`internal/audit/README.md`](./internal/audit/README.md)
- [`internal/auth/README.md`](./internal/auth/README.md)
- [`internal/category/README.md`](./internal/category/README.md)
- [`internal/chat/README.md`](./internal/chat/README.md)
- [`internal/eventoutbox/README.md`](./internal/eventoutbox/README.md)
- [`internal/platform/README.md`](./internal/platform/README.md)
- [`internal/platform/config/README.md`](./internal/platform/config/README.md)
- [`internal/platform/fxapp/README.md`](./internal/platform/fxapp/README.md)
- [`internal/platform/httpclient/README.md`](./internal/platform/httpclient/README.md)
- [`internal/platform/observability/README.md`](./internal/platform/observability/README.md)
- [`internal/platform/ports/README.md`](./internal/platform/ports/README.md)
- [`internal/platform/server/README.md`](./internal/platform/server/README.md)
- [`internal/platform/server/http/README.md`](./internal/platform/server/http/README.md)
- [`internal/platform/server/http/generic/README.md`](./internal/platform/server/http/generic/README.md)
- [`internal/platform/server/http/middleware/README.md`](./internal/platform/server/http/middleware/README.md)
- [`internal/platform/server/http/middleware/servicetoken/README.md`](./internal/platform/server/http/middleware/servicetoken/README.md)
- [`internal/platform/server/http/ports/README.md`](./internal/platform/server/http/ports/README.md)
- [`internal/platform/server/http/router/README.md`](./internal/platform/server/http/router/README.md)
- [`internal/platform/server/http/utils/README.md`](./internal/platform/server/http/utils/README.md)
- [`internal/platform/server/http/utils/cookies/README.md`](./internal/platform/server/http/utils/cookies/README.md)
- [`internal/platform/server/http/utils/httpresponse/README.md`](./internal/platform/server/http/utils/httpresponse/README.md)
- [`internal/platform/server/http/utils/sharederrors/README.md`](./internal/platform/server/http/utils/sharederrors/README.md)
- [`internal/realtime/README.md`](./internal/realtime/README.md)
- [`internal/record/README.md`](./internal/record/README.md)
- [`internal/shared/README.md`](./internal/shared/README.md)
- [`internal/shared/constants/README.md`](./internal/shared/constants/README.md)
- [`internal/tag/README.md`](./internal/tag/README.md)
- [`internal/user/README.md`](./internal/user/README.md)

### agents
- [`agents/personas/README.md`](./agents/personas/README.md)
- [`agents/playbooks/README.md`](./agents/playbooks/README.md)
- [`agents/review/README.md`](./agents/review/README.md)
- [`agents/standards/README.md`](./agents/standards/README.md)

### makefiles
- [`makefiles/README.md`](./makefiles/README.md)

### tests
- [`tests/coverage/README.md`](./tests/coverage/README.md)
- [`tests/setup/README.md`](./tests/setup/README.md)

</details>
<!-- docs-index:end -->

## License

This repository is source-available but proprietary.

- no right to use, copy, modify, distribute, deploy, or create derivative works is granted without prior written authorization
- see [LICENSE](./LICENSE) for the binding terms
