# Betrayal DB Setup (No Docker)

## Environment Variables (.env)

- `DATABASE_URL` - prod or Railway DB URL (e.g. for local dev, for prod, or Railway proxy; string supplied by you)
- `TEST_DATABASE_URL` - separate Railway DB for tests/CI to isolate and avoid clobbering prod

Example:
```
DATABASE_URL=postgres://user:pass@host:port/db
TEST_DATABASE_URL=postgres://user:pass@host:port/test_db
PORT=8080
```

Ensure `.env` contains both, and `.env.example` is updated as reference/for on-boarding.

## Migrations
To run all migrations (assuming golang-migrate installed):
```
migrate -path internal/db/migrations -database "$DATABASE_URL" up
```

**Alternative (preferred for teams/CI/local):**
Extract the DB URL directly from `.env` using grep/cut. This keeps everything in sync and makes scripts portable:

Normal migrations (dev/prod):
```
migrate -database $(cat .env | grep DATABASE_URL | cut -d '=' -f2) -path internal/db/migrations up
```
Test database (run before `go test`):
```
migrate -database $(cat .env | grep TEST_DATABASE_URL | cut -d '=' -f2) -path internal/db/migrations up
```
- This pattern works no matter the order of lines in `.env`
- Swap the grep key for any DB env var as needed.
- The Makefile automates both patterns for convenience.

Or, for tests (before running Go tests):
```
migrate -path internal/db/migrations -database "$TEST_DATABASE_URL" up
```

(You must run `migrate ... down` to reset between test runs/or before running up.)

## Running DB Tests
- Go unit/integration tests look for `TEST_DATABASE_URL` first (then fall back to `DATABASE_URL`).
- CI should use a test DB endpoint with production schema but non-prod data.
- Running `go test ./internal/db/...` will exercise DB create/get logic and verify migrations.

## Tips
- Railway's web console can reset DB or view logs if needed.
- Use different databases/roles for production and test when possible.
- Never use a production database as your test or development target.
