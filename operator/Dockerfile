# Use multistage build for smaller final image
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o operator .

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/operator .

# Set binary as entrypoint
ENTRYPOINT ["./operator"]