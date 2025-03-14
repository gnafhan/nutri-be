include .env
OS := $(shell uname 2>/dev/null || echo Windows)

ifeq ($(OS), Windows)
	# Perintah untuk Windows (PowerShell)
	GET_ENV_VARS = powershell -Command "(Get-Content .env) -replace '=.*',''"
else
	# Perintah untuk Linux/macOS (sed)
	GET_ENV_VARS = sed 's/=.*//' .env
endif

start:
	@go run src/main.go
lint:
	@golangci-lint run
tests:
	@go test -v ./test/...
tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')
testsum:
	@cd test && gotestsum --format testname
swagger:
	@cd src && swag init