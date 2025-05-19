# store-it WMS — backend

[![E2E Tests](https://github.com/lets-store-it/backend/actions/workflows/e2e-tests.yml/badge.svg)](https://github.com/lets-store-it/backend/actions/workflows/e2e-tests.yml) [![Build and Publish Docker Image](https://github.com/lets-store-it/backend/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/lets-store-it/backend/actions/workflows/docker-publish.yml)

## Структура проекта

### Основные директории

#### .github/

Конфигурация GitHub Actions для CI/CD:

- `e2e-tests.yml` — запуск интеграционных тестов
- `docker-publish.yml` — сборка и публикация Docker образа

#### api/

Содержит исходный код OpenAPI спецификации REST API сервиса, на основе которой, с помощью `ogen`, генерируется код для HTTP обработчиков и моделей данных.

#### cmd/

Точки входа в приложение:

- `init_db/` - утилита для инициализации базы данных
- `server/` - HTTP сервер

#### config/

Код ответственный за конфигурацию приложения.

#### generated/

Автоматически генерируемый код:

- Скопилированная конфигурация `OpenAPI`
- Сгенерированные интерфейсы для HTTP обработчиков (`ogen`)
- Клиент базы данных (`sqlc`).

### Внутренняя структура (internal/)

#### database/

Компоненты для работы с базой данных, ошибки, функции-хелперы.

#### server/

Конфигурация и запуск приложения и веб-сервера.

#### models/

Модели данных и бизнес-сущностей используемые для передачи данных между слоями `handlers`<->`usecases`<->`services`.

#### handlers/

Реализации HTTP обработчиков запросов — подготовка, валидация, запросов для передачи в `usecases` и выполнения бизнес-логики.

#### usecases/

Представляет пользовательские бизнес-сценарии. Инкапсулирует работу с сервисами.

#### services/

Код отвечающий за взаимодействие с внешними сервисами.

#### telemetry/

Работа с метриками и трейсингом Opentelemetry.


#### utils/

Вспомогательные функции и общего назначения.

### Тестирование (tests/)

Комплексные тесты приложения:

- `k6/` — нагрузочные тесты с использованием `k6`
- `e2e/` — интеграционные тесты на основе `pytest`

### Дополнительные файлы

- `schema.sql` — SQL схема базы данных с описанием всех таблиц и связей.
- `query.sql` — SQL запросы для генерации типизированного клиента базы данных на основе `sqlc`.
- `copose.yml` — манифест для запуска базы данных и приложения с помощью Docker.
- `Dockerfile` — манифест для сборки Docker образа приложения.

## Запуск

### Запуск с помощью Docker

Приложение полностью контейнеризировано и может быть запущено с помощью Docker. Для запуска необходимо выполнить следующие шаги:

#### Сборка образа приложения

```bash
docker build -t storeit-backend .
```

#### Запуск базы данных и приложения

Создайте файл `.env` с необходимыми переменными окружения (см. раздел "Конфигурация приложения") и выполните:

```bash
docker compose up -d postgres

docker run --rm --network host storeit-backend /app/init_db -schema /app/schema.sql

docker run --rm -d \
  --network host \
  --env-file .env \
  --name storeit-backend \
  storeit-backend
```

#### Проверка работоспособности

Приложение должно быть доступно по адресу http://localhost:8080

#### Остановка приложения

```bash
docker stop storeit-backend
docker compose down -v
```

### Запуск тестов

Для запуска тестов необходимо предварительно собрать и запустить приложение, после чего:

#### e2e тесты

```bash
cd tests/e2e
uv run pytest . 
```

#### k6 тесты

```bash
cd tests/k6
./run.sh
```

## Конфигурация приложения

Приложение может быть настроено с помощью переменных окружения. Все параметры также могут быть указаны в файле `.env` в PWD директории.

### Основные параметры

- `SERVICE_NAME` - название сервиса, используется в Opentelemetry (по умолчанию: "storeit-backend")

### Параметры сервера

- `LISTEN_ADDRESS` - адрес и порт для прослушивания (по умолчанию: "0.0.0.0:8080")
- `CORS_ORIGINS` - разрешенные источники для CORS, через запятую (по умолчанию: "http://localhost:3000,http://localhost:8080,http://localhost,https://store-it.ru,https://www.store-it.ru,http://store-it.ru,http://www.store-it.ru")

### Параметры базы данных

- `DB_HOST` - хост базы данных (по умолчанию: "localhost")
- `DB_PORT` - порт базы данных (по умолчанию: "5432")
- `DB_NAME` - название базы данных (по умолчанию: "storeit")
- `DB_USER` - пользователь базы данных (по умолчанию: "storeit")
- `DB_PASSWORD` - пароль базы данных (по умолчанию: "storeit")

### Параметры авторизации Яндекс

- `YANDEX_OAUTH_CLIENT_ID` - идентификатор приложения в Яндекс OAuth
- `YANDEX_OAUTH_CLIENT_SECRET` - секретный ключ приложения в Яндекс OAuth

### Параметры Kafka

- `KAFKA_ENABLED` - включение интеграции с Kafka (по умолчанию: false)
- `KAFKA_BROKERS` - список брокеров Kafka через запятую (по умолчанию: "localhost:9092")
- `KAFKA_AUDIT_TOPIC` - название топика для отправки событий изменений (по умолчанию: "audit.object-changes")

Пример файла `.env`:

```env
SERVICE_NAME=storeit-backend
LISTEN_ADDRESS=0.0.0.0:8080
DB_HOST=postgres
DB_PASSWORD=mysecretpassword
YANDEX_OAUTH_CLIENT_ID=your_client_id
YANDEX_OAUTH_CLIENT_SECRET=your_client_secret
KAFKA_ENABLED=true
KAFKA_BROKERS=kafka1:9092,kafka2:9092
```
