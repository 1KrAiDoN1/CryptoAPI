# Этап сборки
FROM golang:1.23.1-alpine AS builder

# Создаем и переходим в рабочую директорию
WORKDIR /cryptoapi
# Создаем структуру директорий
RUN mkdir -p /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/

# Копируем только необходимые файлы для зависимостей
COPY go.mod go.sum ./
RUN go mod download

COPY internal/database/DB_Config.env ./internal/database/DB_Config.env
COPY internal/database/secretHash.env ./internal/database/secretHash.env

# Копируем шаблоны
COPY internal/database/ /Users/pavelvasilev/Desktop/CryptoAPI/internal/database/
COPY pkg/templates/ /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/

# Проверяем, что файлы скопировались
RUN ls -la /Users/pavelvasilev/Desktop/CryptoAPI/pkg/templates/

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go


CMD ["./main"]






