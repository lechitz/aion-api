#!/usr/bin/env bash
# Check if Ollama model exists and download if needed
set -euo pipefail

MODEL_NAME="qwen2.5:7b-instruct-q4_K_M"
CONTAINER_NAME="aion-dev-ollama"

echo "🔍 Checking Ollama model availability..."

# Check if ollama container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "❌ Ollama container '${CONTAINER_NAME}' is not running"
    echo "   Starting docker compose first..."
    exit 1
fi

# Wait for Ollama to be ready (max 30 seconds)
echo "⏳ Waiting for Ollama to be ready..."
for i in {1..30}; do
    if docker exec "${CONTAINER_NAME}" ollama list >/dev/null 2>&1; then
        echo "✅ Ollama is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Timeout waiting for Ollama to be ready"
        exit 1
    fi
    sleep 1
done

# Check if model exists
echo "🔍 Checking if model '${MODEL_NAME}' exists..."
if docker exec "${CONTAINER_NAME}" ollama list 2>/dev/null | grep -q "${MODEL_NAME}"; then
    echo "✅ Model '${MODEL_NAME}' is already installed"
    exit 0
fi

# Model not found, download it
echo "📥 Model not found. Downloading '${MODEL_NAME}'..."
echo "   This may take 5-10 minutes (~3-5GB download)"
echo ""

if docker exec "${CONTAINER_NAME}" ollama pull "${MODEL_NAME}"; then
    echo ""
    echo "✅ Model '${MODEL_NAME}' successfully downloaded!"
else
    echo ""
    echo "❌ Failed to download model"
    exit 1
fi

# Verify model is working
echo ""
echo "Testing model..."
if docker exec "${CONTAINER_NAME}" ollama run "${MODEL_NAME}" "Responda apenas: OK" 2>/dev/null | grep -q "OK"; then
    echo "✅ Model is working correctly!"
else
    echo "⚠️  Model downloaded but test failed (this is usually OK)"
fi

echo ""
echo "✅ Ollama setup complete!"
