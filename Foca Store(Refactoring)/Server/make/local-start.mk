CONAINER_NAME=voca_store

# APP
start:
	cd .. && go run cmd/main.go

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