# Documentação do Projeto: API de Gerenciamento de Livros com Docker

Esta API foi projetada para gerenciar um catálogo de livros, permitindo autenticação via JWT (JSON Web Tokens) e controle de acesso baseado em papéis (roles). Entre as funcionalidades principais, estão:
- Middleware para autenticação JWT e validação de papéis de usuário;
- Operações CRUD para livros;
- Adição e cálculo da média de avaliações de livros;
- Implementação simples e extensível utilizando o framework Gin para Go.

Abaixo, encontram-se os passos para configurar o ambiente e utilizar os endpoints da API.

## Introdução
Esta documentação detalha como configurar, rodar e interagir com a API de gerenciamento de livros. A API utiliza autenticação JWT para proteger as rotas e oferece funcionalidades como listagem de livros, avaliação, adição e consulta de média de avaliações.

A configuração do ambiente é feita utilizando **Docker**, simplificando o processo de criação e execução do projeto.

---

## Configuração do Docker

O Docker é utilizado neste projeto para simplificar o processo de implantação. Ele permite criar um ambiente isolado e replicável, garantindo que a API funcione de forma consistente em diferentes sistemas operacionais e máquinas. Com o Docker, você pode facilmente construir, empacotar e executar a aplicação sem se preocupar com dependências externas.

### 1. Construir a Imagem Docker

Execute o seguinte comando para construir a imagem Docker com base no `Dockerfile`:

```bash
docker build -t meu-projeto-jwt .
```

**Descrição:**
- `docker build`: Constrói a imagem Docker.
- `-t meu-projeto-jwt`: Nomeia a imagem como `meu-projeto-jwt`.
- `.`: Indica que o contexto de construção é o diretório atual.

### 2. Rodar o Container

Após criar a imagem, execute este comando para iniciar o container:

```bash
docker run -p 8080:8080 meu-projeto-jwt
```

**Descrição:**
- `docker run`: Inicia um container.
- `-p 8080:8080`: Mapeia a porta do container para a porta 8080 do seu computador.
- `meu-projeto-jwt`: Nome da imagem criada.

---

## Comandos da API

### 1. Fazer Login para Obter um Token JWT

Este endpoint é usado para autenticar um usuário e retornar um token JWT. Esse token é necessário para acessar as rotas protegidas da API.

**Requisição:**
```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"usuario_id": "joanderson", "senha": "123456"}'
```

**Resposta esperada:**
```json
{
  "token": "SEU_TOKEN_JWT_AQUI"
}
```

---

### 2. Listar Todos os Livros Disponíveis

**Requisição:**
```bash
curl -X GET http://localhost:8080/api/livros \
-H "Authorization: Bearer SEU_TOKEN_JWT_AQUI"
```

**Resposta esperada:**
```json
[
  {
    "id": "1",
    "title": "1984",
    "author": "George Orwell",
    "price": 29.9,
    "rating": 0,
    "Ratings": [
      4.5,
      4.8,
      5
    ]
  },
  ...
]
```

---

### 3. Avaliar um Livro

**Requisição:**
```bash
curl -X POST http://localhost:8080/api/livro/1/rate \
-H "Authorization: Bearer SEU_TOKEN_JWT_AQUI" \
-H "Content-Type: application/json" \
-d '{"rating": 5}'
```

**Resposta esperada:**
```json
{
  "livro": {
    "id": "1",
    "title": "1984",
    "author": "George Orwell",
    "price": 29.9,
    "rating": 4.825,
    "Ratings": [4.5, 4.8, 5, 5]
  },
  "message": "Avaliação registrada"
}
```

---

### 4. Consultar a Média de Avaliações de um Livro

**Requisição:**
```bash
curl -X GET http://localhost:8080/api/livros/1 \
-H "Authorization: Bearer SEU_TOKEN_JWT_AQUI"
```

**Resposta esperada:**
```json
{
  "livro_id": "1",
  "average_rating": 4.4
}
```

---

### 5. Adicionar um Novo Livro (Somente para Admins)

**Nota:** A rota `/api/livros` exige privilégios de administrador para adicionar um novo livro. O middleware `autenticarRole` é utilizado no código para verificar o papel do usuário antes de permitir o acesso a essa funcionalidade. Apenas usuários autenticados com o papel `admin` têm permissão para executar esta operação.

**Login como Admin:**
```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"usuario_id": "admin", "senha": "senha123"}'
```

**Adicionar Livro:**
```bash
curl -X POST http://localhost:8080/api/livros \
-H "Authorization: Bearer SEU_TOKEN_JWT_ADMIN" \
-H "Content-Type: application/json" \
-d '{"id": "13", "title": "A Teoria de Tudo", "author": "Stephen Hawking", "price": 45.99, "rating": 0}'
```

**Resposta esperada:**
```json
{
  "id": "13",
  "title": "A Teoria de Tudo",
  "author": "Stephen Hawking",
  "price": 45.99,
  "rating": 0,
  "Ratings": []
}
```

---

## Notas Finais

- Certifique-se de que o Docker esteja instalado e funcionando corretamente antes de executar os comandos.
- Substitua `SEU_TOKEN_JWT_AQUI` pelo token obtido durante o login.
- Este projeto suporta operações básicas de CRUD e utiliza autenticação JWT para proteção de rotas sensíveis.
- Para erros comuns, como o servidor não iniciar ou falhas de autenticação, verifique se as dependências foram instaladas corretamente e se o arquivo `.env` está configurado com as variáveis apropriadas.
- Use logs de erro exibidos no console para identificar problemas e depurar o servidor de maneira eficiente.

Para mais informações, consulte a documentação da API interna ou os códigos-fontes no repositório.

