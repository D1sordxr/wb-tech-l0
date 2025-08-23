FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/api ./cmd/api/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/api /app/api

COPY --from=builder /app/configs ./configs

RUN addgroup -S service && adduser -S service -G service
USER service

EXPOSE 8080

CMD ["/app/api"]