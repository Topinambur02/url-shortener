# URL Shortener Service

Сервис для сокращения ссылок, написанный на **Go**. Предоставляет REST API для создания уникальных коротких URL и получения оригинальных ссылок по их идентификатору.

## Используемые технологии
- Go (версии 1.26+)
- PostgreSQL
- Swagger
- Make
- Pgx
- Golang-migrate
- Logrus
- net/http
- testcontainers
- testify

## Особенности

* **Уникальность:** На один оригинальный URL генерируется строго одна короткая ссылка (1:1).
* **Формат:** Длина сокращенной ссылки составляет ровно 10 символов.
* **Алфавит:** Используются символы латинского алфавита (верхний и нижний регистр), цифры и символ подчеркивания `[a-zA-Z0-9_]`.
* **Хранилище на выбор:** Поддержка **PostgreSQL** и самописного **In-Memory** хранилища (выбирается параметром при запуске).
* **Контейнеризация:** Полностью упаковано в Docker и готово к запуску через `docker-compose`.
* **Тестирование:** Бизнес-логика и основные компоненты покрыты Unit-тестами.

---

## Быстрый старт

### Требования
* [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/) (для удобства использования команд)
* [psql](https://www.postgresql.org/docs/current/app-psql.html) (опционально, для локального подключения к БД)

### Запуск проекта (без Docker)

1. Склонируйте репозиторий:
```bash
   git clone https://github.com/topinambur02/url-shortener.git
   cd url-shortener
```
2. Создайте .env файл со следующим содержимым (настройки по умолчанию):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=url_shortener_db
```
3. Запустите файл run.sh командой:
```bash
chmod +x ./scripts/run.sh
./scripts/run.sh
```

### Запуск проекта (с Docker)
1. Склонируйте репозиторий:
```bash
   git clone https://github.com/topinambur02/url-shortener
   cd url-shortener
```
2. Создайте .env файл со следующим содержимым (Опционально):
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=url_shortener_db
```
3. Запустите контейнера командой:
```bash
docker compose up -d
```

## API Endpoints

### 1. Создание короткой ссылки

Сохраняет оригинальный URL и возвращает сокращенный.

HTTP Запрос:
```HTTP
POST /api
Content-Type: application/json

{
  "original_url": "[https://finance.ozon.ru](https://finance.ozon.ru)"
}
```
HTTP Ответ (201 Created):
```JSON
{
  "short_url": "http://localhost:8080/LXKcFE4gj6"
}
```

### 2. Получение оригинальной ссылки

Принимает короткий идентификатор и возвращает оригинальный URL.

HTTP Запрос:
```HTTP
GET /api/LXKcFE4gj6
```
HTTP Ответ (200 OK):
```JSON
{
  "original_url": "[https://finance.ozon.ru](https://finance.ozon.ru)"
}
```

### 3. Редирект на оригинальную ссылку

Принимает короткий идентификатор и перенаправляет на оригинальный URL.

HTTP Запрос:
```HTTP
GET /LXKcFE4gj6
```
HTTP Ответ ```307 Temporary Redirect```

## Управление проектом

В проекте используется Makefile для инкапсуляции рутинных задач. Основные команды:

| Команда | Описание |
| :---: | ---: |
| make up | Собрать и запустить контейнеры (App, Postgres) в фоне |
| make down | Остановить контейнеры |
| make restart | Перезапустить контейнеры |
| make build | Скомпилировать бинарный файл |
| make logs | Вывести логи запущенных контейнеров |
| make lint | Запустить golangci-lint для проверки кода |
| make format | Отформатировать код с помощью goimports |
| make run | Локальный запуск сервиса (с in-memory хранилищем) |
| make migrate | Применить миграции к базе данных (Postgres) |
| make test | Запустить все тесты (Unit + Integration) |
| make mocks | Сгенерировать моки для интерфейсов с помощью mockery |
| make swagger | Сгенерировать Swagger документацию |
| make help | Вывести список всех доступных команд |

## Стресс-тестирование и нагрузочное тестирование

Было проведено стресс-тестирование для Postgres и In-memory хранилищ. Результаты:
<img width="812" height="432" alt="Postgres stress" src="https://github.com/user-attachments/assets/147759f5-8b08-40d5-b3f4-2c973f2bb479" />
<img width="812" height="432" alt="In-memory stress" src="https://github.com/user-attachments/assets/25548a58-31fb-40d8-8de7-aef5142c93bd" />

Также проведено нагрузочное тестирование:
<img width="828" height="413" alt="Screenshot 2026-06-06 at 16 29 37" src="https://github.com/user-attachments/assets/30c4e689-9d34-46ce-a45f-1b0933b539e5" />

## Используемые решения для оптимизации
- Использование драйвера pgx для оптимальной работы с PostgreSQL.
- Кастомное решение для формирования короткой ссылки (FNV хэш + приведение к base63).
- Расширение максимального количества подключений к базе данных.
- Использование Map в качестве In-Memory хранилища (амортизированное время поиска и сохранения — O(1)).

## Что можно улучшить

Согласно ТЗ, метод GET /api/{short_url} возвращает оригинальный URL в теле ответа. Однако в реальных production-системах основную нагрузку берет на себя эндпоинт, который сразу отдает HTTP-статус 302 Found (или 307 Temporary Redirect) с заголовком Location: <оригинальный_URL>. Это позволяет браузеру автоматически перенаправлять пользователя.

Использование статусов 302/307 вместо 301 Moved Permanently предпочтительнее, так как они не кэшируются браузером навсегда, что позволяет владельцу сервиса корректно собирать аналитику переходов (считать клики). В данном проекте этот подход реализован через отдельный эндпоинт GET /{short_url}.

## Структура проекта

```
├── Dockerfile            # Файл сборки образа приложения
├── README.md             # Текущий файл документации
├── cmd      
│   ├── app               # Точка входа в приложение
│   └── migrate           # Утилита для применения миграций
├── docker-compose.yml    # Конфигурация для запуска инфраструктуры
├── docs                  # Swagger документация
├── go.mod / go.sum       # Управление зависимостями
├── internal      
│   ├── config            # Парсинг переменных окружения
│   ├── db                # Подключение к PostgreSQL
│   ├── dto               # Data Transfer Objects (запросы и ответы API)
│   ├── handler           # HTTP-обработчики
│   ├── middleware        # HTTP-мидлвари
│   ├── model             # Доменные модели (сущности)
│   ├── repository        # Слой доступа к данным (реализации In-Memory и Postgres)
│   └── service           # Бизнес-логика
├── logs                  # Директория для хранения файлов логов
├── makefile              # Утилиты автоматизации сборки и разработки
├── migrations            # SQL-файлы миграций
├── pkg
│   ├── constants         # Глобальные константы
│   ├── exceptions        # Кастомные ошибки приложения
│   ├── logging           # Настройка системы логирования
│   ├── shutdown          # Утилиты для Graceful Shutdown
│   └── utils             # Вспомогательные функции
└── scripts
    ├── entrypoint.sh     # bash-скрипт для инициализации в Docker
    └── run.sh            # bash-скрипт для запуска приложения локально
```
