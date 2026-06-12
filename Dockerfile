# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bagogo-api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bagogo-migrate ./cmd/migrate

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

# Copy binary from builder
COPY --from=builder /app/bagogo-api .
COPY --from=builder /app/bagogo-migrate .

# Create uploads directory
RUN mkdir -p ./uploads

# Create non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 8001

# Health check (HTTPS with self-signed cert support)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --no-check-certificate --tries=1 --spider https://localhost:8001/health || exit 1

ENTRYPOINT ["./bagogo-api"]
