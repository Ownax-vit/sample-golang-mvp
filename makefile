.PHONY: help swagger migrate-up migrate-down migrate-create build run clean test deps install-tools

# Переменные
BINARY_NAME=app
MAIN_PATH=./cmd/main.go
MIGRATIONS_PATH=./migrations
SWAGGER_DIR=./docs

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

help: ## Показать справку по командам
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# Swagger команды
swagger:
	@echo "$(GREEN)Генерация Swagger документации...$(NC)"
	@swag init -g $(MAIN_PATH) -o $(SWAGGER_DIR)
	@echo "$(GREEN)Swagger документация сгенерирована в $(SWAGGER_DIR)$(NC)"

swagger-clean: ## Очистить сгенерированные файлы Swagger
	@echo "$(YELLOW)Очистка Swagger файлов...$(NC)"
	@rm -rf $(SWAGGER_DIR)/docs.go $(SWAGGER_DIR)/swagger.json $(SWAGGER_DIR)/swagger.yaml

# Миграции
migrate-up: ## Применить все миграции вверх
	@echo "$(GREEN)Применение миграций...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$${POSTGRES_URL}" up
	@echo "$(GREEN)Миграции применены$(NC)"

migrate-down: ## Откатить последнюю миграцию
	@echo "$(YELLOW)Откат последней миграции...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$${POSTGRES_URL}" down 1

migrate-down-all: ## Откатить все миграции
	@echo "$(YELLOW)Откат всех миграций...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$${POSTGRES_URL}" down -all

migrate-force: ## Принудительно установить версию миграции (использовать: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(YELLOW)Ошибка: укажите VERSION. Пример: make migrate-force VERSION=1$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Принудительная установка версии $(VERSION)...$(NC)"
	@migrate -path $(MIGRATIONS_PATH) -database "$${POSTGRES_URL}" force $(VERSION)

migrate-version: ## Показать текущую версию миграции
	@migrate -path $(MIGRATIONS_PATH) -database "$${POSTGRES_URL}" version

migrate-create: ## Создать новую миграцию (использовать: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "$(YELLOW)Ошибка: укажите NAME. Пример: make migrate-create NAME=create_users_table$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Создание миграции $(NAME)...$(NC)"
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
	@echo "$(GREEN)Миграция создана в $(MIGRATIONS_PATH)$(NC)"

# Билд
build: ## Собрать бинарный файл
	@echo "$(GREEN)Сборка приложения...$(NC)"
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Бинарный файл создан: $(BINARY_NAME)$(NC)"

build-linux: ## Собрать бинарный файл для Linux
	@echo "$(GREEN)Сборка приложения для Linux...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "$(GREEN)Бинарный файл создан: $(BINARY_NAME)-linux$(NC)"

build-windows: ## Собрать бинарный файл для Windows
	@echo "$(GREEN)Сборка приложения для Windows...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_PATH)
	@echo "$(GREEN)Бинарный файл создан: $(BINARY_NAME).exe$(NC)"

# Запуск
run: ## Запустить приложение
	@echo "$(GREEN)Запуск приложения...$(NC)"
	@go run $(MAIN_PATH)

run-build: build ## Собрать и запустить приложение
	@echo "$(GREEN)Запуск собранного приложения...$(NC)"
	@./$(BINARY_NAME)

# Очистка
clean: ## Очистить сгенерированные файлы и бинарники
	@echo "$(YELLOW)Очистка...$(NC)"
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-linux $(BINARY_NAME).exe
	@echo "$(GREEN)Очистка завершена$(NC)"

clean-all: clean swagger-clean ## Очистить все (бинарники + swagger)
	@echo "$(GREEN)Полная очистка завершена$(NC)"

# Тестирование
test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	@go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии сохранен в coverage.html$(NC)"

# Зависимости
deps: ## Скачать зависимости
	@echo "$(GREEN)Загрузка зависимостей...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Зависимости обновлены$(NC)"

deps-update: ## Обновить все зависимости
	@echo "$(GREEN)Обновление зависимостей...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)Зависимости обновлены$(NC)"

# Установка инструментов
install-tools: ## Установить необходимые инструменты (swag, migrate)
	@echo "$(GREEN)Установка инструментов...$(NC)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)Инструменты установлены$(NC)"

# Комплексные команды
dev: swagger run ## Сгенерировать swagger и запустить приложение

build-all: swagger build ## Сгенерировать swagger и собрать приложение

setup: install-tools deps swagger ## Полная настройка проекта (установка инструментов, зависимостей, генерация swagger)

# Проверка кода
lint: ## Запустить линтер (требует golangci-lint)
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint не установлен. Установите: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

fmt: ## Форматировать код
	@echo "$(GREEN)Форматирование кода...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Код отформатирован$(NC)"

vet: ## Запустить go vet
	@echo "$(GREEN)Проверка кода go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Проверка завершена$(NC)"

# По умолчанию показываем справку
.DEFAULT_GOAL := help