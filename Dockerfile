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

# Build the application - changed build path
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src

# Debug: verify binary exists
RUN ls -la

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder stage
COPY --from=builder /app/main .

# Debug: verify binary exists in final image
RUN ls -la

# Ensure binary is executable
RUN chmod +x ./main

COPY .env .env

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]