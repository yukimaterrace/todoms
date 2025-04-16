.PHONY: migrate-up migrate-down run test test-verbose test-coverage

# Define variables
DB_URL=postgres://admin:admin@localhost:5432/todoms?sslmode=disable
MIGRATION_PATH=file://migration

# Run database migrations up
migrate-up:
	@echo "Running database migrations up..."
	docker run --rm --network host \
		-v $(shell pwd)/migration:/migration \
		migrate/migrate \
		-path=/migration -database "$(DB_URL)" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	docker run --rm --network host \
		-v $(shell pwd)/migration:/migration \
		migrate/migrate \
		-path=/migration -database "$(DB_URL)" down

# Start the application
run:
	@echo "Starting the application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -count=1 ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
