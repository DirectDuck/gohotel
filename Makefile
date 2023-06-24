BASE_GO_COMMAND := @go

build:
	${BASE_GO_COMMAND} build -o bin/app

run:
	${BASE_GO_COMMAND} run main.go

test:
	${BASE_GO_COMMAND} test -v ./...

tidy:
	${BASE_GO_COMMAND} mod tidy
