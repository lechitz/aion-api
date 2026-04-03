#!/usr/bin/env bash
set -euo pipefail

echo "🔍 DIAGNÓSTICO AUTOMÁTICO DO AION + OLLAMA"
echo "==========================================="
echo ""

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Verificar Ollama
echo "1️⃣  Verificando Ollama..."
if docker ps | grep -q aion-dev-ollama; then
    echo -e "${GREEN}✅ Ollama está rodando${NC}"

    # Verificar se está healthy
    if docker ps | grep aion-dev-ollama | grep -q "healthy"; then
        echo -e "${GREEN}✅ Ollama está healthy${NC}"
    else
        echo -e "${YELLOW}⚠️  Ollama não está healthy ainda${NC}"
    fi
else
    echo -e "${RED}❌ Ollama NÃO está rodando${NC}"
    echo "   Solution: run 'make dev' from repository root"
    exit 1
fi

# 2. Verificar Modelo
echo ""
echo "2️⃣  Verificando modelo Qwen..."
if docker exec aion-dev-ollama ollama list | grep -q "qwen2.5:7b-instruct-q4_K_M"; then
    echo -e "${GREEN}✅ Modelo Qwen instalado${NC}"
else
    echo -e "${RED}❌ Modelo Qwen NÃO instalado${NC}"
    echo "   Baixando agora (~5-10 min, 4.7GB)..."
    docker exec aion-dev-ollama ollama pull qwen2.5:7b-instruct-q4_K_M
    echo -e "${GREEN}✅ Modelo baixado${NC}"
fi

# 3. Verificar Aion Chat
echo ""
echo "3️⃣  Verificando Aion Chat..."
if docker ps | grep -q aion-dev-chat; then
    echo -e "${GREEN}✅ Aion Chat está rodando${NC}"
else
    echo -e "${RED}❌ Aion Chat NÃO está rodando${NC}"
    echo "   Solution: run 'make dev' from repository root"
    exit 1
fi

# 4. Verificar Conectividade
echo ""
echo "4️⃣  Testando conectividade..."
if docker exec aion-dev-chat curl -s http://aion-dev-ollama:11434/api/version > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Aion Chat consegue acessar Ollama${NC}"
else
    echo -e "${RED}❌ Aion Chat NÃO consegue acessar Ollama${NC}"
    echo "   Verificando redes..."

    OLLAMA_NET=$(docker inspect aion-dev-ollama -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}')
    CHAT_NET=$(docker inspect aion-dev-chat -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}')

    echo "   Ollama: $OLLAMA_NET"
    echo "   Chat: $CHAT_NET"

    if [ "$OLLAMA_NET" != "$CHAT_NET" ]; then
        echo "   Conectando à mesma rede..."
        docker network connect rede_local aion-dev-ollama 2>/dev/null || true
        echo -e "${GREEN}✅ Conectado${NC}"
    fi
fi

# 5. Verificar Variáveis de Ambiente
echo ""
echo "5️⃣  Verificando configuração..."
LLM_PROVIDER=$(docker exec aion-dev-chat printenv LLM_PROVIDER 2>/dev/null || echo "não definido")
OLLAMA_URL=$(docker exec aion-dev-chat printenv OLLAMA_BASE_URL 2>/dev/null || echo "não definido")

echo "   LLM_PROVIDER: $LLM_PROVIDER"
echo "   OLLAMA_BASE_URL: $OLLAMA_URL"

if [ "$LLM_PROVIDER" = "ollama" ]; then
    echo -e "${GREEN}✅ Configurado para usar Ollama${NC}"
else
    echo -e "${YELLOW}⚠️  Configurado para usar $LLM_PROVIDER${NC}"
    echo "   Esperado: ollama"
fi

if [ "$OLLAMA_URL" = "http://aion-dev-ollama:11434" ]; then
    echo -e "${GREEN}✅ URL do Ollama correta${NC}"
else
    echo -e "${YELLOW}⚠️  URL do Ollama: $OLLAMA_URL${NC}"
    echo "   Esperado: http://aion-dev-ollama:11434"
fi

# 6. Verificar Logs
echo ""
echo "6️⃣  Verificando logs do Aion Chat..."
if docker logs aion-dev-chat --tail 50 2>&1 | grep -q "Ollama client initialized"; then
    echo -e "${GREEN}✅ Aion Chat inicializou cliente Ollama${NC}"
else
    echo -e "${RED}❌ Aion Chat NÃO inicializou cliente Ollama${NC}"
    echo ""
    echo "   Últimas linhas dos logs:"
    docker logs aion-dev-chat --tail 10 2>&1 | sed 's/^/   /'
fi

# 7. Resumo e Ações
echo ""
echo "📋 RESUMO"
echo "=========="
echo ""

OLLAMA_OK=$(docker ps | grep -q aion-dev-ollama && echo "✅" || echo "❌")
MODEL_OK=$(docker exec aion-dev-ollama ollama list 2>/dev/null | grep -q qwen && echo "✅" || echo "❌")
CHAT_OK=$(docker ps | grep -q aion-dev-chat && echo "✅" || echo "❌")
CONN_OK=$(docker exec aion-dev-chat curl -s http://aion-dev-ollama:11434/api/version > /dev/null 2>&1 && echo "✅" || echo "❌")

echo "$OLLAMA_OK Ollama rodando"
echo "$MODEL_OK Modelo instalado"
echo "$CHAT_OK Aion Chat rodando"
echo "$CONN_OK Conectividade OK"

echo ""

# Detectar problema e sugerir solução
if [ "$OLLAMA_OK" = "❌" ] || [ "$CHAT_OK" = "❌" ]; then
    echo -e "${YELLOW}🔧 AÇÃO NECESSÁRIA:${NC}"
    echo "   Run from repository root:"
    echo "   make dev"
    exit 1
fi

if [ "$MODEL_OK" = "❌" ]; then
    echo -e "${YELLOW}🔧 AÇÃO NECESSÁRIA:${NC}"
    echo "   docker exec aion-dev-ollama ollama pull qwen2.5:7b-instruct-q4_K_M"
    echo "   docker restart aion-dev-chat"
    exit 1
fi

if [ "$CONN_OK" = "❌" ]; then
    echo -e "${YELLOW}🔧 AÇÃO NECESSÁRIA:${NC}"
    echo "   docker network connect rede_local aion-dev-ollama"
    echo "   docker restart aion-dev-chat"
    exit 1
fi

# Tudo OK!
echo -e "${GREEN}✅ TUDO FUNCIONANDO!${NC}"
echo ""
echo "Reiniciando aion-chat para garantir..."
docker restart aion-dev-chat
echo ""
echo "Aguarde 10 segundos..."
sleep 10
echo ""
echo "✅ Pronto! Teste em: http://localhost:5000"
