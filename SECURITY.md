# Security policy

Before 1.0, security fixes target the latest released minor version.

## Report a vulnerability

Use GitHub's private [security advisory form](https://github.com/pawnkit/pawn-api/security/advisories/new). If it is unavailable, contact a maintainer listed in `CODEOWNERS` or the PawnKit organization profile.

Include the affected version, impact, and a small reproduction when possible. Do not open a public issue before a fix is available.

## Scope

The CLI decodes repository-controlled JSON. The public `LoadEntries` API may also receive data supplied by another program. Panics or excessive resource use caused by malformed input are in scope.

The runtime library does not execute Pawn code, load AMX files, spawn subprocesses, or make network requests. Data research may use the network, but generation and loading do not.

Direct dependencies are `pawnkit-core` and `jsonschema/v5`. CI should continue to scan them for published vulnerabilities.
