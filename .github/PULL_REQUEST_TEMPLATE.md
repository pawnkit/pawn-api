## Summary

<!-- What does this change do, and why? -->

## Scope check

- [ ] This change does not duplicate source ranges, diagnostics, edits,
      project discovery, or semantic logic owned by another PawnKit
      repository (see REPOSITORY-BOUNDARIES.md).
- [ ] This change does not import any higher-level PawnKit tool
      (`pawnlint`, `pawnlsp`, `pawnmigrate`, etc.).
- [ ] This change does not encode tool-specific severity or lint policy
      into the data model.
- [ ] Any new external dependency is justified in this PR's description.

## If this changes `data/source/*.json`

- [ ] Every new/changed entry's signature was verified against a primary
      source (not memory/guesswork); `source.commit`/`source.path` point
      at what was actually checked.
- [ ] `source.license` is correct for the upstream repository (see
      docs/licence-report.md — do not assume MIT/Apache-2.0 by analogy).
- [ ] `confidence` honestly reflects how thoroughly this was checked.
- [ ] `pawnapi generate` was run and its output is included in this PR
      (`pawnapi generate --check` passes).
- [ ] `documentationSummary`, if set, was extracted from the source's own
      documentation, and docs/licence-report.md's treatment of that field
      still applies (or this PR updates it for a new upstream source).

## Testing

- [ ] `gofmt -l .` is empty.
- [ ] `go vet ./...` passes.
- [ ] `go test ./...` passes (`CGO_ENABLED=1 go test -race ./...` if
      concurrent code changed).
- [ ] New/changed public behaviour has table-driven tests.
- [ ] Any fixed bug has a regression test.

## Breaking changes

<!-- None, or describe the break and the migration path (see
     docs/compatibility.md's version policy). -->
