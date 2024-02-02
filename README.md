# Sistema de autenticação e autorização para Golang com e-mail e senha

Esta aplicação em GoLang foi desenvolvida com o propósito de fornecer um exemplo prático de implementação de um sistema de autenticação e autorização, incluindo validação de e-mail. A aplicação utiliza e-mail e senha como credenciais para autenticação.

## Tecnologias usadas

- [Go](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [SQLC](https://sqlc.dev/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)

## Requisitos

Certifique-se de ter os seguintes itens instalados antes de iniciar:

- Go (versão utilizada no projeto: 1.21.6)
- Docker
- Imagem do PostgreSQL baixada
- SQLC (foi instalado na maquina)

## Instalação e Execução

Clone o repositório:

```bash
git clone https://github.com/dariomatias-dev/go_auth

cd go_auth
```

Abra o arquivo `.env_example`, remova do seu nome `_example`, depois preencha os campos que estão faltando.

Instale as depedências:

```bash
go mod download
```

Crie o container com o banco de dados:

```bash
docker-compose up -d
```

Obs.: Se estiver no Linux e tiver dado erro, provavelmente irá precisar colocar `sudo` antes do comando.

Rode a aplicação:

```bash
go run main.go
```

## Endpoints

A API fornece os seguintes endpoints:

### Rotas de autenticação

- **POST (login)**: Login
- **GET (refresh)**: atualização dos tokens
- **Post (validate-email)**: validação de email

### Rotas de usuário

- **POST (user)**: criação de usuário
- **GET (user/:id)**: obtenção de usuário
- **GET (users)**: obtenção dos usuários
- **PATCH (user/:id)**: atualização de usuário
- **DELETE (user/:id)**: exclusão de usuário

## Licença

Este projeto é licenciado sob a Licença MIT.
