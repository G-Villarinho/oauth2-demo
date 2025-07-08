FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Construir a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Criar diretório para chaves
RUN mkdir -p /app/keys

# Imagem final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiar o binário da aplicação
COPY --from=builder /app/main .

COPY ecdsa_private.pem ecdsa_public.pem /app/

# Expor a porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 