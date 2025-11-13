.PHONY: help generate-schema

help:
	@echo Available targets:
	@echo generate-schema - Generate OpenAPI types and server from openapi.yaml

generate-schema:
	@echo Generating OpenAPI schema code...
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate "types,echo-server,spec" -o internal/generated/openapi/openapi.gen.go -package generated api/openapi.yaml
	@echo Successfully generated OpenAPI schema in generated/openapi
