postgres:
	docker run --name postgres17 --network bank-network -p 5433:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres17 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:vOe5Go2oTKtBIFSE60I7@simple-bank.cfaqkm6kysnc.eu-west-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 

migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/SunnyChugh99/banking_management_golang/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1


