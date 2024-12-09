# Usar uma imagem base do Golang
FROM golang:1.20

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos do projeto para o diretório de trabalho
COPY . .

# Baixar as dependências e compilar o projeto
RUN go mod tidy && go build -o main .

# Expor a porta que o servidor usa
EXPOSE 8080

# Comando para rodar o servidor
CMD ["./main"]
