.PHONY: build run stop clean test

build:
	docker-compose build

up:
	docker-compose up -d

run:
	docker-compose up

down:
	docker-compose down

clean:
	docker-compose down --rmi all

logs:
	docker-compose logs -f

test:
	go test ./api_test

test-coverage:
	go test ./api_test -coverprofile=coverage.out
	go tool cover -html=coverage.out

dev:
	go run main.go

docker-build:
	docker build -t stasshander/lead-management:latest .