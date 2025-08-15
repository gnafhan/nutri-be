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

# Ensure .env file exists (copy from example if not)
RUN if [ ! -f .env ]; then cp .env.example .env || echo "No .env.example found"; fi

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

# Copy necessary files
# COPY --from=builder /app/.env .

# Switch to non-root user
USER appuser

# # Expose port
# EXPOSE 9097

# Command to run the application
ENTRYPOINT ["./nutribox-api"]