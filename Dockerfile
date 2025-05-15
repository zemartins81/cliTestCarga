# Estágio 1: Imagem de construção
# Usando a imagem oficial do Go com Alpine como base para construir o aplicativo
FROM golang:1.24-alpine3.21 AS builder

# Configurar o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos do projeto para o diretório de trabalho
COPY . .

# Baixar e instalar todas as dependências necessárias do projeto
# go mod tidy remove dependências não utilizadas e adiciona as que faltam
RUN go mod tidy

# NOTA: Esta linha é redundante, pois já copiamos todos os arquivos acima
# Pode ser removida para otimizar o build
COPY . .

# Compilar o aplicativo com flags específicas:
# - CGO_ENABLED=0: desativa CGO para criar um binário estático
# - GOOS=linux: compila especificamente para Linux
# - ldflags "-s -w": remove informações de debug para reduzir o tamanho do binário
# - Salva o binário compilado como 'loadtester'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o /app/loadtester .

# Estágio 2: Imagem de produção
# Usando Alpine como imagem base mínima para executar o aplicativo
FROM alpine:latest

# Configurar o diretório de trabalho para a imagem de produção
WORKDIR /root/

# Copiar apenas o binário compilado do estágio anterior (builder)
# Isso reduz significativamente o tamanho da imagem final
COPY --from=builder /app/loadtester .

# Definir o comando que será executado quando o container iniciar
ENTRYPOINT ["/root/loadtester"]