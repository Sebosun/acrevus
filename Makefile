POSTGRES_PORT := 5433
POSTGRES_ADDRESS := 127.0.0.1:$(POSTGRES_PORT)


APP := acrevus
.PHONY: db-up db-migrate db-status build bra server

db-up:
	@mapped="$$(docker compose port postgres 5432 2>/dev/null || :)"; \
	if [ "$$mapped" != "$(POSTGRES_ADDRESS)" ] && ss -H -ltn "sport = :$(POSTGRES_PORT)" | read -r _; then \
		printf '%s\n' "Host port $(POSTGRES_PORT) is already in use. Stop its owner before starting Acrevus."; \
		exit 1; \
	fi
	docker compose up -d --wait postgres
	@mapped="$$(docker compose port postgres 5432 2>/dev/null || :)"; \
	if [ "$$mapped" != "$(POSTGRES_ADDRESS)" ]; then \
		printf '%s\n' "Postgres is running without the expected $(POSTGRES_ADDRESS) mapping. Run 'docker compose down' (without -v), then 'make db-up'."; \
		exit 1; \
	fi

db-migrate: db-up
	goose up

db-status: db-up
	goose status

build: 
	go build -o bin/$(APP) .

server:
	go build -o bin/$(APP) . && ./bin/${APP} server

# run mock
bra:
	go build -o bin/$(APP) . && ./bin/${APP} r
