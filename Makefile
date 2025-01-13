GO         := $(shell which go)
ROOT_DIR   := $(shell pwd)
DEPLOY_DIR = deploy
# Определяем переменную для хранения первого параметра
PARAM := $(word 2, $(MAKECMDGOALS))

up:
	@echo "Starting migrate"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=facebook sslmode=disable password=postgres port=5432 host=localhost" GOOSE_MIGRATION_DIR="migrations" goose up
	@echo "migrate complected"

down:
	@echo "Starting down"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=facebook sslmode=disable password=postgres port=5432 host=localhost" GOOSE_MIGRATION_DIR="migrations" goose down
	@echo "migrate down complected"

status:
	@echo "Starting status"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=facebook sslmode=disable password=postgres port=5432 host=localhost" GOOSE_MIGRATION_DIR="migrations" goose status
	@echo "status complected"
create:
	@echo "Starting create $(PARAM)"
	GOOSE_MIGRATION_DIR="migrations" goose create $(PARAM) sql
	@echo "create complected"
run:
	@echo "Starting run"
	@echo "Starting deploy"
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml up -d
	@echo "end deploy"