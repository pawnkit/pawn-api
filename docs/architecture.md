# Architecture

`pawn-api` owns machine-readable facts about SA-MP and open.mp APIs. Parsing,
project resolution, and lint policy belong elsewhere.

## Data files

Contributors edit `data/source/*.json`. The generator produces two views of that source:

```text
data/source/*.json
        |
        v
 pawnapi generate
        |
        +-- pawnapi/data/pawn-api.full.json
        +-- data/generated/pawn-api.json
        +-- data/generated/manifest.json
```

The embedded full file contains every field used by the Go API. The interchange file contains the subset allowed by `pawnkit-spec` schema version 1. The manifest records counts and hashes so CI can detect stale output.

Do not edit generated files by hand.

## Why there are two views

`pawnapi.Entry` includes aliases, confidence, constraints, callback context, and documentation summaries. Schema version 1 cannot represent all of those fields and rejects unknown properties.

`Entry.ToSchema` creates a strict schema-compatible projection for non-Go consumers. `pawnapi.Load` gives Go callers the full model. Both files come from the same source entries, so their shared fields cannot drift independently.

The exact field gap is listed in [compatibility.md](compatibility.md).

## Packages

| Path | Responsibility |
|---|---|
| `pawnapi` | Public entry types, indexes, loading, and validation |
| `internal/schema` | Vendored schema and schema validation |
| `internal/generator` | Source loading and deterministic output |
| `cmd/pawnapi` | Validation, generation, snapshots, and diffs |
| `data/source` | Reviewed source entries |
| `data/generated` | Checked-in interchange output |

## Validation errors

Dataset validation identifies entries by stable ID rather than source ranges.
It returns `pawnapi.ValidationErrors` because there is no honest diagnostic
span to report.

The generated manifest uses `pawnkit-core/hash`, which is the shared cache-key and content-hash contract.

## Extending the model

- Add a new entity kind to `pawnkit-spec` before adding it to `pawnapi.Kind`.
- New profile names need data changes but no Go enum change.
- `pawnapi snapshot` accepts full-model JSON or extracts declarations from a Pawn include.
- `pawnapi diff` reports additions, removals, signature changes, and metadata changes.
