// Package lethe is the Go face of the cross-language lineage contract:
// the JSON Schema and the canonical fixture graph, embedded so any Go
// consumer (the mnemonica ports' test suites, CI validators, code
// generators) reads exactly the bytes that ship with this module — no
// copied files that drift apart. The contract itself (what the format
// means, when it may change) lives in the repository's README and
// AGENTS.md; this package only distributes it.
package lethe

import _ "embed"

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
