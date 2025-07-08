# Aetheris ID API

API de autenticação e autorização para o ecossistema Aetheris.

## 🚀 Início Rápido

### Pré-requisitos

- Docker e Docker Compose
- Go 1.24+ (para desenvolvimento local)

### Comandos Principais

```bash
# Ver todos os comandos disponíveis
make help

# Ambiente de desenvolvimento
make dev          # Sobe MongoDB + App
make dev-logs     # Mostra logs em tempo real
make down         # Para todos os serviços

# Testes
make test         # Executa todos os testes (unit + e2e)
make test-unit    # Apenas testes unitários
make test-e2e     # Apenas testes E2E
make test-e2e-quick # E2E sem rebuild (mais rápido)

# Docker
make build        # Constrói imagens
make up           # Sobe serviços
make logs         # Logs da aplicação
make logs-mongo   # Logs do MongoDB
make status       # Status dos containers
make restart      # Reinicia aplicação

# Desenvolvimento
make generate-mocks # Gera mocks para interfaces
make check         # Verificações de código
make clean         # Limpa cache e arquivos temporários
make docker-clean  # Limpa Docker (containers, volumes, imagens)

# Utilitários
make shell         # Shell no container da app
make mongo-shell   # Shell no MongoDB
```

## 🧪 Testes

### Testes Unitários

```bash
make test-unit
```

### Testes E2E

Os testes E2E rodam em containers Docker com banco real:

```bash
make test-e2e
```

**O que acontece:**

1. O Makefile executa o script `./scripts/run-e2e-tests.sh`
2. Se o script não existir, executa diretamente os comandos Docker
3. Para containers existentes
4. Constrói e sobe MongoDB + App
5. Aguarda serviços estarem prontos
6. Executa testes E2E
7. Para todos os serviços

### Testes E2E Rápidos

Se os containers já estiverem rodando:

```bash
make test-e2e-quick
```

## 🏗️ Estrutura do Projeto

```
├── cmd/api/                 # Ponto de entrada da aplicação
├── internal/
│   ├── bootstrap/          # Configuração de dependências
│   ├── configs/           # Configurações
│   ├── domain/            # Entidades e regras de negócio
│   ├── handlers/          # Handlers HTTP
│   ├── middlewares/       # Middlewares
│   ├── models/            # DTOs e modelos
│   ├── repositories/      # Acesso a dados
│   ├── server/            # Configuração do servidor
│   └── services/          # Lógica de negócio
├── pkg/                   # Pacotes públicos
├── test/e2e/             # Testes end-to-end
├── mocks/                # Mocks gerados
├── scripts/              # Scripts de automação
├── Dockerfile            # Imagem da aplicação
├── docker-compose.yml    # Orquestração de containers
└── Makefile              # Comandos de automação
```

## 🔧 Desenvolvimento

### Ambiente Local

```bash
# Sobe ambiente completo
make dev

# Acessa a aplicação
curl http://localhost:8080/health

# Logs em tempo real
make dev-logs
```

### Geração de Mocks

```bash
make generate-mocks
```

### Verificações de Código

```bash
make check
```

## 🐳 Docker

### Containers

- **app**: Aplicação Go (porta 8080)
- **mongodb**: Banco de dados (porta 27017)
- **e2e-tests**: Container para testes E2E

### Comandos Docker

```bash
make build        # Constrói todas as imagens
make up           # Sobe serviços
make down         # Para e remove containers
make logs         # Logs da aplicação
make status       # Status dos containers
```

## 🔐 Segurança

### Chaves ECDSA

O projeto usa chaves ECDSA para assinatura de JWT. As chaves são geradas automaticamente:

```bash
# Gerar chaves (se necessário)
openssl ecparam -genkey -name prime256v1 -noout -out ecdsa_private.pem
openssl ec -in ecdsa_private.pem -pubout -out ecdsa_public.pem
```

**⚠️ Importante:** Nunca faça commit das chaves privadas. Elas já estão no `.gitignore`.

### Variáveis de Ambiente

- Para desenvolvimento: use arquivo `.env` (não versionado)
- Para produção: use variáveis de ambiente do sistema

## 📝 API Endpoints

### Autenticação

- `POST /api/v1/auth/login` - Inicia processo de login
- `POST /api/v1/auth/authenticate` - Autentica com código OTP

### Health Check

- `GET /health` - Status da aplicação

## 🧹 Limpeza

### Limpeza de Desenvolvimento

```bash
make clean        # Limpa cache Go e arquivos temporários
```

### Limpeza Docker

```bash
make docker-clean # Remove containers, volumes e imagens não utilizados
```

## 📚 Comandos Úteis

### Desenvolvimento

```bash
make dev          # Ambiente completo
make dev-logs     # Logs em tempo real
make shell        # Shell na aplicação
make mongo-shell  # Shell no MongoDB
```

### Testes

```bash
make test         # Todos os testes
make test-unit    # Testes unitários
make test-e2e     # Testes E2E completos
```

### Manutenção

```bash
make build        # Reconstruir imagens
make restart      # Reiniciar aplicação
make status       # Ver status
make down         # Parar tudo
```

## 🤝 Contribuindo

1. Faça fork do projeto
2. Crie uma branch para sua feature
3. Execute os testes: `make test`
4. Faça commit das mudanças
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT.
