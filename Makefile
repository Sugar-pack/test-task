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

test-coverage: ## Run go test with coverage
	go test ./... -coverprofile=coverage.out `go list ./...`

gen-repo-mock:
	@docker run -v `pwd`:/src -w /src vektra/mockery --case snake --dir internal/repository --output internal/mocks/repository --outpkg repository --all

gen-qualifier-mock:
	@docker run -v `pwd`:/src -w /src vektra/mockery --case snake --dir internal/api --output internal/mocks/qualifier --outpkg qualifier --all


