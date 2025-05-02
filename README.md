<div align="center">
  <h1>Crypto Market</h1>
  <p>Веб-приложение для отслеживания данных о криптовалютах с аутентификацией пользователей и интеграцией с CoinCap API</p>
  
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/PostgreSQL-13+-4169E1?style=for-the-badge&logo=postgresql" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Docker-2.5+-2496ED?style=for-the-badge&logo=docker" alt="Docker">
</div>

## 🚀 О проекте

Crypto Market — это веб-приложение для отслеживания данных о криптовалютах с:
- Аутентификацией пользователей
- Управлением избранными активами 
- Интеграцией с CoinCap API
- Удобным интерфейсом для анализа рынка
- Персонализированными функциями для пользователей

## 🛠 Технологический стек

### Бэкенд
<div>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Swagger-85EA2D?style=flat-square&logo=swagger&logoColor=black" alt="Swagger">
</div>

### Фронтенд
<div>
  <img src="https://img.shields.io/badge/HTML5-E34F26?style=flat-square&logo=html5&logoColor=white" alt="HTML5">
  <img src="https://img.shields.io/badge/CSS3-1572B6?style=flat-square&logo=css3&logoColor=white" alt="CSS3">
</div>

### Инфраструктура
<div>
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/Docker_Compose-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker Compose">
</div>

### Основные библиотеки
- github.com/swaggo/swag/cmd/swag - Генерация Swagger документации
- github.com/golang-jwt/jwt - Работа с JWT токенами
- github.com/joho/godotenv - Загрузка переменных окружения
- github.com/jackc/pgx/v4 - Драйвер PostgreSQL для Go

## 🌟 Возможности

### 🔐 Аутентификация
- Регистрация и вход по email/паролю
- JWT токены (access + refresh)
- Автоматическое обновление токенов
- Защищенные маршруты через middleware
- Выход с инвалидацией сессии

### 💰 Работа с криптовалютами
- Просмотр списка криптовалют
- Детальная информация по каждой монете
- Добавление в избранное
- Персональный портфель

### 🖥 Пользовательский интерфейс
- Личный кабинет
- Дата регистрации
- Форматирование чисел (капитализация, объемы)
- Адаптивный дизайн

### 🔐 Безопасность
- JWT токены в HttpOnly cookies
- Хеширование паролей (SHA1 + секретный ключ)
- Валидация токенов на каждом запросе


### Запуск
```bash
# Сборка и запуск
make build && make up

# Остановка
make down

# Просмотр логов
make logs

# Применить миграции
make migrate-up

# Откатить миграции
make migrate-dwon

После запуска приложение будет доступно по адресу:
http://localhost:8080

Swagger документация:
http://localhost:8080/swagger/index.html
