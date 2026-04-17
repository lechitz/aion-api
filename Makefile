# ============================================================
#                   GLOBAL VARIABLES & CONFIG
# ============================================================

APPLICATION_NAME := aion-api

COMPOSE_FILE_DEV  := infrastructure/docker/environments/dev/docker-compose-dev.yaml
ENV_FILE_DEV      := infrastructure/docker/environments/dev/.env.dev
COMPOSE_FILE_MY   := infrastructure/docker/environments/my/docker-compose-my.yaml
ENV_FILE_MY       := infrastructure/docker/environments/my/.env.my
COMPOSE_FILE_PROD := infrastructure/docker/environments/prod/docker-compose-prod.yaml
ENV_FILE_PROD     := infrastructure/docker/environments/prod/.env.prod

MY_POSTGRES_CONTAINER ?= aion-my-postgres
MY_POSTGRES_USER      ?= aion
MY_POSTGRES_DB        ?= aion-api_my
MY_MIGRATION_DB       ?= postgres://aion:aion123@localhost:5432/aion-api_my?sslmode=disable
MY_BACKUP_DIR         ?= ../backups/aion-api/my

COVERAGE_DIR = tests/coverage

# --- MIGRATION CONFIG ---
MIGRATION_PATH := infrastructure/db/migrations
MIGRATION_DB   ?= $(DB_URL)
MIGRATE_BIN    := $(shell command -v migrate 2> /dev/null)

# ============================================================
#                HELP & TOOLING SECTION
# ============================================================

.PHONY: all help tools-install tools.check

all: help

