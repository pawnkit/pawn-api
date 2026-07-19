# Performance

`pawn-api` loads a small embedded dataset into memory. The useful measurements are startup cost and indexed lookup time.

Run the benchmarks with allocation reporting:

```sh
go test ./pawnapi/... -run=^$ -bench=. -benchmem
```

The suite measures full loading, index construction, dataset validation, ID lookup, and deterministic serialization. `Index.ByID` should remain a map lookup; callers such as the language server should load one index and reuse it.

There is no hard CI threshold yet. Compare repeated results when changing loading, indexing, or validation. Add a larger fixture when the dataset grows enough that the current 79-entry seed no longer represents normal use.
