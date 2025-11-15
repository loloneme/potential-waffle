.PHONY: help generate-schema docker-build build run generate-mocks test test-docker-up test-docker-down

help:
	@echo Available targets:
	@echo generate-schema - Generate OpenAPI types and server from openapi.yaml
	@echo test - Run all tests (including integration tests with Docker)

generate-schema:
	@echo Generating OpenAPI schema code...
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate "types,echo-server,spec" -o internal/generated/openapi/openapi.gen.go -package generated api/openapi.yaml
	@echo Successfully generated OpenAPI schema in generated/openapi

docker-build:
	@echo Building Docker Image...
	docker-compose up --build -d

docker-test-build:
	@echo Building Docker Test Containers...
	@docker-compose -f docker-compose.test.yml up --build -d
	@timeout /t 5 >nul
	@echo Test Docker containers are ready

build:
	@echo Building service...
	go build -o bin/reviewers_app.exe ./cmd/service

run:
	@echo Starting service...
	go run ./cmd/service/main.go

generate-mocks:
	@echo Generating mocks...
	go generate ./...
	@echo Mocks generated successfully

docker-test-down:
	@echo Stopping test Docker containers...
	@docker-compose -f docker-compose.test.yml down -v
	@echo Test Docker containers stopped

test: docker-test-build
	@echo Running all tests...
	@go test ./... -v -count=1 -p 1
	@make docker-test-down