help:
	@echo ""
	@echo ""
	@echo "┃==================================================================================================================┃"
	@echo "┃                                            AION API - CLI COMMANDS                                               ┃"
	@echo "┃==================================================================================================================┃"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ TOOLING ┃"
	@echo ""
	@echo "     tools-install            →  Install all development tools (goimports, golines, gofumpt, golangci-lint)"
	@echo "     graph-projection-export  →  Export graph-projection-v1 JSON (vars: GRAPH_PROJECTION_USER_ID, DATE, CATEGORY_ID, TAG_IDS, OUTPUT)"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ DOCKER ENVIRONMENT COMMANDS ┃"
	@echo ""
	@echo "  - [MY - Personal Environment (Isolated DB)]"
	@echo ""
	@echo "     my                       →  Build + start Personal MY stack"
	@echo "                                  • Separate PostgreSQL database (aion-api_my)"
	@echo "                                  • Separate Redis instance (aion-my-redis)"
	@echo "                                  • Shares Ollama with dev (saves ~5GB+ per model)"
	@echo "                                  • ⚠️  Requires aion-dev-ollama running (make ollama-up)"
	@echo "                                  • 🔥 HOT RELOAD enabled for all projects"
	@echo "                                  • Config in: infrastructure/docker/environments/my/"
	@echo "     my-fast                  →  Start Personal stack without rebuilding"
	@echo "     my-down                  →  Stop Personal services (preserves volumes)"
	@echo "     clean-my                 →  ⚠️ Remove Personal containers/volumes/images (requires CONFIRM_CLEAN_MY=YES)"
	@echo ""
	@echo "  - [DEV - Docker Full Stack]"
	@echo ""
	@echo "     build-dev                →  Build the development images used by the local stack"
	@echo "     dev                      →  Start FULL STACK without rebuilding (default/fast)"
	@echo "     dev-fast                 →  Same as make dev (compatibility alias)"
	@echo "     rebuild-dev              →  Build ALL images + start FULL STACK"
	@echo "                                  Rebuilds: aion-api, aion-chat, aion-web"
	@echo "                                  Preserves: Ollama models + PostgreSQL data (no re-download!)"
	@echo "                                  🔥 HOT RELOAD enabled for all projects:"
	@echo "                                     • Go: Air auto-recompile (~3-5s)"
	@echo "                                     • Python: Uvicorn --reload (~1-2s)"
	@echo "                                     • TypeScript: Vite HMR (<1s)"
	@echo "     dev-down                 →  Stop services (keeps Ollama running, preserves all volumes)"
	@echo "     clean-dev                →  ⚠️ Remove dev containers/volumes/images (keeps Ollama running)"
	@echo ""
	@echo "  - [DEV - Logs (individual services)]"
	@echo ""
	@echo "     logs-api                 →  Show aion-api logs"
	@echo "     logs-api-publisher       →  Show aion-api-outbox-publisher logs"
	@echo "     logs-chat                →  Show aion-chat logs"
	@echo "     logs-ingest              →  Show aion-ingest logs"
	@echo "     logs-streams             →  Show aion-streams logs"
	@echo "     logs-streams-worker      →  Show aion-streams-worker logs"
	@echo "     logs-dashboard           →  Show aion-web logs"
	@echo "     logs-postgres            →  Show PostgreSQL logs"
	@echo "     logs-redis               →  Show Redis logs"
	@echo "     logs-kafka               →  Show Kafka logs"
	@echo "     logs-ollama              →  Show Ollama logs"
	@echo "     logs-jaeger              →  Show Jaeger logs"
	@echo "     logs-otel                →  Show OpenTelemetry Collector logs"
	@echo "     logs-prometheus          →  Show Prometheus logs"
	@echo "     logs-grafana             →  Show Grafana logs"
	@echo "     logs-all                 →  Show all services logs"
	@echo ""
	@echo "  - [DEV - Ollama LLM Management]"
	@echo ""
	@echo "     ollama-up                →  Start Ollama service (if not running)"
	@echo "     ollama-down              →  Stop Ollama (models preserved in volume)"
	@echo "     ollama-status            →  Show Ollama status and installed models"
	@echo "     ollama-restart           →  Restart Ollama service"
	@echo "     ollama-logs              →  View Ollama logs"
	@echo "     ollama-models            →  List installed models"
	@echo "     ollama-pull              →  Download default model (or use MODEL=name)"
	@echo "     ollama-clean             →  ⚠️ Remove Ollama volumes (deletes all models!)"
	@echo ""
	@echo "  - [DEV - Rebuild Individual Services (removes old image before rebuild)]"
	@echo ""
	@echo "     rebuild-dashboard        →  Stop + remove old image + rebuild aion-web"
	@echo "     rebuild-chat             →  Stop + remove old image + rebuild aion-chat"
	@echo "     rebuild-api              →  Stop + remove old image + rebuild aion-api"
	@echo ""
	@echo "  - [DEV - Local (Hot-Reload)]"
	@echo ""
	@echo "     dev-local                →  Run API locally with Air hot-reload (installs Air if needed)"
	@echo "     dev-local-deps           →  Start only dependencies (postgres, redis, etc) for local dev"
	@echo "     dev-local-full           →  Run API locally without Air (go run)"
	@echo "     dev-local-stop           →  Stop local dev dependencies (preserve containers)"
	@echo "     dev-local-down           →  ⚠️ Remove local dev dependencies completely"
	@echo "     air-install              →  Install Air (hot-reload tool)"
	@echo ""
	@echo "  - [PROD]"
	@echo ""
	@echo "     build-prod               →  Build the production Docker image"
	@echo "     prod-up                  →  Start the production environment"
	@echo "     prod-down                →  Stop and remove prod environment containers/volumes"
	@echo "     clean-prod               →  ⚠️ Clean all prod containers, volumes, and images"
	@echo ""
	@echo "  - [CLEANUP & DIAGNOSTICS]"
	@echo ""
	@echo "     cache-reset              →  ⚠️ Flush Redis cache (dev)"
	@echo "     docker-disk              →  Show Docker disk usage (quick diagnostic)"
	@echo "     docker-prune-aion        →  Clean ONLY aion-api images/containers (safe, preserves volumes)"
	@echo "     docker-prune-dangling    →  ⚠️ Remove dangling images (safe)"
	@echo "     docker-prune-build-cache →  ⚠️ Clear Docker build cache"
	@echo "     docker-prune-full        →  Full aion-api cleanup (containers + images + cache, preserves volumes)"
	@echo "     docker-clean-all         →  ⚠️ Remove ALL Docker containers, volumes, and images (DESTRUCTIVE)"
	@echo ""
	@echo "  - [NETWORK/BANDWIDTH OPTIMIZATION]"
	@echo ""
	@echo "     network-audit            →  📊 Audit bandwidth consumption and optimizations"
	@echo "     network-images           →  Show cached Docker images"
	@echo "     images-update            →  🔄 Manually update base images (only when YOU want)"
	@echo "     ollama-update            →  🔄 Update only Ollama image (~800 MB)"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ CODE GENERATION ┃"
	@echo ""
	@echo "     graphql                  →  Generate GraphQL files with gqlgen"
	@echo "     mocks                    →  Generate all GoMock mocks"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ CODE QUALITY ┃"
	@echo ""
	@echo "     format                   →  Format Go code using goimports/golines/gofumpt"
	@echo "     lint                     →  Run golangci-lint (static code analysis)"
	@echo "     lint-fix                 →  Run golangci-lint with --fix (auto-fix where possible)"
	@echo "     verify                   →  Run full pre-commit pipeline (format, mocks, lint, tests, coverage, codegen)"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ MIGRATIONS ┃"
	@echo ""
	@echo "  - [General (requires MIGRATION_DB env)]"
	@echo ""
	@echo "     migrate-install          →  Install golang-migrate CLI"
	@echo "     migrate-up               →  Run all migrations (up)"
	@echo "     migrate-down             →  Rollback the last migration"
	@echo "     migrate-force VERSION=X  →  Force DB to specific version"
	@echo "     migrate-new              →  Create new migration (with prompt)"
	@echo ""
	@echo "  - [DEV Environment (uses localhost:5432)]"
	@echo ""
	@echo "     migrate-dev-up           →  Apply all migrations to dev DB"
	@echo "     migrate-dev-down         →  Rollback last migration on dev DB"
	@echo "     migrate-dev-status       →  Show current migration version"
	@echo "     migrate-dev-reset        →  ⚠️ Drop all and re-apply migrations"
	@echo "     migrate-my-up            →  Apply migrations to personal MY DB without resetting data"
	@echo "     migrate-my-status        →  Show current migration version for personal MY DB"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ BACKUP / RESTORE ┃"
	@echo ""
	@echo "     backup-dev               →  Create a local custom-format backup in ../backups/aion-api/dev/"
	@echo "     backup-my                →  Create a local custom-format backup in ../backups/aion-api/my/"
	@echo "     restore-dev              →  Restore a backup with CONFIRM_RESTORE=YES BACKUP_FILE=..."
	@echo "     restore-my               →  Restore MY backup with CONFIRM_RESTORE=YES BACKUP_FILE=..."
	@echo ""
	@echo ""
	@echo " 🔶 ┃ SEEDS ┃"
	@echo ""
	@echo "  - [Quick Start]"
	@echo ""
	@echo "     db-full                  →  🚀 ONE COMMAND: reset DB + migrations + essential + realistic 3-month profile"
	@echo "     db-reset                 →  ⚠️ Reset DB + re-apply migrations"
	@echo ""
	@echo "  - [Essential Data]"
	@echo ""
	@echo "     seed-essential           →  Seed roles + admin user only"
	@echo "     seed-roles               →  Seed system roles (owner, admin, user, blocked)"
	@echo "     seed-admin               →  Seed admin user (username: aion)"
	@echo ""
	@echo "  - [Test Data]"
	@echo ""
	@echo "     seed-test                →  🧪 Realistic profile: legacy+new taxonomy + metrics/goals + ~3 months (50-60/day)"
	@echo "     seed-clean-test          →  Remove ONLY realistic test profile (testuser), keeps admin"
	@echo "     hash-gen PASS='pwd'      →  🔐 Generate bcrypt hash for password (use to create seeds)"
	@echo ""
	@echo "  - [Cleanup & Reset]"
	@echo ""
	@echo "     reset-user-data          →  🔄 Delete ALL users + data (keeps roles) ⚠️ Removes admin too!"
	@echo "     seed-clean-all           →  Truncate seeded tables (dev only)"
	@echo ""
	@echo "  - [Legacy/Alternative Seeds]"
	@echo ""
	@echo "     seed-users               →  Seed the users table"
	@echo "     seed-categories          →  Seed the categories table"
	@echo "     seed-all                 →  Run all seed scripts"
	@echo "     seed-user1-all           →  Seed full dataset for default user (id=1)"
	@echo "     seed-everybody           →  Alias for seed-all"
	@echo "     seed-api-caller          →  Gera dados via chamadas HTTP/GraphQL (modo estrito, sem criar usuário)"
	@echo "     seed-api-caller-bootstrap  →  Gera dados via API e cria usuário se login falhar"
	@echo "     seed-api-caller-clean     →  Limpa registros via API e roda modo estrito"
	@echo "     seed-caller              →  Gera via API para N usuários (cria se faltar) - use N=9 ou n=9"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ TESTING ┃"
	@echo ""
	@echo "     test                     →  Run unit tests"
	@echo "     test-cover               →  Run tests with coverage summary (package + %)"
	@echo "     test-cover-detail        →  Run tests with coverage report (excludes mocks)"
	@echo "     test-html-report         →  Generate HTML test report (requires go-test-html-report)"
	@echo "     regression-gate-draft    →  Run cross-repo draft-flow regression gate (dashboard + chat + api)"
	@echo "     e2e-draft-smoke          →  Run dashboard Playwright smoke for DC-08/DC-09 (host deps required)"
	@echo "     dc15-correlate           →  Correlate UI action logs across dashboard + api + chat (vars: SINCE, DRAFT_ID)"
	@echo "     mcp-smoke                →  Run MCP smoke test through aion-chat against the current local stack"
	@echo "     mcp-smoke-readonly       →  Run MCP smoke test in read-only mode"
	@echo "     record-projection-smoke  →  Run cross-repo smoke for record -> outbox -> kafka -> projection"
	@echo "     realtime-record-smoke    →  Run SSE smoke for record projection readiness notifications"
	@echo "     record-projection-page-smoke  →  Run cursor pagination smoke for derived record projections"
	@echo "     ingest-event-smoke       →  Run cross-repo smoke for aion-ingest -> kafka envelope publication"
	@echo "     outbox-diagnose          →  Inspect outbox backlog and sample pending/failed rows"
	@echo "     load-test-auth-login     →  Run versioned load scenario for POST /auth/login"
	@echo "     load-test-record-projections  →  Run versioned load scenario for recordProjectionsLatest"
	@echo "     load-test-dashboard-snapshot  →  Run versioned load scenario for dashboardSnapshot"
	@echo "     load-test-realtime-record-created  →  Run versioned load scenario for realtime SSE delivery"
	@echo "     load-test-baseline       →  Run the committed local load baseline across auth + GraphQL + realtime"
	@echo "     load-test                →  Alias for load-test-baseline"
	@echo "     event-backbone-gate-preflight  →  Check local stack and repo prerequisites for the v2 gate"
	@echo "     event-backbone-gate      →  Run the full v2 records gate across api, streams, ingest, and dashboard"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ API DOCS (SWAGGER) ┃"
	@echo ""
	@echo "     docs.gen                 →  Generate Swagger artifacts (docs.go, swagger.json/yaml)"
	@echo "     docs.check-dirty         →  Fail if Swagger artifacts are out-of-date"
	@echo "     docs.clean               →  Remove generated Swagger artifacts"
	@echo ""
	@echo ""
	@echo " 🔶 ┃ API CALLS ┃"
	@echo ""
	@echo "     call-login               →  POST /auth/login (vars: USER, PASS, SAVE_TOKEN=true to cache)"
	@echo "     call-health              →  GET /aion/health (also available at /aion/api/v1/health)"
	@echo "     call-me                  →  GET /user/me (needs TOKEN or .cache/aion/token)"
	@echo "     call-chat                →  POST /chat (needs MESSAGE + token)"
	@echo "     call-graphql             →  POST /graphql (vars: QUERY or QUERY_FILE)"
	@echo ""
	@echo ""
	@echo "┃==================================================================================================================┃"
	@echo ""

