# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o viacep-cloud-run ./cmd/api

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/viacep-cloud-run .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/viacep-cloud-run"]
