# Changelog

Notable changes are recorded here. Breaking pre-1.0 changes are called out
explicitly.

## 0.1.1 - 2026-07-23

### Fixed

- Replaced the parser pseudo-version with `pawn-parser v1.1.9`.

## [Unreleased]

### Added

- Go package for loading, validating, indexing, and serializing API entries.
- CLI commands for validation, generation, snapshots, and classified diffs.
- Seed dataset: 79 hand-verified entries (53 natives, 18 callbacks, 5
  constants, 3 tags) extracted live from `openmultiplayer/omp-stdlib` at
  commit `689c824f6cc558b1ca1b36cfd1aae7f1cda16d65`, covering player and
  vehicle functions plus core lifecycle callbacks.
- Deterministic source generator with a content-hash manifest.
- Vendored schema and JSON Schema Draft 2020-12 validation.
- Licence report covering omp-stdlib facts and documentation summaries.

[Unreleased]: https://github.com/pawnkit/pawn-api/commits/main
