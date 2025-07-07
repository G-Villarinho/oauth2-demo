# Makefile para gerenciar containers Docker do MongoDB

.PHONY: help up down start stop restart logs status clean clean-all

# Comando padrão - mostra ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make up      - Subir containers em background"
	@echo "  make start   - Subir containers em background"
	@echo "  make down    - Parar e remover containers"
	@echo "  make stop    - Parar containers"
	@echo "  make restart - Reiniciar containers"
	@echo "  make logs    - Ver logs dos containers"
	@echo "  make status  - Ver status dos containers"
	@echo "  make clean   - Parar containers e remover volumes"
	@echo "  make clean-all - Parar containers, remover volumes e imagens"

# Subir containers em background
up:
	@echo "🚀 Subindo containers MongoDB..."
	docker-compose up -d
	@echo "✅ Containers iniciados!"
	@echo "📊 MongoDB disponível em: localhost:27017"
	@echo "🌐 MongoDB Express disponível em: http://localhost:8081"

# Alias para up
start: up

# Parar e remover containers
down:
	@echo "🛑 Parando containers..."
	docker-compose down
	@echo "✅ Containers parados e removidos!"

# Parar containers (sem remover)
stop:
	@echo "⏸️  Parando containers..."
	docker-compose stop
	@echo "✅ Containers parados!"

# Reiniciar containers
restart:
	@echo "🔄 Reiniciando containers..."
	docker-compose restart
	@echo "✅ Containers reiniciados!"

# Ver logs dos containers
logs:
	@echo "📋 Logs dos containers:"
	docker-compose logs -f

# Ver logs apenas do MongoDB
logs-mongo:
	@echo "📋 Logs do MongoDB:"
	docker-compose logs -f mongodb

# Ver logs apenas do MongoDB Express
logs-express:
	@echo "📋 Logs do MongoDB Express:"
	docker-compose logs -f mongo-express

# Ver status dos containers
status:
	@echo "📊 Status dos containers:"
	docker-compose ps

# Limpar containers e volumes (cuidado - apaga dados)
clean:
	@echo "🧹 Limpando containers e volumes..."
	docker-compose down -v
	@echo "✅ Containers e volumes removidos!"

# Limpeza completa (containers, volumes e imagens)
clean-all:
	@echo "🧹 Limpeza completa..."
	docker-compose down -v --rmi all
	@echo "✅ Tudo removido!"

# Verificar se containers estão rodando
check:
	@echo "🔍 Verificando status dos containers..."
	@if docker-compose ps | grep -q "Up"; then \
		echo "✅ Containers estão rodando!"; \
		echo "📊 MongoDB: localhost:27017"; \
		echo "🌐 MongoDB Express: http://localhost:8081"; \
	else \
		echo "❌ Containers não estão rodando. Execute 'make up' para iniciar."; \
	fi

# Conectar ao MongoDB via shell
shell:
	@echo "🐚 Conectando ao MongoDB shell..."
	docker-compose exec mongodb mongosh -u admin -p password123

# Backup do banco
backup:
	@echo "💾 Criando backup do MongoDB..."
	mkdir -p backups
	docker-compose exec mongodb mongodump --username admin --password password123 --authenticationDatabase admin --db aetheris --out /data/backup
	docker cp mongodb:/data/backup ./backups/$(shell date +%Y%m%d_%H%M%S)
	@echo "✅ Backup criado em ./backups/"

# Restaurar backup
restore:
	@echo "📥 Restaurando backup do MongoDB..."
	@if [ -z "$(BACKUP_DIR)" ]; then \
		echo "❌ Especifique o diretório do backup: make restore BACKUP_DIR=./backups/20231201_120000"; \
		exit 1; \
	fi
	docker cp $(BACKUP_DIR) mongodb:/data/restore
	docker-compose exec mongodb mongorestore --username admin --password password123 --authenticationDatabase admin /data/restore
	@echo "✅ Backup restaurado!" 