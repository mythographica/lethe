# The `.tactica` format

`tactica/hierarchy.schema.json`, `tactica/definitions.schema.json`,
`tactica/collections.schema.json` — the JSON Schema (draft 2020-12)
contract for the three type-graph files every mnemonica generator writes
into a `.tactica/` directory: tactica (TypeScript), and later the Python
(`mnemonica.stubgen`) and Go (`mnemonica-gen`) ports. A generator's test
suite validates its fresh output against these schemas — a mismatch fails
a test instead of silently breaking the graph a reader builds.

`tactica/usages.schema.json`, `tactica/flow.schema.json`,
`tactica/eds.schema.json` and `tactica/control.schema.json` are the same
contract for the four analysis files (usage sites, native-instance flow,
dive wrap sites, control flow), drafted from what tactica writes today.
Their entry `kind` fields are open strings — the observed values are
documented per schema, and readers tolerate values they do not know rather
than failing.

`version` is required and is `"1.2"`; the schemas accept nothing else.
An older `.tactica/` is regenerated, not read. Incompatible changes mean
a new version string, never a silent edit.

## The seven files

| file | content | required files? |
|---|---|---|
| `hierarchy.json` | the type tree: roots and nested `children` | no — a source may ship only what it has |
| `definitions.json` | per-type declaration metadata, keyed by fullPath | no |
| `collections.json` | the types-collection manifest | no |
| `usages.json` | per-type usage sites (lookup/instantiation/propertyAccess), keyed by short type name | no |
| `flow.json` | native-instance flow sites (instantiation/passAsArg/propertyRead/...), keyed by native type name | no |
| `eds.json` | dive wrap sites + context consumes, keyed by owning type name or `"unknown"`; absent where the source has no dive | no |
| `control.json` | control flow: `callers` (the call paths from entry points to constructor invocations) and plugin-supplied `points` | no |

A `.tactica/` directory may contain any subset (a graph-only consumer
reads `hierarchy.json` alone). Any other file a generator writes into
`.tactica/` — language bindings, internal working data — is NOT part of
this contract, and readers never depend on it: everything they need is in
the seven files above (the scopeIds that reach readers are
self-describing, see `scopeId` below).

### `control.json`

Data flow is the runtime lineage — which instance came from which.
`control.json` is the control flow that makes it happen, found statically:

- `callers.anchors` — every constructor invocation site (`new X`), pinned
  to the scope holding it (`holderScopeId`), with the constructed type's
  fullPath (`typePath`) and, when known, the constructor text as written,
  the variable bound to the new instance, and that variable's first
  reassignment (`terminatedAt`). `rooted: true` marks an invocation at a
  file's top level.
- `callers.nodes` — the scopes holding an invocation or sitting on a call
  path to one, walked outward until no caller remains; those carry
  `starter: true` — the application's entry points, often scopes a
  framework calls. Node `kind` is an open string: tactica writes
  `module`/`function`/`method`/`arrow`; other languages add their own.
- `callers.edges` — `caller → callee`, the callee closer to the
  invocation. Paths may cycle; readers must not assume a tree.
- `points` — framework steps around the entry points (guards,
  interceptors, pipes, filters, middleware for the nestjs plugin),
  supplied only by generator plugins: absent or empty without one. `kind`
  and `scope` are open strings, so each language's plugins bring their
  own.

`scopeId` is a scope's identity: the file path for a file's top level,
`file:line:col` of the scope start for a function, method or closure. It
is the same vocabulary in `usages.json` (`holderScopeId`), `eds.json`
(`scopeId`, `callbackScopeId`) and `control.json`.

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
colons, relative to the **project root** (the directory the generator
analyzes), with forward slashes. Files outside the project root may stay
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
- `control.json`: every edge's `caller` and `callee`, and every anchor's
  `holderScopeId`, is a `callers.nodes` `scopeId`. Every anchor's
  `location` is the `location` of an `instantiation` entry in
  `usages.json`, and every such entry has its anchor; every usage's
  `holderScopeId` is a node. `eds.json` scopeIds use the same vocabulary
  but need not be nodes — a wrap site may sit off every invocation path.

### Required vs optional

- Top-level `version` is **required** in all seven files; `generatedAt`
  (ISO-8601) is optional and informational — readers that use it treat
  it as a freshness hint, not a validity gate.
- `hierarchy.json`: every node carries all four fields (`name`,
  `fullPath`, `location`, `children`); `children` is always present,
  possibly empty.
- `definitions.json`: every entry carries all six fields (`name`,
  `location`, `kind`, `parent`, `strictChain`, `blockErrors`).
  `kind` is `"define"` (define()/lazy()) or `"decorate"`
  (@decorate()/Type.decorate()) in generator output; the schema keeps it
  an open string because readers tolerate other values rather than
  failing. An entry may also carry `fields` (optional): the type's own
  fields in declaration order, each `{ name, type, optional }`, where
  `type` is source-language text for display. "Own" means declared ahead
  of time by the type's definition — its `this` annotation, class body
  or constructor assignments — not what its instance layer holds at
  runtime; a field the parent also declares is listed again. Absent
  `fields` means not recorded; an empty array means the type has no
  fields. An entry may also carry `args` (optional): the constructor's parameters in declaration order —
  what it consumes to produce its fields — each
  `{ name, type, optional?, kind? }`. The receiver (`this`, the parent
  instance) is not listed; `parent` already names its type. `optional`
  absent means false; a rest parameter carries none, since it may always
  receive nothing. `kind` absent means an ordinary parameter; generators
  write `"rest"` (TS `...xs`, Go `xs ...T`, Python `*args`, Java
  `T... xs`), `"keywordRest"` (Python `**kwargs`), `"keywordOnly"` (a
  Python parameter after `*`) or `"positionalOnly"` (a Python parameter
  before `/`); readers tolerate other values. Absent `args` means not
  recorded; an empty array means the constructor takes none.
- `collections.json`: `id`, `name`, `location` and `language` are
  always present (`id`/`location` may be `null`, only for the default
  entry). `language` is `"typescript" | "python" | "go"` — the language
  of that collection's codebase; there is no default. Each language
  codebase writes its own `.tactica/`; readers open one at a time.

## Samples

`testdata/tactica/` holds recorded real-world output (trimmed to
representative subsets, paths unchanged, `version` set to `"1.2"`):

- `multi-collection/` — default + two named collections (from
  mnemographica's own `.tactica/`), plus trimmed `usages.json`/`flow.json`;
- `decorate/` — `kind: "decorate"` chains (from test-core/todo-app);
- `hierarchy-only/` — a directory shipping just `hierarchy.json`;
- `eds/` — trimmed dive wrap sites (from tactica-nestjs's `.tactica/`);
- `nestjs-control/` — `control.json` trimmed from tactica-nestjs's
  `.tactica/` (two anchors with every caller path up to the starters, one
  point per kind);
- `python/` — a hand-made sample of a Python collection;
- `invalid/` — hand-made samples that must be rejected.

`npm test` compiles the schemas and checks every sample: valid ones
pass, invalid ones fail.
