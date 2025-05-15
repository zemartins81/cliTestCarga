FROM golang:1.24-alpine3.21

# Configurar o diretório de trabalho
WORKDIR /app

# Copiar os arquivos do projeto
COPY . .

# Baixe as dependências do módulo
RUN go mod tidy

# Copie o restante do código-fonte da aplicação para o diretório de trabalho
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o /app/loadtester .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/loadtester .

ENTRYPOINT ["/root/loadtester"]