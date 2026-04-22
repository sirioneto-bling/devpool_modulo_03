.PHONY: run build test test-e2e test-coverage lint swagger up-db down-db dev-init db-migrate db-shell clean

APP_NAME=devpool-base-web-api
APP_SRC=./cmd/web_api/
BIN_DIR=bin
APP_BIN=$(BIN_DIR)/web_api
DB_SERVICE=mysql
DB_USER=devpool
DB_PASS=devpool123
DB_NAME=devpool

# ========================= Development =========================

dev-init: up-db
	@echo "Copying example .env..."
	@cp -n .github/envs/localhost.env .env 2>/dev/null || true
	@echo "Downloading Go dependencies..."
	@go mod download
	@$(MAKE) db-migrate
	@echo "Dev environment ready! Run 'make run' to start the server."

run:
	@echo "Starting $(APP_NAME)..."
	@go run $(APP_SRC)

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(APP_BIN) $(APP_SRC)
	@echo "Binary created at $(APP_BIN)"

# ========================= Testing =========================

test:
	@echo "Running tests..."
	@go test ./... -v

test-e2e:
	@echo "Running e2e tests..."
	@go test ./internal/presentation/... -v

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# ========================= Documentation =========================

swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/web_api/main.go -o docs/swagger --parseInternal

# ========================= Quality =========================

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# ========================= Database =========================

up-db:
	@echo "Starting MySQL container..."
	@docker compose up -d --wait $(DB_SERVICE)

down-db:
	@echo "Stopping MySQL container..."
	@docker compose down

db-migrate:
	@echo "Applying database scripts..."
	@for f in scripts_db/*.sql; do \
		echo "  -> $$f"; \
		docker compose exec -T $(DB_SERVICE) mysql -u$(DB_USER) -p$(DB_PASS) $(DB_NAME) < $$f; \
	done

db-shell:
	@docker compose exec $(DB_SERVICE) mysql -u$(DB_USER) -p$(DB_PASS) $(DB_NAME)

# ========================= Cleanup =========================

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) coverage.out coverage.html
