# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Deco is a Go CLI for managing game design documents as structured, validated YAML. See [SPEC.md](SPEC.md) for full specification.

This project uses **bd** (beads) for issue tracking.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Session Completion

Before ending a session: close finished issues, create issues for remaining work, then push.

```bash
bd close <id>           # Close completed work
bd sync --flush-only    # Export beads
git add -A && git commit -m "..." && git push
```

Work is not complete until `git push` succeeds.
