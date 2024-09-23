include .env
export $(shell sed 's/=.*//' .env)

build:
	go build ./cmd/main.go -o generator.exe

create:
	psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/create_tables.sql

constraints:
	psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/add_constraints.sql

copy:
	psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/copy_data.sql

drop:
	psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f ./sql/drop_tables.sql

clear: build
	./generator.exe -d

generate: 
	./generator.exe -c $(RECORDS_PER_TABLE)

generate-csv:
	./generator.exe -c $(RECORDS_PER_TABLE) -csv ./data/output.csv