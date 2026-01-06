# Build stage
FROM dhi.io/golang:1.25-alpine3.23-dev AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM dhi.io/alpine-base:3.23

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=nonroot:nonroot /app/main .

# Expose port
EXPOSE 8080

USER nonroot

# Run the application
CMD ["./main"] 