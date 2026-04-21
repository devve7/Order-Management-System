.PHONY: build run docker-up docker-down logs test migrate-up migrate-down

APP_NAME=oms-app

# ---------- Local ----------

build:
	go build -o $(APP_NAME) ./cmd/app

run:
	go run ./cmd/app/main.go

# ---------- Docker ----------

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# ---------- Tests ----------

test:
	go test ./...

# ---------- Migrations ----------

migrate-up:
	migrate -path migrations -database "$$DB_DSN" up

migrate-down:
	migrate -path migrations -database "$$DB_DSN" down

# ---------- Cleanup ----------

clean:
	rm -f $(APP_NAME)