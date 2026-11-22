# promy-go-kit — Yokai-aware shared Go code for the Promy fleet

[![CI](https://github.com/tclavelloux/promy-go-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/tclavelloux/promy-go-kit/actions/workflows/ci.yml)

<!-- TOC -->
* [Overview](#overview)
* [Architecture](#architecture)
* [Boundaries](#boundaries)
* [Resources](#resources)
<!-- TOC -->

## Overview

Go library of shared code that depends on the [Yokai](https://github.com/ankorstore/yokai) framework, imported by the `promy-*` services. Framework-free shared code stays in [promy-event-bus](https://github.com/tclavelloux/promy-event-bus).

**Module:** `github.com/tclavelloux/promy-go-kit`

```bash
go get github.com/tclavelloux/promy-go-kit
```

### Implemented

- `fxotlphttp`: exports Yokai traces over OTLP/HTTP when `modules.trace.processor.type: otlp-http`. Endpoint, headers and protocol come only from the standard `OTEL_EXPORTER_OTLP_*` variables. Decision: [promy-steering#16](https://github.com/tclavelloux/promy-steering/issues/16).

### Planned

- Nothing yet. Add a package only when two or more services need the same Yokai-aware code.

## Architecture

### Project Structure

```
fxotlphttp/
├── module.go        # Module, DecorateTracerProviderFactory, Factory
├── module_test.go   # fx app tests against an httptest OTLP receiver
└── testdata/        # Yokai config for the tests
```

### fxotlphttp

- `Module` is a root-level `fx.Options` holding one `fx.Decorate` on `trace.TracerProviderFactory`. Not an `fx.Module`: fx scopes a decorator to its declaring module, and Yokai's trace module consumes the factory.
- Yokai maps the unknown `otlp-http` type to its noop processor; the decorated factory appends an `otlptracehttp` batch span processor next to it.
- `otlp-http` with no `OTEL_EXPORTER_OTLP_ENDPOINT` or `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`: logs a warning, exports nothing, never fails boot.
- Any other processor type: factory returned unchanged.
- Shutdown: Yokai's trace module flushes and shuts down every registered processor.

## Boundaries

### Public API

| Symbol | Purpose |
|---|---|
| `fxotlphttp.Module` | fx option that enables OTLP/HTTP export |
| `fxotlphttp.ProcessorType` | `otlp-http`, the config value that enables it |
| `fxotlphttp.Factory` | decorated `trace.TracerProviderFactory` |

### Downstream (repos that import this module)

- [promy-template-go](https://github.com/tclavelloux/promy-template-go), [promy-product](https://github.com/tclavelloux/promy-product), [promy-user](https://github.com/tclavelloux/promy-user), [promy-crm](https://github.com/tclavelloux/promy-crm), [promy-identifier](https://github.com/tclavelloux/promy-identifier): rollout tracked in [promy-steering#16](https://github.com/tclavelloux/promy-steering/issues/16).

### Not owned by this library

- Span, propagation and logging conventions: `~/.claude/rules/yokai-go-observability.md`.
- Framework-free shared code: [promy-event-bus](https://github.com/tclavelloux/promy-event-bus).

## Resources

- [HOWTO.md](HOWTO.md): wiring, configuration, development
- [CHANGELOG.md](CHANGELOG.md)
- [Fleet architecture](https://github.com/tclavelloux/promy-steering/blob/main/ARCHITECTURE.md)
