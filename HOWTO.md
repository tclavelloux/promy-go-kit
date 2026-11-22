# HOWTO — TEMPLATE_PACKAGE_NAME

## Code Conventions

- **Domain module pattern** — one package per domain under `internal/domain/<domain>/`, split into `model/` (entities) and `repository/` (persistence). Handlers live separately under `internal/api/handler/<domain>/`.
- **Typed responses** — every handler defines `{Action}Response` and `{Action}Filters` structs. `Filters` mirrors the accepted query parameters; pagination endpoints carry `Limit` and `Offset` inside it. Never return `map[string]any`.
- **FX registration** — new providers are registered in [`internal/register.go`](internal/register.go). A constructor that isn't registered there is never wired.
- **Tests** — use Yokai's `internal.RunTest()` to boot a test application rather than constructing dependencies by hand. Handler and repository tests sit beside the code they cover.
- **Migrations** — pre-launch, edit [`db/migrations/0001_create_tables.sql`](db/migrations/0001_create_tables.sql) in place. goose will not re-apply an edited migration to an existing volume, so run `make full-fresh` afterwards to reseed.
- **Commits and PRs** — [Conventional Commits](https://www.conventionalcommits.org/); commits are signed. The PR title is the commit that lands on `main` (squash merge only), so it must be conventional. One PR per business purpose: a whole feature, fix or chore end to end, however many commits it takes.

## Getting Started

### Prerequisites

- Go — version in [`go.mod`](go.mod)
- Docker & Docker Compose
- Make
- [pre-commit](https://pre-commit.com/) — `brew install pre-commit`; install it before `make setup`

### Install

```bash
make setup
```

Runs `pre-commit install --install-hooks`. Hook types come from `default_install_hook_types` in [`.pre-commit-config.yaml`](.pre-commit-config.yaml): `pre-commit`, `commit-msg`, `pre-push`. `core.hooksPath` stays unset; the hooks live in `.git/hooks/`.

- `make setup` fails if `pre-commit` is not on PATH.
- A pre-existing `.git/hooks/pre-commit` is moved to `pre-commit.legacy` and still chains. Delete it after setup.

### Run

```bash
make up          # start the stack
make logs        # follow container logs
make down        # stop
make fresh       # rebuild without cache
make full-fresh  # rebuild and drop volumes (reseeds the database)
```

`make up` copies `.env.example` to `.env` on first run.

Migrations:

```bash
make migrate            # apply
make migrate-status     # inspect
make migrate-down       # roll back one
make migrate-reset      # roll back all
```

### Test

```bash
make test
```

Runs `go test` with `-race -cover -count=1 -failfast`. Coverage is gated separately by `make check-coverage` (the pre-push hook) and in CI against the thresholds in [`.testcoverage.yml`](.testcoverage.yml) — the single source of truth for both.

### Lint

```bash
make lint          # golangci-lint at the pinned version
make lint-install  # install that version without linting
```

Both route through [`scripts/golangci-lint.sh`](scripts/golangci-lint.sh), which reads `.golangci-version`, installs that exact version on demand and execs it — so the binary on your PATH is never used implicitly and local findings match CI. Enabled linters are in [`.golangci.yml`](.golangci.yml).

### Build / Deploy

```bash
go build ./...
```

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) runs on pull requests to `main` only (no `push` trigger) and skips draft PRs. Jobs:
- `lint`, `vuln`, `docker`, `coverage` — call the shared [`go-lint.yml`, `go-vuln.yml`, `go-docker.yml` and `go-coverage.yml`](https://github.com/tclavelloux/promy-github-workflows) workflows, so tooling versions are decided fleet-wide rather than here. `docker` pins the Dockerfile's Go version against `go.mod`.
- `test` — race suite; uploads the coverage profile that `coverage` consumes.
- `gitleaks` — secret scan.

[`.github/workflows/pr-title.yml`](.github/workflows/pr-title.yml) validates the PR title (Conventional Commits) on every PR, drafts included. The repo allows squash merge only, so the PR title becomes the commit on `main`.

Releases are cut by [Release Please](https://github.com/googleapis/release-please) from the PR titles that squash-merge leaves on `main`; see [`release-please-config.json`](release-please-config.json) and [`.github/workflows/release-please.yml`](.github/workflows/release-please.yml).

**Railway** (configured in [`.railway/railway.ts`](.railway/railway.ts), named partial in project `main`):
- New service: rename the partial and the service name, then declare its `source` and live variables (`NAME: preserve()`), plus `sleepApplication` if App Sleeping is on; an omitted field plans as unset
- Pre-deploy: `/app migrate-database up` runs the migrations before each deploy (timeout 300s); the app does not migrate at boot
  - Before migrating, `migrate-database` waits for the database: connection errors (refused, timeout, DNS, handshake) are retried with backoff 1 s → 8 s for up to 90 s; auth, unknown-database and SQL errors fail at once. `create`, `fix` and `validate` do not wait.
- Health check: `GET /healthz` (timeout: 150s)
- Restart policy: on failure (max 10 retries)
- Apply a change: Railway CLI >= 5.42.1 and Node >= 22, `cd .railway && npm ci`, `railway link -p main -e production`, `railway config plan`, `railway config apply`

## Contributing

Branch from `main`; the shared `no-direct-commit-to-main` hook (`tclavelloux/promy-github-workflows`, pinned by `rev:`) blocks commits on it.

| Hook | Enforces |
|---|---|
| `pre-commit` | `no-direct-commit-to-main`, then formatting, lint, gitleaks, `go mod tidy` |
| `commit-msg` | Conventional Commits |
| `pre-push` | `check-coverage` only: `make check-coverage` runs the race suite and gates it on `.testcoverage.yml`, matching CI |

```bash
git checkout -b feat/my-feature
# stage by explicit path, never `git add .`
git add <paths>
git commit -S -m "feat(domain): add the thing"
```

Check everything before opening a PR:

```bash
pre-commit run --all-files
pre-commit run --hook-stage pre-push --all-files
```
