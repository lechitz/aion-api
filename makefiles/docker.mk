# ============================================================
#                DOCKER ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-dev dev-up dev-down dev dev-fast rebuild-dev dev-full dev-clean clean-dev
.PHONY: build-my my-up my-down my my-fast my-clean clean-my
.PHONY: rebuild-dashboard rebuild-chat rebuild-api
.PHONY: build-prod prod-up prod-down prod clean-prod
.PHONY: docker-clean-all

APPLICATION_NAME := aion-api

# ============================================================
#         REBUILD INDIVIDUAL SERVICES (without full restart)
# ============================================================

rebuild-dashboard:
	@echo "Stopping and removing old aion-web container..."
	@docker stop aion-dev-web 2>/dev/null || true
	@docker rm aion-dev-web 2>/dev/null || true
	@echo "Removing old aion-web image..."
	@docker rmi aion-web:dev 2>/dev/null || true
	@echo "Rebuilding aion-web..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) build aion-web
	@echo "Auto-cleanup: removing dangling images..."
	@docker image prune -f > /dev/null 2>&1 || true
	@echo "Starting aion-web..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d aion-web
	@echo "✅ aion-web rebuilt and restarted!"
	@echo "   → http://localhost:5000"

rebuild-chat:
	@echo "Stopping and removing old aion-chat container..."
	@docker stop aion-dev-chat 2>/dev/null || true
	@docker rm aion-dev-chat 2>/dev/null || true
	@echo "Removing old aion-chat image..."
	@docker rmi aion-chat:dev 2>/dev/null || true
	@echo "Rebuilding aion-chat..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) build aion-chat
	@echo "Auto-cleanup: removing dangling images..."
	@docker image prune -f > /dev/null 2>&1 || true
	@echo "Starting aion-chat..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d aion-chat
	@echo "✅ aion-chat rebuilt and restarted!"
	@echo "   → http://localhost:8000/health"

rebuild-api:
	@echo "Stopping and removing old aion-api container..."
	@docker stop aion-dev-api 2>/dev/null || true
	@docker rm aion-dev-api 2>/dev/null || true
	@echo "Removing old aion-api image..."
	@docker rmi $(APPLICATION_NAME):dev 2>/dev/null || true
	@echo "Rebuilding aion-api..."
	@DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="" -f infrastructure/docker/environments/dev/Dockerfile.dev -t $(APPLICATION_NAME):dev .
	@echo "Auto-cleanup: removing dangling images..."
	@docker image prune -f > /dev/null 2>&1 || true
	@echo "Starting aion-api..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d aion-api
	@echo "✅ aion-api rebuilt and restarted!"
	@echo "   → http://localhost:5001/aion/api/v1/health"

build-dev:
	@echo "[BUILD-DEV] Building DEV images used by the local compose stack..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) build \
			aion-api \
			aion-api-outbox-publisher \
			aion-ingest \
			aion-streams \
			aion-streams-worker

dev-up: dev-down
	@echo "[DEV-UP] Starting DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-down:
	@echo "[DEV-DOWN] Stopping DEV environment (preserving volumes)..."
	@echo "      Ollama will be kept RUNNING (models preserved)"
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) stop aion-api aion-api-outbox-publisher aion-chat aion-ingest aion-streams aion-streams-worker aion-web postgres redis kafka localstack jaeger otel-collector prometheus grafana loki fluent-bit 2>/dev/null || true
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) rm -f aion-api aion-api-outbox-publisher aion-chat aion-ingest aion-streams aion-streams-worker aion-web postgres redis kafka localstack jaeger otel-collector prometheus grafana loki fluent-bit 2>/dev/null || true
	@echo ""
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "✅ Services stopped (Ollama still running)"; \
		echo ""; \
		echo "💡 Ollama is still RUNNING to preserve models"; \
		echo "   • To stop Ollama: make ollama-down"; \
		echo "   • To check status: make ollama-status"; \
	else \
		echo "✅ Services stopped"; \
		echo ""; \
		echo "ℹ️  Ollama is not running"; \
		echo "   • To start Ollama: make ollama-up"; \
	fi

