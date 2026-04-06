APP_NAME=voca-store
DOCKER_COMPOSE=docker-compose --env-file .env

CONTAINER_NAME=voca-store

.PHONY: up down restart logs build rebuild rebuild-nocache ps sh migrate seed clean redis-start

# Containers
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
	$(DOCKER_COMPOSE) logs -f api
sh:
	$(DOCKER_COMPOSE) exec api sh

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
	docker system prune -f
delete-all:
	docker system prune -a --volumes

# APP
start:
	go run cmd/main.go

# Redis
redis-start:
	docker run -d --name $(CONTAINER_NAME) -p 6379:6379 redis:7
redis-run:
	docker start $(CONTAINER_NAME)
redis-stop:
	docker stop $(CONTAINER_NAME)
redis-remove:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)