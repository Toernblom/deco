# PRD Template (Beads Issues)

Use this template in the **description** of a Beads issue. This repo tracks work in `.beads/issues.jsonl` via `bd`.

## Quick Start
1. `bd create --title="<feature>" --type=task --priority=2`
2. Paste the template below into the issue description.
3. Update status as you work: `bd update <id> --status in_progress`

## Issue Description Template
**Summary**
[One paragraph: what is being built and why]

**Problem / Opportunity**
[What is broken or missing? Who is impacted?]

**Goals**
- [Goal 1]
- [Goal 2]

**Non-goals**
- [Non-goal 1]
- [Non-goal 2]

**User Stories**
- As a [user], I want [capability], so that [outcome].

**Scope**
- In scope: [List the concrete deliverables]
- Out of scope: [Explicit exclusions]

**Constraints**
- [Performance, compatibility, API stability, timeline, etc.]

**Proposed Approach**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**TDD Plan (Required)**
- Tests to write first (red): [List failing test cases and files]
- Minimal implementation (green): [Smallest change to pass tests]
- Refactor: [Cleanup, naming, structure, additional coverage]

**Acceptance Criteria**
- [ ] [Testable outcome 1]
- [ ] [Testable outcome 2]
- [ ] [Testable outcome 3]

**Test Plan**
- [Commands or suites to run]
- [Edge cases to cover]

**Dependencies**
- [Issue IDs, if any]
- If blocked, add deps: `bd dep add <this-issue> <depends-on>`

**Files / Areas Touched**
- [Path(s) or modules]

**Risks**
- [Risk 1]
- [Risk 2]

**Open Questions**
- [Question 1]
- [Question 2]
