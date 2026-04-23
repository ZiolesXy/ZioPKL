APP_NAME=voca-store

# Docker Compose
ENV_FILE=.env
DOCKER_COMPOSE=docker-compose --env-file $(ENV_FILE)

# Containers
API_CONTAINER=api
REDIS_CONTAINER=voca-store

.PHONY: \
	up down restart logs ps sh \
	build rebuild rebuild-nocache \
	clean delete-all \
	start \
	redis-start redis-run redis-stop redis-remove

# Docker Compose Commands
setup :
	$(DOCKER_COMPOSE) build
	$(DOCKER_COMPOSE) up -d

image-setup :
	docker pull redis:7
	docker pull postgres:15

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up -d

ps:
	$(DOCKER_COMPOSE) ps

logs:
	$(DOCKER_COMPOSE) logs -f $(API_CONTAINER)

sh:
	$(DOCKER_COMPOSE) exec $(API_CONTAINER) sh

# Build
build:
	$(DOCKER_COMPOSE) build

rebuild:
	$(DOCKER_COMPOSE) up -d --build

rebuild-nocache:
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

# Cleanup
clean:
	$(DOCKER_COMPOSE) down -v
	docker system prune --volumes -f

delete-all:
	docker system prune -a --volumes

# App Local Run (without Docker)
start:
	go run cmd/main.go

# Redis (Standalone Redis, not inside compose)
infra-up:
	docker compose up -d

infra-down:
	docker compose down
	docker volume prune -f