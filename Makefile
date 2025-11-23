
.PHONY: all
all: native docker

.PHONY: native
native:
	@echo "Сборка нативной версии..."
	@. ./scripts/build-native_reivew.sh

.PHONY: docker
docker:
	@echo "Сборка Docker образа..."
	@. ./scripts/build-docker_review.sh

.PHONY: clean
clean:
	@echo "Очистка..."
	@rm -rf build/bin/review-*
	@docker images -q review:* 2>/dev/null | xargs -r docker rmi -f

.PHONY: version
version:
	@git describe --tags --abbrev=0 2>/dev/null || echo "unknown"

.PHONY: test
test:
	@echo "Запуск тестов..."
	@go test ./...

.PHONY: lint
lint:
	@echo "Проверка кодстайла..."
	@test -x $(shell which golangci-lint) || (echo "Установите golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run

.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make native    - Собрать нативную версию"
	@echo "  make docker    - Собрать Docker образ"
	@echo "  make lint      - Запустить линтер"
	@echo "  make test      - Запустить тесты"
	@echo "  make all       - Собрать обе версии"
	@echo "  make clean     - Очистить собранные файлы"
	@echo "  make version   - Показать текущую версию"
	@echo "  make help      - Показать эту справку"

.DEFAULT_GOAL := help
