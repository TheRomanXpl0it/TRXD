END		= \033[0m
YELLOW	= \033[33m


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

py-test:
	@go build .
	@cd ../tests/backend/ && (./run.sh || (kill $$(pidof trxd) && exit 1))

lint:
	@go vet ./...
	@ineffassign ./...
	@gocyclo -ignore "sqlc|_test|test_" -over 20 .
	@misspell -error .
	@golangci-lint run

cyclo:
	@gocyclo -ignore "sqlc|_test|test_" -top 15 .

pipe: lint test py-test
