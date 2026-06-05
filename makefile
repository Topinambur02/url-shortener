DC = docker-compose

.PHONY: up down restart build logs lint format help run migrate test mocks swagger

help:
	@echo "Доступные команды:"
	@echo "  up              Поднять контейнеры"
	@echo "  down            Остановить и удалить контейнеры"
	@echo "  restart         Перезапустить контейнеры"
	@echo "  logs            Посмотреть логи"
	@echo "  run             Запустить приложение локально"
	@echo "  migrate         Запустить накатывание миграций"
	@echo "  build           Собрать бинарный файл приложения"
	@echo "  lint            Проверить код линтером"
	@echo "	 format			 Отформатировать код"
	@echo "  test            Запустить тесты"
	@echo "  mocks           Сгенерировать моки для интерфейсов"
	@echo "	 swagger		 Сгенерировать swagger документацию

up:
	$(DC) up -d --build

down:
	$(DC) down

restart: 
	down up

logs:
	$(DC) logs -f

run:
	go run cmd/app/main.go

migrate:
	go run cmd/migrate/main.go up

build:
	go build -o bin/app cmd/app/main.go

lint:
	golangci-lint run ./...

format:
	goimports -w .

test:
	go test -v -cover ./...

mocks:
	go generate ./...

swagger:
	swag init -g cmd/app/main.go