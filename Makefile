.PHONY: help generate-schema

help:
	@echo Available targets:
	@echo generate-schema - Generate OpenAPI types and server from openapi.yaml

generate-schema:
	@echo Generating OpenAPI schema code...
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate "types,echo-server,spec" -o internal/generated/openapi/openapi.gen.go -package generated api/openapi.yaml
	@echo Successfully generated OpenAPI schema in generated/openapi

docker-build:
	@echo Building Docker Image...
	docker-compose up --build -d

build:
	@echo "Building service..."
	go build -o bin/reviewers_app.exe ./cmd/service

run:
	@echo "Starting service..."
	go run ./cmd/service/main.go