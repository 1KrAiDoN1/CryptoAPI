# # Этап сборки
# FROM golang:1.23.1-alpine AS builder

# # Создаем и переходим в рабочую директорию
# WORKDIR /cryptoapi
# # Создаем структуру директорий
# RUN mkdir -p /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/
# RUN mkdir -p /Users/pavelvasilev/Desktop/CryptoAPI/internal/database/
# # Копируем только необходимые файлы для зависимостей
# COPY go.mod go.sum ./
# RUN go mod download
# COPY *.env ./
# COPY internal/database/DB_Config.env ./internal/database/DB_Config.env
# COPY internal/database/SecretHash.env ./internal/database/SecretHash.env


# # Копируем шаблоны
# COPY internal/database/ /Users/pavelvasilev/Desktop/CryptoAPI/internal/database/
# COPY pkg/templates/ /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/
# # COPY --from=builder /cryptoapi/migrations ./migrations

# # Проверяем, что файлы скопировались
# RUN ls -la /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/

# # Копируем остальные файлы
# COPY . .

# # Собираем приложение
# RUN go build -o main ./cmd/main.go
# # RUN chmod +x wait-for.sh


# CMD ["./main"]

# Этап сборки
FROM golang:1.23.1-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /cryptoapi

# Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем необходимые файлы окружения (убедитесь, что они есть в проекте)
COPY internal/database/DB_Config.env ./internal/database/DB_Config.env
COPY internal/database/SecretHash.env ./internal/database/SecretHash.env

# Копируем остальные исходники и шаблоны
COPY internal/database/ ./internal/database/
COPY pkg/templates/ ./pkg/templates/
COPY . .

# Собираем бинарь
RUN go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /cryptoapi

# Копируем бинарь из builder
COPY --from=builder /cryptoapi/main .

# Копируем необходимые файлы окружения и шаблоны, если нужны в runtime
COPY --from=builder /cryptoapi/internal/database/DB_Config.env ./internal/database/DB_Config.env
COPY --from=builder /cryptoapi/internal/database/SecretHash.env ./internal/database/SecretHash.env
COPY --from=builder /cryptoapi/pkg/templates/ ./pkg/templates/

# Устанавливаем необходимые зависимости для запуска, если нужны (например, ca-certificates)
RUN apk add --no-cache ca-certificates

CMD ["./main"]





