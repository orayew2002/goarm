FROM golang:1.24.2 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
copy . .

# Build the Go app
RUN go build -o main cmd/app/main.go

# --- Run Stage ---
FROM alpine:latest

# Set working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose application port (change if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
