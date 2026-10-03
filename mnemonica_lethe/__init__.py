"""mnemonica-lethe: the cross-language lineage contract, packaged.

Exposes `lineage.schema.json` and `testdata/lineage/fixture.json` — the
SAME files, via symlinks, through importlib.resources — so every port can
validate its exports against the contract and reproduce the fixture.
"""

from importlib.resources import files
from typing import Any, cast

__all__ = [
    "FIXTURE_RESOURCE",
    "SCHEMA_RESOURCE",
    "fixture",
    "fixture_path",
    "fixture_text",
    "schema",
    "schema_path",
    "schema_text",
]

SCHEMA_RESOURCE = "lineage.schema.json"
FIXTURE_RESOURCE = "data/fixture.json"


def schema_path() -> str:
    """Filesystem path of the schema (the same file, not a copy)."""
    result = str(files(__name__).joinpath(SCHEMA_RESOURCE))
    return result


def fixture_path() -> str:
    """Filesystem path of the shared fixture (the same file, not a copy)."""
    result = str(files(__name__).joinpath(FIXTURE_RESOURCE))
    return result


def schema_text() -> str:
    """The schema as text."""
    result = files(__name__).joinpath(SCHEMA_RESOURCE).read_text(encoding="utf-8")
    return result


def fixture_text() -> str:
    """The fixture as text."""
    result = files(__name__).joinpath(FIXTURE_RESOURCE).read_text(encoding="utf-8")
    return result


def schema() -> dict[str, Any]:
    """The schema as a parsed JSON object."""
    import json

    result = cast(dict[str, Any], json.loads(schema_text()))
    return result


def fixture() -> dict[str, Any]:
    """The fixture graph as a parsed JSON object."""
    import json

    result = cast(dict[str, Any], json.loads(fixture_text()))
    return result
