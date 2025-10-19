migration_create:
	goose -s --dir ./migrations create $(word 2, $(MAKECMDGOALS)) sql

swag-gen:
	swag init -g cmd/app/main.go -o docs