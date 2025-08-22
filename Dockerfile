# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install dependencies for build (if needed)
RUN apk add --no-cache git

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/api ./cmd/api/main.go

# Final lightweight image
FROM alpine:3.18

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/api /app/api

# Copy configs
COPY --from=builder /app/configs ./configs

# Create a non-root user and switch to it
RUN addgroup -S service && adduser -S service -G service
USER service

EXPOSE 8080

CMD ["/app/api"]