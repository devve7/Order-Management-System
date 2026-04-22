# Order Management System

Backend сервис для управления заказами и товарами, реализованный на Go с использованием архитектурного подхода Domain-Driven Design (DDD).

Проект демонстрирует разработку масштабируемого backend-сервиса с чистой архитектурой, разделением ответственности и вниманием к инженерным практикам.

---

## 🚀 Основные возможности

- Управление заказами (создание, обработка)
- Работа с товарами (интеграция через отдельный сервис внутри монолита)
- REST API для взаимодействия с системой
- Валидация входных данных
- Логирование запросов
- Централизованная обработка ошибок
- Graceful shutdown
- Работа с PostgreSQL
- Миграции базы данных
- Контейнеризация через Docker

---

## 🏗 Архитектура

Проект построен с использованием принципов **DDD** и близок к **Hexagonal Architecture**.

### Слои:

- **domain** — бизнес-сущности и инварианты  
- **application** — use cases и бизнес-логика  
- **infrastructure** — работа с базой данных и внешними зависимостями  
- **transport (HTTP)** — обработка HTTP-запросов  

### Ключевые особенности:

- Слабая связность компонентов через **ports (интерфейсы)**
- Разделение сервисов:
  - Order service
  - Product service
- Возможность расширения (например, управление складом уже предусмотрено в домене)

---

## ⚙️ Технологии

- Go
- PostgreSQL
- Docker / docker-compose
- REST API
- SQL миграции
- Context propagation
- Error wrapping

---

## 🧠 Инженерные решения

- Использование `context.Context` во всех слоях приложения
- Централизованная обработка ошибок на транспортном уровне
- Error wrapping для сохранения контекста ошибок
- Транзакционная работа с базой данных
- Разделение логики через use cases
- Repository pattern для работы с данными
- Middleware для логирования HTTP-запросов
- Graceful shutdown сервера

---

## 📦 Структура проекта

```
cmd/app              # точка входа
internal/
  domain/            # доменная модель
  application/       # use cases
  infrastructure/    # база данных, репозитории
  transport/http/    # HTTP слой (handlers, middleware)
migrations/          # SQL миграции
```

---

## 🛠 Makefile

Проект поддерживает удобное управление через Makefile.

### Основные команды:

```bash
make docker-up      # сборка и запуск сервиса (app + postgres)
make docker-down    # остановка контейнеров
make docker-logs    # просмотр логов

make run            # запуск приложения локально (без Docker)
make build          # сборка бинарника

make test           # запуск всех тестов

make migrate-up     # применить миграции
make migrate-down   # откатить миграции
```

---

## 🐳 Run with Docker

```bash
docker-compose up --build
```

Сервис будет доступен по адресу:

http://localhost:9091

---

## 📦 API

### Orders

---

#### Create order
POST /orders

``` json 
{
  "customer_id": 1
}
```

---

#### Get order 
GET /orders/{id}

---

#### Get all orders
GET /orders

---

#### Add item to order
POST /orders/{id}/items

``` json 
{
  "product_id": 1,
  "quantity": 2
}
```

--- 

#### Remove item from order
DELETE /orders/{id}/items/{item_id}

--- 

#### Pay order
POST /orders/{id}/pay

--- 

#### Ship order
POST /orders/{id}/ship

---

#### Cancel order
POST /orders/{id}/cancel

----

### Products

---

#### Create product
POST /products

```json
{
  "name": "Laptop",
  "price_cents": 100000,
  "stock": 10
}
```

---

#### Get product
GET /products/{id}

---

#### Get all products
GET /products

---

#### Change price
PATCH /products/{id}/price

```json
{
  "price": 120000
}
```

---

#### Add stock
POST /products/{id}/stock/add

```json
{
  "stock_amount": 5
}
```

---

#### Remove stock
POST /products/{id}/stock/remove

```json
{
  "stock_amount": 3
}
```

---

#### Activate product
POST /products/{id}/activate

---

#### Deactivate product
POST /products/{id}/deactivate

---


## 🧪 Тестирование

Проект покрыт unit и integration тестами.

```bash
go test ./...
```

---

## 🔮 Дальнейшее развитие

- Управление складскими остатками
- Кэширование (Redis)
- Асинхронная обработка заказов
- Расширение продуктового сервиса
- Внедрение очередей (Kafka / RabbitMQ)

---

## 📬 Контакты

GitHub: https://github.com/devve7 
Email: kreator445@gmail.com
Telegram: @devve7
