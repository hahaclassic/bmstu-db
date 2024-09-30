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

clear: build
	./generator.exe -d

generate: 
	./generator.exe -c $(RECORDS_PER_TABLE)

generate-csv:
	./generator.exe -c $(RECORDS_PER_TABLE) -csv ./data/output.csv