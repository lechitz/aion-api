# ============================================================
#                         MIGRATIONS
# ============================================================

# Default dev database URL (used by migrate-dev-* commands)
DEV_MIGRATION_DB ?= postgres://aion:aion123@localhost:5432/aion-api?sslmode=disable

.PHONY: migrate-up migrate-down migrate-force migrate-new migrate-dev-up migrate-dev-down migrate-dev-status migrate-my-up migrate-my-status migrate-install

# Install golang-migrate CLI
migrate-install:
	@echo "📦 Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ golang-migrate installed"

migrate-up:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "Running all migrations (up)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" up

migrate-down:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "↩️  Rolling back the last migration (1 step)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" down 1

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION not provided. Use 'make migrate-force VERSION=X'"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "🚨 Forcing DB schema version to $(VERSION)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" force "$(VERSION)"

migrate-new:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Migration name is required"; \
		exit 1; \
	fi; \
	$(MIGRATE_BIN) create -ext sql -dir "$(MIGRATION_PATH)" "$$name"

# ============================================================
#                 DEV ENVIRONMENT MIGRATIONS
# ============================================================
# These commands use DEV_MIGRATION_DB by default (localhost:5432)
# Requires: postgres container running (make dev or make dev-fast)

migrate-dev-up:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "🚀 Running all migrations on DEV database..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(DEV_MIGRATION_DB)" up
	@echo "✅ Migrations applied successfully"

migrate-dev-down:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "↩️  Rolling back last migration on DEV database..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(DEV_MIGRATION_DB)" down 1

migrate-dev-status:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "📊 Migration status on DEV database..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(DEV_MIGRATION_DB)" version

migrate-dev-reset:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "⚠️  Resetting DEV database (dropping all and re-applying)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(DEV_MIGRATION_DB)" drop -f || true
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(DEV_MIGRATION_DB)" up
	@$(MAKE) seed-roles
	@echo "✅ DEV database reset complete"

# ============================================================
#                 MY ENVIRONMENT MIGRATIONS
# ============================================================
# These commands use MY_MIGRATION_DB by default.
# They intentionally do not provide a reset target because MY is for durable
# personal data.

migrate-my-up:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "🚀 Running all migrations on MY database..."
	@set +e; \
		output="$$( $(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MY_MIGRATION_DB)" up 2>&1 )"; \
		status="$$?"; \
		set -e; \
		echo "$$output" | grep -v "no change" || true; \
		if [ "$$status" -ne 0 ] && ! echo "$$output" | grep -q "no change"; then \
			exit "$$status"; \
		fi
	@echo "✅ MY migrations checked/applied"

migrate-my-status:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Run: make migrate-install"; \
		exit 1; \
	fi
	@echo "📊 Migration status on MY database..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MY_MIGRATION_DB)" version
