include .env
export

local:
	@echo "Setting MYSQL_DB_URL to LOCAL_DB_URL"
	@echo MYSQL_DB_URL=$(LOCAL_DB_URL) >> .env
	@make connect

hosted:
	@echo "Setting MYSQL_DB_URL to HOSTED_DB_URL"
	@echo MYSQL_DB_URL=$(HOSTED_DB_URL) >> .env
	@make connect

connect:
	@go run main.go
