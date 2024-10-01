include .env
export $(shell sed 's/=.*//' .env)

build:
	go build -o generator.exe ./cmd/generator/main.go

create:
	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/create_tables.sql

constraints:
	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/add_constraints.sql

copy:
	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/copy_data.sql

drop:
	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/drop_tables.sql

delete-all: build
	./generator.exe -d

generate: build
	./generator.exe -c $(RECORDS_PER_TABLE)

generate-csv: build
	./generator.exe -c $(RECORDS_PER_TABLE) -csv ./data/

quick-start: create constraints

start-db:
	docker compose up

delete-db:
	docker rm music-service-db
	docker volume rm music-service-db-volume

clear:
	rm -f ./generator.exe

