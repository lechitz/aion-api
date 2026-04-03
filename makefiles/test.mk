# ============================================================
#                         TESTING
# ============================================================

GO_CACHE := $(CURDIR)/.cache/go-build
MCP_SMOKE_USER_ID ?= 999
LOAD_TEST_REQUESTS ?= 60
LOAD_TEST_CONCURRENCY ?= 6
LOAD_TEST_WARMUP ?= 5
LOAD_TEST_TIMEOUT ?= 5s
LOAD_TEST_RECORD_LIMIT ?= 20
LOAD_TEST_DASHBOARD_TIMEZONE ?= UTC
LOAD_TEST_REALTIME_REQUESTS ?= 20
LOAD_TEST_REALTIME_CONCURRENCY ?= 4
LOAD_TEST_REALTIME_WARMUP ?= 2
LOAD_TEST_REALTIME_TIMEOUT ?= 30s

.PHONY: test test-cover test-cover-detail test-html-report test-ci test-clean test-checks mcp-smoke mcp-smoke-readonly record-projection-smoke realtime-record-smoke record-projection-page-smoke ingest-event-smoke outbox-diagnose event-backbone-gate event-backbone-gate-preflight load-test load-test-auth-login load-test-record-projections load-test-dashboard-snapshot load-test-realtime-record-created load-test-baseline

# Execute unit tests
test:
	@echo "Running unit tests with race detector..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go test ./... -v -race

# Run tests with coverage and print package summary plus:
# - total project coverage (from coverprofile)
# - number of source files without tests (by directory heuristic)
# This is intentionally lightweight and meant for quick feedback.
test-cover:
	@echo "Running tests with coverage (summary)..."
	@mkdir -p $(GO_CACHE) $(COVERAGE_DIR)
	@tmpout=$$(mktemp); \
	GOCACHE=$(GO_CACHE) go test ./... -race -coverprofile=$(COVERAGE_DIR)/coverage_summary_tmp.out -count=1 2>&1 | tee $$tmpout; \
	status=$$?; \
	echo ""; \
	echo "--- Coverage summary (package + %) ---"; \
	grep -E "^ok\s+|^\?\s+" $$tmpout || true; \
	rm -f $$tmpout; \
	echo ""; \
	echo "--- Total project coverage (from coverprofile) ---"; \
	if [ -f $(COVERAGE_DIR)/coverage_summary_tmp.out ]; then \
		go tool cover -func=$(COVERAGE_DIR)/coverage_summary_tmp.out | tail -n 1 | awk '{print "TOTAL:", $$3}'; \
		mv -f $(COVERAGE_DIR)/coverage_summary_tmp.out $(COVERAGE_DIR)/coverage.out; \
	else \
		echo "TOTAL: n/a (coverprofile not generated; likely due to packages with no test files)"; \
	fi; \
	echo ""; \
	echo "--- Files without tests (dir heuristic) ---"; \
	echo "Counting .go files (excluding *_test.go, mocks, generated) in directories that have zero *_test.go"; \
	notest_files=$$( \
		find . -type f -name "*.go" \
			! -name "*_test.go" \
			! -path "./tests/mocks/*" \
			! -path "./contracts/openapi/*" \
			! -path "./vendor/*" \
			! -path "./.gomodcache/*" \
			! -name "*.gen.go" \
			! -name "*_gen.go" \
			! -name "mock_*.go" \
			! -name "*_mock.go" \
			-print0 | \
		xargs -0 -I {} sh -c 'd=$$(dirname "{}" ); if ! ls "$$d"/*_test.go >/dev/null 2>&1; then echo "{}"; fi' | wc -l | tr -d ' ' \
	); \
	echo "FILES_WITHOUT_TESTS: $$notest_files"; \
	exit $$status

# Run tests with coverage, filter mocks, and generate HTML coverage report (detailed).
test-cover-detail:
	@echo "Running tests with coverage report..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go test ./... -race -coverprofile=$(COVERAGE_DIR)/coverage_tmp.out -v
	@echo "Filtering out mock files from coverage..."
	grep -v "Mock" $(COVERAGE_DIR)/coverage_tmp.out > $(COVERAGE_DIR)/coverage.out
	@rm -f $(COVERAGE_DIR)/coverage_tmp.out
	@echo "Generating HTML coverage report..."
	GOCACHE=$(GO_CACHE) go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

