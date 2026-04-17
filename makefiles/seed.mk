# ============================================================
#                          SEEDS
# ============================================================
# All seeds use *_generate.sql (dynamic, parametrizable via N=count)
# Usage: make seed-all N=10 | make populate N=100 | make seed-caller N=5
# ============================================================
.PHONY: seed-users seed-categories seed-all seed-tags seed-records seed-roles seed-user-roles seed-admin seed-user1-all seed-everybody seed-clean-users seed-clean-categories seed-clean-tags seed-clean-records seed-clean-roles seed-clean-user-roles seed-clean-registration-sessions seed-clean-all seed-helper seed-setup seed-quick seed-api-caller seed-api-caller-bootstrap seed-api-caller-clean seed-caller populate reset-user-data seed-test-timeline seed-clean-test-timeline seed-essential seed-my-essential db-full db-reset

POSTGRES_CONTAINER := aion-dev-postgres
POSTGRES_USER := aion
POSTGRES_DB := aion-api

SEED_DEFAULT_PASSWORD := testpassword123
SEED_DEFAULT_PASSWORD_HASH := $$2a$$10$$BIv0nYxelFEGDods46gtuuIpGH8NCThM1frbbhG5Ro/UqQ80ziwXS

# Default count for seeds (can be overridden: make seed-all N=100)
N ?= 10
ifdef n
N := $(n)
endif

# Build the seed-helper tool
seed-helper:
	@echo "Building seed-helper..."
	@go build -o bin/seed-helper ./hack/tools/seed-helper

# Generate .env.local with all seed variables (interactive setup)
seed-setup: seed-helper
	@echo "Setting up seed environment..."
	@read -p "Number of users to generate (default 10): " count; \
	count=$${count:-10}; \
	./bin/seed-helper generate-env $$count
	@echo ""
	@echo "✅ Setup complete! Now you can run: make seed-quick"

# Quick seed using .env.local (must run seed-setup first)
seed-quick:
	@echo "Quick seeding with .env.local..."
	@if [ ! -f infrastructure/db/seed/.env.local ]; then \
		echo "❌ .env.local not found. Run 'make seed-setup' first."; \
		exit 1; \
	fi
	@export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	$(MAKE) seed-all

# Seeds use *_generate.sql files (dynamic, parametrizable)
# Usage: make seed-users N=10 (or SEED_USER_COUNT=10)
seed-users:
	@echo "Seeding users..."
	@if [ -f infrastructure/db/seed/.env.local ]; then \
		export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	fi; \
	count=$${SEED_USER_COUNT:-$(N)}; \
	count=$${count:-10}; \
	echo "Generating $$count users using pgcrypto (password: $${DEV_PASSWORD:-$(SEED_DEFAULT_PASSWORD)})..."; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-v seed_count="$$count" \
		-v user_seed_password_plain="$${DEV_PASSWORD:-$(SEED_DEFAULT_PASSWORD)}" \
		< infrastructure/db/seed/user_generate.sql

seed-categories:
	@echo "Seeding categories..."
	@if [ -f infrastructure/db/seed/.env.local ]; then \
		export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	fi; \
	count=$${SEED_USER_COUNT:-$(N)}; \
	count=$${count:-10}; \
	echo "Generating categories for $$count users..."; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-v seed_count="$$count" \
		< infrastructure/db/seed/category_generate.sql

seed-tags:
	@echo "Seeding tags..."
	@if [ -f infrastructure/db/seed/.env.local ]; then \
		export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	fi; \
	count=$${SEED_USER_COUNT:-$(N)}; \
	count=$${count:-10}; \
	echo "Generating tags for $$count users..."; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-v seed_count="$$count" \
		< infrastructure/db/seed/tags_generate.sql

seed-records:
	@echo "Seeding records..."
	@if [ -f infrastructure/db/seed/.env.local ]; then \
		export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	fi; \
	count=$${SEED_USER_COUNT:-$(N)}; \
	count=$${count:-10}; \
	days=$${SEED_DAYS:-7}; \
	echo "Generating records for $$count users ($$days days of data)..."; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-v seed_count="$$count" \
		-v days="$$days" \
		< infrastructure/db/seed/records_generate.sql

# --- Admin user seed ---
seed-admin:
	@echo "Seeding admin user 'aion'..."
	@if [ -z "$${USER_TOKEN_TEST:-}" ]; then \
		echo "USER_TOKEN_TEST not set; using default hash for password '$(SEED_DEFAULT_PASSWORD)'"; \
		USER_TOKEN_TEST='$(SEED_DEFAULT_PASSWORD_HASH)'; \
	fi; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v user_seed_password_hash="$${USER_TOKEN_TEST}" < infrastructure/db/seed/admin_user.sql

