BASE_GO_COMMAND := @go

run:
	${BASE_GO_COMMAND} run main.go

test:
	${BASE_GO_COMMAND} test -v ./... -count=1

tidy:
	${BASE_GO_COMMAND} mod tidy

seed:
	${BASE_GO_COMMAND} run scripts/seed.go
