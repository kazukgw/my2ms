.DEFAULT_GOAL := help

help: ## show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## run test
	@go test

reset: ## reset test suite
	@echo rm mysql/app/var/master.info
	@if [[ -f mysql/app/var/master.info ]]; then rm mysql/app/var/master.info; fi
	docker-compose kill
	docker-compose rm
	docker-compose up -d
	@echo
	@echo sleep 10s ...
	@sleep 10
	@echo
	go run main.go models.go -action reset

insert: ## insert test data
	go run main.go models.go -action insert

update:
	go run main.go models.go -action update

delete:
	go run main.go models.go -action delete
