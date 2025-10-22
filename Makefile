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
# Freemium Trial System Tests
test-freemium:
	@echo "🧪 Running Freemium Trial System Tests..."
	@echo "Note: Some tests require database connection"
	@go test -v ./test/unit/middleware/freemium_or_access_test.go
	@echo "⚠️  Service and Integration tests require database - run with database connected"
test-freemium-unit:
	@echo "🔬 Running Freemium Unit Tests (No Database Required)..."
	@go test -v ./test/unit/middleware/freemium_or_access_test.go -run TestFreemiumOrAccess
test-freemium-integration:
	@echo "🔗 Running Freemium Integration Tests (Requires Database)..."
	@echo "⚠️  Make sure PostgreSQL is running on localhost:5432"
	@go test -v ./test/integration/freemium_test.go -run TestFreemiumFlow
test-freemium-middleware:
	@echo "🛡️ Running Freemium Middleware Tests (No Database Required)..."
	@go test -v ./test/unit/middleware/freemium_or_access_test.go -run TestFreemiumOrAccess
test-freemium-service:
	@echo "⚙️ Running Freemium Service Tests (Requires Database)..."
	@echo "⚠️  Make sure PostgreSQL is running on localhost:5432"
	@go test -v ./test/unit/service/subscription_service_test.go -run TestCreateFreemiumSubscription
test-freemium-no-db:
	@echo "🚀 Running Freemium Tests That Don't Require Database..."
	@go test -v ./test/unit/middleware/freemium_or_access_test.go -run TestFreemiumOrAccess
test-help:
	@echo "📋 Available Test Commands:"
	@echo "  make tests                    - Run all tests"
	@echo "  make test-freemium           - Run freemium tests (middleware only)"
	@echo "  make test-freemium-no-db     - Run freemium tests without database"
	@echo "  make test-freemium-middleware - Run middleware tests only"
	@echo "  make test-freemium-unit      - Run unit tests (middleware only)"
	@echo "  make test-freemium-service   - Run service tests (requires database)"
	@echo "  make test-freemium-integration - Run integration tests (requires database)"
	@echo ""
	@echo "💡 Database Required Tests:"
	@echo "  - Service tests need PostgreSQL running on localhost:5432"
	@echo "  - Integration tests need PostgreSQL running on localhost:5432"
	@echo "  - Middleware tests run without database (using mocks)"
swagger:
	@cd src && swag init