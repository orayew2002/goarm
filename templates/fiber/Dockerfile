# --- Build Stage ---
FROM golang:1.24.2 AS builder

# Disable CGO and set target OS/architecture
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy module files and download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o cmd/bin/main cmd/app/main.go

# --- Run Stage ---
FROM alpine:latest

# Set working directory inside the runtime container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app .

# Expose the application's port (update if needed)
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./cmd/bin/main"]
