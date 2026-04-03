# ============================================================
#                   OLLAMA LIFECYCLE MANAGEMENT
# ============================================================
# 
# Ollama é gerenciado separadamente porque:
# - Modelos são grandes (3-5GB) e demoram para baixar
# - Não mudam frequentemente
# - Devem persistir entre rebuilds das outras aplicações
#
# ============================================================

.PHONY: ollama-up ollama-down ollama-status ollama-restart ollama-logs
.PHONY: ollama-models ollama-pull ollama-clean

# ============================================================
#                   BASIC COMMANDS
# ============================================================

ollama-up:
	@echo "🚀 Starting Ollama service..."
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" | grep -q aion-dev-ollama; then \
		echo "✅ Ollama is already running"; \
	else \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) up -d ollama; \
		echo "⏳ Waiting for Ollama to be healthy..."; \
		for i in $$(seq 1 30); do \
			if docker exec aion-dev-ollama ollama list >/dev/null 2>&1; then \
				echo "✅ Ollama is ready!"; \
				break; \
			fi; \
			if [ $$i -eq 30 ]; then \
				echo "⚠️  Timeout waiting for Ollama"; \
			fi; \
			sleep 1; \
		done; \
	fi
	@echo ""
	@echo "📍 Ollama API: http://localhost:11434"

ollama-down:
	@echo "🛑 Stopping Ollama service..."
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" | grep -q aion-dev-ollama; then \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) stop ollama; \
		echo "✅ Ollama stopped (volumes preserved)"; \
	else \
		echo "ℹ️  Ollama is not running"; \
	fi

ollama-restart:
	@$(MAKE) ollama-down
	@$(MAKE) ollama-up

ollama-status:
	@echo "🔍 Ollama Status"
	@echo "=================="
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" | grep -q aion-dev-ollama; then \
		echo "✅ Container: Running"; \
		echo ""; \
		echo "📦 Installed models:"; \
		docker exec aion-dev-ollama ollama list 2>/dev/null || echo "   (error listing models)"; \
		echo ""; \
		echo "💾 Volume info:"; \
		docker volume inspect dev_ollama-models --format "   Location: {{ .Mountpoint }}" 2>/dev/null || echo "   Volume not found"; \
		docker volume inspect dev_ollama-models --format "   Created: {{ .CreatedAt }}" 2>/dev/null || true; \
	else \
		echo "❌ Container: Not running"; \
		echo ""; \
		echo "💾 Volume status:"; \
		if docker volume ls | grep -q "dev_ollama-models"; then \
			echo "   ✅ Volume exists (models preserved)"; \
			docker volume inspect dev_ollama-models --format "   Location: {{ .Mountpoint }}" 2>/dev/null; \
		else \
			echo "   ❌ Volume does not exist"; \
		fi; \
	fi

ollama-logs:
	@echo "📋 Ollama Logs (Ctrl+C to exit)"
	@echo "==============================="
	@docker logs -f aion-dev-ollama

# ============================================================
#                   MODEL MANAGEMENT
# ============================================================

ollama-models:
	@echo "📦 Installed Ollama Models"
	@echo "============================"
	@if docker ps --filter "name=aion-dev-ollama" --filter "status=running" | grep -q aion-dev-ollama; then \
		docker exec aion-dev-ollama ollama list; \
		echo ""; \
		echo "Commands:"; \
		echo "   make ollama-pull MODEL=<name>  → Download a specific model"; \
		echo "   make ollama-pull               → Download default (qwen2.5:7b-instruct-q4_K_M)"; \
	else \
		echo "❌ Ollama is not running. Start it with: make ollama-up"; \
	fi

ollama-pull:
	@MODEL_NAME="$${MODEL:-qwen2.5:7b-instruct-q4_K_M}"; \
	echo "📥 Downloading model: $$MODEL_NAME"; \
	echo "   This may take 5-15 minutes depending on model size..."; \
	echo ""; \
	if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" | grep -q aion-dev-ollama; then \
		echo "⚠️  Ollama is not running. Starting..."; \
		$(MAKE) ollama-up; \
		echo ""; \
	fi; \
	docker exec aion-dev-ollama ollama pull "$$MODEL_NAME" && \
	echo "" && \
	echo "✅ Model downloaded successfully!" && \
	echo "" && \
	docker exec aion-dev-ollama ollama list

ollama-clean:
	@echo "⚠️  DANGER: Remove Ollama volumes (including all models)"
	@echo "========================================================"
	@echo "This will:"
	@echo "   • Stop Ollama container"
	@echo "   • Remove the container"
	@echo "   • DELETE the volume with all downloaded models (~4-5GB)"
	@echo ""
	@read -p "Are you SURE? Type 'yes' to confirm: " confirmation; \
	if [ "$$confirmation" = "yes" ]; then \
		echo ""; \
		echo "🛑 Stopping Ollama..."; \
		docker stop aion-dev-ollama 2>/dev/null || true; \
		echo "🗑️  Removing container..."; \
		docker rm aion-dev-ollama 2>/dev/null || true; \
		echo "🗑️  Removing volume..."; \
		docker volume rm dev_ollama-models 2>/dev/null || true; \
		echo ""; \
		echo "✅ Ollama cleaned. Next 'make ollama-up' will start fresh."; \
		echo "   Models will need to be downloaded again."; \
	else \
		echo "❌ Cancelled"; \
	fi

# ============================================================
#                   HELPER FUNCTIONS
# ============================================================

# Check if Ollama is running and warn if not
.PHONY: _check-ollama-running
_check-ollama-running:
	@if ! docker ps --filter "name=aion-dev-ollama" --filter "status=running" -q | grep -q .; then \
		echo ""; \
		echo "⚠️  NOTICE: Ollama is NOT running"; \
		echo "   Chat service will not work without Ollama."; \
		echo ""; \
		echo "   Start it with: make ollama-up"; \
		echo ""; \
	fi
