# ============================================================
#                         BACKUP / RESTORE
# ============================================================

.PHONY: backup-dev restore-dev backup-my restore-my

BACKUP_DIR ?= ../backups/aion-api/dev

backup-dev:
	@mkdir -p "$(BACKUP_DIR)"
	@backup_file="$(BACKUP_FILE)"; \
		if [ -z "$$backup_file" ]; then \
			backup_file="$(BACKUP_DIR)/aion-api-dev-$$(date -u +%Y%m%dT%H%M%SZ).dump"; \
		fi; \
		echo "Creating dev database backup: $$backup_file"; \
		docker exec $(POSTGRES_CONTAINER) pg_dump -U $(POSTGRES_USER) -d $(POSTGRES_DB) --format=custom --no-owner --no-privileges > "$$backup_file"; \
		echo "Backup created: $$backup_file"

restore-dev:
	@if [ "$(CONFIRM_RESTORE)" != "YES" ]; then \
		echo "Restore is destructive. Re-run with CONFIRM_RESTORE=YES BACKUP_FILE=../backups/aion-api/dev/<file>.dump"; \
		exit 1; \
	fi
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "BACKUP_FILE is required. Example: make restore-dev CONFIRM_RESTORE=YES BACKUP_FILE=../backups/aion-api/dev/aion-api-dev-20260416T120000Z.dump"; \
		exit 1; \
	fi
	@if [ ! -f "$(BACKUP_FILE)" ]; then \
		echo "Backup file not found: $(BACKUP_FILE)"; \
		exit 1; \
	fi
	@echo "Restoring dev database from $(BACKUP_FILE)"
	@cat "$(BACKUP_FILE)" | docker exec -i $(POSTGRES_CONTAINER) pg_restore -U $(POSTGRES_USER) -d $(POSTGRES_DB) --clean --if-exists --no-owner --no-privileges
	@echo "Restore completed from $(BACKUP_FILE)"

backup-my:
	@mkdir -p "$(MY_BACKUP_DIR)"
	@backup_file="$(BACKUP_FILE)"; \
		if [ -z "$$backup_file" ]; then \
			backup_file="$(MY_BACKUP_DIR)/aion-api-my-$$(date -u +%Y%m%dT%H%M%SZ).dump"; \
		fi; \
		echo "Creating MY database backup: $$backup_file"; \
		docker exec $(MY_POSTGRES_CONTAINER) pg_dump -U $(MY_POSTGRES_USER) -d $(MY_POSTGRES_DB) --format=custom --no-owner --no-privileges > "$$backup_file"; \
		echo "MY backup created: $$backup_file"

restore-my:
	@if [ "$(CONFIRM_RESTORE)" != "YES" ]; then \
		echo "MY restore is destructive. Re-run with CONFIRM_RESTORE=YES BACKUP_FILE=../backups/aion-api/my/<file>.dump"; \
		exit 1; \
	fi
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "BACKUP_FILE is required. Example: make restore-my CONFIRM_RESTORE=YES BACKUP_FILE=../backups/aion-api/my/aion-api-my-20260416T120000Z.dump"; \
		exit 1; \
	fi
	@if [ ! -f "$(BACKUP_FILE)" ]; then \
		echo "Backup file not found: $(BACKUP_FILE)"; \
		exit 1; \
	fi
	@echo "Restoring MY database from $(BACKUP_FILE)"
	@cat "$(BACKUP_FILE)" | docker exec -i $(MY_POSTGRES_CONTAINER) pg_restore -U $(MY_POSTGRES_USER) -d $(MY_POSTGRES_DB) --clean --if-exists --no-owner --no-privileges
	@echo "MY restore completed from $(BACKUP_FILE)"
