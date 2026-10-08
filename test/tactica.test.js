'use strict';
// The .tactica contract tests: the six schemas compile, every recorded
// real-world sample validates against the schema it ships with, and
// hand-made invalid samples are rejected. Generators (tactica, the Go and
// Python ports) run the same validation against their own fresh output.
const { test } = require('node:test');
const { readFileSync, readdirSync, existsSync } = require('node:fs');
const { join } = require('node:path');
const Ajv2020 = require('ajv/dist/2020');

const root = join(__dirname, '..');
const schemas = {
	hierarchy: JSON.parse(readFileSync(join(root, 'tactica', 'hierarchy.schema.json'), 'utf8')),
	definitions: JSON.parse(readFileSync(join(root, 'tactica', 'definitions.schema.json'), 'utf8')),
	collections: JSON.parse(readFileSync(join(root, 'tactica', 'collections.schema.json'), 'utf8')),
	usages: JSON.parse(readFileSync(join(root, 'tactica', 'usages.schema.json'), 'utf8')),
	flow: JSON.parse(readFileSync(join(root, 'tactica', 'flow.schema.json'), 'utf8')),
	eds: JSON.parse(readFileSync(join(root, 'tactica', 'eds.schema.json'), 'utf8')),
};
const sampleRoot = join(root, 'testdata', 'tactica');

const compile = (schema) => new Ajv2020({ allErrors: true, strict: true }).compile(schema);
const parse = (path) => JSON.parse(readFileSync(path, 'utf8'));

const samples = readdirSync(sampleRoot, { withFileTypes: true })
	.filter((d) => d.isDirectory() && d.name !== 'invalid')
	.map((d) => d.name);

test('the three .tactica schemas compile', () => {
	for (const [name, schema] of Object.entries(schemas)) {
		compile(schema);
	}
});

for (const sample of samples) {
	const dir = join(sampleRoot, sample);
	for (const kind of Object.keys(schemas)) {
		const file = join(dir, `${kind}.json`);
		if (!existsSync(file)) { continue; }
		test(`${sample}/${kind}.json validates against ${kind}.schema.json`, () => {
			const validate = compile(schemas[kind]);
			const valid = validate(parse(file));
			if (!valid) {
				throw new Error(`${sample}/${kind}.json failed: ` + JSON.stringify(validate.errors, null, 2));
			}
		});
	}
}

test('hand-made invalid samples are rejected', () => {
	const invalidDir = join(sampleRoot, 'invalid');
	for (const file of readdirSync(invalidDir)) {
		const kind = file.split('.')[1];
		const validate = compile(schemas[kind]);
		if (validate(parse(join(invalidDir, file)))) {
			throw new Error(`${file} unexpectedly validated against ${kind}.schema.json`);
		}
	}
});
