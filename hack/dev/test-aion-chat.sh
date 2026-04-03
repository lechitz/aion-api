#!/usr/bin/env bash
#
# SCRIPT DE VALIDAÇÃO COMPLETA - Aion Chat
# Testa os 4 casos críticos após correções
#

set -euo pipefail

if [[ -z "${AION_CHAT_SERVICE_KEY:-}" ]]; then
    echo "AION_CHAT_SERVICE_KEY is required"
    exit 1
fi

SERVICE_KEY="${AION_CHAT_SERVICE_KEY}"
BASE_URL="http://localhost:8000/internal/process"

echo "🧪 INICIANDO TESTES DE VALIDAÇÃO AION CHAT"
echo "==========================================="
echo ""

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Contador de testes
PASSED=0
FAILED=0

test_case() {
    local test_name="$1"
    local message="$2"
    local expected_pattern="$3"
    local should_not_contain="${4:-}"

    echo "───────────────────────────────────────────"
    echo "📋 Teste: $test_name"
    echo "💬 Mensagem: $message"
    echo ""

    response=$(curl -s -X POST "$BASE_URL" \
      -H "Content-Type: application/json" \
      -H "X-Service-Key: $SERVICE_KEY" \
      -d "{\"user_id\":1,\"message\":\"$message\",\"history\":[]}")

    echo "📤 Resposta:"
    echo "$response" | jq -r '.response' | head -n 5
    echo ""

    # Check expected pattern
    if echo "$response" | jq -r '.response' | grep -qi "$expected_pattern"; then
        echo -e "${GREEN}✅ PASSOU${NC} - Contém padrão esperado: $expected_pattern"
        ((PASSED++))
    else
        echo -e "${RED}❌ FALHOU${NC} - NÃO contém padrão esperado: $expected_pattern"
        ((FAILED++))
    fi

    # Check should NOT contain
    if [ -n "$should_not_contain" ]; then
        if echo "$response" | jq -r '.response' | grep -qi "$should_not_contain"; then
            echo -e "${RED}❌ FALHOU${NC} - Contém padrão indesejado: $should_not_contain"
            ((FAILED++))
            ((PASSED--))  # Desconta o passe anterior
        else
            echo -e "${GREEN}✅ OK${NC} - NÃO contém padrão indesejado: $should_not_contain"
        fi
    fi

    echo ""
}

echo "🔧 Aguardando containers ficarem prontos..."
sleep 5

echo ""
echo "═══════════════════════════════════════════"
echo "  TESTE 1: Anti-Alucinação (Nome de Usuário)"
echo "═══════════════════════════════════════════"
test_case \
    "Nome do usuário (não deve inventar)" \
    "qual meu nome?" \
    "não tenho\|não sei\|não possuo" \
    "joão\|maria\|carlos\|ana"

echo "═══════════════════════════════════════════"
echo "  TESTE 2: Listar Categorias (Dados Reais)"
echo "═══════════════════════════════════════════"
test_case \
    "Listar categorias reais" \
    "quais minhas categorias?" \
    "pessoal\|saude_fisica\|saude_mental" \
    ""

echo "═══════════════════════════════════════════"
echo "  TESTE 3: Listar Tags de Categoria"
echo "═══════════════════════════════════════════"
test_case \
    "Listar tags de saude_fisica" \
    "quais tags de saude_fisica?" \
    "Run\|Stretching" \
    "Mindfulness\|Terapia\|Yoga"  # Não deve inventar tags

echo "═══════════════════════════════════════════"
echo "  TESTE 4: Create Record (Fluxo Completo)"
echo "═══════════════════════════════════════════"
# Este teste verifica se IA executa create_record ao invés de só perguntar
test_case \
    "Criar registro com tag Stretching" \
    "registre alongamento 22min na tag Stretching categoria saude_fisica" \
    "registrei\|pronto\|criado\|sucesso" \
    "<tool_call>\|qual categoria\|qual tag"  # Não deve pedir mais info

echo ""
echo "═══════════════════════════════════════════"
echo "           RESULTADO FINAL"
echo "═══════════════════════════════════════════"
echo ""
echo -e "✅ Testes Passados: ${GREEN}$PASSED${NC}"
echo -e "❌ Testes Falhados: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM!${NC}"
    echo ""
    echo "Sistema está funcionando conforme esperado:"
    echo "  ✅ Não alucina dados"
    echo "  ✅ Lista dados reais corretamente"
    echo "  ✅ Executa create_record quando tem todos dados"
    echo ""
    exit 0
else
    echo -e "${RED}⚠️  ALGUNS TESTES FALHARAM${NC}"
    echo ""
    echo "Verifique os logs:"
    echo "  docker logs -f aion-dev-chat"
    echo ""
    exit 1
fi
