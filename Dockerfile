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

# Build the application with explicit output path
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./src

# Debug: verify binary exists and show directory structure
RUN ls -la /app && echo "Binary exists: $(test -f /app/api-server && echo YES || echo NO)"

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder stage with explicit paths
COPY --from=builder /app/api-server /app/api-server

# Debug: verify binary exists in final image
RUN ls -la /app && echo "Binary exists: $(test -f /app/api-server && echo YES || echo NO)"

# Ensure binary is executable
RUN chmod +x /app/api-server

# Copy environment file
COPY .env /app/.env

# Expose port
EXPOSE 3000

# Run the application with explicit path
CMD ["/app/api-server"]