GO         := $(shell which go)
ROOT_DIR   := $(shell pwd)
DEPLOY_DIR = deploy
# Определяем переменную для хранения первого параметра
PARAM := $(word 2, $(MAKECMDGOALS))
DB_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=6432 host=localhost"
DB_SHARD_STRING := "user=postgres dbname=facebook sslmode=disable password=postgres port=7432 host=localhost"

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
run:
	@echo "Starting run"
	@echo "Starting deploy"
	POSTGRES_PASSWORD=postgres docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml up --scale worker=2 --scale app=3 -d
	@echo "end deploy"
stop:
	@echo "Starting down"
	docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml down
	@echo "end down"
build:
	@echo "Starting build"
	docker-compose -f $(DEPLOY_DIR)/app.yml -f $(DEPLOY_DIR)/db.yml -f $(DEPLOY_DIR)/cache.yml -f $(DEPLOY_DIR)/kafka.yml -f $(DEPLOY_DIR)/sre.yml build
	@echo "end build"