# @mnemonica/lethe

Lethe — the river of forgetfulness. This package is the part of mnemonica
that survives it.

## What lethe is

When live instances are gone — the process exited, the language changed,
the heap was collected — what remains of a construction is its
**lineage**: which types were declared, who built what from whom, which
data belonged to which level. Lethe is the contract for exporting that
graph across the boundaries where nothing live can cross:

- **process boundaries** — a crash handler ships the lineage of the
  failing instances out with the report;
- **language boundaries** — the Go, JavaScript, and Python ports speak one
  format, so a graph exported in one runtime can be checked against
  another;
- **time boundaries** — a fixture recorded today must still validate in
  ten years.

Three files make the contract:

- `lineage.schema.json` — the JSON Schema (draft 2020-12) every export
  must satisfy. `version` is `"1"`; any incompatible change means a new
  version and a migration, never a silent edit.
- `testdata/lineage/fixture.json` — the canonical graph. Every port
  reproduces it **byte-for-byte** from the construction script in the
  README beside it, mapping instance ids 1:1 in first-encounter order —
  ids are implementation-specific (a per-realm random prefix plus a
  counter) and are never compared across processes or languages.
- `testdata/lineage/README.md` — the language-neutral recipe (the types,
  the constructions, the export) and the id-order comparison rule. This
  is the part a new port implements first.

## What an export says

`heads` lists the exported instances' ids in argument order. `nodes`
holds every reachable instance, deduplicated at any depth, keyed by id.
Each node carries:

- `type` — `{ collection, path }`: where the type was DECLARED, not just
  its name (names collide across levels and collections).
- `own` — the fields this level set itself, so shadowed values stay with
  their writer. A field holding another mnemonica instance exports as
  `{ "$ref": <id> }` and the referenced instance joins `nodes`.
- `parent` — the parent node's id, `null` at the root.
- `args` / `props` — construction arguments and metadata (e.g.
  timestamp), present only when the caller opts in.
- Values JSON cannot carry (functions, symbols, NaN, infinities, cycles,
  …) export as tagged placeholders `{ "$mnemonica": "unsupported",
  "kind": … }` — never an error.

## The .tactica contract

Next to the lineage contract, lethe holds a second one: the
`tactica/` directory with JSON Schemas (draft 2020-12) for the three
type-graph files every mnemonica generator writes into a `.tactica/`
directory — `hierarchy.json`, `definitions.json`, `collections.json`
(format versions "1.0" and "1.1"). `tactica/README.md` carries the
rules a schema cannot say (how `fullPath` is built, the
`collection_N::` prefix, `file:line:col` locations, cross-file joins),
and `testdata/tactica/` holds recorded real-world samples the schemas
must accept. Generators validate their fresh output against these
schemas in their own test suites.

## Using lethe

Implementations depend on lethe as a devDependency, validate their export
against `lineage.schema.json`, and reproduce `fixture.json` in their test
suite. Nothing here is published to be imported at runtime — the schema
and fixture are the deliverables; `npm test` in this package only checks
that the fixture still validates against the schema it ships with.