# Generate JUnit XML test report via gotestsum
test-html-report:
	@echo "Running tests and generating JUnit XML report..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) gotestsum --junitfile $(COVERAGE_DIR)/junit-report.xml -- -race ./...
	@echo "✅ JUnit report generated at $(COVERAGE_DIR)/junit-report.xml"

# CI target: tests with coverage but no HTML UI
test-ci:
	@echo "Running CI tests with coverage output..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go test ./... -race -coverprofile=$(COVERAGE_DIR)/coverage.out -v

# Check tests for common anti-patterns and fail early with actionable message
# Currently checks for uses of context.Background() inside _test.go files.
# Prefer t.Context() (or suite.T().Context()) so test cancellations/timeouts propagate.
test-checks:
	@echo "Checking tests for discouraged patterns..."
	@matches=$$(grep -R --line-number "context.Background()" --include="*_test.go" . \
		--exclude-dir=.git --exclude-dir=vendor --exclude-dir=tests --exclude-dir=.gomodcache --exclude-dir=.cache || true); \
	if [ -n "$$matches" ]; then \
		echo "Found context.Background() usages in tests. Prefer using t.Context() or the suite's context. Matches:"; \
		echo "$$matches"; \
		exit 1; \
	else \
		echo "No discouraged test patterns found."; \
	fi

# Cleanup coverage artifacts and test reports
test-clean:
	@echo "Cleaning up coverage reports and test artifacts..."
	@rm -f \
		$(COVERAGE_DIR)/coverage.out \
		$(COVERAGE_DIR)/coverage_tmp.out \
		$(COVERAGE_DIR)/coverage.html \
		$(COVERAGE_DIR)/junit-report.xml
	@echo "✅ Cleanup complete!"

mcp-smoke:
	@echo "Running MCP smoke test via aion-chat..."
	@if docker ps --filter "name=aion-dev-chat" --filter "status=running" -q | grep -q .; then \
		if docker exec aion-dev-chat test -f /app/scripts/mcp_smoke_test.py >/dev/null 2>&1; then \
			docker exec aion-dev-chat python /app/scripts/mcp_smoke_test.py --user-id $(MCP_SMOKE_USER_ID); \
		else \
			echo "⚠️  MCP smoke script not found inside aion-dev-chat."; \
			echo "   Falling back to host repo execution."; \
			cd ../aion-chat && AION_API_GRAPHQL_URL=http://localhost:5001/aion/api/v1/graphql .venv/bin/python scripts/mcp_smoke_test.py --user-id $(MCP_SMOKE_USER_ID) --env-file infrastructure/docker/environments/dev/.env.dev; \
		fi; \
	else \
		echo "⚠️  aion-dev-chat is not running. Falling back to host repo execution."; \
		cd ../aion-chat && AION_API_GRAPHQL_URL=http://localhost:5001/aion/api/v1/graphql .venv/bin/python scripts/mcp_smoke_test.py --user-id $(MCP_SMOKE_USER_ID) --env-file infrastructure/docker/environments/dev/.env.dev; \
	fi

mcp-smoke-readonly:
	@echo "Running MCP smoke test (read-only) via aion-chat..."
	@if docker ps --filter "name=aion-dev-chat" --filter "status=running" -q | grep -q .; then \
		if docker exec aion-dev-chat test -f /app/scripts/mcp_smoke_test.py >/dev/null 2>&1; then \
			docker exec aion-dev-chat python /app/scripts/mcp_smoke_test.py --read-only --user-id $(MCP_SMOKE_USER_ID); \
		else \
			echo "⚠️  MCP smoke script not found inside aion-dev-chat."; \
			echo "   Falling back to host repo execution."; \
			cd ../aion-chat && AION_API_GRAPHQL_URL=http://localhost:5001/aion/api/v1/graphql .venv/bin/python scripts/mcp_smoke_test.py --read-only --user-id $(MCP_SMOKE_USER_ID) --env-file infrastructure/docker/environments/dev/.env.dev; \
		fi; \
	else \
		echo "⚠️  aion-dev-chat is not running. Falling back to host repo execution."; \
		cd ../aion-chat && AION_API_GRAPHQL_URL=http://localhost:5001/aion/api/v1/graphql .venv/bin/python scripts/mcp_smoke_test.py --read-only --user-id $(MCP_SMOKE_USER_ID) --env-file infrastructure/docker/environments/dev/.env.dev; \
	fi

