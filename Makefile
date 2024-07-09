build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))
	
migrate-up:
	@echo "Applying database migrations..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	@go run cmd/migrate/main.go down

migrate-force:
	@echo "Forcing migration version..."
	@go run cmd/migrate/main.go force $(version)
