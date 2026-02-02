# Issue Template for AI Agents

This template ensures issues contain enough context for an AI agent with **zero prior knowledge** to complete the work independently. When creating issues, fill sections thoroughly—the agent working this issue will likely have no access to your current conversation or research.

**Workflow:** Research → Create detailed issue → Context cleared → Fresh agent claims and completes work

---

## Issue Types

| Type | Use When |
|------|----------|
| **bug** | Deviates from spec or broke something that worked |
| **task** | Discrete work: refactor, config, cleanup, migration |
| **feature** | New capability that doesn't exist yet |
| **epic** | Container for related work (not directly workable) |
| **question** | Research needed before work can be defined |
| **docs** | Documentation-only, no code changes |

## Priority Levels

| Priority | Meaning |
|----------|---------|
| **P0** | Drop everything (broken build, security, data loss) |
| **P1** | Do soon (blocks other work, significant bugs) |
| **P2** | Normal (default for most work) |
| **P3** | Nice to have (do when convenient) |
| **P4** | Backlog (may never happen) |

---

## Required Sections

### Context / Problem
What's broken, missing, or suboptimal? Include:
- **Current behavior** (what happens now)
- **Why it matters** (impact on users, system, or other work)
- **Relevant files/functions** (specific paths, line numbers if known)

### Goal / Outcome
What does "done" look like? Be specific enough that success is unambiguous.

### Scope
- **In:** What this issue covers
- **Out:** What's explicitly NOT included (prevents scope creep)

### TDD Plan
- **Tests to write first:** Specific test cases that should fail before implementation
- **Expected failures:** What the test output should show pre-implementation

### Acceptance Criteria
Checkboxes that must all pass for the issue to close:
```
- [ ] Unit tests pass
- [ ] Behavior matches spec
- [ ] No regressions in related functionality
```

---

## Supporting Sections

### Proposed Approach (optional but recommended)
Outline the implementation strategy. Helps agents avoid wrong paths:
- Key steps or phases
- Technical decisions already made
- Alternatives considered and why rejected

### Test / Verification Plan
How to verify the work is complete:
- **Commands:** Exact commands to run (e.g., `go test ./internal/...`)
- **Expected outcome:** What success looks like

### Risks / Edge Cases
Potential pitfalls the agent should watch for:
- Edge cases that need handling
- Backwards compatibility concerns
- Performance implications

### Dependencies
- **Blocks:** Issues that cannot start until this completes
- **Blocked by:** Issues that must complete before this can start

### Notes / Links
Supporting context: related issues, docs, specs, external references, prior conversations summarized.

---

## Before Creating the Issue

Ask yourself: Could an agent with no prior context complete this work?

- [ ] **Context is self-contained** - No references to "the conversation above" or assumed knowledge
- [ ] **Files/locations identified** - Agent knows where to look
- [ ] **Success is measurable** - Clear acceptance criteria, not vague goals
- [ ] **TDD path is clear** - Agent knows what tests to write first
- [ ] **Scope is bounded** - Explicit in/out prevents drift