rebuild-dev: build-dev
	@echo "[REBUILD-DEV] Building images + starting FULL STACK (detached)..."
	@echo "      → aion-api + aion-chat + dashboard + infrastructure"
	@echo "      ℹ️  Volumes preserved (Ollama models + Database)"
	@echo "      💡 Use 'make dev' for fast startup without forced rebuild"
	@echo ""
	@echo "Starting/restarting services (preserving volumes)..."
	@if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "Starting Ollama..."; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d ollama; \
		sleep 2; \
	fi
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d --build --no-recreate 2>/dev/null || \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d --build
	@echo ""
	@echo "⏳ Waiting for Database to be ready..."
	@for i in $$(seq 1 30); do \
		if docker exec aion-dev-postgres pg_isready -U aion -d aion-api >/dev/null 2>&1; then \
			echo "✅ Database is ready!"; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "⚠️  Timeout waiting for Database"; \
		fi; \
		sleep 1; \
	done
	@echo ""
	@echo "🗄️  Running database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path infrastructure/db/migrations -database "postgres://aion:aion123@localhost:5432/aion-api?sslmode=disable" up && \
		echo "✅ Migrations applied successfully"; \
	else \
		echo "⚠️  'migrate' CLI not found. Run: make migrate-install"; \
		echo "   Then run: make migrate-dev-up"; \
	fi
	@echo ""
	@echo "🔧 Checking Ollama model setup..."
	@bash hack/dev/check-and-setup-ollama.sh || echo "⚠️  Ollama setup had issues but continuing..."
	@echo ""
	@echo "⏳ Waiting for services to be healthy..."
	@for i in $$(seq 1 60); do \
		if docker compose -f $(COMPOSE_FILE_DEV) ps | grep -q "aion-chat.*healthy"; then \
			echo "✅ aion-chat is healthy!"; \
			break; \
		fi; \
		if [ $$i -eq 60 ]; then \
			echo "⚠️  Timeout waiting for aion-chat health check"; \
		fi; \
		sleep 1; \
	done
	@echo ""
	@echo "✅ All services started in background"
	@echo ""
	@echo "📍 Service URLs:"
	@echo "   • Dashboard:  http://localhost:5000"
	@echo "   • API:        http://localhost:5001/aion/api/v1/health"
	@echo "   • Chat AI:    http://localhost:8000/health"
	@echo "   • Grafana:    http://localhost:3000"
	@echo "   • Logs (Loki): http://localhost:3000/explore"
	@echo "   • Jaeger:     http://localhost:16686"
	@echo "   • Loki:       http://localhost:3100"
	@echo ""
	@echo "Quick commands:"
	@echo "   make dev                   → Start without rebuilding ANY images (default/fast)"
	@echo "   make dev-fast              → Same as make dev (kept for compatibility)"
	@echo "   make rebuild-dev           → Build all images + start full stack"
	@echo "   make dev-down              → Stop services (Ollama stays running)"
	@echo "   make rebuild-api           → Rebuild only aion-api (smart rebuild)"
	@echo "   make dev-attach            → Attach to aion-api logs"
	@echo "   make ollama-status         → Check Ollama status and models"
	@echo "   make ollama-down           → Stop Ollama (models preserved)"
	@echo "   make docker-prune-dangling → Clean temp images (run weekly)"

dev-fast:
	@echo "[DEV-FAST] Starting services WITHOUT rebuilding images..."
	@echo "      ⚡ Use this when you haven't changed ANY code (fastest option)"
	@echo ""
	@if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "Starting Ollama..."; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d ollama; \
		sleep 2; \
	fi
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d
	@echo ""
	@echo "Waiting for Database..."
	@for i in $$(seq 1 20); do \
		if docker exec aion-dev-postgres pg_isready -U aion -d aion-api >/dev/null 2>&1; then \
			break; \
		fi; \
		sleep 1; \
	done
	@echo " Applying migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path infrastructure/db/migrations -database "postgres://aion:aion123@localhost:5432/aion-api?sslmode=disable" up 2>&1 | grep -v "no change" || true; \
	fi
	@echo ""
	@echo " Services started (using existing images)"
	@echo ""
	@echo "📍 Service URLs:"
	@echo "   • Dashboard:  http://localhost:5000"
	@echo "   • API:        http://localhost:5001/aion/api/v1/health"
	@echo "   • Chat AI:    http://localhost:8000/health"
	@echo "   • Grafana:    http://localhost:3000"
	@echo "   • Logs (Loki): http://localhost:3000/explore"
	@echo "   • Jaeger:     http://localhost:16686"
	@echo "   • Loki:       http://localhost:3100"

