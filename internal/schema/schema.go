// Package schema validates pawn-api interchange documents.
package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed pawn-api.schema.json
var Raw []byte

const SchemaID = "https://schemas.pawnkit.dev/pawn-api/v1/schema.json"

var (
	compileOnce sync.Once
	compiled    *jsonschema.Schema
	compileErr  error
)

// Compiled returns the cached pawn-api schema.
func Compiled() (*jsonschema.Schema, error) {
	compileOnce.Do(func() {
		c := jsonschema.NewCompiler()
		c.Draft = jsonschema.Draft2020
		if err := c.AddResource(SchemaID, bytes.NewReader(Raw)); err != nil {
			compileErr = fmt.Errorf("schema: registering vendored pawn-api.schema.json: %w", err)
			return
		}
		sch, err := c.Compile(SchemaID)
		if err != nil {
			compileErr = fmt.Errorf("schema: compiling vendored pawn-api.schema.json: %w", err)
			return
		}
		compiled = sch
	})
	return compiled, compileErr
}

// ValidateDocument validates a JSON interchange document.
func ValidateDocument(raw []byte) error {
	sch, err := Compiled()
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var inst any
	if err := dec.Decode(&inst); err != nil {
		return fmt.Errorf("schema: decoding document: %w", err)
	}

	if err := sch.Validate(inst); err != nil {
		return fmt.Errorf("schema: document does not conform to pawn-api.schema.json: %w", err)
	}
	return nil
}
