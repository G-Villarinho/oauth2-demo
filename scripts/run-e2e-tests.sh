#!/bin/bash

# Script para executar testes E2E
# Uso: ./scripts/run-e2e-tests.sh

set -e

echo "🚀 Iniciando testes E2E..."

# Verificar se o Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor, inicie o Docker e tente novamente."
    exit 1
fi

# Parar containers existentes (se houver)
echo "🛑 Parando containers existentes..."
docker-compose down --remove-orphans

# Construir e subir os serviços
echo "🔨 Construindo e subindo serviços..."
docker-compose up --build -d mongodb app

# Aguardar serviços estarem prontos
echo "⏳ Aguardando serviços estarem prontos..."
sleep 10

# Executar testes E2E
echo "🧪 Executando testes E2E..."
docker-compose run --rm e2e-tests

# Capturar o código de saída dos testes
TEST_EXIT_CODE=$?

# Parar todos os serviços
echo "🛑 Parando serviços..."
docker-compose down

# Verificar resultado dos testes
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✅ Todos os testes E2E passaram!"
    exit 0
else
    echo "❌ Alguns testes E2E falharam!"
    exit 1
fi 