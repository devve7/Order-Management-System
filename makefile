include .env
export

.PHONY: build run docker-up docker-down docker-logs test migrate-up migrate-down seed clean

# ---------- Local ----------

build:
	@go build -o ${APP_NAME} ./cmd/app

run:
	@go run ./cmd/app/main.go

# ---------- Docker ----------

docker-up:
	@docker-compose up --build

docker-down:
	@docker-compose down

docker-logs:
	@docker-compose logs -f

# ---------- Tests ----------

test:
	@go test ./...

# ---------- Migrations ----------

migrate-up:
	@migrate -path migrations -database "${DB_DSN}" up

migrate-down:
	@migrate -path migrations -database "${DB_DSN}" down

# ---------- Seed ----------

seed:
	@docker-compose exec -T postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < seeds/001_seed_dev_data.sql

# ---------- Cleanup ----------

clean:
	@rm -f ${APP_NAME}