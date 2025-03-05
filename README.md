# API Gateway для дипломного проекта "Интеллектуальная система оценки спроса на продукт"

## Обзор

API Gateway является центральным компонентом системы, который обеспечивает:
- Единую точку входа для всех клиентских запросов
- Аутентификацию пользователей через Auth Service
- Проксирование запросов к внутренним микросервисам
- Валидацию токенов с использованием кэширования для повышения производительности

## Структура проекта

```
api-gateway/
├── cmd/
│   └── main.go              # Точка входа в приложение
├── internal/
│   ├── handlers/            # HTTP обработчики
│   │   ├── auth.go          # Обработчики для /login и /register
│   │   └── analytics.go     # Обработчик для /analytics/demand
│   ├── middlewares/         # Middleware для аутентификации
│   │   └── auth.go
│   ├── grpc/                # gRPC клиент для Auth Service
│   │   └── client.go
│   └── cache/               # Кэш токенов
│       └── token_cache.go
├── proto/
│   ├── auth.proto           # Protobuf схема
│   ├── auth.pb.go           # Сгенерированный код (предполагается, что уже сгенерирован)
│   └── auth_grpc.pb.go      # Сгенерированный gRPC код (предполагается, что уже сгенерирован)
├── go.mod                   # Файл модуля
├── Dockerfile               # Конфигурация Docker
└── docker-compose.yml       # Конфигурация Docker Compose
```

## Требования

- Go 1.21 или выше
- Docker и Docker Compose (для запуска в контейнере)
- Доступ к Auth Service для аутентификации и валидации токенов

## Запуск приложения

### Локальный запуск

1. Клонировать репозиторий:
```bash
git clone https://github.com/yourusername/api-gateway.git
cd api-gateway
```

2. Создать файл `.env` со следующим содержимым:
```
AUTH_HTTP_ADDR=http://localhost:8080
AUTH_GRPC_ADDR=localhost:50051
API_GATEWAY_PORT=8081
```

3. Генерация gRPC-кода (если proto-файлы изменились):
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth.proto
```

4. Установка зависимостей и запуск:
```bash
go mod tidy
go run cmd/main.go
```

### Запуск с Docker Compose

```bash
docker-compose up --build
```

## Endpoints API

### Публичные эндпоинты

#### POST /register
Регистрация нового пользователя через Auth Service.

Пример запроса:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### POST /login
Аутентификация пользователя и получение JWT-токена.

Пример запроса:
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Пример ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Защищенные эндпоинты (требуют токен)

#### GET /analytics/demand
Получение данных о спросе на продукт.

Заголовок запроса:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Пример ответа:
```json
{
  "data": "Some demand data"
}
```

## Масштабирование и производительность

- Простая реализация кэша токенов использует `map`, что подходит для MVP
- Для продакшена рекомендуется использовать Redis для распределенного кэширования
- Для повышения отказоустойчивости добавить обработку ошибок с ретраями и circuit breaker

## Расширение

Для добавления новых защищенных эндпоинтов используйте группу `protected` в файле `cmd/main.go`.

## Логирование

В текущей реализации используется базовое логирование. Для продакшена рекомендуется:
- Использовать структурированный логгер, например, `zap`
- Настроить уровни логирования
- Добавить трассировку запросов

## Тестирование

Для тестирования API используйте прилагаемый файл `api-tests.rest` (см. документацию ниже).
