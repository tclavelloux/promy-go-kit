# CLAUDE.md

## What this repo is

- Go **library** of Yokai-aware shared code, imported by `promy-*` services. No binary, no HTTP server, no database.
- Module: `github.com/tclavelloux/promy-go-kit`.
- Scope: code that depends on Yokai or fx and is needed by two or more services. Framework-free shared code goes to `promy-event-bus`.

## Commands

```bash
make test             # race suite
make check-coverage   # coverage gate (pre-push hook, same as CI)
make lint             # golangci-lint
make vuln             # govulncheck
go test -run TestModule_ExportsSpansOverHTTPWithHeaders -v ./fxotlphttp/...
```

## Rules

- Load the `yokai-go` skill before writing Go; follow `~/.claude/rules/yokai-go-*.md`.
- Shared modules decorating Yokai dependencies are root-level `fx.Options`, never `fx.Module`: fx scopes decorators to their declaring module.
- Never fail service boot from telemetry code: log a warning and degrade.
- `fxotlphttp` reads only the standard `OTEL_EXPORTER_OTLP_*` variables. No sink-specific code or endpoint (decision: promy-steering#16).
- Tests start a real fx app with `fxconfig`, `fxlog` and `fxtrace`; config lives in `<package>/testdata/`.
- `.github/workflows/` and `.github/dependabot.yml` are governed by `promy-github-workflows`: never edit them here.
