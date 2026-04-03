#!/usr/bin/env bash
# Script para testar o Aion Chat
set -euo pipefail

if [[ -z "${AION_CHAT_SERVICE_KEY:-}" ]]; then
    echo "AION_CHAT_SERVICE_KEY is required"
    exit 1
fi

echo "🤖 Testando Aion Chat..."
echo ""

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Teste 1: Health check do aion-chat
echo "1️⃣  Verificando health do aion-chat..."
if curl -s http://localhost:8000/health | grep -q "healthy"; then
    echo -e "${GREEN}✓ aion-chat está saudável${NC}"
else
    echo -e "${RED}✗ aion-chat não está respondendo${NC}"
    exit 1
fi

echo ""

# Teste 2: Verificar modelo Ollama
echo "2️⃣  Verificando modelo Ollama..."
if docker exec aion-dev-ollama ollama list | grep -q "qwen2.5:7b-instruct-q4_K_M"; then
    echo -e "${GREEN}✓ Modelo qwen2.5:7b-instruct-q4_K_M está instalado${NC}"
else
    echo -e "${RED}✗ Modelo não está instalado${NC}"
    echo "Execute: docker exec aion-dev-ollama ollama pull qwen2.5:7b-instruct-q4_K_M"
    exit 1
fi

echo ""

# Teste 3: Teste direto no aion-chat (Python)
echo "3️⃣  Testando endpoint interno do aion-chat..."
RESPONSE=$(curl -s -X POST http://localhost:8000/internal/process \
  -H "Content-Type: application/json" \
  -H "X-Service-Key: ${AION_CHAT_SERVICE_KEY}" \
  --max-time 120 \
  -d '{
    "user_id": 1,
    "message": "Olá! Qual é o seu nome?"
  }')

if echo "$RESPONSE" | grep -q "response"; then
    echo -e "${GREEN}✓ aion-chat respondeu com sucesso${NC}"
    echo ""
    echo "Resposta do Aion:"
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
else
    echo -e "${RED}✗ Erro na resposta do aion-chat${NC}"
    echo "$RESPONSE"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Todos os testes passaram!${NC}"
echo ""
echo "📌 Acesse o dashboard em: http://localhost:5000"
echo "📌 API disponível em: http://localhost:5001/aion/api/v1"
echo "📌 Chat Python em: http://localhost:8000"
