# Clean Architecture Challenge

Este projeto implementa um sistema utilizando Clean Architecture, com suporte para REST, gRPC e GraphQL.

## Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.19 ou superior
- Ferramenta de migração de banco de dados (`migrate`)
  - Você pode instalar o `migrate` com o comando:
    ```bash
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

## Configuração e Execução

Siga os passos abaixo para rodar o projeto:

1. **Subir o ambiente com Docker**:

   Na raiz do projeto, execute o comando para subir os serviços com Docker:

   ```bash
   docker compose up -d

Execute o seguinte comando para rodar as migrações:
migrate -path=internal/infra/database/migrations -database="mysql://root:root@tcp(localhost:3306)/orders" -verbose up

execute o seguinte comando para iniciar a aplicação:
cd cmd/ordersystem 

>go run main.go wire_gen.go
