.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:
	@GOOS=linux GOARCH=amd64 go build -o my2ms

build_for_test:
	@GOOS=linux GOARCH=amd64 go build -o test/mysql/app/my2ms

test:
	@go test
