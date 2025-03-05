FROM golang:1.21-alpine AS builder

WORKDIR /app

# Только копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o api-gateway ./cmd/main.go

# Финальный образ на основе scratch (минимальный)
FROM alpine:latest

WORKDIR /app

# Копируем SSL-сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

# Копируем только скомпилированное приложение
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/.env .

# Создаём непривилегированного пользователя
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 8081

CMD ["./api-gateway"]
