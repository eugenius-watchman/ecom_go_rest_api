build:
	@go build -o bin/ecom_go_rest_api cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom_go_rest_api

# migration:
# 	@migration create

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down