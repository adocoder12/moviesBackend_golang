# moviesBackend_golang

A small REST API for movies, built in Go with SQLite. It is a learning project: the goal is to practice idiomatic Go, clean project structure and the basics of `database/sql`.

## Stack

| Tool                                | Purpose                                                             |
| ----------------------------------- | ------------------------------------------------------------------- |
| Go `net/http`                       | HTTP server and routing (Go 1.22+ patterns like `GET /movies/{id}`) |
| `log/slog`                          | Structured logging from the standard library                        |
| `database/sql` + `mattn/go-sqlite3` | SQLite access (needs CGO / a C compiler)                            |
| `joho/godotenv`                     | Loads `.env` into environment variables                             |
| Air (`air-verse/air`)               | Live reload during development                                      |

## Project structure

```
moviesBackend_golang/
├── cmd/api/main.go        # wiring only: load env, open DB, build app, start server
├── internal/
│   ├── db/                # opens the connection, creates the schema
│   ├── model/             # Movie struct (with JSON tags)
│   ├── repository/        # SQL queries (only place that knows SQL)
│   ├── services/          # business logic, defines the storage interface it needs
│   └── handler/           # HTTP layer: routes, handlers, error helpers
├── .air.toml              # live reload config
├── .env                   # local config (never committed)
└── go.mod
```

### How a request flows

```
HTTP request
  -> handler   (parse request, write response)
  -> service   (business rules)
  -> repository (SQL)
  -> SQLite
```

`main.go` builds the chain from the bottom up and passes each piece into the next:

```
conn := db.Connect(...)
repo := repository.NewMoviesRepository(conn)
svc  := services.NewMoviesService(repo, logger)
app  := handler.NewApplication(logger, svc)
server.Handler = app.SetupRoutes()
```

## Running it

```
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin     # add to ~/.zshrc

go mod tidy
air
```

`.env` example:

```
PORT=8080
APP_ENV=development
DB_PATH=movies.db
```

Commit a `.env.example` with fake values, and keep `.env`, `tmp/`, `*.db` and `build-errors.log` in `.gitignore`.

### Try it

```
curl -i localhost:8080/health
curl -i localhost:8080/movies
curl -i localhost:8080/banana        # 404
curl -i -X POST localhost:8080/health   # 405, route is GET only
```

## Endpoints

| Method | Path      | Description                |
| ------ | --------- | -------------------------- |
| GET    | `/`       | Home message               |
| GET    | `/health` | Health check, returns `ok` |
| GET    | `/movies` | List all movies as JSON    |

---

# What I learned

Notes from building this step by step.

## Tooling

- **Air** is a Go CLI tool, installed with `go install`, not a library. Do not add it to `go.mod`. Make sure `$(go env GOPATH)/bin` is on your PATH.
- Air's `cmd` in `.air.toml` must point at the folder containing `main.go` (`go build -o ./tmp/main ./cmd/api`). Pointing at `.` gives "no Go files".
- Add `"env"` to `include_ext` so Air restarts when `.env` changes.
- "address already in use" usually means an old process is still running. Find it with `lsof -i :8080` and stop it with `kill <PID>`. Graceful shutdown prevents most of these.

## Config and environment variables

- Env var names are **case-sensitive**: `os.Getenv("Port")` never matches `PORT`.
- `godotenv.Load()` reads `.env` from the current working directory. In production the variables come from the server, so a missing file is only a warning.
- Never commit `.env`. It holds secrets.

## Logging with slog

- `log/slog` is in the standard library (Go 1.21+). No package needed.
- One `*slog.Logger` is enough. Levels (`Debug`, `Info`, `Warn`, `Error`) are built in, so the old two-logger pattern is not necessary.
- Use a message plus key/value pairs: `logger.Error("server failed", "err", err)`.
- Text output is easier to read in development, JSON is better in production. The level should come from config, not be hardcoded.
- `slog.Error` does not exit. Use `os.Exit(1)` after it for fatal startup errors. `os.Exit` skips deferred calls.

