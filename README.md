# goProcessClient

API em Go para autenticação e roteamento de dados  para  (modelo genérico para apps mobile).

## Tecnologias

- Go 1.24+
- MySQL (`github.com/go-sql-driver/mysql`)
- JWT (`github.com/golang-jwt/jwt/v5`)

## Estrutura

- `cmd/api` — ponto de entrada da API (main)
- `internal/config` — configurações (endereço do servidor, DSN, segredo JWT)
- `internal/data` — acesso a dados (conexão MySQL)
- `internal/domain` — modelos de domínio (usuário, login, config, etc.)
- `internal/http/handlers/auth` — handlers de autenticação (login, health, JWT)

## Rodando localmente

```bash
go run ./cmd/api
