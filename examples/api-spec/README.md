# API Spec Example

A non-game example showing how to document a REST API with Deco. It demonstrates
custom block validation, schema rules, cross-node references, and constraints.

## Highlights

- **Custom blocks**: `endpoint`, `schema_field`, and `error_case`
- **Schema rules**: required custom fields for `endpoint`, `schema`, and `system`
- **Cross-node constraint**: endpoint rate limits must not exceed the platform default
- **Contracts**: login flow contract scenario
- **Issues**: resolved auth workflow question

## Layout

```
api-spec/
└── .deco/
    ├── config.yaml
    ├── history.jsonl
    └── nodes/
        ├── systems/
        ├── endpoints/
        ├── schemas/
        └── glossaries/
```

## Try It

```bash
cd examples/api-spec

deco validate

deco show endpoints/users

deco graph --format mermaid
```
