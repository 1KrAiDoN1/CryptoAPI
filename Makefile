# Makefile для проекта CryptoAPI

# Переменные
DOCKER_COMPOSE = docker-compose


build: ## Собрать контейнеры
	$(DOCKER_COMPOSE) build --no-cache

up: ## Запустить контейнеры
	$(DOCKER_COMPOSE) up -d

down: ## Остановить контейнеры
	$(DOCKER_COMPOSE) down

logs: ## Показать логи приложения
	$(DOCKER_COMPOSE) logs -f cryptoapi

migrate-up:
    migrate -path ./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable" up

migrate-down:
    migrate -path ./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable" down	

