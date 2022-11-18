POSTGRES_CONTAINER_NAME = balance-app-db
POSTGRES_PASSWORD = qwerty
PORTS = 5436:5432

build:
	docker-compose build balance-app
	
run:
	docker-compose up balance-app

postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -p $(PORTS) -d --rm postgres

stop-postgres:
	docker stop $(POSTGRES_CONTAINER_NAME)

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' down

swag:
	~/go/bin/swag init -g cmd/main/go