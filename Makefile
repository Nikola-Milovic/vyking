# Load environment variables from .env file
include .env
export

.PHONY: help
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: dev-setup
dev-setup: clean ## Setup development environment
	docker-compose up -d
	./scripts/mysql_isready.sh
	make migrate-up

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	docker run --rm -v $(PWD)/migrations:/migrations --network vyking_internal migrate/migrate \
		-path=/migrations \
		-database="mysql://$(DB_USER):$(DB_PASSWORD)@tcp(vyking-mysql:3306)/$(DB_NAME)" \
		up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	docker run --rm -v $(PWD)/migrations:/migrations --network vyking_internal migrate/migrate \
		-path=/migrations \
		-database="mysql://$(DB_USER):$(DB_PASSWORD)@tcp(vyking-mysql:3306)/$(DB_NAME)" \
		down

.PHONY: migrate-force
migrate-force: ## Force migration to specific version
	docker run --rm -v $(PWD)/migrations:/migrations --network vyking_internal migrate/migrate \
		-path=/migrations \
		-database="mysql://$(DB_USER):$(DB_PASSWORD)@tcp(vyking-mysql:3306)/$(DB_NAME)" \
		force $(VERSION)

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Please provide a migration name: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%s); \
	@SEQ=$$(ls -1 migrations/*.up.sql 2>/dev/null | wc -l | xargs); \
	@SEQ=$$(printf "%03d" $$((SEQ + 1))); \
	touch migrations/$${SEQ}_$(NAME).up.sql; \
	touch migrations/$${SEQ}_$(NAME).down.sql; \
	echo "Created migration files:"; \
	echo "  migrations/$${SEQ}_$(NAME).up.sql"; \
	echo "  migrations/$${SEQ}_$(NAME).down.sql"

.PHONY: up
up: ## Start all services with docker-compose
	docker-compose up -d

.PHONY: down
down: ## Stop all services
	docker-compose down

.PHONY: logs
logs: ## Show logs from all services
	docker-compose logs -f

.PHONY: clean
clean: ## Stop services and remove volumes
	docker-compose down -v

.PHONY: mysql
mysql: ## Connect to MySQL database 
	docker exec -it vyking-mysql mysql -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -cover ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: generate-mocks
generate-mocks: ## Generate mocks for interfaces
	@mkdir -p internal/clients/mock
	mockgen -source=internal/domain/apiclients.go -destination=internal/clients/mock/mocks.go -package=mock

.PHONY: install-tools
install-tools: ## Install required tools
	go install github.com/golang/mock/mockgen@v1.6.0