# New default behavior: fast startup without rebuild.
dev: dev-fast

# Alias for explicit naming.
dev-full: rebuild-dev

# ============================================================
#                     LOGS COMMANDS
# ============================================================

.PHONY: logs-api logs-api-publisher logs-chat logs-ingest logs-streams logs-streams-worker logs-dashboard logs-postgres logs-redis logs-kafka logs-ollama logs-jaeger logs-otel logs-prometheus logs-grafana logs-all

logs-api:
	@echo "📋 aion-api logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-api

logs-api-publisher:
	@echo "📋 aion-api-outbox-publisher logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-api-outbox-publisher

logs-chat:
	@echo "📋 aion-chat logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-chat

logs-ingest:
	@echo "📋 aion-ingest logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-ingest

logs-streams:
	@echo "📋 aion-streams logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-streams

logs-streams-worker:
	@echo "📋 aion-streams-worker logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-streams-worker

logs-dashboard:
	@echo "📋 aion-web logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-web

logs-postgres:
	@echo "📋 PostgreSQL logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f postgres

logs-redis:
	@echo "📋 Redis logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f redis

logs-kafka:
	@echo "📋 Kafka logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f kafka

logs-ollama:
	@echo "📋 Ollama logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f ollama

logs-jaeger:
	@echo "📋 Jaeger logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f jaeger

logs-otel:
	@echo "📋 OpenTelemetry Collector logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f otel-collector

logs-prometheus:
	@echo "📋 Prometheus logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f prometheus

logs-grafana:
	@echo "📋 Grafana logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Container still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f grafana

logs-all:
	@echo "📋 All services logs (Ctrl+C to exit)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Containers still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f

# Backwards compatibility aliases
dev-attach: logs-api
dev-logs: logs-all

# Alias for backwards compatibility (use clean-dev instead)
dev-clean: clean-dev

clean-dev:
	@echo "[CLEAN-DEV] Cleaning DEV containers, volumes, images..."
	@echo "      ⚠️  This will remove:"
	@echo "         • PostgreSQL data"
	@echo "         • Redis cache"
	@echo "         • Localstack assets"
	@echo "         • aion-api:dev image"
	@echo "         • aion-chat:dev image"
	@echo "         • aion-web:dev image"
	@echo ""
	@echo "      ℹ️  Ollama will be kept RUNNING (models preserved)"
	@echo ""
	@echo "→ Stopping and removing services (except Ollama)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) stop aion-api aion-api-outbox-publisher aion-chat aion-ingest aion-streams aion-streams-worker aion-web postgres redis kafka localstack jaeger otel-collector prometheus grafana loki fluent-bit 2>/dev/null || true
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) rm -f -v aion-api aion-api-outbox-publisher aion-chat aion-ingest aion-streams aion-streams-worker aion-web postgres redis kafka localstack jaeger otel-collector prometheus grafana loki fluent-bit 2>/dev/null || true
	@echo "→ Removing dev images..."
	@docker images --filter "reference=$(APPLICATION_NAME):dev" -q | xargs -r docker rmi -f || true
	@docker images --filter "reference=aion-chat:dev" -q | xargs -r docker rmi -f || true
	@docker images --filter "reference=aion-web:dev" -q | xargs -r docker rmi -f || true
	@echo "→ Removing dev volumes (Postgres/Redis/Observability)..."
	@docker volume ls -q | grep -E '(^|[-_])postgres-data-dev$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])grafana-data$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])loki-data-dev$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])fluentbit-data$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])localstack-data$$' | xargs -r docker volume rm -f || true
	@echo ""
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "✅ Cleanup complete (Ollama still running)"; \
		echo ""; \
		echo "💡 Ollama is still RUNNING to preserve models (~4-5GB)"; \
		echo "   To remove Ollama: make ollama-clean"; \
	else \
		echo "✅ Cleanup complete"; \
	fi
	@echo ""