# ============================================================
#                 CONSOLIDATED .PHONY TARGETS
# ============================================================

-include makefiles/*.mk

.PHONY: graphql mocks docs.gen docs.validate docs.check-dirty docs-verify lint test test-cover test-cover-detail test-ci test-clean

# Short aliases
.PHONY: install-tools
install-tools: tools-install

.PHONY: \
	help tools-install tools.check \
	graph-projection-export \
	build-dev rebuild-dev dev-full dev-up dev-down dev dev-fast dev-attach dev-logs dev-clean clean-dev \
	logs-api logs-chat logs-dashboard logs-postgres logs-redis logs-ollama logs-jaeger logs-otel logs-prometheus logs-grafana logs-all \
	ollama-up ollama-down ollama-status ollama-restart ollama-logs ollama-models ollama-pull ollama-clean \
	dev-local dev-local-deps dev-local-full dev-local-stop dev-local-down air-install \
	build-prod prod-up prod-down prod clean-prod \
	docker-clean-all docker-disk docker-prune-aion docker-prune-dangling docker-prune-build-cache docker-prune-full \
	cache-reset \
	graphql mocks \
	format lint lint-fix verify \
	test test-cover test-cover-detail test-html-report test-ci test-clean \
	regression-gate-draft e2e-draft-smoke dc15-correlate \
	mcp-smoke mcp-smoke-readonly \
	migrate-up migrate-down migrate-force migrate-new migrate-install \
	migrate-dev-up migrate-dev-down migrate-dev-status migrate-dev-reset migrate-my-up migrate-my-status \
	backup-dev restore-dev backup-my restore-my \
	docs.gen docs.check-dirty docs.clean docs.validate docs-verify

docs-serve:
	@.venv-docs/bin/python -m mkdocs serve

docs-build:
	@.venv-docs/bin/python -m mkdocs build

docs-verify:
	@.venv-docs/bin/python -m mkdocs build --strict

regression-gate-draft:
	@./hack/regression-gate-draft-flow.sh

e2e-draft-smoke:
	@cd ../aion-web && npm run test:e2e:draft

dc15-correlate:
	@./hack/dc15-correlate-ui-action.sh --since "$${SINCE:-45m}" $${DRAFT_ID:+--draft-id "$$DRAFT_ID"}

# Include debug makefile (opt-in to avoid overriding targets and noisy warnings)
# Usage: make INCLUDE_DEBUG_MK=1 debug-roles
ifeq ($(INCLUDE_DEBUG_MK),1)
-include makefiles/debug.mk
endif
