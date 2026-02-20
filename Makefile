.PHONY: help up down migrate test build logs

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start all services (db, redis, api, worker)
	docker compose up -d --build

down: ## Stop all services
	docker compose down

migrate: ## Run database migrations
	docker compose run --rm migrate

test: ## Run all Go tests
	cd backend && go test ./... -v

build: ## Build Go binaries locally
	cd backend && go build ./...

logs: ## Tail all service logs
	docker compose logs -f

db-shell: ## Open psql shell
	docker compose exec db psql -U assistant -d assistant