# ============================================================
#               PERSONAL/MY ENVIRONMENT (Isolated DB)
# ============================================================

build-my:
	@echo "[BUILD-MY] Building MY image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):my .

my-up: my-down
	@echo "[MY-UP] Starting MY personal environment..."
	export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_MY) rm -f -v postgres-my
	export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_MY) up

my-down:
	@echo "[MY-DOWN] Stopping MY environment (preserving volumes)..."
	@echo "      ℹ️  Ollama will be kept RUNNING (models preserved)"
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_MY) stop aion-api aion-chat aion-web postgres-my redis-aion-my localstack-my jaeger otel-collector prometheus-my grafana-my loki fluent-bit 2>/dev/null || true
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_MY) rm -f aion-api aion-chat aion-web postgres-my redis-aion-my localstack-my jaeger otel-collector prometheus-my grafana-my loki fluent-bit 2>/dev/null || true
	@echo ""
	@if docker ps --filter "name=ollama-my" --filter "status=running" -q | grep -q .; then \
		echo "✅ Services stopped (Ollama still running)"; \
		echo ""; \
		echo "💡 Ollama is still RUNNING to preserve models"; \
		echo "   • To stop Ollama: make ollama-down"; \
		echo "   • To check status: make ollama-status"; \
	else \
		echo "✅ Services stopped"; \
		echo ""; \
	fi

my:
	@echo "=================================================="
	@echo "  BUILDING + STARTING PERSONAL ENVIRONMENT (MY)"
	@echo "=================================================="
	@echo ""
	@echo "This will:"
	@echo "  ✓ Build aion-api:my, aion-chat:dev, aion-web:dev"
	@echo "  ✓ Start ALL services (postgres-my, redis-my, etc.)"
	@echo "  ✓ Use SEPARATE database (aion-api_my)"
	@echo "  ✓ Share Ollama with dev environment (saves resources)"
	@echo "  ✓ Enable hot-reload for all projects"
	@echo ""
	@echo "⚠️  Prerequisites:"
	@echo "   • Ollama-dev must be running (shared between dev/my)"
	@echo "   • Run 'make ollama-up' if not running"
	@echo ""
	@if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "❌ ERROR: aion-dev-ollama is not running!"; \
		echo ""; \
		echo "Start it with: make ollama-up"; \
		echo ""; \
		exit 1; \
	fi
	@echo "✅ aion-dev-ollama is running"
	@echo ""
	@echo "⏳ Building aion-api:my..."
	@DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):my . || { echo "❌ Build failed"; exit 1; }
	@echo ""
	@echo "⏳ Starting services..."
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_MY) build
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_MY) up -d
	@echo ""
	@echo "✅ Personal environment started!"
	@echo ""
	@echo "📍 Services:"
	@echo "   • aion-api:       http://localhost:5001/aion/api/v1/health"
	@echo "   • Dashboard:     http://localhost:5000"
	@echo "   • Aion-Chat:     http://localhost:8000/health"
	@echo "   • GraphQL:       http://localhost:5001/aion/graphql"
	@echo "   • PostgreSQL:    localhost:5432 (DB: aion-api_my)"
	@echo "   • Redis:         localhost:6379"
	@echo "   • Ollama:        http://localhost:11434 (shared with dev)"
	@echo "   • Jaeger UI:     http://localhost:16686"
	@echo "   • Prometheus:    http://localhost:9090"
	@echo "   • Grafana:       http://localhost:3001"
	@echo ""
	@echo "📋 Logs: make logs-all"
	@echo "⏹️  Stop:  make my-down"
	@echo ""

my-fast:
	@echo "=================================================="
	@echo "🚀  STARTING PERSONAL ENVIRONMENT (MY) - NO REBUILD"
	@echo "=================================================="
	@echo ""
	@echo "⚠️  Prerequisites:"
	@echo "   • Ollama-dev must be running (shared between dev/my)"
	@echo ""
	@if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo "❌ ERROR: aion-dev-ollama is not running!"; \
		echo ""; \
		echo "Start it with: make ollama-up"; \
		echo ""; \
		exit 1; \
	fi
	@echo "✅ aion-dev-ollama is running"
	@echo ""
	@echo "⏳ Starting services (using existing images)..."
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_MY) up -d
	@echo ""
	@echo "✅ Personal environment started!"
	@echo ""
	@echo "📍 Services:"
	@echo "   • aion-api:       http://localhost:5001/aion/api/v1/health"
	@echo "   • Dashboard:     http://localhost:5000"
	@echo "   • PostgreSQL:    localhost:5432 (DB: aion-api_my)"
	@echo ""
	@echo "📋 Logs: make logs-all"
	@echo "⏹️  Stop:  make my-down"
	@echo ""

