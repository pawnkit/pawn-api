# pawn-api

`pawn-api` is PawnKit's versioned database of SA-MP and open.mp natives, callbacks, constants, and tags. It includes a Go package for querying the data and a CLI for maintaining it.

Tools use this data to answer questions such as whether `SetPlayerPos` exists for a profile, which parameters it accepts, and whether it is deprecated.

## Status

This is a pre-1.0 dataset. It currently has 1,146 reviewed entries covering
the public declarations in the pinned omp-stdlib source.

## Install

```sh
go get github.com/pawnkit/pawn-api
```

Requires Go 1.26 or later.

```sh
go install github.com/pawnkit/pawn-api/cmd/pawnapi@latest
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/pawnkit/pawn-api/pawnapi"
)

func main() {
	ix, err := pawnapi.Load()
	if err != nil {
		panic(err)
	}

	e, ok := ix.ByID("native:SetPlayerPos")
	fmt.Println(ok, e.Name, e.Signature.ReturnTag)

	for _, cb := range ix.ByKind(pawnapi.KindCallback) {
		fmt.Println(cb.Name)
	}

	for _, e := range ix.ByProfile("samp-037") {
		fmt.Println(e.ID, "is available under samp-037")
	}
}
```

## Maintain the dataset

```sh
pawnapi validate
pawnapi generate
pawnapi generate --check
pawnapi coverage pawnapi/data/pawn-api.full.json /path/to/omp-stdlib/omp_*.inc
```

`generate --check` is intended for CI. It fails when the embedded and interchange files do not match the source data.
`coverage` compares a dataset with declarations extracted from one or more
includes. It exits with status 1 when declarations are missing.

## Current coverage

| Kind | Count | Examples |
|---|---:|---|
| `native` | 936 | `SetTimerEx`, `TextDrawCreate`, `DB_Open` |
| `callback` | 91 | `OnPlayerConnect`, `OnVehicleSpawn`, `OnNPCSpawn` |
| `function` | 3 | `SetTeamCount`, `EnableTirePopping` |
| `constant` | 67 | `INVALID_TEXT_DRAW`, `HTTP_GET`, `HTTP_ERROR_BAD_HOST` |
| `define` | 2 | `INVALID_TIMER` |
| `tag` | 47 | `Float`, `PLAYER_STATE`, `VEHICLE_MODEL_INFO`, `NPC_MOVE_MODE` |

Each entry records its upstream repository, file, commit, licence, confidence,
and review status. [docs/licence-report.md](docs/licence-report.md) explains
how that provenance is collected.

## Limits

- The reviewed dataset covers every public omp-stdlib subsystem.
- At the pinned omp-stdlib commit, 1,055 of the 1,061 declarations extracted
  from `omp_*.inc` have matching entries. The other six are typed constants
  represented as `define` declarations by the importer.
- `pawnapi snapshot` accepts full-model JSON or a Pawn include. Include import
  covers natives, forwards/callbacks, named tags, and literal defines.
- SA-MP availability is reviewed against a pinned 0.3.7 include source.
- Schema version 1 does not contain every field used by the Go model. The gap
  is listed in [docs/compatibility.md](docs/compatibility.md).

Third-party library APIs are not bundled here. PawnKit reads them from the
project's installed includes, so completion and diagnostics match the version
the project actually uses.

## Links

- PawnKit organisation: <https://github.com/pawnkit>
- Architecture: [docs/architecture.md](docs/architecture.md)
- Compatibility and schema gap: [docs/compatibility.md](docs/compatibility.md)
- Licence report: [docs/licence-report.md](docs/licence-report.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
- Security policy: [SECURITY.md](SECURITY.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
