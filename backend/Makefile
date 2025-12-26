END				= \033[0m

GREY			= \033[30m
RED				= \033[31m
GREEN			= \033[32m
YELLOW			= \033[33m
BLUE			= \033[34m
PURPLE			= \033[35m
CYAN			= \033[36m

HIGH_RED		= \033[91m


default:
	@echo "${YELLOW}[Warn] This makefile is not to built, but for testing purposes only.${END}"

sqlc:
	@sqlc generate
sql: sqlc

run:
	@go run . -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	@go run . -t
	@go run .

test:
	@REDIS_DISABLE=true go test -v ./...

lint:
	@go vet ./...
	@ineffassign ./...
	@gocyclo -ignore "sqlc|_test|test_" -over 20 .
	@misspell -error .
	@golangci-lint run

cyclo:
	@gocyclo -ignore "sqlc|_test|test_" -top 15 .