# Alias for backwards compatibility
my-clean: clean-my

clean-my:
	@echo "[CLEAN-MY] Cleaning MY containers, volumes, images..."
	@echo "      ⚠️  This will remove:"
	@echo "         • PostgreSQL data (aion-api_my)"
	@echo "         • Redis cache (redis-my)"
	@echo "         • Localstack assets"
	@echo "         • aion-api:my image"
	@echo "         • aion-chat:dev image"
	@echo "         • aion-web:dev image"
	@echo ""
	@echo "      ℹ️  Ollama will be kept RUNNING (models preserved)"
	@echo ""
	@echo "→ Stopping and removing services (except Ollama)..."
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_MY) stop aion-api aion-chat aion-web postgres-my redis-aion-my localstack-my jaeger otel-collector prometheus-my grafana-my loki fluent-bit 2>/dev/null || true
	@export $$(cat $(ENV_FILE_MY) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_MY) rm -f -v aion-api aion-chat aion-web postgres-my redis-aion-my localstack-my jaeger otel-collector prometheus-my grafana-my loki fluent-bit 2>/dev/null || true
	@echo "→ Removing my images..."
	@docker images --filter "reference=$(APPLICATION_NAME):my" -q | xargs -r docker rmi -f || true
	@docker images --filter "reference=aion-chat:dev" -q | xargs -r docker rmi -f || true
	@docker images --filter "reference=aion-web:dev" -q | xargs -r docker rmi -f || true
	@echo "→ Removing my volumes (Postgres/Redis/Observability)..."
	@docker volume ls -q | grep -E '(^|[-_])postgres-data-my$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])grafana-data$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])loki-data-my$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])fluentbit-data$$' | xargs -r docker volume rm -f || true
	@docker volume ls -q | grep -E '(^|[-_])localstack-data$$' | xargs -r docker volume rm -f || true
	@echo ""
	@if docker ps --filter "name=ollama-my" --filter "status=running" -q | grep -q .; then \
		echo "✅ Cleanup complete (Ollama still running)"; \
		echo ""; \
		echo "💡 Ollama is still RUNNING to preserve models (~4-5GB)"; \
		echo "   To remove Ollama: make ollama-clean"; \
	else \
		echo "✅ Cleanup complete"; \
	fi
	@echo ""

# ============================================================
#                PRODUCTION BUILD
# ============================================================


build-prod:
	@echo "[BUILD-PROD] Building PROD image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="-s -w" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):prod .

prod-up: prod-down
	@echo "[PROD-UP] Starting PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) up

prod-down:
	@echo "[PROD-DOWN] Stopping PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) down -v

prod: build-prod prod-up

clean-prod:
	@echo "[CLEAN-PROD] Cleaning PROD containers, volumes, images..."
	@docker ps -a --filter "name=prod" -q | xargs -r docker rm -f
	@docker volume ls --filter "name=prod" -q | xargs -r docker volume rm
	@docker images --filter "reference=*prod*" -q | xargs -r docker rmi -f

docker-clean-all:
	@echo "[CLEAN-ALL] Removing ALL containers, volumes, images..."
	@docker ps -a -q | xargs -r docker rm -f
	@docker volume ls -q | xargs -r docker volume rm
	@docker images -a -q | xargs -r docker rmi -f

# ============================================================
#         DOCKER CLEANUP & DIAGNOSTICS (Disk Management)
# ============================================================

.PHONY: docker-disk docker-prune-aion docker-prune-dangling docker-prune-build-cache cache-reset

cache-reset:
	@echo "🧹 Flushing Redis cache (dev)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) exec -T redis redis-cli FLUSHALL
	@echo "✅ Redis cache cleared"

