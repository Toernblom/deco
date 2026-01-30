# Deco / Design Compilator

Deco is a **design engine and CLI** for building games (and game systems) from **structured design**, not a swamp of Markdown files and forgotten references.

Instead of treating a GDD as a pile of prose, Deco treats it like **source code**:

- Design is stored as **typed YAML nodes** with stable IDs.
- Nodes reference each other through explicit links, not “remember to update page 17”.
- Deco validates the whole design graph (schemas, references, constraints).
- Deco compiles design into **contracts** (BDD / test specs) and other artifacts.
- AI agents can generate or propose changes as **patch operations**, while Deco enforces correctness.

The point is simple: **design changes should propagate predictably**, and tests/contracts should prove nothing drifted.

---

## Why Deco exists

If you’ve ever maintained a big GDD, you know how it ends:

- “Loose ends” everywhere (TBDs, contradictory rules, outdated pages).
- Updating one system means hunting references across 30 documents.
- An LLM helps write, then loses context, and now you have two truths.

Deco replaces that with a model that is:

- **indexable**
- **diffable**
- **refactorable**
- **compilable**
- **validated**

So your GDD becomes something you can actually trust.

---

## Core idea

### Structured design nodes are the source of truth

A node represents a system, feature, mechanic, entity, item, etc.

Each node has:
- stable `id` (the anchor of truth)
- `refs` to other nodes
- structured `sections` with blocks (rules, tables, mechanics, parameters)
- `issues` for TBDs/questions so “loose ends” are tracked
- optional `contracts` describing behavior in testable form

Humans can still read and write it, but the engine can also reason about it.

---

## What Deco does ( TBD, NOT YET DONE AND NEED THOROUGH INVESTIGATION )

### ✅ Design graph management
- Parse all YAML nodes into a graph
- Validate schemas
- Validate references (no dangling links)
- Build reverse links (“used-by”) automatically
- Track issues and TODOs as first-class objects

### ✅ CLI-driven updates (no manual doc whack-a-mole)
- Patch nodes with `--set`, `--append`, `--unset`
- Rename/move nodes with automatic ref updates
- Show impact of changes across the graph

### ✅ Contracts-first design (optional but powerful)
Deco can store and/or compile contracts from design:

- BDD-style scenarios
- contract specs (YAML/JSON)
- generated `.feature` files (Gherkin) for existing runners
- test scaffolds for game code verification

### ✅ AI-friendly, engine-controlled
Non-technical users can dump text into a chat UI.

- LLM converts text into structured YAML **or patch ops**
- Deco validates + applies changes deterministically
- The engine, not the LLM, is the memory and source of truth

---

## Repository layout (suggested)

```txt
/design
  /systems
    /settlement
      colonists.yaml
      housing.yaml
/contracts
  /systems
    /settlement
      colonists.contract.yaml   # optional, if contracts are separate
/build
  /docs                          # generated human docs (optional)
  /gherkin                       # generated .feature files (optional)
  /context                       # generated LLM context bundles (optional)
deco.schema.json                 # schema or CUE definition
