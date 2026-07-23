# Changelog

Notable changes are recorded here. Breaking pre-1.0 changes are called out
explicitly.

## 0.19.0 - 2026-07-23

### Added

- Added review status to API entries and generated metadata.

### Changed

- Exported confidence through the shared schema projection.

## 0.18.0 - 2026-07-23

### Added

- Added 489 reviewed player, vehicle, network, database, and NPC entries.
- Completed coverage of the public declarations in the pinned omp-stdlib source.

## 0.17.0 - 2026-07-23

### Added

- Added 116 reviewed global and per-player text-draw entries.
- Recorded query APIs, typed identifiers, limits, and compatibility aliases.

## 0.16.0 - 2026-07-23

### Added

- Added 114 reviewed core, timer, console, and variable entries.
- Recorded callable defaults, variadic arguments, and compatibility aliases.

### Fixed

- Preserved defaults and variadic tags when importing include declarations.

## 0.15.0 - 2026-07-23

### Added

- Added 40 reviewed global and per-player 3D text-label entries.
- Recorded legacy aliases, defaults, tags, and invalid identifiers.

## 0.14.0 - 2026-07-23

### Added

- Added 14 reviewed HTTP native, tag, and constant entries.
- Recorded the response callback shape and HTTP error codes.

## 0.13.0 - 2026-07-23

### Added

- Added 43 reviewed global and per-player gang-zone entries.
- Recorded visibility checks, callbacks, and colour-name aliases.

## 0.12.0 - 2026-07-23

### Added

- Added 35 reviewed global and per-player pickup entries.
- Recorded pickup visibility, streaming callbacks, and update defaults.

## 0.11.0 - 2026-07-23

### Added

- Added eight reviewed player-class functions.
- Recorded open.mp class queries and weapon defaults.

## 0.10.0 - 2026-07-23

### Added

- Added 28 reviewed object query, attachment, and custom-model entries.
- Recorded SA-MP 0.3.DL functions outside the `samp-037` profile.

## 0.9.0 - 2026-07-23

### Added

- Added 37 reviewed per-player object natives and callbacks.
- Recorded open.mp query functions, defaults, and compatibility names.

### Fixed

- Marked the legacy object editing and camera-collision names as deprecated.

## 0.8.0 - 2026-07-23

### Added

- Added 45 reviewed object material and editing entries.
- Recorded typed constants, open.mp editing names, and material text defaults.

## 0.7.0 - 2026-07-23

### Added

- Added 20 reviewed global-object natives and constants.
- Recorded legacy availability, open.mp aliases, and optional movement values.

## 0.6.0 - 2026-07-23

### Added

- Added 23 reviewed menu natives, callbacks, constants, and tags.
- Recorded open.mp formatting arguments and the SA-MP menu sentinel value.

## 0.5.0 - 2026-07-23

### Added

- Added 15 reviewed dialog natives, callbacks, constants, and tags.
- Recorded the open.mp formatting extension and deprecated dialog alias.

## 0.4.0 - 2026-07-23

### Added

- Added 15 reviewed checkpoint natives, callbacks, and tags.
- Recorded SA-MP and open.mp availability for the checkpoint API.

## 0.3.0 - 2026-07-23

### Added

- Added 25 reviewed actor natives, callbacks, and constants.
- Recorded SA-MP and open.mp availability for the actor API.

## 0.2.0 - 2026-07-23

### Added

- Added `pawnapi coverage` for comparing reviewed data with include declarations.

### Documented

- Recorded current omp-stdlib coverage and the boundary for third-party APIs.

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
