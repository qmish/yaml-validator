# Сборка (версия Go должна соответствовать go.mod)
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o yaml-validator .

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /workspace
COPY --from=builder /app/yaml-validator /usr/local/bin/yaml-validator
ENTRYPOINT ["yaml-validator"]
