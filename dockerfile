# Usando uma imagem base do Debian mais recente
FROM debian:bookworm-slim AS builder

# Instalar dependências
RUN apt-get update && apt-get install -y wget gcc make

# Baixar e instalar o Go
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz && \
    rm go1.23.0.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"

# Criar diretório de trabalho
WORKDIR /app

# Copiar arquivos do projeto para o contêiner
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar o binário
RUN go build -o main .

# Fase final para otimizar o contêiner
FROM debian:bookworm-slim

# Diretório de trabalho
WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /app/main .

# Variáveis de ambiente
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

# Porta que a aplicação escuta
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
