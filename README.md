# template-svc

Шаблон микросервиса на Go для быстрого старта новых сервисов.

В репозитории уже есть:
- слоистая архитектура (Clean / Hexagonal);
- PostgreSQL + автоматические миграции (Goose);
- graceful shutdown;
- Swagger;
- пример **CRUD** для сущности `Item`.

---

## Быстрый старт

### 1. Требования

- Go **1.24+**
- PostgreSQL **14+**

### 2. Настройка окружения

```bash
cp .env.example .env
```

Создайте базу данных с именем из `DB_NAME` и проверьте параметры подключения `DB_*`.

### 3. Запуск

```bash
make run-local
```

Сервис поднимется на порту из `HTTP_PORT` (по умолчанию `8080`).

### 4. Проверка

```bash
# health check
curl http://localhost:8080/actuator/health

# создать item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"title":"First","description":"Example item"}'
```

### 5. Swagger

```bash
make swag
```

Откройте: http://localhost:8080/swagger/index.html

---

## Архитектура

### Принципы

1. **Зависимости направлены внутрь** — HTTP и БД не «просачиваются» в домен.
2. **Интерфейсы на границах слоёв** — `delivery` знает только `Usecase`, `usecase` — только `Service`, и т.д.
3. **Composition root в одном месте** — все зависимости собираются в `internal/app/app.go`.
4. **Инфраструктура в `pkg/`** — код, который можно переиспользовать между сервисами.
5. **Бизнес-модули в `internal/app/`** — каждый модуль = отдельный bounded context.

### Общая схема

```
                    ┌─────────────────────┐
                    │     cmd/main.go     │
                    │  env · logger · run │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │  internal/app/app   │  ← composition root (DI)
                    │  DB · Router · deps │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
┌────────▼────────┐   ┌────────▼────────┐   ┌───────▼───────┐
│  someModule     │   │  yourModule     │   │  env/         │
│  (пример CRUD)  │   │  (новый модуль) │   │  конфигурация │
└────────┬────────┘   └─────────────────┘   └───────────────┘
         │
         │  delivery → usecase → service → repository → PostgreSQL
         │
┌────────▼────────────────────────────────────────────────────┐
│  pkg/  — инфраструктура                                     │
│  db · http · conductor · cache · middleware · i18n          │
└─────────────────────────────────────────────────────────────┘
```

### Жизненный цикл приложения

```
main()
  ├─ godotenv.Load()              # .env (опционально)
  ├─ env.New()                    # парсинг переменных окружения
  ├─ app.New(config, logger)      # DI: DB → repo → service → usecase → router
  ├─ регистрация /actuator/health, /swagger
  ├─ pkg/http.New(...)            # HTTP-сервер
  └─ conductor.Run / Shutdown     # graceful stop по SIGINT / SIGTERM
```

Сборка зависимостей в `internal/app/app.go`:

```go
sqlxDB  := db.NewPostgresDB(config.DbConfig)
repo    := postgres.NewItemRepository(sqlxDB)
svc     := service.New(repo)
uc      := usecase.New(svc)
router  := gin.New()
delivery.AttachRoutes(router, uc)
```

---

## Слои модуля (на примере `someModule`)

Каждый бизнес-модуль живёт в `internal/app/<moduleName>/`:

```
someModule/
├── interface.go          # контракты Repository, Service, Usecase
├── errors.go             # доменные ошибки (ErrNotFound и т.д.)
├── model/                # сущности без зависимостей от Gin/SQL
├── repository/postgres/  # SQL-запросы
├── service/              # операции над сущностями
├── usecase/              # бизнес-правила и валидация
└── delivery/
    ├── handler.go        # HTTP handlers
    ├── routes.go         # регистрация маршрутов
    └── adapter/          # request/response DTO
```

| Слой | Пакет | Что делает | Чего не делает |
|------|--------|------------|----------------|
| **Delivery** | `delivery/` | HTTP, JSON, коды ответа, Swagger-аннотации | SQL, бизнес-правила |
| **Usecase** | `usecase/` | Валидация, лимиты, orchestration | Знание о Gin и SQL |
| **Service** | `service/` | CRUD-операции, маппинг в `model` | HTTP-форматы |
| **Repository** | `repository/postgres/` | SQL, работа с `sqlx` | HTTP, валидация |
| **Model** | `model/` | Структуры домена | Зависимости от фреймворков |

