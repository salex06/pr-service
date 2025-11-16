.PHONY: full-setup, quick-setup, fmt, lint, compose-up, compose-down

full-setup: fmt lint compose-up
quick-setup: compose-up

fmt:	
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted successfully!"

lint:
	@echo "Running linters..."
	golangci-lint run ./...

compose-up:
	@echo "Starting all services with docker-compose..."
	docker-compose up -d --build

compose-down:
	@echo "Stopping all services..."
	docker-compose down