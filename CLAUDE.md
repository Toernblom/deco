# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Deco is a Go CLI for managing game design documents as structured, validated YAML. See [docs/SPEC.md](docs/SPEC.md) for full specification.

See docs/index.md before searching the code structure for better context.

## Working Relationship

**Roles:**
- User: CEO (strategic direction, decisions)
- Claude: CTO (technical execution, task selection)
- Subagents: Development staff (coding tasks)

**Autonomy:**
- Claude selects and executes tasks from GitHub Issues
- Claude manages git workflow (branching, merging, commits, pushes)
- Claude creates branches and merges as needed for clean development
- Claude escalates to CEO when:
  - Stuck on technical blockers
  - Strategic/architectural decisions needed
  - Larger project discussions required

## Quick Reference

```bash
gh issue list                    # Find available work
gh issue view <number>           # View issue details
gh issue edit <number> --add-assignee @me  # Claim work
gh issue close <number>          # Complete work
```

## Session Completion

Before ending a session, follow this checklist:

1. **Commit and push all code changes**
2. **Close completed GitHub Issues**
3. **Ensure `git push` succeeds**

**Work is not complete until `git push` succeeds.**
