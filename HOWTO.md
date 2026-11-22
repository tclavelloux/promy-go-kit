# HOWTO — promy-go-kit

## Use fxotlphttp in a service

Add the module to the bootstrapper, next to Yokai's modules:

```go
import "github.com/tclavelloux/promy-go-kit/fxotlphttp"

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	fxhttpserver.FxHttpServerModule,
	fxotlphttp.Module,
	// service modules...
)
```

Select the processor in `configs/config.prod.yaml`:

```yaml
modules:
  trace:
    processor:
      type: otlp-http
```

Set no endpoint in YAML. The exporter reads the standard variables:

| Variable | Prod (Railway) | Dev (docker compose) |
|---|---|---|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | injected by Railway | `http://promy-jaeger:4318` |
| `OTEL_EXPORTER_OTLP_HEADERS` | injected by Railway | unset |

- `service.name` always comes from `app.name`: Yokai's trace resource ignores `OTEL_SERVICE_NAME`.

Enable Railway tracing in the service's `.railway/railway.ts`:

```ts
tracing: { enabled: true },
```

- A variable set on the service wins over Railway's: setting your own `OTEL_EXPORTER_OTLP_ENDPOINT` stops Railway injecting its tracing variables.
- Switch sinks (for example to Grafana Cloud) by setting `OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_EXPORTER_OTLP_HEADERS` and removing the `tracing` block. No code change.
- No endpoint at boot: one warning `otlp-http trace processor without OTEL_EXPORTER_OTLP_ENDPOINT, spans are not exported`, then the service runs without exporting.

## Getting Started

### Prerequisites

- Go (version in `go.mod`)
- `pre-commit` (`brew install pre-commit`)

### Install

```bash
make setup   # installs pre-commit, commit-msg and pre-push hooks
```

### Test

```bash
make test             # race suite
make check-coverage   # race suite gated on .testcoverage.yml (also the pre-push hook)
make coverage         # coverage.html
```

### Lint

```bash
make lint   # golangci-lint, version pinned in .golangci-version
make vuln   # govulncheck, same pin as CI
```

## Contributing

- Branch from `main`; the `no-direct-commit-to-main` hook blocks commits on `main`.
- Conventional Commits, enforced by the `commit-msg` hook and the PR title check.
- Squash-merge; release-please turns each PR title into a changelog entry and cuts the release.
- CI callers and `dependabot.yml` are governed by [promy-github-workflows](https://github.com/tclavelloux/promy-github-workflows): change them there.
