ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: migrateUp
migrateUp:
	docker run --rm -v ./migrations:/migrations --network appnet migrate/migrate -path=./migrations \
		-database 'postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=$(POSTGRES_SSLMODE)' up

.PHONY: migrateDown
migrateDown:
	yes | docker run --rm -i -v ./migrations:/migrations --network appnet migrate/migrate -path=./migrations \
		-database 'postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=$(POSTGRES_SSLMODE)' down

.PHONY: migrateDrop
migrateDrop:
	yes | docker run --rm -i -v ./migrations:/migrations --network appnet migrate/migrate -path=./migrations \
		-database 'postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=$(POSTGRES_SSLMODE)' drop

.PHONY: up
up:
	docker compose up -d --build

.PHONY: down
down:
	docker compose down -v