# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking:** the Ollama health check moved from `GET /health/llm` to `GET /llm/health`, to match the other `/llm/*` routes. The old path now returns `404`.
- `POST /llm/stream` now answers failures that happen before the stream starts with `500` and a JSON body `{"error":"..."}`, the same as `POST /llm/generate`. It used to return plain text.
