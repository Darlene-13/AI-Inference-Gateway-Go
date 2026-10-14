# AI Inference Gateway (Go)

A lightweight HTTP gateway that sits in front of LLM inference backends (Ray Serve, Ollama, OpenAI-compatible APIs) and exposes them behind a single OpenAI-style endpoint. Clients ask for a **model alias** such as `fast`, and the gateway routes the request to the real provider and model configured for that alias, with fallback targets tried in order.

> **Status:** early development. The configuration layer is done and tested. The HTTP server (`cmd/gateway`) is next.

## Features

### Done
- **YAML configuration** with strict parsing: unknown keys such as typos are rejected at startup.
- **Environment variable expansion**: `${NAME}` and `${NAME:-default}`, so secrets and hosts never live in the file.
- **Sensible defaults** for server timeouts, body size, logging and provider timeouts.
- **Validation that reports every problem at once**, not just the first one.
- Unit tests for parsing, env expansion, defaults and validation.

### Planned
- OpenAI-compatible `/v1/chat/completions` proxy
- Model aliases with ordered fallback targets
- Health and readiness endpoints (`/healthz`, `/readyz`) with graceful shutdown and drain delay
- Structured logging (text or JSON)

## Project layout

```
.
├── internal/
│   └── config/                  # load, expand, default and validate the config
│       ├── config.go            # config types
│       ├── load.go              # Load / Parse, env expansion, defaults
│       ├── validate.go          # Validate and ValidationError
│       ├── config_test.go       # unit tests
│       ├── testdata/            # YAML fixtures used by the tests
│       ├── gateway.example.yaml # every setting with its default
│       └── local.yaml           # local two-machine setup (Ray + Ollama)
├── .github/workflows/           # CI (lint, test, build) and release
├── .golangci.yml                # linter configuration
├── Makefile                     # common dev commands
└── go.mod
```

## Configuration

Every setting is documented in [`internal/config/gateway.example.yaml`](internal/config/gateway.example.yaml). A minimal config looks like this:

```yaml
providers:
  ray:
    base_url: ${RAY_URL:-http://localhost:8000}
    request_timeout: 120s
  ollama:
    base_url: ${OLLAMA_URL:-http://localhost:11434}

models:
  fast:                      # the alias clients send as "model"
    - provider: ray          # tried first
      model: llama
    - provider: ollama       # fallback
      model: "llama3.2:1b"
```

| Section     | Key                   | Default  | Notes                                     |
|-------------|-----------------------|----------|-------------------------------------------|
| `server`    | `addr`                | `:8080`  | listen address                            |
|             | `read_header_timeout` | `5s`     |                                           |
|             | `read_timeout`        | `30s`    |                                           |
|             | `idle_timeout`        | `120s`   |                                           |
|             | `shutdown_timeout`    | `30s`    | time in-flight requests get on shutdown   |
|             | `drain_delay`         | `0s`     | use ~5s behind a load balancer            |
|             | `max_body_bytes`      | `1048576`| 1 MiB                                     |
| `log`       | `level`               | `info`   | `debug`, `info`, `warn`, `error`          |
|             | `format`              | `text`   | `text` or `json`                          |
| `providers` | `base_url`            | required | must be an http(s) URL                    |
|             | `api_key`             |          | read it from the environment              |
|             | `connect_timeout`     | `2s`     |                                           |
|             | `request_timeout`     | `60s`    |                                           |
| `models`    | `provider`, `model`   | required | provider must exist under `providers`     |

**Environment variables:** `${NAME}` is replaced by `NAME`. `${NAME:-default}` uses `default` when `NAME` is unset or empty. Other `$` forms are left untouched.

## Development

Requirements: Go (version in `go.mod`) and [golangci-lint](https://golangci-lint.run/) v2.

```bash
make help        # list the commands
make test        # run all unit tests
make test-race   # tests with the race detector, like CI
make cover       # coverage per function
make lint        # run linters
make fmt         # format code
make tidy        # tidy go.mod / go.sum
```

You can also run the config tests directly:

```bash
go test -v ./internal/config/
```

## CI

GitHub Actions runs lint, tests with the race detector and coverage, and a build on every push to `main` and on every pull request. See [`.github/workflows/ci.yml`](.github/workflows/ci.yml).

## License

[Apache License 2.0](LICENSE)
