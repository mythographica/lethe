# The `.tactica` format

`tactica/hierarchy.schema.json`, `tactica/definitions.schema.json`,
`tactica/collections.schema.json` — the JSON Schema (draft 2020-12)
contract for the three type-graph files every mnemonica generator writes
into a `.tactica/` directory: tactica (TypeScript), and later the Python
(`mnemonica.stubgen`) and Go (`mnemonica-gen`) ports. A generator's test
suite validates its fresh output against these schemas — a mismatch fails
a test instead of silently breaking the graph a reader builds.

`tactica/usages.schema.json`, `tactica/flow.schema.json` and
`tactica/eds.schema.json` are the same contract for the three analysis
files (usage sites, native-instance flow, dive wrap sites), drafted from
what tactica writes today. Their entry `kind` fields are open strings —
the observed values are documented per schema, and readers tolerate values
they do not know rather than failing.

`version` is `"1.0"` or `"1.1"`; readers accept both. Absent `version`
means 1.0 (pre-versioning output is identical in shape to 1.0).
Incompatible changes mean a new version string, never a silent edit.

## The six files

| file | content | required files? |
|---|---|---|
| `hierarchy.json` | the type tree: roots and nested `children` | no — a source may ship only what it has |
| `definitions.json` | per-type declaration metadata, keyed by fullPath | no |
| `collections.json` | the types-collection manifest | no |
| `usages.json` | per-type usage sites (lookup/instantiation/propertyAccess), keyed by short type name | no |
| `flow.json` | native-instance flow sites (instantiation/passAsArg/propertyRead/...), keyed by native type name | no |
| `eds.json` | dive wrap sites + context consumes, keyed by owning type name or `"unknown"`; absent where the source has no dive | no |

A `.tactica/` directory may contain any subset (a TypeScript-less source
writes none of them; a graph-only consumer reads `hierarchy.json` alone).
The remaining `.tactica/` files (`types.ts`, `instrumentation.json`,
`scopes.json`, `modules.json`) are NOT part of this contract: `types.ts`
is TypeScript-only by definition, `instrumentation.json` `points` are
parked (see plan boundary-cleanup step 2), and scopes/modules keep their
own evolution cadence.

## Rules a schema cannot say

### `fullPath` — the one identity of a type

A type's fullPath is its **only** identity. There are no numeric node
ids; every join between the files (and into usages/flow/eds data) is a
string match on the fullPath.

- Built from dotted declared names: a subtype's fullPath is
  `<parentFullPath>.<name>`.
- The `collection_N::` prefix appears **only at collection roots**, and
  only when the type belongs to a minted `createTypesCollection()`
  collection: `collection_1::Backend.User`. Subtypes inherit the prefix
  through the parent's fullPath.
- The **default collection** (types from the global registry — plain
  `define()`, `lazy()`, `mnemonica.decorate()`) never carries a prefix.
  The absence of the prefix IS its identity; it is never written as
  `default::` or similar.
- Collection ids are minted in first-encounter order during analysis:
  `collection_1`, `collection_2`, … Aliases of one collection variable
  reuse its id.

### `location`

Always `file:line:col` with **1-based** line and column, joined by
colons, relative to the **project root** (the tsconfig directory for
tactica), with forward slashes. Files outside the project root may stay
absolute. The single exception: the default-collection entry in
`collections.json` has `location: null` — the global registry has no
declaration site. Readers parse with `^(.+):(\d+):(\d+)$`; a location
that does not match is treated as unknown, never as fatal.

### Cross-file joins

- `definitions.json` keys are fullPaths and must equal the corresponding
  `hierarchy.json` node's `fullPath` — for both files the hierarchy's
  `parent` pointer and the definition's `parent` field hold the parent's
  fullPath (`null` at roots).
- `collections.json` `id` is the join key for the `collection_N::`
  prefix. A fullPath with no prefix belongs to the entry with `id: null`
  (readers key it as `"defaultTypes"`).
- `hierarchy.json` nodes without a `definitions.json` entry (and vice
  versa) are tolerated by readers; the contract expects generators to
  emit both for every type they discover.

### Required vs optional

- Top-level `version` and `generatedAt` are **optional** in all three
  files. `generatedAt` (ISO-8601) is informational; readers that use it
  treat it as a freshness hint, not a validity gate. Hand-written and
  pre-versioning files in the wild omit both.
- `hierarchy.json`: every node carries all four fields (`name`,
  `fullPath`, `location`, `children`); `children` is always present,
  possibly empty.
- `definitions.json`: every entry carries all six fields (`name`,
  `location`, `kind`, `parent`, `strictChain`, `blockErrors`).
  `kind` is `"define"` (define()/lazy()) or `"decorate"`
  (@decorate()/Type.decorate()) in generator output; the schema keeps it
  an open string because readers tolerate other values rather than
  failing.
- `collections.json`: `id`, `name`, `location` are always present
  (`id`/`location` may be `null`, only for the default entry);
  `registryInterface` and `language` are optional.

## Version 1.1 additions

- `collections[].language`: `"typescript" | "python" | "go"` —
  the language of that collection's codebase. Absent means `"typescript"`
  (all 1.0 output is TypeScript). Each language codebase is its own
  collection set; one `.tactica/` graph may mix languages.
- `collections[].registryInterface` becomes **optional** (a TypeScript
  interface name; meaningless for other languages). When present it is a
  non-empty string; when absent it is omitted, never `null`.

## Decorated collection roots

The schema has always required `definitions.json` keys to equal the
hierarchy fullPaths — including the `collection_N::` prefix on
`@MyCollection.decorate()` roots. An early tactica build keyed those
definitions without the prefix (the root never joined between the two
files); that was a generator bug, fixed 1.0-compatibly — the format did
not change.

## Samples

`testdata/tactica/` holds recorded real-world output (trimmed to
representative subsets, paths unchanged):

- `multi-collection/` — default + two named collections with
  `registryInterface` (from mnemographica's own `.tactica/`), plus trimmed
  `usages.json`/`flow.json`;
- `decorate/` — `kind: "decorate"` chains (from test-core/todo-app);
- `hierarchy-only/` — a directory shipping just `hierarchy.json`;
- `eds/` — trimmed dive wrap sites (from tactica-nestjs's `.tactica/`);
- `v1.1-python/` — a hand-made 1.1 sample showing `language` and an
  absent `registryInterface`;
- `invalid/` — hand-made samples that must be rejected.

`npm test` compiles the three schemas and checks every sample: valid
ones pass, invalid ones fail.
