# Makefile para Aetheris ID API
# Uso: make <comando>

.PHONY: help build up down logs test test-e2e test-unit clean docker-clean generate-mocks

# Comandos principais
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker Compose
build: ## Constrói as imagens Docker
	docker-compose build --no-cache

up: ## Sobe os serviços (MongoDB + App)
	docker-compose up -d mongodb app

down: ## Para e remove todos os containers
	docker-compose down --remove-orphans

logs: ## Mostra logs da aplicação
	docker-compose logs -f app

logs-mongo: ## Mostra logs do MongoDB
	docker-compose logs -f mongodb

# Testes
test: ## Executa todos os testes (unit + e2e)
	@echo "🧪 Executando todos os testes..."
	$(MAKE) test-unit
	$(MAKE) test-e2e

test-unit: ## Executa testes unitários
	@echo "🧪 Executando testes unitários..."
	go test -v ./...

test-e2e: ## Executa testes E2E usando o script
	@echo "🚀 Executando testes E2E via script..."
	@if [ -f "./scripts/run-e2e-tests.sh" ]; then \
		chmod +x ./scripts/run-e2e-tests.sh; \
		./scripts/run-e2e-tests.sh; \
	else \
		echo "❌ Script de testes E2E não encontrado: ./scripts/run-e2e-tests.sh"; \
		echo "💡 Executando testes E2E diretamente..."; \
		$(MAKE) test-e2e-direct; \
	fi

test-e2e-direct: ## Executa testes E2E diretamente (sem script)
	@echo "🚀 Iniciando testes E2E..."
	@echo "🛑 Parando containers existentes..."
	docker-compose down --remove-orphans
	@echo "🔨 Construindo e subindo serviços..."
	docker-compose up --build -d mongodb app
	@echo "⏳ Aguardando serviços estarem prontos..."
	@sleep 10
	@echo "🧪 Executando testes E2E..."
	docker-compose run --rm e2e-tests
	@echo "🛑 Parando serviços..."
	docker-compose down

test-e2e-quick: ## Executa testes E2E sem rebuild (mais rápido)
	@echo "🚀 Executando testes E2E (modo rápido)..."
	docker-compose run --rm e2e-tests

# Desenvolvimento
dev: ## Sobe ambiente de desenvolvimento
	@echo "🚀 Iniciando ambiente de desenvolvimento..."
	docker-compose up --build -d mongodb app
	@echo "✅ Ambiente pronto! Acesse http://localhost:8080"

dev-logs: ## Mostra logs em modo desenvolvimento
	docker-compose logs -f

# Limpeza
clean: ## Limpa arquivos temporários e cache
	go clean -cache -testcache
	find . -name "*.tmp" -delete
	find . -name "*.log" -delete

docker-clean: ## Limpa containers, volumes e imagens não utilizados
	@echo "🧹 Limpando Docker..."
	docker system prune -f
	docker volume prune -f
	docker image prune -f

# Geração de código
generate-mocks: ## Gera mocks para interfaces
	@echo "🔧 Gerando mocks..."
	mockery --dir=services --name=AuthService --output=mocks
	mockery --dir=services --name=OTPService --output=mocks
	mockery --dir=services --name=JWTService --output=mocks
	mockery --dir=repositories --name=UserRepository --output=mocks
	mockery --dir=repositories --name=OTPRepository --output=mocks
	mockery --dir=middlewares --name=CookieMiddleware --output=mocks
	@echo "✅ Mocks gerados!"

# Verificações
check: ## Executa verificações de código
	@echo "🔍 Executando verificações..."
	go vet ./...
	gofmt -d .
	golint ./... || true

# Produção
prod-build: ## Constrói para produção
	@echo "🏗️ Construindo para produção..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api

# Utilitários
status: ## Mostra status dos containers
	docker-compose ps

restart: ## Reinicia a aplicação
	docker-compose restart app

shell: ## Abre shell no container da aplicação
	docker-compose exec app sh

mongo-shell: ## Abre shell no MongoDB
	docker-compose exec mongodb mongosh

# Comando padrão
.DEFAULT_GOAL := help