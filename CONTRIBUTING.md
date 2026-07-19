# Contributing to pawn-api

PawnKit is maintained by volunteers, so reviews may take a little time.

Corrections, missing entries, and better source references are welcome. A small,
well-sourced update is more useful than a large import nobody can verify.

One incorrect signature can affect completion, linting, migration, and generated documentation. Check every entry against a primary source; do not fill gaps from memory.

## Add or update an entry

Edit `data/source/*.json`. Never edit `data/generated` or `pawnapi/data/pawn-api.full.json` by hand.

For an omp-stdlib entry:

1. Check the declaration in the upstream `.inc` file.
2. Record the exact commit, source path, and licence.
3. Use low confidence and a clear note if part of the entry cannot be verified. Leave it out if a plausible guess would be misleading.
4. Make deprecation replacements point to an existing entry ID.
5. Read [the licence report](docs/licence-report.md) before copying a documentation summary.

Then regenerate and test:

```sh
go run ./cmd/pawnapi validate
go run ./cmd/pawnapi generate
go test ./...
go vet ./...
```

Commit the source change and generated files together. CI runs `pawnapi generate --check` and rejects stale output.

## What validation proves

Validation checks the schema, duplicate IDs and names, version ranges, and deprecation references. It cannot prove that a signature matches upstream. The reviewer and contributor must check that source.

Facts belong in this repository. Lint severity and fix policy belong in `pawnlint`.

## Tests

Dataset-wide tests cover new entries automatically. Add a fixture under `testdata/broken` when fixing a validation gap. Generator output must remain byte-for-byte deterministic.

The module requires Go 1.26 or later. Before opening a pull request, run the commands above and `go fmt ./...`.

Breaking public API or generated-format changes need a note in [CHANGELOG.md](CHANGELOG.md). Report vulnerabilities through [SECURITY.md](SECURITY.md).
