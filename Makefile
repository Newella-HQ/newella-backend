ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: migrateUp
migrateUp:
	docker run -v ./migrations:/migrations --network appnet migrate/migrate -path=./migrations \
		-database 'postgres://postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=$(POSTGRES_SSLMODE)' up

.PHONY: migrateDown
migrateDown:
	docker run -v ./migrations:/migrations --network appnet migrate/migrate -path=./migrations \
		-database 'postgres://postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=$(POSTGRES_SSLMODE)' down
