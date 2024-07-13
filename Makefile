.PHONY: build run clean

build:
	go build -o directory_logger cmd/main.go

run:
	go run cmd/main.go

depend:
	go mod download && go mod verify

lint:
	golangci-lint run

migrate:
	goose -dir ./migrations mysql "root:mysql_password1@tcp(mysql:3306)/directory_logger" up

down:
	goose -dir ./migrations mysql "root:mysql_password1@tcp(mysql:3306)/directory_logger" down

reset:
	goose -dir ./migrations mysql "root:mysql_password1@tcp(mysql:3306)/directory_logger" reset

clean: reset
	rm directory_logger
