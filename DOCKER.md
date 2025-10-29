# Docker Deployment Guide

Complete guide for running ODT Image Replacer API in Docker with multi-stage builds and production-ready configuration.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

### Using Docker CLI

```bash
# Build the image
docker build -t odt-api:latest .

# Run the container
docker run -d \
  --name odt-api \
  -p 8080:8080 \
  -e GIN_MODE=release \
  odt-api:latest

# View logs
docker logs -f odt-api

# Stop the container
docker stop odt-api
docker rm odt-api
```

---

## Dockerfile Architecture

### Multi-Stage Build

The Dockerfile uses a **two-stage build** for optimal size and security:

#### Stage 1: Builder (golang:1.25.2-alpine)
- Compiles the Go application
- Downloads dependencies
- Creates static binary
- ~300MB+ image size

#### Stage 2: Production (alpine:latest)
- Runs the compiled binary
- Minimal dependencies
- Non-root user
- ~15MB final image size

### Build Process

```dockerfile
# Build Stage: Compiles Go code
FROM golang:1.25.2-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o odt-api ./cmd/api/main.go

# Production Stage: Runs the binary
FROM alpine:latest
COPY --from=builder /build/odt-api /app/odt-api
USER appuser
ENTRYPOINT ["/app/docker-entrypoint.sh"]
```

---

## Entrypoint Script

The `docker-entrypoint.sh` script handles:

1. ✅ **Environment validation** - Checks binary exists
2. ✅ **Configuration display** - Shows runtime settings
3. ✅ **Graceful shutdown** - Handles SIGTERM/SIGINT
4. ✅ **Process management** - Monitors API process
5. ✅ **Error handling** - Returns proper exit codes

### Features

```bash
#!/bin/sh
# Validates binary
# Sets up temp directories
# Traps signals for graceful shutdown
# Starts API server
# Waits for completion
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GIN_MODE` | `release` | Gin framework mode: `debug`, `release`, or `test` |
| `PORT` | `8080` | API server port |
| `HOST` | `0.0.0.0` | API server host |
| `TZ` | `UTC` | Timezone |

### Setting Environment Variables

**Docker CLI:**
```bash
docker run -d \
  -e GIN_MODE=debug \
  -e PORT=3000 \
  -e TZ=America/New_York \
  odt-api:latest
```

**Docker Compose:**
```yaml
environment:
  - GIN_MODE=debug
  - PORT=3000
  - TZ=America/New_York
```

---

## Port Configuration

### Default Port: 8080

**Change port mapping:**

```bash
# Map to port 3000 on host
docker run -p 3000:8080 odt-api:latest

# Map to port 80 (requires root/privileges)
docker run -p 80:8080 odt-api:latest
```

**Docker Compose:**
```yaml
ports:
  - "3000:8080"  # Host:Container
```

**Environment variable (`.env` file):**
```bash
API_PORT=3000
```

Then:
```yaml
ports:
  - "${API_PORT:-8080}:8080"
```

---

## Health Checks

The container includes a built-in health check:

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
```

**Check health status:**
```bash
docker ps  # Shows health status
docker inspect odt-api | grep -A 10 Health
```

**Docker Compose health check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 10s
```

---

## Resource Limits

Docker Compose includes resource limits for production:

```yaml
deploy:
  resources:
    limits:
      cpus: '1.0'        # Max 1 CPU core
      memory: 512M       # Max 512MB RAM
    reservations:
      cpus: '0.25'       # Min 0.25 CPU core
      memory: 128M       # Min 128MB RAM
```

**Adjust for your needs:**
- Light load: 256MB RAM, 0.5 CPU
- Medium load: 512MB RAM, 1.0 CPU
- Heavy load: 1GB RAM, 2.0 CPU

---

## Security Features

### 1. Non-Root User
```dockerfile
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
```

### 2. No New Privileges
```yaml
security_opt:
  - no-new-privileges:true
```

### 3. Read-Only Root Filesystem (Optional)
```yaml
read_only: true
tmpfs:
  - /tmp:size=100M,mode=1777
```

