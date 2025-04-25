# Makefile для проекта CryptoAPI

# Переменные
DOCKER_COMPOSE = docker-compose
APP_NAME = cryptoapi-app

.PHONY: help build up down logs clean db-shell app-shell

help: ## Показать справку
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собрать контейнеры
	$(DOCKER_COMPOSE) build --no-cache

up: ## Запустить контейнеры
	$(DOCKER_COMPOSE) up -d

down: ## Остановить контейнеры
	$(DOCKER_COMPOSE) down

logs: ## Показать логи приложения
	$(DOCKER_COMPOSE) logs -f app

clean: ## Очистить все Docker-артефакты
	$(DOCKER_COMPOSE) down -v --rmi all

db-shell: ## Подключиться к консоли PostgreSQL
	$(DOCKER_COMPOSE) exec db psql -U postgres -d registration

app-shell: ## Подключиться к контейнеру приложения
	$(DOCKER_COMPOSE) exec app sh

migrate: ## Выполнить миграции (автоматически через volume)
	@echo "Миграции выполняются автоматически при старте контейнера БД"