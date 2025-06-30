# # Этап сборки
# FROM golang:1.23.1-alpine AS builder

# # Создаем и переходим в рабочую директорию
# WORKDIR /cryptoapi
# # Создаем структуру директорий
# RUN mkdir -p /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/

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

# Создаем и переходим в рабочую директорию
WORKDIR /cryptoapi

# Копируем только необходимые файлы для зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Создаем структуру директорий
RUN mkdir -p internal/database pkg/templates

# Копируем env файлы
COPY internal/database/DB_Config.env ./internal/database/
COPY internal/database/SecretHash.env ./internal/database/

# Копируем шаблоны и остальные файлы
COPY pkg/templates/ ./pkg/templates/
COPY internal/ ./internal/
COPY cmd/ ./cmd/

# Собираем приложение
RUN go build -o main ./cmd/main.go

CMD ["./main"]




