# Repository guidance

## Entrypoints and flow
- Run from the root: `go run . fetch <url>`, `go run . analyze <url>`
  (alias `a`), `go run . ui`, or `go run . server`.
- `main.go` initializes article storage before dispatching any Cobra command,
  including help and server startup.
- Fetching uses `fetcher.generalParserAnalyze` → `analyzer.ParseContentDensity`
  → `storage.SaveArticle`. The other site parsers are not the active fetch path.
- `analyze` extracts content without adding it to the saved article collection.
  `repl/` is the Bubble Tea reader for that collection.

## Runtime gotchas
- Article storage is hardcoded to `$HOME/.local/share/acrevus/`; it does not
  honor `XDG_DATA_HOME`. Metadata lives in `entries.json`, HTML in `articles/`.
  Use an isolated HOME when exercising storage destructively.
- Fetch and analyze use Rod with a live browser and network access.
  Both write `./temp.html` through the analyzer; this file is not gitignored.
- User routes are registered under `/api/v1` via `ApiConfig.registerUserRoutes`.
  `PATCH /users/:id` preserves omitted fields atomically in SQL; DELETE is soft.
  Public responses use `UserResponse`, not the generated model (which has a hash).

## Database workflow
- See `README.md` for environment variables and tool installation. Compose,
  Goose, and the server load root `.env`; the server requires `DATABASE_URL` and `PORT`.
- sqlc targets `database/sql`; use `*sql.DB` with the `lib/pq` driver and
  `sql.Open("postgres", ...)`. Generated transactions use `*sql.Tx`.
- Start Postgres with `docker compose up -d --wait postgres`, then `goose up`.
  Create numbered migrations with `goose -s create <name> sql`.
- `db/migrations/` is the schema source for Goose and sqlc; edit queries in
  `db/queries/`, then run `sqlc generate`. Commit generated `internal/database/`
  files; do not hand-edit them. Generation does not require a live database.
- User queries exclude soft-deleted rows and maintain `updated_at`; direct SQL
  writes must maintain it too. `password` stores a hash supplied by the caller.

## Verification
- Use `go build ./...` and `go test ./...` for baseline checks;
  `go test ./server` runs user input-validation tests without Postgres.
- `TestUserLifecycle` requires `TEST_DATABASE_URL` pointing to a migrated Postgres
  database. It uses a session-local temporary copy of `public.users`; keep its
  connection pool limited to one connection so every query sees the same table.
- `cmd/fetch.go` currently contains an invalid `%i` printf verb, which can
  cause the vet phase of `go test ./...` to fail.
- For database changes, also run `goose validate` and `sqlc vet`; exercise Goose
  up/down against a disposable database when changing migrations.
