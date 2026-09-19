# TODOS

- [ ] Instance admin tooling for deleting articles etc.

- [ ] Parser:
  - Penalize comment sections
  - [ ] Refactor with goquery for speed improvements n shit
  - [x] Articles should also return:
    - [x] Author
    - [x] Description
    - [x] Publish Date
    - [ ] Last Updated Date
    - [ ] Estimated time to read
    - [x] Original URL

- [ ] Attaching tags to articles/content

- [ ] Queueing the article fetching work
- [ ] Adding option to add RSS feed

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
make db-migrate
sqlc generate
```

Compose starts PostgreSQL 17 on `127.0.0.1:5433` by default and keeps its data in
the `postgres_data` named volume. `make db-up` verifies that the host port is
available and that Docker published it before `make db-migrate` applies migrations.

### Environment variables

`.env.example` is the tracked template; `.env` and other `.env.*` files are ignored.
Docker Compose and Goose automatically read `.env` from the working directory.
The server also loads `.env`; its configuration requires `DATABASE_URL` and
`PORT`. When using the template, add the server port to your `.env`:

```dotenv
PORT=8080
```

| Variable | Default / purpose |
| --- | --- |
| `POSTGRES_USER` | `acrevus`, database role created on first startup |
| `POSTGRES_PASSWORD` | `acrevus_dev`, local development password |
| `POSTGRES_DB` | `acrevus`, database created on first startup |
| `POSTGRES_PORT` | `5433`, dedicated Acrevus host port |
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

Host ports are shared by every Docker project. Acrevus reserves `5433`; do not
reuse it in another Compose project. If Docker reports Postgres as healthy but
`make db-up` says its mapping is missing, run `docker compose down` (without
`-v`) and then `make db-up`. This recreates the container while preserving the
database volume.

Postgres initialization variables only apply to an empty data volume; editing
credentials in `.env` does not change an existing database role.

### Migrations and query generation

```sh
make db-status                  # Start and verify Postgres, then show migration status.
make db-migrate                 # Start and verify Postgres, then apply pending migrations.
goose -s create add_example sql  # Create the next numbered migration.
# Fill in both the +goose Up and +goose Down sections.
goose validate                  # Check migration file structure.
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

Generated queries include `CreateUser`, `GetUser`, `GetUserByEmail`, `ListUsers`,
`UpdateUser`, and `SoftDeleteUser`. Reads and updates exclude soft-deleted users. Update and
soft-delete queries set `updated_at`; direct SQL writers must also maintain it.
Hash passwords before passing them to create/update queries. The HTTP handlers
hash passwords automatically.

### User API

Start the server with `go run . server`. All user routes use `/api/v1`:

| Method | Path | Behavior |
| --- | --- | --- |
| `POST` | `/api/v1/users` | Create from required `name`, `email`, and `password`; returns 200 |
| `GET` | `/api/v1/users?limit=20&offset=0` | List active users in ID order; returns 200 |
| `GET` | `/api/v1/users/:id` | Get an active user; returns 200 |
| `PATCH` | `/api/v1/users/:id` | Update any of `name`, `email`, or `password`; returns 200 |
| `DELETE` | `/api/v1/users/:id` | Soft-delete the user; returns 204 with no body |

IDs must be positive integers. List pagination defaults to `limit=20`,
`offset=0`; the maximum limit is 100. `PATCH` requires at least one non-null field;
omitted fields are preserved, and supplied fields must be non-empty. Email values
must be valid email addresses. Passwords must be at most 72 bytes for bcrypt.

Create/get/update return `{"user": {...}}`; lists return
`{"users": [...], "limit": 20, "offset": 0}` (an empty list is `[]`). User responses
contain `id`, `name`, `email`, `created_at`, `updated_at`, and `deleted_at` (`null`
for active users), with no password/hash field. Errors use `{"error": "..."}`:
400 for invalid input, 404 for missing/deleted users, 409 for duplicate email,
and 500 for internal failures.

For example, updating a name does not require resubmitting the email or password:

```sh
curl -X PATCH http://localhost:8080/api/v1/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Updated Name"}'
curl -X DELETE http://localhost:8080/api/v1/users/1
```

Substitute your configured `PORT`. The Bruno create request in
`api-tests/create-user.yaml` uses its environment's `baseUrl`.

### User handler tests

`go test ./server` runs input-validation tests without Postgres. To also run the
full create/list/get/update/delete lifecycle, set `TEST_DATABASE_URL` to a migrated
Postgres database and run:

```sh
TEST_DATABASE_URL='postgres://acrevus:acrevus_dev@localhost:5432/acrevus?sslmode=disable' \
  go test ./server -run TestUserLifecycle -v
```

The lifecycle test copies `public.users` into a connection-local temporary table,
exercises the real generated queries, and drops the temporary data when its
connection closes. Without `TEST_DATABASE_URL`, this test is skipped.

### Stopping and checking

```sh
docker compose down             # Stop services, keeping database data.
docker compose down -v          # Delete the local database volume and all its data.
go build ./...
go test ./...
```
