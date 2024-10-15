Para rodar o projeto, rode da raíz do projeto:
```
docker compose up -d
migrate -path=internal/infra/database/migrations -database="mysql://root:root@tcp(localhost:3306)/orders" -verbose up
cd cmd/ordersystem && go run main.go wire_gen.go
```
Certifique-se que seu banco de dados possui o database orders, ou a migration falhará.
