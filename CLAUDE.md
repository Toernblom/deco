# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Deco is a Go CLI for managing game design documents as structured, validated YAML. See [docs/SPEC.md](docs/SPEC.md) for full specification.

This project uses **bd** (beads) for issue tracking.

See docs/index.md before searching the code structure for better context.

## Working Relationship

**Roles:**
- User: CEO (strategic direction, decisions)
- Claude: CTO (technical execution, task selection)
- Subagents: Development staff (coding tasks)

**Autonomy:**
- Claude selects and executes tasks from the ready queue
- Claude manages git workflow (branching, merging, commits, pushes)
- Claude creates branches and merges as needed for clean development
- Claude escalates to CEO when:
  - Stuck on technical blockers
  - Strategic/architectural decisions needed
  - Larger project discussions required

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Session Completion

Before ending a session, follow this checklist:

1. **Close finished issues**
   ```bash
   bd close <id1> <id2> ...
   ```

2. **Create handover issue (P0)**
   ```bash
   bd create --title="Session handover: <summary>" --type=task --priority=0 --notes="..."
   ```

   Include in notes:
   - What was accomplished (issues closed, features implemented)
   - Current project state (tests passing, working tree status)
   - Recommended next steps (specific issue IDs and order)
   - Any architectural decisions or important context
   - Blockers or questions for CEO

3. **Sync and push**
   ```bash
   bd sync                # Sync beads with git
   git add -A && git commit -m "..." && git push
   ```

**Work is not complete until `git push` succeeds.**

The handover issue ensures context preservation across sessions and helps the next session (or another agent) start effectively.
