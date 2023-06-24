BASE_GO_COMMAND := go

run:
	@${BASE_GO_COMMAND} run main.go

tidy:
	@${BASE_GO_COMMAND} mod tidy
