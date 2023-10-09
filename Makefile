start:
	nodemon --exec go run main.go serve --signal SIGTERM

test:
	go test -v -cover ./...

mocks:
	mockery --all --keeptree --recursive=true --outpkg=mocks --output ./mocks

linter:
	golangci-lint run --out-format html > golangci-lint.html

migrate-file:
	migrate create -ext sql -dir migrations $(name)

build:
	go build -o order-service