# Show Docker disk usage (quick diagnostic)
docker-disk:
	@echo "Docker Disk Usage Summary"
	@echo "=============================="
	@docker system df
	@echo ""
	@echo "Aion-related images:"
	@docker images | grep -E "(aion|dashboard)" || echo "   (none found)"
	@echo ""
	@echo "Aion-related containers (all states):"
	@docker ps -a | grep -E "(aion|dashboard)" || echo "   (none found)"
	@echo ""
	@echo "Tips:"
	@echo "   make docker-prune-aion       → Clean only aion-api images/containers"
	@echo "   make docker-prune-dangling   → Remove dangling images (safe)"
	@echo "   make docker-prune-build-cache → Clear Docker build cache"

# Clean ONLY aion-api-related images and containers (safe, preserves volumes)
docker-prune-aion:
	@echo "Cleaning aion-api-related images and containers..."
	@echo "   ⚠️  This will NOT delete volumes (PostgreSQL data, Ollama models)"
	@echo ""
	@echo "→ Stopping Aion containers..."
	@docker stop aion-dev-api aion-dev-chat aion-dev-web 2>/dev/null || true
	@echo "→ Removing Aion containers..."
	@docker rm aion-dev-api aion-dev-chat aion-dev-web 2>/dev/null || true
	@echo "→ Removing Aion images..."
	@docker rmi aion-api:dev aion-chat:dev aion-web:dev 2>/dev/null || true
	@echo ""
	@echo "→ Removing dangling images from Aion builds..."
	@docker images --filter "dangling=true" -q | xargs -r docker rmi 2>/dev/null || true
	@echo ""
	@echo "✅ aion-api cleanup complete!"
	@echo "   Next 'make dev' will rebuild images from scratch (using cache)"

# Remove dangling images only (safe, no data loss)
docker-prune-dangling:
	@echo "Removing dangling images..."
	@docker image prune -f
	@echo ""
	@echo "✅ Dangling images removed!"

# Clear Docker build cache (frees significant space after many rebuilds)
docker-prune-build-cache:
	@echo "Clearing Docker build cache..."
	@echo "   ⚠️  Next builds will be slower (no cache)"
	@docker builder prune -f
	@echo ""
	@echo "✅ Build cache cleared!"

# Full prune for aion-api stack (containers + images + build cache, preserves volumes)
docker-prune-full:
	@echo "Full aion-api Docker cleanup..."
	@echo "   ⚠️  This will remove:"
	@echo "      • All Aion containers"
	@echo "      • All Aion images"
	@echo "      • Docker build cache"
	@echo "   ✅ This will PRESERVE:"
	@echo "      • PostgreSQL data"
	@echo "      • Ollama models"
	@echo ""
	@read -p "Continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(MAKE) docker-prune-aion; \
		$(MAKE) docker-prune-build-cache; \
		$(MAKE) docker-prune-dangling; \
		echo ""; \
		echo "📊 Space after cleanup:"; \
		docker system df; \
	else \
		echo "❌ Cancelled"; \
	fi

# ============================================================
#               NETWORK/BANDWIDTH DIAGNOSTICS
# ============================================================

.PHONY: network-audit network-images ollama-update images-update

