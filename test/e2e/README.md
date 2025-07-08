# Testes End-to-End (E2E)

Este diretório contém testes E2E reais que testam a aplicação de ponta a ponta, incluindo:

- Aplicação rodando em container
- Banco de dados MongoDB real
- Requisições HTTP reais
- Criação e limpeza de dados de teste

## Estrutura

```
test/e2e/
├── auth_login_e2e_test.go  # Testes E2E para o endpoint de login
└── README.md               # Este arquivo
```

## Como Executar

### Opção 1: Usando o script automatizado (Recomendado)

```bash
./scripts/run-e2e-tests.sh
```

Este script:

1. Para containers existentes
2. Constrói e sobe MongoDB e aplicação
3. Executa os testes E2E
4. Para todos os serviços
5. Reporta o resultado

### Opção 2: Manual

1. **Subir serviços:**

   ```bash
   docker-compose up --build -d mongodb app
   ```

2. **Aguardar serviços estarem prontos:**

   ```bash
   sleep 10
   ```

3. **Executar testes:**

   ```bash
   docker-compose run --rm e2e-tests
   ```

4. **Parar serviços:**
   ```bash
   docker-compose down
   ```

### Opção 3: Local (sem Docker)

Se você quiser rodar os testes localmente (com aplicação e banco rodando):

```bash
# Certifique-se que a aplicação está rodando em localhost:8080
# e o MongoDB em localhost:27017

go test ./test/e2e -v
```

## Testes Disponíveis

### `TestLoginE2E`

Testa o endpoint `/api/v1/auth/login` com cenários:

- ✅ Login com email válido (usuário existe)
- ❌ Login com email inválido (usuário não existe)
- ❌ Email com formato inválido
- ❌ Email vazio
- ❌ Payload JSON inválido
- ❌ Método HTTP não permitido (GET)

### `TestLoginE2EWithRepository`

Testa integração com repositório:

- ✅ Cria usuário via repositório
- ✅ Faz login via HTTP
- ✅ Verifica resposta completa

### `BenchmarkLoginE2E`

Benchmark de performance do endpoint de login.

## Configuração

Os testes usam as seguintes variáveis de ambiente (definidas no `docker-compose.yml`):

- `MONGODB_CONNECTION_URI`: URI de conexão com MongoDB
- `MONGODB_DATABASE_NAME`: Nome do banco de dados de teste
- `API_BASE_URL`: URL base da API
- Outras configurações de JWT, CORS, etc.

## Banco de Dados de Teste

- **Database**: `aetheris_test`
- **Coleções**: `users`, `otps`, `clients`, `authorization_codes`, `refresh_tokens`
- **Limpeza**: Automática após cada teste
- **Isolamento**: Cada teste roda em ambiente isolado

## Dicas

1. **Logs**: Use `docker-compose logs app` para ver logs da aplicação
2. **Debug**: Use `docker-compose logs e2e-tests` para ver logs dos testes
3. **Banco**: Use MongoDB Express em `http://localhost:8081` para inspecionar dados
4. **Timeout**: Os testes têm timeout de 30 segundos por padrão

## Próximos Passos

Para adicionar mais testes E2E:

1. Crie um novo arquivo `test/e2e/[feature]_e2e_test.go`
2. Use as funções helper existentes (`setupTestDatabase`, `createTestUser`, etc.)
3. Siga o padrão de nomenclatura e estrutura
4. Adicione o novo teste ao `docker-compose.yml` se necessário

## Troubleshooting

### Erro: "connection refused"

- Verifique se o MongoDB está rodando
- Aguarde mais tempo para os serviços inicializarem

### Erro: "user not found"

- Verifique se o usuário foi criado corretamente
- Verifique logs da aplicação

### Erro: "timeout"

- Aumente o timeout no `docker-compose.yml`
- Verifique se a aplicação está respondendo
