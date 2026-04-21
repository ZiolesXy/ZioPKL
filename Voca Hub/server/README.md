# Voca Hub Server

## Stack

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- MinIO
- WebSocket
- Clerk
- Docker

## Setup

1. Copy `.env.example` to `.env`
2. Update Clerk values, including `CLERK_SECRET_KEY`
3. Set `CORS_ALLOWED_ORIGINS` with frontend origin list separated by commas when needed
3. Choose one workflow:
4. Full docker stack: `make docker-setup`
5. Local app with Redis and MinIO only: `make setup`
6. Run migration: `make migrate`
7. Run app locally when using local workflow: `make run`

## Make Targets

- `make setup`: start local Redis and MinIO only
- `make run`: run the Go server locally
- `make up`, `make down`, `make restart`, `make logs`, `make ps`, `make clear`: manage local Redis and MinIO
- `make build`, `make tidy`: local Go utility targets
- `make docker-setup`, `make docker-up`, `make docker-down`, `make docker-restart`, `make docker-build`, `make docker-logs`, `make docker-ps`, `make docker-clear`: manage the full docker stack
- `make migrate`, `make seed`, `make drop`, `make reset`: database and seeder commands

## Local Dependency Defaults

- Local MinIO credentials are fixed to `minioadmin` / `minioadmin`
- Local workflow only starts Redis and MinIO from `make/local.compose.yml`

## Seeder Notes

- `go run ./seeders/seed.go` and `go run ./seeders/reset.go` resolve `clerk_id` from Clerk using the configured user emails.
- The seed emails must already exist as users in your Clerk instance.

## Endpoints

- `GET /health`
- `POST /api/friends/request`
- `POST /api/friends/:id/accept`
- `POST /api/friends/:id/reject`
- `GET /api/friends`
- `GET /api/chat/history/:user_id`
- `GET /api/chat/ws`
- `POST /api/games/upload`
- `GET /api/games`
- `GET /api/games/:id`
- `GET /api/games/:id/play`
- `GET /api/admin/dashboard`
- `GET /api/admin/users`
- `GET /api/admin/games`
- `POST /api/admin/games/:id/approve`
- `POST /api/admin/games/:id/reject`

## CORS

- Configure allowed browser origins with `CORS_ALLOWED_ORIGINS`
- Use comma-separated values, for example: `http://localhost:3000,http://localhost:5173`
- Use `*` to allow every origin
