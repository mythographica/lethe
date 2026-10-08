# AGENTS.md — @mnemonica/lethe

Guidance for AI agents working on the lethe package.

## What this package is

The cross-language lineage contract: `lineage.schema.json` +
`testdata/lineage/fixture.json` + the recipe README. It exists so every
mnemonica port (Go, JavaScript, Python, …) exports instance lineage in
one format. See [`README.md`](./README.md) for the shape and meaning.

A second contract lives in `tactica/`: the JSON Schemas for the six
`.tactica/` files — the three type-graph files (`hierarchy.json`,
`definitions.json`, `collections.json`) plus the analysis files
(`usages.json`, `flow.json`, `eds.json`) — described by
`tactica/README.md` and exercised against recorded samples in
`testdata/tactica/`. The same cross-language rule applies — any schema
change must be writable by every generator (tactica, `mnemonica.stubgen`,
`mnemonica-gen`).

## The contract rules (non-negotiable)

1. **The schema is a cross-language contract.** Any change to
   `lineage.schema.json` — a required field, a new kind, a loosened
   constraint — must be reproducible by EVERY implementation. Before
   changing the schema, list the implementations that must follow and
   confirm the migration path; a schema edit without ported follow-ups
   breaks the contract.
2. **`version` is `"1"`.** Incompatible changes mean a new version field
   and a new fixture, never a silent edit of the current one.
3. **The fixture is canonical.** `fixture.json` and its README recipe
   change only with a schema version bump, and every port's byte-for-byte
   reproduction must be updated in the same breath.
4. **Instance ids are implementation-specific.** Never put real ids in
   the fixture or the docs — the semantic placeholders (`s`, `a1`, `u`,
   `a2`) and the first-encounter mapping rule are the contract.
5. **Tagged placeholders, never exceptions.** A value JSON cannot carry
   exports as `{ "$mnemonica": "unsupported", "kind": … }`; an export
   that throws on a value is a bug in the exporting implementation.

## Docs rules

- Present state only: no dated records, no changelog entries, no
  name attributions. Why-a-design-is-so may stay when it explains the
  present contract.
- No new dependencies beyond the test runner's `ajv` (devDependency).
  This package must stay trivially portable — it is a contract, not a
  runtime.

## Git

The owner runs git. Agents do not commit, push, or tag here.