## HTTP server

- Always set timeouts on `http.Server` (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
- Use your own `http.NewServeMux()` instead of the global default mux.
- `HandleFunc("/", ...)` is a catch-all. Use `"GET /{$}"` to match only the exact root.
- A method in the pattern (`"GET /health"`) makes other methods return `405 Method Not Allowed`.
- A handler has the signature `func(w http.ResponseWriter, r *http.Request)`. With `net/http`, set the header, then use `json.NewEncoder(w).Encode(...)`.
- Headers and status must be set **before** writing the body. After the first write they are locked, so an encode error can only be logged.
- `ErrServerClosed` is expected after `server.Shutdown()`, so it is not treated as a failure.

## Structuring the app

- `cmd/api` holds the entrypoint, and `main.go` only wires things together.
- `internal/` cannot be imported by other modules, which keeps the code private to the project.
- An `Application` struct holds shared dependencies (logger, services). Handlers are methods on it, so they reach dependencies through `app.` instead of globals.
- Helpers keep error handling consistent:
  - `serverError`: logs the real error with method and URL, sends a generic 500.
  - `clientError`: sends 400, 404 and similar, no logging needed.
  - `notFound`: shortcut for a 404.
- Do not split into packages until it hurts. Split when a file gets crowded.

## Dependency injection

- Whoever needs to own a lifecycle creates the thing. `main` opens the DB, so it can fail fast and close it on shutdown.
- Pass dependencies through constructors and struct fields. No framework is needed.
- Avoid global `DB` variables: they hide dependencies and make tests harder.

## Interfaces

- Go interfaces are satisfied implicitly. There is no "implements" keyword.
- **The package that uses an interface defines it**, not the package that implements it. Here `services` defines `MovieStore`, and `repository.MoviesRepository` satisfies it without importing anything.
- Do not use a pointer to an interface (`*MyInterface`). An interface value already holds a pointer inside.
- No `Interface` suffix in names. Only add an interface when you need one, for example to write a fake in tests.

## Database

- `sql.Open` does not connect. Call `Ping()` to check right away.
- `*sql.DB` is a connection pool, safe for concurrent use. Create it once and share it.
- Driver names must match the driver: `"sqlite3"` for `mattn/go-sqlite3`, `"sqlite"` for `modernc.org/sqlite`. Pure Go (`modernc`) avoids CGO, which makes builds and Docker easier.
- Use `Query` for many rows and `QueryRow` for exactly one.
- Scan takes one pointer per column, in order. Prefer explicit columns over `SELECT *`.
- `defer rows.Close()` right after the error check, then check `rows.Err()` after the loop.
- Start with `movies := []model.Movie{}` so an empty result encodes as `[]`, not `null`.
- `NOT NULL` columns avoid NULL scan errors on plain `int` fields.
- Close the connection on failure paths (`conn.Close()` after a failed ping).

## Errors

- Wrap with context using `%w`: `fmt.Errorf("open database: %w", err)`.
- Wrap in one place only, otherwise messages are duplicated.
- Log an error once, at the top (in `main` or the handler), not at every layer.
- Never send raw errors to clients.

## Migrations

- `CREATE TABLE IF NOT EXISTS` is fine for one table while learning. It does not change an existing table.
- When the schema starts changing, use numbered SQL files (`001_create_movies.sql`) with a tool such as goose. Never edit a migration that has already run.
- Do not `Exec` a goose file directly, because its `Down` section would run too.

## Context

- `context.Context` carries cancellation and deadlines through a call chain.
- It starts at `r.Context()` in the handler and is passed down to `QueryContext`, where the database call stops if the client disconnects.
- It is always the first parameter, named `ctx`. Do not store it in structs.
- Use `context.WithTimeout` for limits, and always `defer cancel()`.
- Practical value with a local SQLite file is small. The habit matters once queries are slow or remote.

## Testing with curl

```
curl -i localhost:8080/health
curl -i -X POST localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"Alien","year":1979}'
```

`-i` shows status and headers, `-v` shows full detail.
