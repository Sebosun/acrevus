# About

Article fetcher / parser

Stack:

- Rod
- Bubbletea (cli text display)
- Parsing cmd args with cobra
- PostgreSQL (Docker Compose), sqlc (typed Go queries), Goose (migrations)

By default all data fetched is saved in
`~/.local/share/acrevus/`

Within `~/.local/share/acrevus/entries.json` you can find JSON containing all grabbed articles.
Postgres holds users; fetched articles currently use the filesystem storage above.

## Local database setup

Requires Go 1.26+, Docker with Compose, and a running Docker daemon. Install the
database tools (versions used for this setup):

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0
go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
export PATH="$(go env GOPATH)/bin:$PATH"
```

Run all commands below from the repository root:

```sh
cp .env.example .env # First-time setup; edit .env for your local settings.
docker compose up -d --wait postgres
goose up
sqlc generate
```

Compose starts PostgreSQL 17 on `127.0.0.1:5432` by default and keeps its data in
the `postgres_data` named volume. Migrations are applied explicitly with Goose,
after Postgres is healthy.

### Environment variables

`.env.example` is the tracked template; `.env` and other `.env.*` files are ignored.
Docker Compose and Goose automatically read `.env` from the working directory.
The server also loads `.env`; its current configuration requires `DB_URL` and
`PORT`. When using the template, add these below `DATABASE_URL` in your `.env`:

```dotenv
DB_URL=${DATABASE_URL}
PORT=8080
```

| Variable | Default / purpose |
| --- | --- |
| `POSTGRES_USER` | `acrevus`, database role created on first startup |
| `POSTGRES_PASSWORD` | `acrevus_dev`, local development password |
| `POSTGRES_DB` | `acrevus`, database created on first startup |
| `POSTGRES_PORT` | `5432`, published host port; change in `.env` if already in use |
| `DATABASE_URL` | Host-side Postgres URL, expanded from the values above |
| `GOOSE_DRIVER` | `postgres` |
| `GOOSE_DBSTRING` | Expanded from `DATABASE_URL`; Goose's connection string |
| `GOOSE_MIGRATION_DIR` | `./db/migrations` |

The example URL uses `localhost` and `sslmode=disable` for host-side local
development. A client running in the same Compose network would use
`postgres:5432`. URL-encode credentials containing URL-special characters in
`DATABASE_URL`. For a different environment file, use
`docker compose --env-file .env.test up -d --wait postgres` and
`goose -env .env.test up`.

Postgres initialization variables only apply to an empty data volume; editing
credentials in `.env` does not change an existing database role.

### Migrations and query generation

```sh
goose status                    # Show applied/pending migrations.
goose -s create add_example sql  # Create the next numbered migration.
# Fill in both the +goose Up and +goose Down sections.
goose validate                  # Check migration file structure.
goose up                        # Apply pending migrations.
sqlc generate                   # Regenerate Go code after SQL changes.
sqlc vet                        # Check SQL/configuration for errors.
goose down                      # Roll back the most recent migration (may drop data).
```

- `db/migrations/` is the schema source for both Goose and sqlc. sqlc understands
  Goose's up/down annotations and does not need a running database to generate code.
- Write named queries in `db/queries/`; `sqlc.yaml` generates the `database`
  package in `internal/database/`, using `database/sql`. Commit generated files along
  with SQL changes; do not edit generated Go files directly.
- Use `sql.Open("postgres", connectionURL)` with the `github.com/lib/pq` driver,
  then pass the resulting `*sql.DB` to `database.New(db)`. Use `queries.WithTx(tx)`
  for a `*sql.Tx` transaction. Call `db.PingContext(ctx)` to verify connectivity;
  `sql.Open` alone does not establish a connection.

### Users

The initial migration creates `users` with:

| Column | Definition |
| --- | --- |
| `id` | Auto-generated `BIGINT` primary key |
| `name` | Required text |
| `email` | Required, unique text (case-sensitive; remains reserved after soft deletion) |
| `password` | Required text for a password hash, never plaintext |
| `created_at` | Required `TIMESTAMPTZ`, defaults to `NOW()` |
| `updated_at` | Required `TIMESTAMPTZ`, defaults to `NOW()` |
| `deleted_at` | Nullable `TIMESTAMPTZ`; `NULL` means active |

Generated queries include `CreateUser`, `GetUser`, `GetUserByEmail`, `UpdateUser`,
and `SoftDeleteUser`. Reads and updates exclude soft-deleted users. Update and
soft-delete queries set `updated_at`; direct SQL writers must also maintain it.
Hash passwords before passing them to create/update queries. The current
`POST /users` handler validates input and returns a placeholder response; it does
not yet insert users.

### Stopping and checking

```sh
docker compose down             # Stop services, keeping database data.
docker compose down -v          # Delete the local database volume and all its data.
go build ./...
go test ./...
```
