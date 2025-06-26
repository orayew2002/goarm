prod: tidy vet linter test
	go run cmd/app/main.go

vet:
	go vet ./...

linter:
	@golangci-lint run ./...

tidy:
	@go mod tidy

test: 
	@go test ./...

run:
	go run cmd/app/main.go
