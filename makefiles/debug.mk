# Debug script for roles issue
# Run: make debug-roles

.PHONY: debug-roles test-user-create

debug-roles:
	@echo "=== DEBUGGING ROLES SYSTEM ==="
	@echo ""
	@echo "1. Checking if roles table exists..."
	@docker exec aion-dev-postgres psql -U aion -d aion-api -c "\d aion_api.roles" || echo "❌ Table does not exist"
	@echo ""
	@echo "2. Counting roles..."
	@docker exec aion-dev-postgres psql -U aion -d aion-api -t -c "SELECT COUNT(*) FROM aion_api.roles;" || echo "❌ Query failed"
	@echo ""
	@echo "3. Listing roles..."
	@docker exec aion-dev-postgres psql -U aion -d aion-api -c "SELECT role_id, name, is_active FROM aion_api.roles ORDER BY role_id;" || echo "❌ Query failed"
	@echo ""
	@echo "4. Checking API health..."
	@curl -sf http://localhost:5001/aion/api/v1/health && echo "✅ API is healthy" || echo "❌ API is not responding"
	@echo ""
	@echo "5. Testing user creation..."
	@$(MAKE) test-user-create

test-user-create:
	@echo "Creating test user..."
	@curl -X POST http://localhost:5001/aion/api/v1/user/create \
		-H "Content-Type: application/json" \
		-d '{"name":"Debug User","username":"debuguser$$(date +%s)","email":"debug$$(date +%s)@test.com","password":"test123"}' \
		-w "\nHTTP Status: %{http_code}\n" \
		2>&1 || echo "❌ Request failed"
