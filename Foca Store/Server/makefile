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
	docker system prune -a

delete-all:
	docker system prune -a --volumes

# App Local Run (without Docker)
start:
	go run cmd/main.go

# Redis (Standalone Redis, not inside compose)
redis-start:
	docker run -d --name $(REDIS_CONTAINER) -p 6379:6379 redis:7

redis-run:
	docker start $(REDIS_CONTAINER)

redis-stop:
	docker stop $(REDIS_CONTAINER)

redis-remove:
	docker stop $(REDIS_CONTAINER) || true
	docker rm $(REDIS_CONTAINER) || true
	docker volume prune -f