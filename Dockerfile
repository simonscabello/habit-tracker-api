# Build stage
FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Definir variáveis de ambiente para o build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação
RUN go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:3.19

# Instalar certificados CA e timezone
RUN apk add --no-cache ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -g '' appuser

WORKDIR /app

# Copiar o binário compilado do stage anterior
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Mudar para o usuário não-root
USER appuser

# Labels para documentação
LABEL maintainer="Simon Scabello" \
      version="1.0" \
      description="Habit Tracker API"

# Expor a porta da aplicação
EXPOSE 8000

# Comando para executar a aplicação
CMD ["./main"]
