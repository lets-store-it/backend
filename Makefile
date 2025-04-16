.PHONY: generate
generate:
	@echo "Bundling OpenAPI specifications..."
	@redocly bundle api/openapi/openapi.yaml -o generated/openapi.yaml
	@echo "Running Go generate..."
	@go generate ./...
	@echo "Running sqlc generate..."
	@sqlc generate
