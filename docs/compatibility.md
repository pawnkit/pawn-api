# Compatibility

## Schema version

The generated interchange file uses `pawn-api.schema.json` version 1. A vendored copy lives at `internal/schema/pawn-api.schema.json`.

When version 2 is published, this repository will keep the version 1 reader during the migration window.

## Compatibility profiles

Profiles are strings matching `^[a-z][a-z0-9-]*$`, not a closed enum. The seed dataset uses two:

| Profile | Meaning |
|---|---|
| `samp-037` | Legacy SA-MP 0.3.7 client/server API surface. |
| `openmp` | open.mp's API surface (a superset of `samp-037` plus open.mp-only additions). |

If a profile is absent from an entry's availability list, the API is not available for that profile.

## Gap between `pawnapi.Entry` and `pawn-api.schema.json` v1

The Go model contains fields that schema version 1 cannot represent. Schema objects reject extra properties, so those fields cannot be added unofficially to the interchange file.

| ARCHITECTURE.md field | In `pawnapi.Entry`? | In schema v1? |
|---|---|---|
| Stable ID / canonical name | Yes | Yes |
| Aliases / spelling variants | Yes (`Entry.Aliases`) | **No** |
| Signature / parameters | Yes | Yes |
| `const`, array dimensions, varargs | Yes (`Parameter.Const/ArrayDimensions/Variadic`) | **No** (parameter is `{name, tag, default}` only) |
| Return semantics (prose) | Yes (`Signature.ReturnSemantics`) | **No** (`returnTag` only) |
| Introduced/deprecated/removed versions | Yes (`Availability.Since/Until`, `Deprecated.Since`) | Yes |
| Replacement/migration guidance | Yes (`Deprecated.Replacement/Reason`) | Yes |
| Availability by profile | Yes | Yes |
| Owning include/plugin/component | Yes (`Entry.OwningInclude`) | **No** (only `source.path`, the file pawn-api extracted from, which is not always the include a script author writes) |
| Constraints, sentinels, side effects | Yes (`Entry.Constraints`, free text) | **No** |
| Callback context / thread notes | Yes (`Entry.CallbackContext`) | **No** |
| Documentation summary (short prose) | Yes (`Entry.DocumentationSummary`) | **No** (`documentationUrl` only, a link) |
| Source provenance | Yes | Yes |
| Confidence | Yes (`Entry.Confidence`) | Yes |
| Review status | Yes (`Entry.ReviewStatus`) | Yes |

`Entry.ToSchema` creates `data/generated/pawn-api.json`; `pawnapi.Load` exposes
the full model to Go callers. A future schema revision should add aliases and
documentation summaries so non-Go consumers can use them too.

## Version policy

The module is pre-1.0. Breaking Go API or generated-format changes are recorded in [CHANGELOG.md](../CHANGELOG.md). Normal semantic versioning applies after 1.0.
