#!/usr/bin/env bash
# Script de Teste GPU - Execute para verificar tudo
set -euo pipefail

echo "======================================"
echo "🎮 TESTE DE GPU - OLLAMA"
echo "======================================"
echo ""

echo "1️⃣  Verificando se GPU está visível..."
if docker exec aion-dev-ollama nvidia-smi > /dev/null 2>&1; then
    echo "✅ GPU detectada!"
    docker exec aion-dev-ollama nvidia-smi | grep "NVIDIA GeForce"
else
    echo "❌ GPU não detectada"
    echo "Execute: sudo systemctl restart docker && make dev"
    exit 1
fi
echo ""

echo "2️⃣  Verificando modelo Qwen..."
if docker exec aion-dev-ollama ollama list | grep -q "qwen2.5"; then
    echo "✅ Modelo Qwen instalado!"
else
    echo "⏳ Baixando modelo (~5 min)..."
    docker exec aion-dev-ollama ollama pull qwen2.5:7b-instruct-q4_K_M
fi
echo ""

echo "3️⃣  Teste de performance..."
echo "   Pergunta: 'Conte até 5'"
echo "   Aguarde..."
START=$(date +%s.%N)
docker exec aion-dev-ollama ollama run qwen2.5:7b-instruct-q4_K_M "Conte até 5" > /tmp/ollama-test.txt 2>&1
END=$(date +%s.%N)
DIFF=$(echo "$END - $START" | bc)

echo ""
echo "⏱️  Tempo de resposta: ${DIFF}s"
echo ""

if (( $(echo "$DIFF < 5.0" | bc -l) )); then
    echo "✅ PERFEITO! GPU está funcionando!"
    echo "   (< 5s = GPU funcionando)"
elif (( $(echo "$DIFF < 15.0" | bc -l) )); then
    echo "⚠️  OK, mas pode estar usando CPU"
    echo "   (5-15s = pode ser CPU ou GPU lenta)"
else
    echo "❌ LENTO! Provavelmente usando CPU"
    echo "   (> 15s = definitivamente CPU)"
fi
echo ""

echo "4️⃣  Resposta do modelo:"
cat /tmp/ollama-test.txt
echo ""

echo "======================================"
echo "📊 RESUMO"
echo "======================================"
echo "GPU Detectada: ✅"
echo "Modelo Instalado: ✅"
echo "Tempo: ${DIFF}s"

if (( $(echo "$DIFF < 5.0" | bc -l) )); then
    echo "Performance: ⚡ EXCELENTE (GPU)"
elif (( $(echo "$DIFF < 15.0" | bc -l) )); then
    echo "Performance: ⚠️  MÉDIO"
else
    echo "Performance: ❌ LENTO (CPU)"
fi

echo ""
echo "🎯 Próximo passo:"
echo "   open http://localhost:5000"
echo "   Digite: 'Olá Aion!'"
echo ""