# seed-all includes everything in correct order:
# 1. roles (must exist first)
# 2. admin user (aion with admin role)
# 3. regular users
# 4. user_roles (assigns 'user' role to all users without roles)
# 5. business data (categories, tags, records)
# Usage: make seed-all N=10 (default: 10 users)
seed-all: seed-roles seed-admin seed-users seed-user-roles seed-categories seed-tags seed-records
	@echo "✅ All seeds applied."

# Convenience target: seeds the full dataset for N=1 user
seed-user1-all:
	@$(MAKE) seed-all N=1

# Alias to seed everything available; easier mnemonic than seed-all.
seed-everybody: seed-all
	@echo "✅ Everyone seeded."

seed-api-caller:
	@echo "Seeding via API (HTTP/GraphQL)..."
	@go run ./hack/tools/seed-caller

seed-api-caller-bootstrap:
	@echo "Seeding via API (bootstrap: cria usuário se necessário)..."
	@API_CALLER_AUTO_CREATE=true go run ./hack/tools/seed-caller

seed-api-caller-clean:
	@echo "Limpando via API (soft delete de records, sem criar nada)..."
	@API_CALLER_CLEAN=true API_CALLER_ONLY_CLEAN=true go run ./hack/tools/seed-caller

seed-caller:
	@echo "Seeding via API (multi-user: count=$(N))..."
	@echo "📋 Step 0/4: Ensuring migrations are applied..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path infrastructure/db/migrations -database "postgres://aion:aion123@localhost:5432/aion-api?sslmode=disable" up 2>&1 | grep -v "no change" || true; \
	else \
		echo "⚠️  'migrate' not found. Assuming DB is already initialized."; \
	fi
	@echo "📋 Step 1/4: Seeding roles..."
	@$(MAKE) seed-roles
	@echo "📋 Step 2/4: Seeding admin user (aion)..."
	@$(MAKE) seed-admin
	@echo "📋 Step 3/4: Waiting for API to be ready..."
	@for i in 1 2 3 4 5; do \
		curl -sf http://localhost:5001/aion/api/v1/health > /dev/null 2>&1 && break || sleep 2; \
	done
	@echo "📋 Step 4/4: Running API seed caller..."
	@API_CALLER_COUNT=$(N) API_CALLER_AUTO_CREATE=true go run ./hack/tools/seed-caller

seed-api-caller-many: seed-caller

# Populate N users (default 10) with categories/tags/records via SQL generators.
# Usage: make populate N=100
populate:
	@echo "Cleaning tables and populating $(N) users (password=$(SEED_DEFAULT_PASSWORD))..."
	@$(MAKE) seed-clean-all
	@$(MAKE) seed-all N=$(N)
	@echo "✅ Populate completed for $(N) users."

# --- Clean helpers (dev-only): truncate seeded tables safely and reset IDs ---
# NOTE: Intended for local dev. These use TRUNCATE ... RESTART IDENTITY CASCADE.
seed-clean-users:
	@echo "Truncating users (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.users RESTART IDENTITY CASCADE;"

seed-clean-categories:
	@echo "Truncating categories (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.categories RESTART IDENTITY CASCADE;"

seed-clean-tags:
	@echo "Truncating tags (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.tags RESTART IDENTITY CASCADE;"

seed-clean-records:
	@echo "Truncating records (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.records RESTART IDENTITY CASCADE;"

seed-clean-user-roles:
	@echo " Truncating user_roles (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.user_roles RESTART IDENTITY CASCADE;"

seed-clean-roles:
	@echo "Truncating roles (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.roles RESTART IDENTITY CASCADE;"

seed-clean-registration-sessions:
	@echo "Truncating registration_sessions (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.registration_sessions RESTART IDENTITY CASCADE;"

seed-clean-all: seed-clean-records seed-clean-tags seed-clean-categories seed-clean-user-roles seed-clean-users seed-clean-roles seed-clean-registration-sessions
	@echo "✅ All seeded tables truncated (dev only)."

# --- Missing seed targets (referenced by seed-all but not defined) ---
seed-roles:
	@echo "📋 Seeding system roles..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/roles.sql
	@echo "✅ Roles seeded (owner, admin, user, blocked)"

seed-user-roles:
	@echo "📋 Assigning default roles to users without roles..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/user_roles.sql
	@echo "✅ User roles assigned"

# --- Reset user data (keeps system data like roles) ---
# Deletes all user-generated data but preserves system configuration
# Perfect for starting fresh without losing roles/permissions setup
reset-user-data: seed-clean-records seed-clean-tags seed-clean-categories seed-clean-user-roles seed-clean-users seed-clean-registration-sessions
	@echo "Resetting cache..."
	@$(MAKE) cache-reset
	@echo "📋 Re-seeding system roles..."
	@$(MAKE) seed-roles
	@echo ""
	@echo "✅ User data reset complete!"
	@echo "   ✓ Users deleted"
	@echo "   ✓ Categories deleted"
	@echo "   ✓ Tags deleted"
	@echo "   ✓ Records deleted"
	@echo "   ✓ Registration sessions deleted"
	@echo "   ✓ Cache cleared"
	@echo "   ✓ System roles preserved"
	@echo ""
	@echo "💡 Next steps:"
	@echo "   → Create new user: make seed-admin (or signup via API)"
	@echo "   → Seed test data: make seed-all N=10"

