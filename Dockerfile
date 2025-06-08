# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"] 