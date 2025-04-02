GO         := $(shell which go)
ROOT_DIR   := $(shell pwd)
DEPLOY_DIR = deploy
# Определяем переменную для хранения первого параметра
PARAM := $(word 2, $(MAKECMDGOALS))
DB_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=6432 host=localhost"
DB_SHARD_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=7432 host=localhost"

up:
	@echo "Starting migrate"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIR="migrations" goose up
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_shard" goose up
	@echo "migrate complected"

down:
	@echo "Starting down"

	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIR="migrations" goose down
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_shard" goose down
	@echo "migrate down complected"

reshard:
	@echo "Starting resharding"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_SHARD_STRING) GOOSE_MIGRATION_DIR="migrations_script" goose -no-versioning up
	@echo "resharding complected"

status:
	@echo "Starting status"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_STRING) GOOSE_MIGRATION_DIR="migrations" goose status
	@echo "status complected"
create:
	@echo "Starting create $(PARAM)"
	GOOSE_MIGRATION_DIR="migrations" goose create $(PARAM) sql
	@echo "create complected"
run:
	@echo "Starting run"
	@echo "Starting deploy"
	POSTGRES_PASSWORD=postgres docker-compose -f $(DEPLOY_DIR)/docker-compose.yml up --scale worker=2 -d
	@echo "end deploy"
stop:
	@echo "Starting down"
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml down
	@echo "end down"
build:
	@echo "Starting build"
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml build
	@echo "end build"