# Audit Docker images and estimate potential bandwidth consumption
network-audit:
	@echo ""
	@echo "┃═══════════════════════════════════════════════════════════════┃"
	@echo "┃              AION STACK - NETWORK/BANDWIDTH AUDIT             ┃"
	@echo "┃═══════════════════════════════════════════════════════════════┃"
	@echo ""
	@echo "📊 DOCKER IMAGES (cached locally):"
	@echo "──────────────────────────────────────────────────────────────────"
	@docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep -E "(ollama|aion|postgres|redis|node|python|grafana|prometheus|jaeger|loki|fluent|otel|localstack)" || echo "   (none found)"
	@echo ""
	@echo "💾 VOLUMES (persistent data):"
	@echo "──────────────────────────────────────────────────────────────────"
	@docker volume ls --format "{{.Name}}" | grep -E "(ollama|postgres|grafana|loki|redis)" | while read vol; do \
		size=$$(docker system df -v 2>/dev/null | grep "$$vol" | awk '{print $$3}' || echo "?"); \
		echo "   $$vol: $$size"; \
	done
	@echo ""
	@echo "🌐 ESTIMATED BANDWIDTH PER COMMAND (if cache miss):"
	@echo "──────────────────────────────────────────────────────────────────"
	@echo "   make dev (full rebuild):"
	@echo "      • Ollama image:         ~800 MB"
	@echo "      • aion-chat deps:       ~1.5 GB (PyTorch, LangChain, Whisper)"
	@echo "      • aion-api deps:        ~200 MB (Go modules)"
	@echo "      • dashboard deps:       ~300 MB (node_modules)"
	@echo "      • Infra images:         ~500 MB (Postgres, Redis, Grafana, etc)"
	@echo "      ────────────────────────────────"
	@echo "      TOTAL (worst case):     ~3.3 GB"
	@echo ""
	@echo "   make dev-fast (no rebuild):"
	@echo "      • Downloads:            0 MB (uses cached images)"
	@echo ""
	@echo "   make ollama-pull (new model):"
	@echo "      • Qwen 7B:              ~4.5 GB"
	@echo "      • Qwen 14B:             ~8.5 GB"
	@echo "      • Llama 3.1 8B:         ~4.7 GB"
	@echo ""
	@echo "💡 OTIMIZAÇÕES APLICADAS:"
	@echo "──────────────────────────────────────────────────────────────────"
	@if grep -q "pull_policy: missing" $(COMPOSE_FILE_DEV) 2>/dev/null; then \
		echo "   ✅ Ollama: pull_policy: missing (não baixa a cada make dev)"; \
	else \
		echo "   ❌ Ollama: pull_policy: always (baixa ~800MB a cada make dev!)"; \
	fi
	@echo "   ✅ Volumes persistentes: ollama-models, postgres-data-dev"
	@echo "   ✅ Go module cache: go-mod-cache volume"
	@echo "   ✅ Dashboard node_modules: volume separado"
	@echo ""
	@echo "🔧 COMANDOS ÚTEIS:"
	@echo "──────────────────────────────────────────────────────────────────"
	@echo "   make dev-fast          → Inicia SEM rebuild (0 MB download)"
	@echo "   make ollama-status     → Verifica modelos já baixados"
	@echo "   make images-update     → Atualiza imagens base manualmente"
	@echo "   make docker-disk       → Mostra uso de disco do Docker"
	@echo ""

# Show current Docker images that could be updated
network-images:
	@echo "📦 Imagens Docker do stack Aion:"
	@echo ""
	@docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}" | grep -E "(ollama|aion|postgres|redis|node|python|grafana|prometheus|jaeger|loki|fluent|otel|localstack|REPOSITORY)" || echo "(none)"

# Manually update base images (only when YOU want, not every make dev)
images-update:
	@echo "🔄 Atualizando imagens base (isso usa banda!)..."
	@echo ""
	@echo "Isso vai baixar:"
	@echo "   • ollama/ollama:latest    (~800 MB)"
	@echo "   • postgres:16             (~400 MB)"
	@echo "   • redis:7.2               (~150 MB)"
	@echo "   • node:20-slim            (~200 MB)"
	@echo "   • python:3.12-slim        (~150 MB)"
	@echo ""
	@read -p "Continuar? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo ""; \
		echo "→ Pulling ollama..."; \
		docker pull ollama/ollama:latest; \
		echo "→ Pulling postgres..."; \
		docker pull postgres:16; \
		echo "→ Pulling redis..."; \
		docker pull redis:7.2; \
		echo "→ Pulling node..."; \
		docker pull node:20-slim; \
		echo "→ Pulling python..."; \
		docker pull python:3.12-slim; \
		echo ""; \
		echo "✅ Imagens atualizadas!"; \
	else \
		echo "❌ Cancelado"; \
	fi

# Alias for ollama update only
ollama-update:
	@echo "🔄 Atualizando apenas Ollama (~800 MB)..."
	@docker pull ollama/ollama:latest
	@echo "✅ Ollama atualizado!"
	@echo ""
	@echo "💡 Reinicie o Ollama para usar a nova versão:"
	@echo "   make ollama-restart"