### 4. Static Binary
- No dynamic linking
- Reduced attack surface
- Smaller image size

---

## Building the Image

### Production Build

```bash
docker build -t odt-api:latest .
```

### Development Build

```bash
docker build \
  --build-arg GIN_MODE=debug \
  -t odt-api:dev .
```

### With Build Args

```bash
docker build \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  --build-arg VERSION=2.0.0 \
  -t odt-api:2.0.0 .
```

### Multi-Platform Build

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t odt-api:latest .
```

---

## Running in Production

### 1. Docker Compose (Recommended)

Create `.env` file:
```bash
API_PORT=8080
GIN_MODE=release
TZ=UTC
```

Start services:
```bash
docker-compose up -d
```

### 2. Docker Swarm

```bash
docker stack deploy -c docker-compose.yml odt-stack
```

### 3. Kubernetes

Create deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: odt-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: odt-api
        image: odt-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
```

---

## Logging

### View Logs

```bash
# Docker CLI
docker logs odt-api
docker logs -f odt-api        # Follow
docker logs --tail 100 odt-api  # Last 100 lines

# Docker Compose
docker-compose logs
docker-compose logs -f odt-api
```

### Log Configuration

Docker Compose includes log rotation:
```yaml
logging:
  driver: json-file
  options:
    max-size: "10m"    # Max 10MB per file
    max-file: "3"      # Keep 3 files
```

---

## Networking

### Bridge Network (Default)

```yaml
networks:
  odt-network:
    driver: bridge
```

### Connect Multiple Services

```yaml
services:
  odt-api:
    networks:
      - odt-network

  nginx:
    networks:
      - odt-network
```

### Access from Host

```bash
curl http://localhost:8080/health
```

### Access from Another Container

```bash
curl http://odt-api:8080/health
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs odt-api

# Check entrypoint script
docker run --rm --entrypoint /bin/sh odt-api:latest

# Verify binary
docker run --rm odt-api:latest ls -l /app/
```

### Health Check Failing

```bash
# Test health endpoint
docker exec odt-api curl http://localhost:8080/health

# Check if service is listening
docker exec odt-api netstat -tlnp
```

### Out of Memory

```bash
# Check current usage
docker stats odt-api

# Increase memory limit
docker run -m 1g odt-api:latest
```

### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Use different port
docker run -p 3000:8080 odt-api:latest
```

---

## Best Practices

### 1. Use Docker Compose
- Easier configuration
- Environment file support
- Service dependencies
- Easy scaling

### 2. Set Resource Limits
```yaml
deploy:
  resources:
    limits:
      memory: 512M
```

### 3. Enable Health Checks
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
```

### 4. Use .env Files
```bash
API_PORT=8080
GIN_MODE=release
```

### 5. Monitor Logs
```bash
docker-compose logs -f --tail=100
```

### 6. Regular Updates
```bash
docker-compose pull
docker-compose up -d
```

---

## Examples

### Example 1: Basic Development

```bash
docker run -d \
  --name odt-api-dev \
  -p 8080:8080 \
  -e GIN_MODE=debug \
  odt-api:latest
```

### Example 2: Production with SSL

```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./certs:/etc/nginx/certs
    depends_on:
      - odt-api

  odt-api:
    build: .
    environment:
      - GIN_MODE=release
```

### Example 3: High Availability

```bash
docker-compose up -d --scale odt-api=3
```

---

## Cleanup

```bash
# Stop and remove containers
docker-compose down

# Remove images
docker rmi odt-api:latest

# Remove everything (including volumes)
docker-compose down -v --rmi all

# Clean up system
docker system prune -a
```

---

## Summary

✅ **Multi-stage build** - Optimized image size
✅ **Entrypoint script** - Graceful startup/shutdown
✅ **Health checks** - Container monitoring
✅ **Non-root user** - Enhanced security
✅ **Resource limits** - Controlled resource usage
✅ **Environment variables** - Flexible configuration
✅ **Docker Compose** - Easy deployment
✅ **Production-ready** - All best practices included

Your API is ready to deploy in Docker! 🐳
