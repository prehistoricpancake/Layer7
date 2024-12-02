# Build stage
FROM golang:1.23-alpine AS builder

# Install git for downloading dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat-server

# Final stage
FROM alpine:latest

# Add timezone data and certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/chat-server .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./chat-server"]