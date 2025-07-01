
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  help          - Show this help message"
	@echo "  check-prereqs - Check that Go, Docker, and docker compose are installed"
	@echo "  env-setup     - Setup the .env file based on .env.sample (with confirmation prompt)"
	@echo "  clean         - Stop docker compose services and remove volumes"
	@echo "  start         - Run the VaultStream project"
	@echo "  tests         - Run tests with Cargo"



.PHONY: quick-start
quick-start: env-setup check-prereqs start

.PHONY: start
start: stop db.seed nats.start
	@go run ./keys-service & \
	go run ./records-service & \
	go run ./signing-service & \
	wait

ifneq (,$(wildcard .env))
    include .env
    export
endif
test: stop db.setup nats.start
	GOPROXY=https://proxy.golang.org,direct go test ./keys-service
	go test ./records-service
	go test ./signing-service


.PHONY: stop
stop:
	docker compose down
	rm -rf ./database/data
	rm -rf ./nats/data

# Environment setup: create or overwrite .env from .env.sample
.PHONY: env-setup
env-setup:
	@if [ -f .env ]; then \
		read -p ".env file exists. Overwrite? (y/N): " ans; \
		if [ "$$ans" = "y" ]; then \
			echo "Overwriting .env file..."; \
			cp .env.sample .env; \
		else \
			echo "Aborted. Please configure your .env manually based on .env.sample."; \
			exit 1; \
		fi; \
	else \
		echo "Creating .env file from .env.sample..."; \
		cp .env.sample .env; \
	fi

.PHONY: db.start
db.start: check-prereqs
	docker compose up -d db
	@echo "Waiting for Postgres' readiness confirmation..."
	@until docker compose exec -T db pg_isready -U "$${DOCKER_COMPOSE_POSTGRES_USER}"; do \
	    echo "Postgres not ready yet..."; \
	    sleep 1; \
	done

.PHONY: nats.start
nats.start: check-prereqs
	docker compose up -d nats
	@echo "Waiting for NATS readiness confirmation..."
	@until curl -s http://localhost:8222/varz | grep 'server_id' ; do \
	    echo "NATS not ready yet..."; \
	    sleep 1; \
	done
	@echo "NATS is ready!"

.PHONY: db.setup
db.setup: db.start 
	go generate ./database

.PHONY: db.seed
db.seed: db.setup
	go run ./seeder



.PHONY: sync-deps
sync-deps:
	@for d in database seeder types; do \
		echo "Tidying module in $$d..."; \
		(cd $$d && go mod tidy) || exit 1; \
	done

check-prereqs:
	@command -v docker >/dev/null 2>&1 || { \
		echo >&2 "Error: Docker not found. Please install Docker from https://docs.docker.com/get-docker/"; exit 1; }
	@docker compose version >/dev/null 2>&1 || { \
		echo >&2 "Error: 'docker compose' not available. Please ensure you have Docker 20.10 or later."; exit 1; }
	@command -v go >/dev/null 2>&1 || { \
		echo >&2 "Error: Go not found. Please install Go from https://golang.org/dl/"; exit 1; }
