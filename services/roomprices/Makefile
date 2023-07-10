BASE_GO_COMMAND := @go

build:
	@set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && ${BASE_GO_COMMAND} build -o bin/app

run:
	${BASE_GO_COMMAND} run .

test:
	${BASE_GO_COMMAND} test -v ./... -count=1

tidy:
	${BASE_GO_COMMAND} mod tidy

compile_proto:
	protoc --go_out=rpc --go_opt=paths=source_relative \
	       --go-grpc_out=rpc --go-grpc_opt=paths=source_relative \
	       roomprices.proto
