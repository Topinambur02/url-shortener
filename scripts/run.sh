#!/bin/bash

if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-url_shortener_db}

echo "Выберите тип хранилища для запуска сервиса:"
echo "1) in-memory (по умолчанию)"
echo "2) postgres"
read -p "Введите 1 или 2: " STORAGE_CHOICE

if [ "$STORAGE_CHOICE" == "2" ] || [ "$STORAGE_CHOICE" == "postgres" ]; then
    echo "Выбран режим: postgres"
    
    if ! command -v psql &> /dev/null; then
        echo "Ошибка: утилита psql не найдена. Она требуется для создания базы данных."
        exit 1
    fi

    echo "Проверка и создание базы данных '$DB_NAME' на $DB_HOST:$DB_PORT..."
    
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME"

    if [ $? -ne 0 ]; then
        echo "Ошибка при создании базы данных. Проверьте данные для подключения к PostgreSQL."
        exit 1
    fi

    echo "Применение миграций к БД..."
    go run ./cmd/migrate/main.go up

    if [ $? -ne 0 ]; then
        echo "Ошибка при выполнении миграций. Запуск сервиса отменен."
        exit 1
    fi

    echo "Запуск сервиса с хранилищем postgres..."
    go run ./cmd/app/main.go -storage=postgres

else
    echo "Выбран режим: in-memory"
    echo "Запуск сервиса..."
    go run ./cmd/app/main.go
fi