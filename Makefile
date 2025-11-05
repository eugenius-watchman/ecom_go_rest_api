build:
	@go build -o bin/ecom_go_rest_api cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom_go_rest_api