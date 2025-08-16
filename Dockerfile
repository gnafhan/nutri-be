# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Ensure .env file exists for Docker
RUN if [ ! -f .env ]; then \
    if [ -f .env.docker ]; then \
        cp .env.docker .env; \
    elif [ -f .env.example ]; then \
        cp .env.example .env; \
    else \
        echo "No .env file found"; \
    fi \
fi

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nutribox-api ./src/main.go

# Final stage
FROM alpine:latest

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" > /etc/timezone

# Create non-root user
RUN adduser -D -g '' appuser

# Create necessary directories
RUN mkdir -p /app/uploads
RUN chown -R appuser:appuser /app
RUN chmod -R 777 /app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/nutribox-api .

# Copy the .env file from builder stage
COPY --from=builder /app/.env .

# Copy necessary files
# COPY --from=builder /app/.env .

# Switch to non-root user
USER appuser

# # Expose port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["./nutribox-api"]