# Aetheris ID - Servidor de Identidade OAuth2

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com/)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.0+-orange.svg)](https://www.mongodb.com/)

OAUTH-DEMO é um servidor de identidade robusto e seguro que implementa o protocolo OAuth 2.0 com suporte a PKCE (Proof Key for Code Exchange), OpenID Connect e autenticação baseada em OTP (One-Time Password).

## 🚀 Características

- **OAuth 2.0 Authorization Code Flow** com PKCE
- **OpenID Connect** para autenticação federada
- **Autenticação OTP** via email
- **Refresh Tokens** para renovação automática
- **JWT Tokens** (Access Token, ID Token, OTP Token)
- **MongoDB** como banco de dados
- **Arquitetura Limpa** com injeção de dependência
- **Testes Unitários** abrangentes
- **Docker Compose** para desenvolvimento

## 🏗️ Arquitetura

O projeto segue uma arquitetura limpa com as seguintes camadas:

```
├── cmd/api/           # Ponto de entrada da aplicação
├── configs/           # Configurações e variáveis de ambiente
├── internal/          # Código interno da aplicação
│   ├── api/          # Definições de rotas
│   ├── bootstrap/    # Configuração de dependências
│   ├── domain/       # Entidades e regras de negócio
│   ├── handlers/     # Controladores HTTP
│   ├── middlewares/  # Middlewares de autenticação
│   ├── models/       # Modelos de dados
│   ├── repositories/ # Acesso a dados
│   ├── server/       # Configuração do servidor
│   └── services/     # Lógica de negócio
├── pkg/              # Pacotes públicos reutilizáveis
└── infra/            # Infraestrutura (banco de dados)
```

## 🔐 Fluxo OAuth2

O Aetheris ID implementa o fluxo OAuth 2.0 Authorization Code com PKCE, conforme ilustrado na imagem abaixo:

![Fluxo OAuth2](public/oauth2-fluxo-example.png)

### Etapas do Fluxo:

1. **Autorização**: O cliente redireciona o usuário para `/oauth/authorize`
2. **Login**: Se não autenticado, o usuário é redirecionado para login
3. **Consentimento**: Após login, o usuário autoriza o acesso
4. **Callback**: O servidor redireciona para o cliente com um código de autorização
5. **Troca de Token**: O cliente troca o código por tokens de acesso
6. **Acesso**: O cliente usa o access token para acessar recursos protegidos

## 🛠️ Tecnologias

- **Go 1.24.2+** - Linguagem principal
- **Echo v4** - Framework web
- **MongoDB** - Banco de dados
- **JWT** - Tokens de autenticação
- **Dig** - Injeção de dependência
- **Docker Compose** - Orquestração de containers

## 📋 Pré-requisitos

- Go 1.24.2 ou superior
- Docker e Docker Compose
- OpenSSL (para geração de chaves)

## 🚀 Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/G-Villarinho/oauth2-demo.git
cd aetheris-id
```

### 2. Gere as chaves ECDSA

```bash
make generate-key
```

### 3. Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Servidor
PORT=5001
ENV=development

# URLs
API_BASE_URL=http://localhost:5001
CLIENT_LOGIN_URL=http://localhost:3000/login

# MongoDB
MONGODB_URI=mongodb://admin:password123@localhost:27017/aetheris?authSource=admin

# Segurança
ACCESS_TOKEN_EXPIRATION_HOURS=1
REFRESH_TOKEN_EXPIRATION_HOURS=24
ID_TOKEN_EXPIRATION_MINUTES=15
OTP_JWT_EXPIRATION_MINUTES=5

# OTP
OTP_EXPIRATION_MINUTES=5
OTP_RESEND_COOLDOWN_MINUTES=1
```

### 4. Inicie o MongoDB

```bash
make up
```

### 5. Execute a aplicação

```bash
go run cmd/api/main.go
```

## 📚 Uso

### Criando um Cliente OAuth2

```bash
curl -X POST http://localhost:5001/api/v1/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Minha Aplicação",
    "description": "Aplicação web de exemplo",
    "redirect_uris": ["http://localhost:3000/callback"],
    "grant_types": ["authorization_code", "refresh_token"]
  }'
```

### Fluxo de Autorização

1. **Redirecionar para autorização**:

```
GET /api/v1/oauth/authorize?
  client_id=seu_client_id&
  redirect_uri=http://localhost:3000/callback&
  response_type=code&
  scope=openid profile email&
  state=xyz&
  code_challenge=challenge&
  code_challenge_method=S256
```

2. **Trocar código por token**:

```bash
curl -X POST http://localhost:5001/api/v1/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "code=authorization_code&code_verifier=verifier&client_id=seu_client_id&redirect_uri=http://localhost:3000/callback"
```

### Autenticação de Usuários

```bash
# Enviar código de verificação
curl -X POST http://localhost:5001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario@exemplo.com"}'

# Autenticar com código
curl -X POST http://localhost:5001/api/v1/auth/authenticate \
  -H "Content-Type: application/json" \
  -d '{"code": "123456"}'
```

## 🧪 Testes

Execute os testes unitários:

```bash
go test ./...
```

Execute testes com cobertura:

```bash
go test -cover ./...
```

## 🐳 Docker

### Desenvolvimento

```bash
# Iniciar MongoDB
make up

# Verificar status
make status

# Ver logs
make logs

# Parar containers
make down
```

### Produção

```bash
# Build da imagem
docker build -t aetheris-id .

# Executar container
docker run -p 5001:5001 aetheris-id
```

## 📖 API Reference

### Endpoints OAuth2

- `GET /api/v1/oauth/authorize` - Iniciar fluxo de autorização
- `POST /api/v1/oauth/token` - Trocar código por token

### Endpoints de Autenticação

- `POST /api/v1/auth/login` - Enviar código de verificação
- `POST /api/v1/auth/authenticate` - Autenticar com código
- `POST /api/v1/auth/register` - Registrar novo usuário
- `POST /api/v1/auth/code/resend` - Reenviar código

### Endpoints de Clientes

- `POST /api/v1/clients` - Criar novo cliente OAuth2

## 🔧 Configuração Avançada

### Variáveis de Ambiente

| Variável                         | Descrição                         | Padrão        |
| -------------------------------- | --------------------------------- | ------------- |
| `PORT`                           | Porta do servidor                 | `5001`        |
| `ENV`                            | Ambiente (development/production) | `development` |
| `MONGODB_URI`                    | URI de conexão MongoDB            | -             |
| `ACCESS_TOKEN_EXPIRATION_HOURS`  | Expiração do access token         | `1`           |
| `REFRESH_TOKEN_EXPIRATION_HOURS` | Expiração do refresh token        | `24`          |
| `ID_TOKEN_EXPIRATION_MINUTES`    | Expiração do ID token             | `15`          |

### Segurança

- **PKCE**: Implementado para prevenir ataques de interceptação
- **JWT**: Tokens assinados com ECDSA
- **OTP**: Códigos de uso único com expiração
- **HTTPS**: Recomendado para produção

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request
