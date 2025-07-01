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
# # COPY internal/database/DB_Config.env ./internal/database/DB_Config.env
# # COPY internal/database/SecretHash.env ./internal/database/SecretHash.env

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



# CMD ["./main"]


FROM golang:1.23.1-alpine AS builder

WORKDIR /cryptoapi

# Создаем нужные директории
RUN mkdir -p ./pkg/templates/
RUN mkdir -p ./internal/database/

# Копируем go.mod и go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные необходимые файлы (без env)
COPY *.env ./

# Копируем исходники и шаблоны
COPY internal/database/ ./internal/database/
COPY pkg/templates/ ./pkg/templates/

# Копируем весь проект
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

CMD ["./main"]









