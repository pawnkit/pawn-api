package pawnapi

// SchemaVersion is the schema version produced by this package.
const SchemaVersion = 1

// SchemaDocument is the cross-tool interchange format.
type SchemaDocument struct {
	SchemaVersion int           `json:"schemaVersion"`
	Entries       []SchemaEntry `json:"entries"`
}

type SchemaParameter struct {
	Name    string   `json:"name"`
	Tag     string   `json:"tag,omitempty"`
	Default *Literal `json:"default,omitempty"`
}

type SchemaSignature struct {
	Parameters []SchemaParameter `json:"parameters,omitempty"`
	ReturnTag  string            `json:"returnTag,omitempty"`
}

type SchemaAvailability struct {
	Profile string  `json:"profile"`
	Since   string  `json:"since,omitempty"`
	Until   *string `json:"until,omitempty"`
}

type SchemaDeprecated struct {
	Since       string `json:"since"`
	Replacement string `json:"replacement,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

type SchemaSource struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Commit     string `json:"commit"`
	License    string `json:"license"`
}

type SchemaEntry struct {
	ID               string               `json:"id"`
	Kind             string               `json:"kind"`
	Name             string               `json:"name"`
	Signature        *SchemaSignature     `json:"signature,omitempty"`
	Value            *Literal             `json:"value,omitempty"`
	Tags             []string             `json:"tags,omitempty"`
	Availability     []SchemaAvailability `json:"availability"`
	Deprecated       *SchemaDeprecated    `json:"deprecated,omitempty"`
	Source           SchemaSource         `json:"source"`
	DocumentationURL string               `json:"documentationUrl,omitempty"`
}

// ToSchema returns the schema-compatible subset of e.
func (e Entry) ToSchema() SchemaEntry {
	out := SchemaEntry{
		ID:               e.ID,
		Kind:             string(e.Kind),
		Name:             e.Name,
		Value:            e.Value,
		Tags:             e.Tags,
		DocumentationURL: e.DocumentationURL,
		Source: SchemaSource{
			Repository: e.Source.Repository,
			Path:       e.Source.Path,
			Commit:     e.Source.Commit,
			License:    e.Source.License,
		},
	}

	if e.Signature != nil {
		sig := &SchemaSignature{ReturnTag: e.Signature.ReturnTag}
		for _, p := range e.Signature.Parameters {
			sig.Parameters = append(sig.Parameters, SchemaParameter{
				Name:    p.Name,
				Tag:     p.Tag,
				Default: p.Default,
			})
		}
		out.Signature = sig
	}

	for _, a := range e.Availability {
		out.Availability = append(out.Availability, SchemaAvailability(a))
	}

	if e.Deprecated != nil {
		out.Deprecated = &SchemaDeprecated{
			Since:       e.Deprecated.Since,
			Replacement: e.Deprecated.Replacement,
			Reason:      e.Deprecated.Reason,
		}
	}

	return out
}

// ToSchemaDocument returns entries sorted by kind and ID.
func ToSchemaDocument(entries []Entry) SchemaDocument {
	sorted := sortedByKindAndID(entries)
	doc := SchemaDocument{SchemaVersion: SchemaVersion}
	for _, e := range sorted {
		doc.Entries = append(doc.Entries, e.ToSchema())
	}
	return doc
}
