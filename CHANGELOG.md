# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Optional `OLLAMA_MODEL` environment variable to choose the Ollama model. It defaults to `llama3`, the previous hard-coded value.
- Integration tests against a real Ollama, behind the `integration` build tag: an Ollama client test and an end-to-end API test in `tests/integration`. Unit tests use the `unit` build tag.
- A unit test proving that a client disconnecting from `POST /llm/stream` cancels the request to Ollama.
- Makefile targets `help`, `test-unit`, `test-integration`, `test` and `test-coverage`, plus `docker-*` targets. The old `local-*` targets still work.
- CI integration job with Ollama as a service container, using the small `qwen2.5:0.5b` model.
- `LICENSE` file (MIT). The README already claimed this license.

### Changed

- **Breaking:** the Ollama health check moved from `GET /health/llm` to `GET /llm/health`, to match the other `/llm/*` routes. The old path now returns `404`.
- `POST /llm/stream` now answers failures that happen before the stream starts with `500` and a JSON body `{"error":"..."}`, the same as `POST /llm/generate`. It used to return plain text.
- `make local-init` (and so `make local-complete-first-start`) no longer overwrites an existing `.env`.
- `make docker-pull-model` no longer uses `docker exec -ti`, so it also works without a terminal.
- Docker Compose: the API waits until Ollama is healthy before starting, and gets `OLLAMA_URL` and `OLLAMA_MODEL` from `.env` or from defaults. `.env` is now optional.
- `make docker-pull-model` and the integration tests use the same `OLLAMA_MODEL` as the API: environment first, then `.env`, then `llama3`.
- README rewritten: API, errors and health tables, a streaming and cancellation section, the real rate limiter and circuit breaker settings, design decisions, and a mermaid architecture diagram.
- The coverage badge is published by CI to a `badges` branch and read through a shields.io endpoint, instead of a percentage committed into the README.

### Removed

- `LOG_LEVEL` and `OTEL_SERVICE_NAME` from `.env.example`. The code never read them.
- The Loki healthcheck in Docker Compose. It always failed, because the image has no `wget`.
- `docs/project.png`, replaced by the mermaid diagram in the README.

### Fixed

- The README now gives the Loki query that works, `{service_name="ai-backend-api"}`, and no longer claims telemetry beyond logs.
