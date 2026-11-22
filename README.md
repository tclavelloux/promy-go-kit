# TEMPLATE_PACKAGE_NAME — <one-line role description>

[![CI](https://github.com/TEMPLATE_ORG/TEMPLATE_PACKAGE_NAME/actions/workflows/ci.yml/badge.svg)](https://github.com/TEMPLATE_ORG/TEMPLATE_PACKAGE_NAME/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8)](go.mod)

<!-- TOC -->
* [Overview](#overview)
* [Architecture](#architecture)
* [Boundaries](#boundaries)
* [Decision Log](#decision-log)
* [Resources](#resources)
<!-- TOC -->

Scaffolded from [promy-template-go](https://github.com/tclavelloux/promy-template-go) — see [TEMPLATE_INSTRUCTION.md](TEMPLATE_INSTRUCTION.md) and delete it once setup is done.

## Overview

<What problem this service solves, for whom, and its role in the promy ecosystem. Max 4 sentences.>

### Implemented

- <Feature> — <one-clause description>

### Planned

- <Feature> — <status: roadmap phase / dependency present / not wired>

## Architecture

### Layers

```
HTTP request
    │
    ▼
internal/router.go ─────────► route table, middleware
    │
    ▼
internal/api/handler/ ──────► decode, validate, typed response structs
    │
    ▼
internal/domain/<domain>/ ──► business rules
    │           ├── model/       entities
    │           └── repository/  persistence
    ▼
MySQL
```

Dependencies are wired through Uber FX in [`internal/register.go`](internal/register.go); [`internal/bootstrap.go`](internal/bootstrap.go) assembles the Yokai application.

### Project Structure

| Path | Contents |
|---|---|
| [`cmd/`](cmd/) | Cobra entry points — `root.go`, `run.go`, `migrate.go` |
| [`internal/api/handler/`](internal/api/handler/) | HTTP handlers, one package per domain |
| [`internal/api/request/`](internal/api/request/) | Shared request decoding helpers |
| [`internal/domain/`](internal/domain/) | Business logic, one package per domain |
| [`configs/`](configs/) | Yokai config per environment |
| [`db/migrations/`](db/migrations/) | goose SQL migrations |
| [`scripts/`](scripts/) | Lint wrapper, tidy check, scaffold script |
| [`.pre-commit-config.yaml`](.pre-commit-config.yaml) | Git hooks (pre-commit, commit-msg, pre-push), installed by `make setup` |

### Conventions

- Handlers return typed `{Action}Response` / `{Action}Filters` structs — never `map[string]any`. List endpoints put `Limit` and `Offset` inside `Filters`.
- Pre-launch: one migration file per service. Edit [`db/migrations/0001_create_tables.sql`](db/migrations/0001_create_tables.sql) in place rather than adding numbered files, then reseed with `make full-fresh`.
- `.golangci-version` pins the linter for local hooks and `make lint`. It must match the `golangci-lint-version` default in [`go-lint.yml`](https://github.com/tclavelloux/promy-github-workflows/blob/v1/.github/workflows/go-lint.yml), which CI calls — the workflow warns on every run if they drift.
- Coverage thresholds live only in [`.testcoverage.yml`](.testcoverage.yml), read by both the pre-push `check-coverage` hook and CI. Ratchet up, never down.
- Squash merge is the only permitted merge strategy. Merge commits and rebase merges are disabled at repo level. This is the arbitration chosen to work on par with release-please: GitHub always writes the PR title into a merge commit's message, so the merge commit and the branch commits it carries each parse as conventional commits, and every change reaches the CHANGELOG twice. Squash collapses a PR into exactly one conventional commit — its title.
- One PR per business purpose: a whole feature, fix or chore end to end, however many commits it takes. Never split one purpose across PRs (no PR per commit or per layer), and never bundle unrelated purposes into one.
- PR titles must be conventional. The squash subject is what release-please parses and what lands in the CHANGELOG; local commit messages are no longer user-facing.

See [HOWTO.md](HOWTO.md) for conventions, setup, and day-to-day commands.

## Boundaries

### Upstream (services this repo calls)

| Service | Role | Auth | Entry point |
|---|---|---|---|
| <service-name> | <what it provides> | <auth method> | [`path/to/client.go`](path) |

### Downstream (services that call this repo)

| Service | What it calls | Auth |
|---|---|---|
| <service-name> | <endpoints> | <auth method> |

### API Surface

OpenAPI spec not yet available — see [`internal/router.go`](internal/router.go) for the endpoint list.

## Decision Log

| Date | Decision | Rationale |
|---|---|---|
| YYYY-MM | <decision> | <why> |
| 2026-09 | Squash merge only; merge commits and rebase merges disabled | GitHub offers no merge-commit message format that omits the PR title, so release-please parses a merge commit and its branch commits as separate conventional commits and duplicates every CHANGELOG entry. Squash yields one conventional commit per PR. The unit of change is the PR: one business purpose per PR. |

## Resources

- [Yokai framework](https://github.com/ankorstore/yokai) — application framework and FX wiring
- [promy-github-workflows](https://github.com/tclavelloux/promy-github-workflows) — shared reusable CI workflows

<!-- readme-updated-at: efa05c6 -->
