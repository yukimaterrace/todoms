.PHONY: migrate-up migrate-down run

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
