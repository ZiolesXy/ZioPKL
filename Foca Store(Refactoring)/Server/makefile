up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f app

restart:
	docker compose restart app

migrate:
	docker compose exec app migrate up

seed:
	docker compose exec app go run scripts/seeder.go