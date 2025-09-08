# Default host (for local dev)
DB_HOST ?= localhost

# Run the app with DB_HOST injected
prod: tidy vet linter test
	DB_HOST=$(DB_HOST) go run cmd/app/main.go

run:
	DB_HOST=$(DB_HOST) go run cmd/app/main.go

build:
	DB_HOST=$(DB_HOST) go build -o bin/app cmd/app/main.go

docker-build:
	docker build \
		--build-arg DB_HOST=host.docker.internal \
		-t myapp:latest .

docker-run:
	docker run --rm \
		-e DB_HOST=host.docker.internal \
		-p 8080:8080 myapp:latest

# Everything else remains the same...
vet:
	go vet ./...

linter:
	@golangci-lint run ./...

tidy:
	@go mod tidy

test:
	@go test -v ./...

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

up:
	docker compose up --build -d

down:
	docker compose down
