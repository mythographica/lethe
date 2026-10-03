package lethe

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestLineageSchemaIsValidJSON: the embedded schema parses, is one JSON
// object, and pins the contract version. Full schema-VALIDATION lives in
// the mnemonica core-go tools module, which owns the validator dependency;
// this module stays stdlib-only.
func TestLineageSchemaIsValidJSON(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(LineageSchema, &schema); err != nil {
		t.Fatalf("lineage.schema.json does not parse: %v", err)
	}
	if schema["version"] != nil {
		t.Error("the schema document itself must not hardcode a graph version")
	}
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("schema has no $defs")
	}
	if _, ok := defs["node"]; !ok {
		t.Error("schema $defs has no node definition")
	}
	if _, ok := defs["value"]; !ok {
		t.Error("schema $defs has no value definition")
	}
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("schema has no properties")
	}
	for _, key := range []string{"version", "heads", "nodes"} {
		if _, ok := properties[key]; !ok {
			t.Errorf("schema properties missing %q", key)
		}
	}
}

// TestLineageFixtureIsValidJSON: the embedded fixture parses as the
// graph shape — version "1", heads array, nodes object.
func TestLineageFixtureIsValidJSON(t *testing.T) {
	var fixture map[string]any
	if err := json.Unmarshal(LineageFixture, &fixture); err != nil {
		t.Fatalf("fixture.json does not parse: %v", err)
	}
	if fixture["version"] != "1" {
		t.Errorf("fixture version = %v, want %q", fixture["version"], "1")
	}
	heads, ok := fixture["heads"].([]any)
	if !ok || len(heads) != 2 {
		t.Errorf("fixture heads = %v, want two", fixture["heads"])
	}
	nodes, ok := fixture["nodes"].(map[string]any)
	if !ok || len(nodes) != 4 {
		t.Errorf("fixture nodes = %v, want four", fixture["nodes"])
	}
}

// TestLineageFixtureREADME accompanies the fixture: non-empty and names
// the id-mapping rule that the byte-for-byte contract depends on.
func TestLineageFixtureREADME(t *testing.T) {
	if len(LineageFixtureREADME) == 0 {
		t.Fatal("the fixture README is empty")
	}
	if !bytes.Contains([]byte(LineageFixtureREADME), []byte("first-encounter order")) {
		t.Error("the fixture README must document the first-encounter id mapping")
	}
}
