# Build stage
FROM golang:1.23-bullseye AS builder
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app

# Install necessary dependencies for Alpine
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

# Make sure the binary is executable
RUN chmod +x /app/main

EXPOSE 8080
CMD ["/app/main"]