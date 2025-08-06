.PHONY: help up up_payment_processor clean test

## help: Show this help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'

## clean: Stop and remove all containers and networks
clean:
	@echo "Parando e removendo containers existentes..."
	@docker compose -f ./payment-processor/docker-compose.yml down --remove-orphans
	@docker compose down --remove-orphans

## up_payment_processor: Start payment processor service
up_payment_processor:
	@docker compose -f ./payment-processor/docker-compose.yml up -d

## up: Clean and start all services
up: clean up_payment_processor
	@docker compose up -d

## test: Run k6 tests
test: up
	@command -v k6 >/dev/null 2>&1 || { \
		echo "k6 não está instalado. acesse https://grafana.com/docs/k6/latest/set-up/install-k6/"; \
		exit 1; \
	}
	k6 run ./rinha-test/rinha.js
