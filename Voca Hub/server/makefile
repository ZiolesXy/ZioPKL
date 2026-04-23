APP_NAME=server
SHELL := powershell.exe

run:
	go run ./cmd

migrate:
	go run ./seeders/migrate.go

seed:
	go run ./seeders/seed.go

seed-games:
	go run ./seeders/seed_games.go

drop:
	go run ./seeders/drop.go

reset:
	go run ./seeders/reset.go

build:
	go build -o $(APP_NAME) ./cmd

tidy:
	go mod tidy

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down

infra-up:
	docker compose up -d
# 	powershell -ExecutionPolicy Bypass -File scripts/setup.ps1

infra-down:
	docker compose down
	docker volume prune -f
# 	powershell -ExecutionPolicy Bypass -File scripts/teardown.ps1

infra-restart:
	powershell -ExecutionPolicy Bypass -File scripts/teardown.ps1
	powershell -ExecutionPolicy Bypass -File scripts/setup.ps1
