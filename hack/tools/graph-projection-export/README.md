# Graph Projection Export

**Path:** `hack/tools/graph-projection-export`

## Purpose

Internal tool for exporting the `graph-projection-v1` payload from canonical Aion entities.

## Usage

```bash
go run ./hack/tools/graph-projection-export --user-id 999
go run ./hack/tools/graph-projection-export --user-id 999 --window WINDOW_90D --output ./tmp/graph.json
go run ./hack/tools/graph-projection-export --user-id 999 --category-id 3 --tag-ids 14,15
make graph-projection-export GRAPH_PROJECTION_USER_ID=999 GRAPH_PROJECTION_WINDOW=WINDOW_30D
```

## Inputs

- required: `--user-id`
- optional: `--window`, `--date`, `--timezone`, `--category-id`, `--tag-ids`, `--output`
- the `make graph-projection-export` target loads the dev env file, forces `DB_HOST=localhost`, and forwards those flags

## Boundary Rules

- this is a dev and debugging tool, not a public contract or runtime endpoint
- graph rules remain owned by `internal/record/core/...`
- the tool bootstraps config, DB, repositories, and canonical mappers; it should not fork business logic

## Validate

```bash
go test ./hack/tools/graph-projection-export/...
make graph-projection-export GRAPH_PROJECTION_USER_ID=999
```
