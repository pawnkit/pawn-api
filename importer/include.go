// Package importer extracts API declarations from Pawn includes.
package importer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pawnkit/pawn-api/pawnapi"
	parser "github.com/pawnkit/pawn-parser"
	"github.com/pawnkit/pawn-parser/token"
)

type Options struct {
	Profile    string
	Repository string
	Path       string
	Commit     string
	License    string
}

// Include extracts native and forward declarations.
func Include(text []byte, opts Options) ([]pawnapi.Entry, error) {
	file := parser.ParseCompact(text, parser.ParseOptions{})
	if file.HasParseErrors() {
		return nil, fmt.Errorf("importer: include contains %d parse diagnostics", len(file.Diagnostics))
	}
	profile := opts.Profile
	if profile == "" {
		profile = pawnapi.ProfileOpenMP
	}
	commit := opts.Commit
	if commit == "" {
		commit = pawnapi.HandAuthoredCommit
	}
	source := pawnapi.Source{Repository: opts.Repository, Path: opts.Path, Commit: commit, License: opts.License}
	var entries []pawnapi.Entry
	declarations := file.Syntax().Declarations()
	for declarations.Next() {
		declaration := declarations.Declaration()
		if enum, ok := parser.AsEnum(declaration); ok {
			if name, ok := enum.Name(); ok {
				entries = append(entries, entry(pawnapi.KindTag, name.Text(), nil, profile, source, opts.Path))
			}
			continue
		}
		function, ok := parser.AsFunction(declaration)
		if !ok {
			continue
		}
		storage, ok := function.Field("storage")
		if !ok || storage.Token().Kind() != token.KwNative && storage.Token().Kind() != token.KwForward {
			continue
		}
		nameNode, ok := function.Name()
		if !ok {
			continue
		}
		name := nameNode.Text()
		kind := pawnapi.KindNative
		if storage.Token().Kind() == token.KwForward {
			kind = pawnapi.KindFunction
			if strings.HasPrefix(name, "On") {
				kind = pawnapi.KindCallback
			}
		}
		signature := &pawnapi.Signature{}
		if tag, ok := function.Field("tag"); ok {
			signature.ReturnTag = strings.TrimSuffix(tag.Text(), ":")
		}
		parameters := function.Parameters()
		for parameters.Next() {
			parameter := parameters.Parameter()
			parameterText := parameter.Text()
			parameterName, ok := parameter.Name()
			if !ok {
				if strings.Contains(parameterText, "...") {
					item := pawnapi.Parameter{Name: "arguments", Variadic: true}
					if tag, hasTag := parameter.Field("tag"); hasTag {
						item.Tag = strings.TrimSuffix(tag.Text(), ":")
					}
					signature.Parameters = append(signature.Parameters, item)
				}
				continue
			}
			item := pawnapi.Parameter{Name: strings.TrimPrefix(parameterName.Text(), "...")}
			if tag, ok := parameter.Field("tag"); ok {
				item.Tag = strings.TrimSuffix(tag.Text(), ":")
			}
			item.Const = strings.Contains(parameterText, "const ")
			item.Reference = strings.Contains(parameterText, "&")
			item.Variadic = strings.Contains(parameterText, "...")
			item.ArrayDimensions = dimensions(parameterText)
			if value, ok := parameterDefault(parameterText); ok {
				item.Default = &value
			}
			signature.Parameters = append(signature.Parameters, item)
		}
		entries = append(entries, entry(kind, name, signature, profile, source, opts.Path))
	}
	children := file.Syntax().Children()
	for children.Next() {
		define := children.Node()
		if define.Kind() != parser.KindDirectiveDefine {
			continue
		}
		if _, parameterized := define.Field("parameters"); parameterized {
			continue
		}
		name, nameOK := define.Field("name")
		value, valueOK := define.Field("value")
		literal, literalOK := parseLiteral(value.Text())
		if !nameOK || !valueOK || !literalOK {
			continue
		}
		item := entry(pawnapi.KindDefine, name.Text(), nil, profile, source, opts.Path)
		item.Value = &literal
		entries = append(entries, item)
	}
	if err := pawnapi.ValidateDataset(entries); err != nil {
		return nil, fmt.Errorf("importer: %w", err)
	}
	return entries, nil
}

func parameterDefault(text string) (pawnapi.Literal, bool) {
	_, value, ok := strings.Cut(text, "=")
	if !ok {
		return pawnapi.Literal{}, false
	}
	value = strings.TrimSpace(value)
	if literal, ok := parseLiteral(value); ok {
		return literal, true
	}
	if value == "" {
		return pawnapi.Literal{}, false
	}
	return pawnapi.StringLiteral(value), true
}

func entry(kind pawnapi.Kind, name string, signature *pawnapi.Signature, profile string, source pawnapi.Source, include string) pawnapi.Entry {
	return pawnapi.Entry{
		ID: string(kind) + ":" + name, Kind: kind, Name: name, Signature: signature,
		Availability: []pawnapi.Availability{{Profile: profile}}, Source: source,
		OwningInclude: include, Confidence: pawnapi.ConfidenceHigh,
	}
}

func parseLiteral(text string) (pawnapi.Literal, bool) {
	text = strings.TrimSpace(text)
	if value, err := strconv.ParseInt(text, 0, 64); err == nil {
		return pawnapi.NumberLiteral(value), true
	}
	if value, err := strconv.Unquote(text); err == nil {
		return pawnapi.StringLiteral(value), true
	}
	if text == "true" || text == "false" {
		return pawnapi.BoolLiteral(text == "true"), true
	}
	return pawnapi.Literal{}, false
}

func dimensions(text string) []int {
	count := strings.Count(text, "[")
	if count == 0 {
		return nil
	}
	return make([]int, count)
}
