prod: tidy vet linter test
	go run cmd/app/main.go

vet:
	go vet ./...

linter:
	@golangci-lint run ./...

tidy:
	@go mod tidy

test:
	@go test -v ./...

run:
	go run cmd/app/main.go

build:
	go build -o bin/app cmd/app/main.go

clean:
	@rm -rf bin

fmt:
	@gofmt -s -w .

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

doc:
	@go doc ./...

deps:
	@go mod download

docker-build:
	docker build -t myapp:latest .

docker-run:
	docker run --rm -p 8080:8080 myapp:latest

up:
	docker compose up --build -d

down:
	docker compose down