### Поток HTTP-запроса

```
POST /api/v1/items
    │
    ▼
delivery/handler.go       ShouldBindJSON → adapter.CreateItemRequest
    │
    ▼
usecase/usecase.go        validateTitle(), TrimSpace
    │
    ▼
service/service.go        model.Item → repo.Create
    │
    ▼
repository/postgres/      INSERT INTO items ...
    │
    ▼
PostgreSQL                таблица items
```

### Правила зависимостей

```
delivery  ──►  usecase  ──►  service  ──►  repository  ──►  pkg/db
   │              │             │              │
   └─ adapter     └─ model      └─ model       └─ model
```

- `model` не импортирует ничего из `delivery`, `usecase`, `service`, `repository`.
- `interface.go` объявляет контракты — удобно для моков в тестах.
- Gin используется **только** в `delivery/` и `internal/app/app.go`.

---

## Инфраструктурные пакеты (`pkg/`)

| Пакет | Назначение |
|-------|------------|
| `pkg/db` | Подключение к PostgreSQL, пул соединений, миграции Goose |
| `pkg/http` | Обёртка над `net/http` с таймаутами и `Shutdown` |
| `pkg/conductor` | Запуск и остановка сервисов по сигналам ОС |
| `pkg/cache` | Абстракция кэша: Redis Sentinel / in-memory |
| `pkg/middleware` | JWT-авторизация, CORS |
| `pkg/i18n` | Локализация сообщений (go-i18n) |

JWT и CORS подключаются в `app.go` или на уровне группы роутов — в шаблоне CORS включается через `CORS_ENABLED=true`.

---

## API — CRUD Items (пример)

Базовый путь: **`/api/v1/items`**

| Метод | Путь | Описание | Ответ |
|-------|------|----------|-------|
| `POST` | `/api/v1/items` | Создать | `201` + item |
| `GET` | `/api/v1/items` | Список | `200` + `{ items, limit, offset }` |
| `GET` | `/api/v1/items/:id` | Получить по ID | `200` + item |
| `PUT` | `/api/v1/items/:id` | Обновить | `200` + item |
| `DELETE` | `/api/v1/items/:id` | Удалить | `204` |

Query-параметры списка: `limit` (по умолчанию 20, макс. 100), `offset` (по умолчанию 0).

### Примеры curl

```bash
# список
curl "http://localhost:8080/api/v1/items?limit=20&offset=0"

# получить
curl http://localhost:8080/api/v1/items/1

# обновить
curl -X PUT http://localhost:8080/api/v1/items/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","description":"New text"}'

# удалить
curl -X DELETE http://localhost:8080/api/v1/items/1
```

### Форматы JSON

**Успех (item):**
```json
{
  "id": 1,
  "title": "First",
  "description": "Example item",
  "created_at": "2026-05-28T12:00:00Z",
  "updated_at": "2026-05-28T12:00:00Z"
}
```

**Ошибка:**
```json
{
  "code": "not_found",
  "message": "item not found"
}
```

Коды ошибок: `bad_request`, `not_found`, `internal_error`, `unauthorized`.

---

## Как добавить новый модуль

Пошаговая инструкция для нового bounded context (например, `orderModule`):

### Шаг 1. Скопировать структуру

```bash
cp -r internal/app/someModule internal/app/orderModule
```

Переименуйте пакеты и типы под свой домен.

### Шаг 2. Описать контракты

В `internal/app/orderModule/interface.go` объявите:

- `Repository` — методы работы с БД;
- `Service` — операции над сущностями;
- `Usecase` — публичный API для delivery.

### Шаг 3. Реализовать слои снизу вверх

1. `model/` — структуры домена;
2. `repository/postgres/` — SQL;
3. `service/` — делегирование в repository;
4. `usecase/` — валидация и бизнес-правила;
5. `delivery/` — handlers, adapter, `AttachRoutes`.

