.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-force

ifneq (,$(wildcard .env))
    include .env
    export
endif

DATABASE_DSN ?= "host=localhost user=postgres password=postgres dbname=travel_db port=5432 sslmode=disable"
MIGRATIONS_DIR = ./migrations


migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: nome é obrigatório. Use 'make migrate-create name=nome_da_migration'"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATIONS_DIR)
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	up_file="$(MIGRATIONS_DIR)/$${timestamp}_$(name).up.sql"; \
	down_file="$(MIGRATIONS_DIR)/$${timestamp}_$(name).down.sql"; \
	echo "-- Migration: $(name) (up)" > $${up_file}; \
	echo "-- Adicione suas consultas SQL aqui" >> $${up_file}; \
	echo "" >> $${up_file}; \
	echo "-- Migration: $(name) (down)" > $${down_file}; \
	echo "-- Adicione suas consultas SQL para reverter a migration aqui" >> $${down_file}; \
	echo "" >> $${down_file}; \
	echo "Migration criada: $${up_file} e $${down_file}"

migrate-up:
	go run cmd/migrate/migrate.go up

migrate-down:
	go run cmd/migrate/migrate.go down

migrate-down-steps:
	@if [ -z "$(steps)" ]; then \
		echo "Error: steps é obrigatório. Use 'make migrate-down-steps steps=N'"; \
		exit 1; \
	fi
	go run cmd/migrate/migrate.go down $(steps)

migrate-status:
	go run cmd/migrate/migrate.go status

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: version é obrigatório. Use 'make migrate-force version=20230101120000'"; \
		exit 1; \
	fi
	go run cmd/migrate/migrate.go force $(version)