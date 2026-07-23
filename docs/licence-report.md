# Licence report

The repository contains original code, public API facts, and short text copied
from upstream documentation. They do not all share the same licence.

## Repository code

The Go source and the structure of the generated JSON are original PawnKit
work released under the [MIT licence](../LICENSE).

## omp-stdlib data

The seed dataset was checked against
`openmultiplayer/omp-stdlib` at commit
`689c824f6cc558b1ca1b36cfd1aae7f1cda16d65`. That repository uses
MPL-2.0. Each derived entry records the commit, source file, and licence.

The dataset contains two kinds of upstream material:

- Names, signatures, tags, and constant values are facts about the public API.
- `documentationSummary` values come from omp-stdlib's pawndoc comments and
  remain under MPL-2.0.

Consumers that redistribute the generated dataset should keep this report or
an equivalent notice. The surrounding Go code and JSON structure remain MIT.

The `Apache-2.0` value previously used for omp-stdlib in pawnkit-spec examples
was a placeholder and is incorrect. pawn-api records `MPL-2.0`.

No complete omp-stdlib source files are stored here. Contributors should copy
only the facts needed for an entry and follow the process in
[CONTRIBUTING.md](../CONTRIBUTING.md).

## Legacy SA-MP availability

Actor, checkpoint, dialog, and menu availability was checked against `pawn-lang/samp-stdlib` at commit
`8ffb055624308b25521665b60e78b5e6e6b3717f`. Other legacy entries use medium
confidence unless their primary source is recorded.

## Vendored schema

`internal/schema/pawn-api.schema.json` is copied from the PawnKit-owned
`pawnkit-spec` repository. A test compares it with a sibling checkout when one
is available.

## Summary

| Content | Licence or status |
|---|---|
| Go source and JSON structure | MIT |
| API names, signatures, and values | Public API facts sourced from omp-stdlib |
| `documentationSummary` text | MPL-2.0 |
| Vendored PawnKit schema | PawnKit-owned |
