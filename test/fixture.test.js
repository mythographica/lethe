'use strict';
// The lethe contract test: the shared fixture must validate against the
// schema it ships with. Ports additionally reproduce the fixture bytes
// (see testdata/lineage/README.md) — that lives in each implementation.
const { test } = require('node:test');
const { readFileSync } = require('node:fs');
const { join } = require('node:path');
const Ajv2020 = require('ajv/dist/2020');

const schema = JSON.parse(readFileSync(join(__dirname, '..', 'lineage.schema.json'), 'utf8'));
const fixture = JSON.parse(readFileSync(join(__dirname, '..', 'testdata', 'lineage', 'fixture.json'), 'utf8'));

test('the shared fixture validates against the lethe schema', () => {
	const ajv = new Ajv2020({ allErrors: true, strict: true });
	const validate = ajv.compile(schema);
	const valid = validate(fixture);
	if (!valid) {
		throw new Error('fixture failed schema validation: ' + JSON.stringify(validate.errors, null, 2));
	}
});