### Шаг 4. Миграция

Создайте файл в `migrations/`:

```
migrations/00002_create_orders.sql
```

Миграции применяются автоматически при старте (`pkg/db`).

### Шаг 5. Подключить в composition root

В `internal/app/app.go`:

```go
orderRepo := orderPostgres.NewOrderRepository(sqlxDB)
orderSvc  := orderService.New(orderRepo)
orderUC   := orderUsecase.New(orderSvc)
orderDelivery.AttachRoutes(router, orderUC)
```

### Шаг 6. Swagger

Добавьте godoc-аннотации в handlers и выполните:

```bash
make swag
```

---

## Миграции

- Каталог: `migrations/`
- Инструмент: [Goose](https://github.com/pressly/goose)
- Применение: автоматически при старте сервиса

Пример миграции:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS items (...);

-- +goose Down
DROP TABLE IF EXISTS items;
```

---

## Переменные окружения

| Переменная | Обязательная | Описание |
|------------|:------------:|----------|
| `LOG_LEVEL` | да | `debug`, `info`, `warn`, `error` |
| `HTTP_PORT` | да | Порт HTTP (`8080`) |
| `HTTP_CLIENT_TIMEOUT` | да | Read timeout сервера (`15s`) |
| `CORS_ENABLED` | нет | `true` — включить CORS middleware |
| `DB_TYPE` | нет | Тип БД, по умолчанию `postgres` |
| `DB_USERNAME` | да | Пользователь БД |
| `DB_PASSWORD` | да | Пароль БД |
| `DB_HOST` | да | Хост БД |
| `DB_PORT` | да | Порт БД |
| `DB_NAME` | да | Имя базы |
| `DB_PARAMS` | нет | Query-параметры DSN |
| `DB_MAX_OPEN_CONNECTIONS` | да | Макс. открытых соединений |
| `DB_MAX_IDLE_CONNECTIONS` | да | Макс. idle-соединений |
| `DB_CONNECTION_MAX_IDLE_TIME` | нет | TTL idle-соединения |

---

## Makefile

| Команда | Действие |
|---------|----------|
| `make run-local` | Запуск сервиса |
| `make test` | Тесты + coverage |
| `make run-tests` | Тесты через gotestsum |
| `make lint` | golangci-lint |
| `make swag` | Генерация Swagger в `docs/` |
| `make gen` | `go generate ./...` |

---

## Структура репозитория

```
cmd/
  main.go                         # точка входа

internal/
  app/
    app.go                        # composition root (DI, router, health)
    someModule/                   # пример модуля с CRUD
      interface.go
      errors.go
      model/
      repository/postgres/
      service/
      usecase/
      delivery/
  env/
    env.go                        # конфигурация из env

migrations/
  00001_create_items.sql          # SQL-миграции Goose

pkg/
  db/                             # Postgres + миграции
  http/                           # HTTP server + options
  conductor/                      # graceful shutdown
  cache/                          # Redis / in-memory
  middleware/                     # JWT, CORS
  i18n/                           # локализация

docs/                             # Swagger (генерируется: make swag)
.env.example                      # шаблон переменных окружения
Makefile
go.mod
```

---

## Graceful shutdown

HTTP-сервер запускается через `pkg/http` и управляется `pkg/conductor`:

1. `conductor.Run()` — стартует сервер в goroutine;
2. ожидание `SIGINT` / `SIGTERM` или ошибки;
3. `conductor.Shutdown()` — вызывает `http.Server.Shutdown()` с таймаутом;
4. `app.Close()` — закрывает соединение с БД.

Это позволяет корректно завершать in-flight запросы при деплое или остановке контейнера.

---

## Рекомендации при форке шаблона

1. Переименуйте module path в `go.mod` (`github.com/your-org/your-svc`).
2. Замените `someModule` на свой домен или удалите, если не нужен.
3. Обновите Swagger `@title` в `cmd/main.go`.
4. Настройте CI: `make lint`, `make test`.
5. Для JWT защитите нужные группы роутов через `middleware.Configure(secret)`.
