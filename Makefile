DB_URL=postgres://root:secret@localhost:5433/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres20 --network bank-network -p 5433:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d $(DB_URL)res:12-alpine

createdb:
	docker exec -it postgres20 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres20 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/SunnyChugh99/banking_management_golang/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml 

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 db_docs db_schema sqlc test server mock proto

