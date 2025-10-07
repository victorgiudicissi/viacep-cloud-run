.PHONY: run

help:
	@echo "Available targets:"
	@echo ""
	@echo "run - Run the application"
	@echo "test - Run tests"

run:
	docker compose up --build

test:
	go test ./... -v
