package schema

import (
	"bytes"
	"os"
	"testing"
)

// TestVendoredSchemaMatchesSource checks an optional sibling spec checkout.
func TestVendoredSchemaMatchesSource(t *testing.T) {
	const sourcePath = "../../../pawnkit-spec/schemas/pawn-api.schema.json"
	upstream, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Skipf("pawnkit-spec not available as a sibling checkout: %v", err)
	}
	if !bytes.Equal(upstream, Raw) {
		t.Fatal("internal/schema/pawn-api.schema.json has drifted from pawnkit-spec/schemas/pawn-api.schema.json; re-vendor it and record the new revision in docs/compatibility.md")
	}
}
