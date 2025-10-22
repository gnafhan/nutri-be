# Multi-stage Dockerfile for Nutribox API (Go + Fiber)
# Optimized for Google Cloud Run

# ---------- Build stage ----------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build deps
RUN apk add --no-cache git ca-certificates tzdata

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build static binary
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /app/server ./src

# ---------- Runtime stage ----------
# Distroless base includes CA certificates and a minimal userspace
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy timezone data and certs (defensive; base already has certs)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Copy built binary
COPY --from=builder /app/server /app/server

# Static files directory (served at /uploads). Note: Cloud Run FS is read-only
# except for /tmp. If your app writes uploads, redirect writes to /tmp/uploads
# or use Cloud Storage. This layer only provides an empty directory for reads.
COPY uploads /app/uploads

# Cloud Run expects the server to listen on $PORT
ENV PORT=8080 \
    APP_ENV=prod \
    APP_HOST=0.0.0.0

# If you use Viper to read from a .env file, mount it to /app/.env at deploy time
# using Cloud Run Secret Manager volumes, e.g.:
#   --mount=type=secret,source=nutribox-env,target=/app/.env,mode=0440

# Use an unprivileged user provided by distroless
USER 65532:65532

EXPOSE 8080

# Start the server
ENTRYPOINT ["/app/server"]
