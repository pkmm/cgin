all: gotool
	@go build -v
gotool:
	gofmt -s -w .
	swag init
help:
	@echo "make - compile the source code"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
.PHONY: gotool help