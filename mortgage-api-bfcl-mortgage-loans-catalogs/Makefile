BUILDPATH=$(CURDIR)

build:
	@echo "Creando Binario ..."
	@mkdir -p $(BUILDPATH)/build/bin
	@go build -ldflags '-s -w' -o $(BUILDPATH)/build/bin/dist ./cmd/api
	@echo "Binario generado en build/bin/dist"

test:
	@echo "Ejecutando tests..."
	@go test ./...

coverage:
	@echo "Coverfile..."
	go test -coverprofile=coverfile_out ./...
	@go tool cover -func coverfile_out
	@go tool cover -func coverfile_out | grep total | grep -o '[0-9]*\.[0-9]*' | cut -d' ' -f1 > coverage.txt

.PHONY: test build coverage

