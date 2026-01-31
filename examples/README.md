# Deco Examples

Example game design documents showcasing deco's features.

## Examples

### Snake (`snake/`)
Classic arcade game demonstrating:
- **Systems**: Core game loop, scoring
- **Items**: Food with different types
- **References**: Nodes linking to each other
- **Contracts**: BDD-style test scenarios
- **Parameters**: Configurable tick rate, points

### Space Invaders (`space-invaders/`)
Retro shooter demonstrating:
- **Entities**: Player ship, alien types
- **Tables**: Alien type definitions with stats
- **Mechanics**: Formation movement, wave progression
- **Rules**: Game flow and constraints

## Usage

```bash
# Navigate to an example
cd examples/snake

# List all nodes
deco list

# Show a specific node
deco show systems/core

# Validate the design
deco validate
```

## Structure

Each example follows the standard deco layout:
```
example/
└── .deco/
    ├── config.yaml      # Project config
    ├── history.jsonl    # Change history
    └── nodes/           # Design documents
        ├── systems/
        ├── items/
        └── entities/
```