record-projection-smoke:
	@echo "Running record -> outbox -> kafka -> projection smoke..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/record-projection-smoke

realtime-record-smoke:
	@echo "Running record projection realtime SSE smoke..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/realtime-record-smoke

record-projection-page-smoke:
	@echo "Running derived record projection pagination smoke..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/record-projection-page-smoke

ingest-event-smoke:
	@echo "Running aion-ingest -> kafka smoke..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/ingest-event-smoke

outbox-diagnose:
	@echo "Running outbox diagnose tool..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/outbox-diagnose

load-test:
	@$(MAKE) load-test-baseline

load-test-auth-login:
	@echo "Running versioned load test: auth-login..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/load-test \
		--scenario auth-login \
		--requests $(LOAD_TEST_REQUESTS) \
		--concurrency $(LOAD_TEST_CONCURRENCY) \
		--warmup $(LOAD_TEST_WARMUP) \
		--timeout $(LOAD_TEST_TIMEOUT)

load-test-record-projections:
	@echo "Running versioned load test: record-projections-latest..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/load-test \
		--scenario record-projections-latest \
		--requests $(LOAD_TEST_REQUESTS) \
		--concurrency $(LOAD_TEST_CONCURRENCY) \
		--warmup $(LOAD_TEST_WARMUP) \
		--timeout $(LOAD_TEST_TIMEOUT) \
		--record-limit $(LOAD_TEST_RECORD_LIMIT)

load-test-dashboard-snapshot:
	@echo "Running versioned load test: dashboard-snapshot..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/load-test \
		--scenario dashboard-snapshot \
		--requests $(LOAD_TEST_REQUESTS) \
		--concurrency $(LOAD_TEST_CONCURRENCY) \
		--warmup $(LOAD_TEST_WARMUP) \
		--timeout $(LOAD_TEST_TIMEOUT) \
		--dashboard-timezone $(LOAD_TEST_DASHBOARD_TIMEZONE)

load-test-realtime-record-created:
	@echo "Running versioned load test: realtime-record-created..."
	@mkdir -p $(GO_CACHE)
	GOCACHE=$(GO_CACHE) go run ./hack/tools/load-test \
		--scenario realtime-record-created \
		--requests $(LOAD_TEST_REALTIME_REQUESTS) \
		--concurrency $(LOAD_TEST_REALTIME_CONCURRENCY) \
		--warmup $(LOAD_TEST_REALTIME_WARMUP) \
		--timeout $(LOAD_TEST_REALTIME_TIMEOUT)

load-test-baseline:
	@echo "Running versioned local load baseline..."
	@$(MAKE) load-test-auth-login LOAD_TEST_REQUESTS=60 LOAD_TEST_CONCURRENCY=6 LOAD_TEST_WARMUP=5 LOAD_TEST_TIMEOUT=5s
	@$(MAKE) load-test-record-projections LOAD_TEST_REQUESTS=80 LOAD_TEST_CONCURRENCY=8 LOAD_TEST_WARMUP=5 LOAD_TEST_RECORD_LIMIT=20 LOAD_TEST_TIMEOUT=5s
	@$(MAKE) load-test-dashboard-snapshot LOAD_TEST_REQUESTS=80 LOAD_TEST_CONCURRENCY=8 LOAD_TEST_WARMUP=5 LOAD_TEST_DASHBOARD_TIMEZONE=UTC LOAD_TEST_TIMEOUT=5s
	@$(MAKE) load-test-realtime-record-created LOAD_TEST_REALTIME_REQUESTS=20 LOAD_TEST_REALTIME_CONCURRENCY=4 LOAD_TEST_REALTIME_WARMUP=2 LOAD_TEST_REALTIME_TIMEOUT=30s

event-backbone-gate-preflight:
	@echo "Running event backbone gate preflight..."
	@bash ./hack/dev/event-backbone-gate-preflight.sh

event-backbone-gate:
	@echo "Running v2 event backbone gate..."
	@$(MAKE) event-backbone-gate-preflight
	@$(MAKE) outbox-diagnose
	@$(MAKE) record-projection-smoke
	@$(MAKE) realtime-record-smoke
	@$(MAKE) record-projection-page-smoke
	@$(MAKE) ingest-event-smoke
	@echo "Running dashboard records smoke..."
	@cd ../aion-web && npm run test:e2e:records
