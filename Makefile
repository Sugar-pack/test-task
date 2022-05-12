.PHONY: docker-run
docker-run:
	@docker-compose up --build -d --remove-orphans

.PHONY: docker-up
docker-up:
	@docker-compose up -d

.PHONY: docker-build
docker-build:
	@docker-compose build

lint: ## Run go lint
	golangci-lint run

test: ## Run tests
	go test ./...