# ============================================================
# Hash generation helper
# ============================================================
.PHONY: hash-gen

hash-gen:
	@if [ -z "$(PASS)" ]; then \
		echo "❌ Error: PASS parameter is required"; \
		echo "Usage: make hash-gen PASS='yourpassword'"; \
		exit 1; \
	fi
	@echo "Generating bcrypt hash for: $(PASS)"
	@echo ""
	@go run -mod=readonly -tags tools golang.org/x/crypto/bcrypt/cmd/bcrypt -cost=10 "$(PASS)" 2>/dev/null || \
	{ \
		echo "Password: $(PASS)"; \
		GO_CODE='package main; import ("fmt"; "os"; "golang.org/x/crypto/bcrypt"); func main() { h, _ := bcrypt.GenerateFromPassword([]byte(os.Args[1]), 10); fmt.Printf("Hash:     %s\n", string(h)); err := bcrypt.CompareHashAndPassword(h, []byte(os.Args[1])); if err == nil { fmt.Println("✅ Hash verified successfully!") } else { fmt.Printf("❌ Hash verification failed: %v\n", err) }}'; \
		echo "$$GO_CODE" | go run - "$(PASS)"; \
	}


# ============================================================
#                   REALISTIC DEMO DATASET
# ============================================================
# Complete demo profile seed
# - Test user (username: testuser, password: Test@123)
# - Canonical categories/tags used by dashboard and records
# - Metric definitions + goal templates
# - ~3 months history (50-60 records/day)
# - Adds current-day records for immediate dashboard visibility
# - Only runs in dev/local (user_id=999)
# ============================================================

.PHONY: seed-test seed-clean-test seed-essential db-full db-test db-reset

# Generate complete test profile
seed-test:
	@echo "Generating realistic demo profile..."
	@echo "   • Test user (testuser / Test@123)"
	@echo "   • Legacy + new categories/tags merged for dashboard + timeline"
	@echo "   • Metric definitions + goal templates"
	@echo "   • ~3 months history (50-60 records/day)"
	@echo "   • Records always include duration + recorded_at"
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		< infrastructure/db/seed/test_data.sql
	@echo "✅ Realistic demo profile generated!"
	@echo ""
	@echo "🔐 Test User Login:"
	@echo "   username: testuser"
	@echo "   password: Test@123"

# Remove complete test profile
seed-clean-test:
	@echo "Removing test profile..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-c "DELETE FROM aion_api.metric_definition_tag_bindings WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.goal_instances WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.goal_templates WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.dashboard_widgets WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.dashboard_views WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.metric_definitions WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.records WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.tags WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.categories WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.user_roles WHERE user_id = 999;" \
		-c "DELETE FROM aion_api.users WHERE user_id = 999;"
	@echo "✅ Test profile removed"

# Seed only essential data (roles + admin user)
seed-essential: seed-roles seed-admin
	@echo "✅ Essential data seeded (roles + admin)"

seed-my-essential:
	@echo "📋 Seeding essential MY data without touching user-entered records..."
	@$(MAKE) seed-essential \
		POSTGRES_CONTAINER=$(MY_POSTGRES_CONTAINER) \
		POSTGRES_USER=$(MY_POSTGRES_USER) \
		POSTGRES_DB=$(MY_POSTGRES_DB)
	@echo "✅ MY essential data checked/applied"

# Full database setup: migrations + essential + test profile + cache flush
db-full: migrate-dev-reset seed-essential seed-test cache-reset
	@echo ""
	@echo "✅ Database ready with realistic demo data!"
	@echo ""
	@echo "📊 Summary:"
	@echo "   ✓ Database reset + migrations applied"
	@echo "   ✓ System roles created"
	@echo "   ✓ Admin user (username: aion, password: testpassword123)"
	@echo "   ✓ Test user with 3 months realistic profile (username: testuser, password: Test@123)"
	@echo ""
	@echo "🔐 Available Logins:"
	@echo "   Admin:  username: aion      password: testpassword123"
	@echo "   Test:   username: testuser  password: Test@123"
	@echo ""
	@echo "Platform ready for realistic analysis!"

# Alias: explicit name for test/demo profile setup
db-test: db-full
	@echo "✅ db-test completed (same as db-full)"

# Reset database and apply migrations
db-reset:
	@echo "⚠️  Resetting database..."
	@$(MAKE) migrate-dev-reset
	@echo "✅ Database reset and migrations applied"
