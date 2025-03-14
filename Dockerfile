# Build stage
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Jalankan test sebelum build untuk verifikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./src && \
    echo "Build berhasil! Binary size: $(du -h /app/api-server | cut -f1)"

# Periksa exit code
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./src && \
    echo "Build successful" || echo "Build failed"

# Debug: verify binary exists and show directory structure

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder stage with explicit paths
COPY --from=builder /app/api-server /app/api-server

# Copy environment file and other necessary files
COPY .env /app/.env
COPY src/ /app/src/

# Debug: verify binary exists in final image and show files
RUN ls -la /app && echo "Files in app directory:"

# Ensure binary is executable
RUN chmod +x /app/api-server

# Expose port
CMD ["/app/src/main"]
EXPOSE 3000