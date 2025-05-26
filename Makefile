GO         := $(shell which go)
ROOT_DIR   := $(shell pwd)
DEPLOY_DIR = deploy
# Определяем переменную для хранения первого параметра
PARAM := $(word 2, $(MAKECMDGOALS))
DB_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=6432 host=localhost"
DB_SHARD_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=7432 host=localhost"
NETWORK_NAME = app_network

migrate_up:
	@echo "Starting migrate"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIR="migrations" goose up
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_shard" goose up
	@echo "migrate complected"

migrate_down:
	@echo "Starting down"

	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIRmigrations" goose down
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_shard" goose down
	@echo "migrate down complected"

migrate_status:
	@echo "Starting status"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIR="migrations" goose status
	@echo "status complected"

migrate_create:
	@echo "Starting create $(PARAM)"
	GOOSE_MIGRATION_DIR="migrations" goose create $(PARAM) sql
	@echo "create complected"

reshard:
	@echo "Starting resharding"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_script" goose -no-versioning up
	@echo "resharding complected"
fake:
	@echo "Starting fake"
	POSTGRES_PASSWORD=postgres docker-compose -f $(DEPLOY_DIR)/fake.yml up -d
	@echo "end deploy"
run: create-network
	@echo "Starting run"
	@echo "Starting deploy"
	NETWORK_NAME=$(NETWORK_NAME) POSTGRES_PASSWORD=postgres docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml up --scale worker=2 --scale app=3 -d
	@echo "end deploy"
down:
	@echo "Starting down"
	docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml down
	@echo "end down"
stop: down clean-network
build:
	@echo "Starting build"
	docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml build
	@echo "end build"
create-network:
	@if [ -z "$$(docker network ls -q -f name=${NETWORK_NAME})" ]; then \
		echo "Создание сети ${NETWORK_NAME}..."; \
		docker network create ${NETWORK_NAME}; \
	else \
		echo "Сеть ${NETWORK_NAME} уже существует"; \
	fi

clean-network:
	@echo "Удаление сети ${NETWORK_NAME}..."
	-docker network rm ${NETWORK_NAME} 2>/dev/null || true