mod:
	@echo "Vendoring..."
	@go mod vendor

test:
	@echo "Ejecutando tests..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

coverage:
	@echo "Coverage..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

.PHONY: clean build test

