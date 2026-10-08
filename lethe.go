// Package lethe is the Go face of the cross-language lineage contract:
// the JSON Schema and the canonical fixture graph, embedded so any Go
// consumer (the mnemonica ports' test suites, CI validators, code
// generators) reads exactly the bytes that ship with this module — no
// copied files that drift apart. The contract itself (what the format
// means, when it may change) lives in the repository's README and
// AGENTS.md; this package only distributes it.
package lethe

import (
	"embed"
	"encoding/json"
)

// LineageSchema is the lineage graph JSON Schema (draft 2020-12). The
// canonical copy lives at the repository root; embedding pins the module
// to it.
//
//go:embed lineage.schema.json
var LineageSchema []byte

// LineageFixture is the canonical fixture graph: the exact bytes every
// port must reproduce from the recipe in LineageFixtureREADME, ids mapped
// 1:1 in first-encounter order.
//
//go:embed testdata/lineage/fixture.json
var LineageFixture []byte

// LineageFixtureREADME is the fixture recipe: the language-neutral
// construction script and the id-mapping rule that make byte-for-byte
// reproduction possible across ports.
//
//go:embed testdata/lineage/README.md
var LineageFixtureREADME string

// TacticaSchemas are the six .tactica type-graph JSON Schemas (draft
// 2020-12): hierarchy, definitions, collections (the required structure
// files) plus usages, flow, eds (the analysis files). Same distribution
// idea as LineageSchema — one embedded source of truth for every
// generator's test suite.
//
//go:embed tactica/*.schema.json
var TacticaSchemas embed.FS

// TacticaSchemaNames lists the .tactica schema bases, in contract order.
var TacticaSchemaNames = []string{
	"hierarchy", "definitions", "collections",
	"usages", "flow", "eds",
}

// TacticaSchema returns one .tactica schema's bytes by base name
// ("hierarchy", "definitions", "collections", "usages", "flow", "eds").
func TacticaSchema(name string) ([]byte, error) {
	return TacticaSchemas.ReadFile("tactica/" + name + ".schema.json")
}

// TacticaSchemaJSON returns one .tactica schema parsed as a generic JSON
// object — a stdlib-only smoke check for consumers without a validator.
func TacticaSchemaJSON(name string) (map[string]any, error) {
	raw, err := TacticaSchema(name)
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}
