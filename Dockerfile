# Multi-stage Dockerfile for ODT Image Replacer API
# Build stage: Compiles the Go application
# Production stage: Runs the application in a minimal container

###################
# Build Stage
###################
FROM golang:1.25.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0: Build static binary
# -ldflags: Strip debug info and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /build/odt-api \
    ./cmd/api/main.go

# Verify the binary was built
RUN ls -lh /build/odt-api

###################
# Production Stage
###################
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/odt-api /app/odt-api

# Copy entrypoint script
COPY --from=builder /build/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

# Create temp directory for ODT processing
RUN mkdir -p /tmp/odt-temp && \
    chown -R appuser:appuser /tmp/odt-temp

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release \
    PORT=8080 \
    HOST=0.0.0.0

# Use entrypoint script
ENTRYPOINT ["/app/docker-entrypoint.sh"]

# Default command (can be overridden)
CMD ["-port=8080", "-host=0.0.0.0", "-mode=release"]
