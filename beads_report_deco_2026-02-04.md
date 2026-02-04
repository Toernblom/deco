# Beads Export

*Generated: Wed, 04 Feb 2026 18:37:14 CET*

## Summary

| Metric | Count |
|--------|-------|
| **Total** | 199 |
| Open | 16 |
| In Progress | 0 |
| Blocked | 0 |
| Closed | 183 |

## Quick Actions

Ready-to-run commands for bulk operations:

```bash
# Close open items (16 total, showing first 10)
bd close deco-1iuf deco-jfye deco-hfu4 deco-kz2i deco-02tq deco-gpli deco-v1vi deco-2g91 deco-6iv deco-23bx

# View high-priority items (P0/P1)
bd show deco-1iuf deco-jfye deco-hfu4 deco-kz2i

```

## Table of Contents

- [🟢 deco-1iuf Session handover: strict block validation](#deco-1iuf-session-handover-strict-block-validation)
- [🟢 deco-jfye Validate node IDs to prevent path traversal](#deco-jfye-validate-node-ids-to-prevent-path-traversal)
- [🟢 deco-hfu4 Add optimistic concurrency controls for writes](#deco-hfu4-add-optimistic-concurrency-controls-for-writes)
- [🟢 deco-kz2i Add LLM context export command](#deco-kz2i-add-llm-context-export-command)
- [🟢 deco-02tq Add propose/review mode for AI patches](#deco-02tq-add-propose-review-mode-for-ai-patches)
- [🟢 deco-gpli Allow explicit actor identity in audit history](#deco-gpli-allow-explicit-actor-identity-in-audit-history)
- [🟢 deco-v1vi Add JSON output for validation and mutation errors](#deco-v1vi-add-json-output-for-validation-and-mutation-errors)
- [🟢 deco-2g91 Add reference/schema discovery commands](#deco-2g91-add-reference-schema-discovery-commands)
- [🟢 deco-6iv Deco GDD: Project-lead vision for end-state product](#deco-6iv-deco-gdd-project-lead-vision-for-end-state-product)
- [🟢 deco-23bx Clean invalid field from snake example](#deco-23bx-clean-invalid-field-from-snake-example)
- [🟢 deco-oxu2 Add tests for Resolve*Path with relative rootDir](#deco-oxu2-add-tests-for-resolve-path-with-relative-rootdir)
- [🟢 deco-yle3 Update CLI help text to reflect configurable nodes_path](#deco-yle3-update-cli-help-text-to-reflect-configurable-nodes-path)
- [🟢 deco-vtn5 Resolve*Path should be absolute or comments updated](#deco-vtn5-resolve-path-should-be-absolute-or-comments-updated)
- [🟢 deco-4p9 GDD playbook documentation](#deco-4p9-gdd-playbook-documentation)
- [🟢 deco-a5n Explain error / suggest fix output](#deco-a5n-explain-error-suggest-fix-output)
- [🟢 deco-9rm Compile GDD to LaTeX output](#deco-9rm-compile-gdd-to-latex-output)
- [⚫ deco-u0jb Session handover: custom block types implemented](#deco-u0jb-session-handover-custom-block-types-implemented)
- [⚫ deco-dag0 Session handover: performance optimizations](#deco-dag0-session-handover-performance-optimizations)
- [⚫ deco-5gxq Session handover: sync command fixes](#deco-5gxq-session-handover-sync-command-fixes)
- [⚫ deco-awz Session handover: Review workflow implementation in progress](#deco-awz-session-handover-review-workflow-implementation-in-progress)
- [⚫ deco-ry9 Session handover: Review workflow design complete](#deco-ry9-session-handover-review-workflow-design-complete)
- [⚫ deco-xqq Session handover: Constraint scope enforcement](#deco-xqq-session-handover-constraint-scope-enforcement)
- [⚫ deco-c4y Session handover: File location in validation errors](#deco-c4y-session-handover-file-location-in-validation-errors)
- [⚫ deco-4y9 Session handover: Block validation complete](#deco-4y9-session-handover-block-validation-complete)
- [⚫ deco-lwl Session handover: contract validation CLI complete](#deco-lwl-session-handover-contract-validation-cli-complete)
- [⚫ deco-v4o Session handover: contract ref validation added (deco-wxh)](#deco-v4o-session-handover-contract-ref-validation-added-deco-wxh)
- [⚫ deco-ykl Session handover: contract validator added (deco-0sk)](#deco-ykl-session-handover-contract-validator-added-deco-0sk)
- [⚫ deco-g50 Session handover: contract parser added (deco-4ci)](#deco-g50-session-handover-contract-parser-added-deco-4ci)
- [⚫ deco-gbk Session handover: audit history added](#deco-gbk-session-handover-audit-history-added)
- [⚫ deco-fre Session handover: stats command added](#deco-fre-session-handover-stats-command-added)
- [⚫ deco-vq3 Session handover: diff command added](#deco-vq3-session-handover-diff-command-added)
- [⚫ deco-4kg Session handover: graph command added](#deco-4kg-session-handover-graph-command-added)
- [⚫ deco-dyi Session handover: 3 CLI commands added](#deco-dyi-session-handover-3-cli-commands-added)
- [⚫ deco-duf Session handover: mv command and CLI epic complete](#deco-duf-session-handover-mv-command-and-cli-epic-complete)
- [⚫ deco-gw6 Session handover: Node rename service complete](#deco-gw6-session-handover-node-rename-service-complete)
- [⚫ deco-8rt Session handover: CLI mutation commands complete](#deco-8rt-session-handover-cli-mutation-commands-complete)
- [⚫ deco-flq Session handover: CLI commands (validate, list, show)](#deco-flq-session-handover-cli-commands-validate-list-show)
- [⚫ deco-934 Session handover: CLI foundation complete](#deco-934-session-handover-cli-foundation-complete)
- [⚫ deco-wxe Session handover: Service layer complete - 67% done](#deco-wxe-session-handover-service-layer-complete-67-done)
- [⚫ deco-tw3 Session handover: Service layer progress - 54% complete](#deco-tw3-session-handover-service-layer-progress-54-complete)
- [⚫ deco-8vg Session handover: Storage and error system enhancements complete](#deco-8vg-session-handover-storage-and-error-system-enhancements-complete)
- [⚫ deco-7de Session handover: Error system and storage layer complete](#deco-7de-session-handover-error-system-and-storage-layer-complete)
- [⚫ deco-c5k Session handover: Foundation complete, error system next](#deco-c5k-session-handover-foundation-complete-error-system-next)
- [⚫ deco-q5pk Validate unknown fields in RefLink and other nested structures](#deco-q5pk-validate-unknown-fields-in-reflink-and-other-nested-structures)
- [⚫ deco-snlr Add strict block field validation (configurable)](#deco-snlr-add-strict-block-field-validation-configurable)
- [⚫ deco-pc3o Content hash uses non-deterministic map ordering causing hash churn](#deco-pc3o-content-hash-uses-non-deterministic-map-ordering-causing-hash-churn)
- [⚫ deco-3o2h apply/rewrite don't write content hash, breaking sync detection](#deco-3o2h-apply-rewrite-don-t-write-content-hash-breaking-sync-detection)
- [⚫ deco-36rv CEL constraint engine missing spec'd capabilities (allNodes, custom fields)](#deco-36rv-cel-constraint-engine-missing-spec-d-capabilities-allnodes-custom-fields)
- [⚫ deco-nd7w Config paths (nodes_path, history_path) are hardcoded, not configurable](#deco-nd7w-config-paths-nodes-path-history-path-are-hardcoded-not-configurable)
- [⚫ deco-h91g sync swallows errors - can exit clean despite failures](#deco-h91g-sync-swallows-errors-can-exit-clean-despite-failures)
- [⚫ deco-0t5 CLI mutations don't record content_hash in history](#deco-0t5-cli-mutations-don-t-record-content-hash-in-history)
- [⚫ deco-jyn deco set panics on nested paths (pointer deref in patcher)](#deco-jyn-deco-set-panics-on-nested-paths-pointer-deref-in-patcher)
- [⚫ deco-86e Block data lost when parsing YAML (blocks[].data null)](#deco-86e-block-data-lost-when-parsing-yaml-blocks-data-null)
- [⚫ deco-5xa Errors: Rust-like error system](#deco-5xa-errors-rust-like-error-system)
- [⚫ deco-7yo CLI: Command implementations](#deco-7yo-cli-command-implementations)
- [⚫ deco-t6q Services: Business logic layer](#deco-t6q-services-business-logic-layer)
- [⚫ deco-elv Storage: Repository implementations](#deco-elv-storage-repository-implementations)
- [⚫ deco-16e Foundation: Project setup and domain types](#deco-16e-foundation-project-setup-and-domain-types)
- [⚫ deco-qw2b Validate examples/snake with deco validate](#deco-qw2b-validate-examples-snake-with-deco-validate)
- [⚫ deco-hond set resets status to draft but append/unset don't - inconsistent workflow](#deco-hond-set-resets-status-to-draft-but-append-unset-don-t-inconsistent-workflow)
- [⚫ deco-122q Ref validation ignores emits_events and vocabulary fields](#deco-122q-ref-validation-ignores-emits-events-and-vocabulary-fields)
- [⚫ deco-epm4 README block examples use key/value but validator requires name/datatype](#deco-epm4-readme-block-examples-use-key-value-but-validator-requires-name-datatype)
- [⚫ deco-7z6u YAML error location infrastructure exists but unused in validate output](#deco-7z6u-yaml-error-location-infrastructure-exists-but-unused-in-validate-output)
- [⚫ deco-chwx Status validation only checks presence, not allowed values](#deco-chwx-status-validation-only-checks-presence-not-allowed-values)
- [⚫ deco-fgs5 refactor: extend mv to update all reference types](#deco-fgs5-refactor-extend-mv-to-update-all-reference-types)
- [⚫ deco-rb08 sync: detect manual file renames and update references](#deco-rb08-sync-detect-manual-file-renames-and-update-references)
- [⚫ deco-oh0q --dry-run always exits 0 even when changes would occur](#deco-oh0q-dry-run-always-exits-0-even-when-changes-would-occur)
- [⚫ deco-cxl sync O(nodes × history) performance - needs single-pass indexing](#deco-cxl-sync-o-nodes-history-performance-needs-single-pass-indexing)
- [⚫ deco-8xo Content hash excludes Glossary, Contracts, LLMContext, Constraints, Custom, Kind](#deco-8xo-content-hash-excludes-glossary-contracts-llmcontext-constraints-custom-kind)
- [⚫ deco-mbz Add deco sync command for detecting and fixing unversioned changes](#deco-mbz-add-deco-sync-command-for-detecting-and-fixing-unversioned-changes)
- [⚫ deco-zn8 Add audit history to all node-modifying commands](#deco-zn8-add-audit-history-to-all-node-modifying-commands)
- [⚫ deco-p75 Backwards-compatible schema migrations](#deco-p75-backwards-compatible-schema-migrations)
- [⚫ deco-a5l AI patch/rewrite safety: validate gate + transactional apply + explicit diff](#deco-a5l-ai-patch-rewrite-safety-validate-gate-transactional-apply-explicit-diff)
- [⚫ deco-79k Enforce constraint scope (node kind/pattern) in validator](#deco-79k-enforce-constraint-scope-node-kind-pattern-in-validator)
- [⚫ deco-87y Custom block types with validation hooks](#deco-87y-custom-block-types-with-validation-hooks)
- [⚫ deco-apc Configurable schema rules (org-level constraints)](#deco-apc-configurable-schema-rules-org-level-constraints)
- [⚫ deco-rs0 Review workflow: approvals, status transitions, changelog notes](#deco-rs0-review-workflow-approvals-status-transitions-changelog-notes)
- [⚫ deco-e4g Expand issues/TBD system: filters, severity, per-node tracking](#deco-e4g-expand-issues-tbd-system-filters-severity-per-node-tracking)
- [⚫ deco-603 Add type-specific validation for blocks (rule/table/param/etc.)](#deco-603-add-type-specific-validation-for-blocks-rule-table-param-etc)
- [⚫ deco-hsx Require content for approved/published nodes (allow drafts without content)](#deco-hsx-require-content-for-approved-published-nodes-allow-drafts-without-content)
- [⚫ deco-cxk Decide and implement block schema (inline fields vs data map)](#deco-cxk-decide-and-implement-block-schema-inline-fields-vs-data-map)
- [⚫ deco-rib Define strict top-level schema with explicit extension mechanism](#deco-rib-define-strict-top-level-schema-with-explicit-extension-mechanism)
- [⚫ deco-0ja deco set cannot update non-string fields](#deco-0ja-deco-set-cannot-update-non-string-fields)
- [⚫ deco-a3o validate ignores unknown top-level keys and structural typos](#deco-a3o-validate-ignores-unknown-top-level-keys-and-structural-typos)
- [⚫ deco-8ov validate does not detect duplicate node IDs](#deco-8ov-validate-does-not-detect-duplicate-node-ids)
- [⚫ deco-3dy Set up CI/CD pipeline](#deco-3dy-set-up-ci-cd-pipeline)
- [⚫ deco-tzg Write tests for deco apply command](#deco-tzg-write-tests-for-deco-apply-command)
- [⚫ deco-7kj Write tests for deco history command](#deco-7kj-write-tests-for-deco-history-command)
- [⚫ deco-747 Write tests for deco mv command](#deco-747-write-tests-for-deco-mv-command)
- [⚫ deco-3vb Write tests for deco unset command](#deco-3vb-write-tests-for-deco-unset-command)
- [⚫ deco-zix Write tests for deco append command](#deco-zix-write-tests-for-deco-append-command)
- [⚫ deco-bhb Write tests for node rename](#deco-bhb-write-tests-for-node-rename)
- [⚫ deco-0jp Write tests for deco set command](#deco-0jp-write-tests-for-deco-set-command)
- [⚫ deco-nop Write tests for QueryEngine search](#deco-nop-write-tests-for-queryengine-search)
- [⚫ deco-vz6 Write tests for deco query command](#deco-vz6-write-tests-for-deco-query-command)
- [⚫ deco-5wv Write tests for error aggregation](#deco-5wv-write-tests-for-error-aggregation)
- [⚫ deco-3jm Write tests for QueryEngine filter](#deco-3jm-write-tests-for-queryengine-filter)
- [⚫ deco-4dj Write tests for suggestion engine](#deco-4dj-write-tests-for-suggestion-engine)
- [⚫ deco-s5c Write tests for deco show command](#deco-s5c-write-tests-for-deco-show-command)
- [⚫ deco-h9q Write tests for Patcher apply operation](#deco-h9q-write-tests-for-patcher-apply-operation)
- [⚫ deco-mvn Write tests for YAML context extraction](#deco-mvn-write-tests-for-yaml-context-extraction)
- [⚫ deco-wlu Write tests for deco list command](#deco-wlu-write-tests-for-deco-list-command)
- [⚫ deco-q41 Write tests for Patcher unset operation](#deco-q41-write-tests-for-patcher-unset-operation)
- [⚫ deco-yu5 Write tests for error formatter](#deco-yu5-write-tests-for-error-formatter)
- [⚫ deco-3co Write tests for Patcher append operation](#deco-3co-write-tests-for-patcher-append-operation)
- [⚫ deco-sbk Write tests for deco validate command](#deco-sbk-write-tests-for-deco-validate-command)
- [⚫ deco-8sd Write tests for error code registry](#deco-8sd-write-tests-for-error-code-registry)
- [⚫ deco-aoo Write tests for Patcher set operation](#deco-aoo-write-tests-for-patcher-set-operation)
- [⚫ deco-sz1 Write tests for deco init command](#deco-sz1-write-tests-for-deco-init-command)
- [⚫ deco-16w Write tests for DecoError structure](#deco-16w-write-tests-for-decoerror-structure)
- [⚫ deco-az1 Write tests for Validator orchestrator](#deco-az1-write-tests-for-validator-orchestrator)
- [⚫ deco-336 Write tests for Cobra root command](#deco-336-write-tests-for-cobra-root-command)
- [⚫ deco-eft Write tests for constraint Validator](#deco-eft-write-tests-for-constraint-validator)
- [⚫ deco-xw8 Write tests for reference Validator](#deco-xw8-write-tests-for-reference-validator)
- [⚫ deco-2t4 Write tests for schema Validator](#deco-2t4-write-tests-for-schema-validator)
- [⚫ deco-im0 Write tests for reverse reference indexing](#deco-im0-write-tests-for-reverse-reference-indexing)
- [⚫ deco-jxy Write tests for GraphBuilder service](#deco-jxy-write-tests-for-graphbuilder-service)
- [⚫ deco-1jc Write tests for YAML line number tracking](#deco-1jc-write-tests-for-yaml-line-number-tracking)
- [⚫ deco-hxq Write tests for HistoryRepository](#deco-hxq-write-tests-for-historyrepository)
- [⚫ deco-77y Write tests for ConfigRepository](#deco-77y-write-tests-for-configrepository)
- [⚫ deco-1er Write tests for AuditEntry domain type](#deco-1er-write-tests-for-auditentry-domain-type)
- [⚫ deco-344 Write tests for file discovery](#deco-344-write-tests-for-file-discovery)
- [⚫ deco-xjv Write tests for Ref domain type](#deco-xjv-write-tests-for-ref-domain-type)
- [⚫ deco-53k Write tests for YAML NodeRepository](#deco-53k-write-tests-for-yaml-noderepository)
- [⚫ deco-ll2 Write tests for Constraint domain type](#deco-ll2-write-tests-for-constraint-domain-type)
- [⚫ deco-0hm Write tests for Issue domain type](#deco-0hm-write-tests-for-issue-domain-type)
- [⚫ deco-shb Write tests for Graph domain type](#deco-shb-write-tests-for-graph-domain-type)
- [⚫ deco-ora Write tests for Node domain type](#deco-ora-write-tests-for-node-domain-type)
- [⚫ deco-c5c Add shell completion generation](#deco-c5c-add-shell-completion-generation)
- [⚫ deco-142 Implement node rename with ref updates](#deco-142-implement-node-rename-with-ref-updates)
- [⚫ deco-0kr Implement deco apply command](#deco-0kr-implement-deco-apply-command)
- [⚫ deco-6j1 Implement deco history command](#deco-6j1-implement-deco-history-command)
- [⚫ deco-oam Implement QueryEngine search](#deco-oam-implement-queryengine-search)
- [⚫ deco-vzo Implement deco mv command](#deco-vzo-implement-deco-mv-command)
- [⚫ deco-0s1 Implement QueryEngine filter](#deco-0s1-implement-queryengine-filter)
- [⚫ deco-31k Implement deco unset command](#deco-31k-implement-deco-unset-command)
- [⚫ deco-7yp Implement Patcher apply operation](#deco-7yp-implement-patcher-apply-operation)
- [⚫ deco-inq Add error documentation generator](#deco-inq-add-error-documentation-generator)
- [⚫ deco-1v8 Implement deco append command](#deco-1v8-implement-deco-append-command)
- [⚫ deco-750 Implement Patcher unset operation](#deco-750-implement-patcher-unset-operation)
- [⚫ deco-5fv Implement error aggregation](#deco-5fv-implement-error-aggregation)
- [⚫ deco-53i Implement deco set command](#deco-53i-implement-deco-set-command)
- [⚫ deco-k3j Implement Patcher append operation](#deco-k3j-implement-patcher-append-operation)
- [⚫ deco-3ip Implement suggestion engine](#deco-3ip-implement-suggestion-engine)
- [⚫ deco-gpu Define HistoryRepository interface](#deco-gpu-define-historyrepository-interface)
- [⚫ deco-m9l Implement YAML context extraction](#deco-m9l-implement-yaml-context-extraction)
- [⚫ deco-1eg Implement deco query command](#deco-1eg-implement-deco-query-command)
- [⚫ deco-7rp Implement Patcher set operation](#deco-7rp-implement-patcher-set-operation)
- [⚫ deco-5ns Define ConfigRepository interface](#deco-5ns-define-configrepository-interface)
- [⚫ deco-4q8 Implement deco show command](#deco-4q8-implement-deco-show-command)
- [⚫ deco-e9t Implement error formatter](#deco-e9t-implement-error-formatter)
- [⚫ deco-2y3 Implement Validator orchestrator](#deco-2y3-implement-validator-orchestrator)
- [⚫ deco-p54 Define NodeRepository interface](#deco-p54-define-noderepository-interface)
- [⚫ deco-qll Implement deco list command](#deco-qll-implement-deco-list-command)
- [⚫ deco-7cp Create error code registry](#deco-7cp-create-error-code-registry)
- [⚫ deco-9ls Implement constraint Validator](#deco-9ls-implement-constraint-validator)
- [⚫ deco-8ca Define AuditEntry domain type](#deco-8ca-define-auditentry-domain-type)
- [⚫ deco-1fm Implement deco validate command](#deco-1fm-implement-deco-validate-command)
- [⚫ deco-03s Define DecoError structure](#deco-03s-define-decoerror-structure)
- [⚫ deco-jna Implement reference Validator](#deco-jna-implement-reference-validator)
- [⚫ deco-rrk Define Ref domain type](#deco-rrk-define-ref-domain-type)
- [⚫ deco-7ba Implement deco init command](#deco-7ba-implement-deco-init-command)
- [⚫ deco-bx5 Implement schema Validator](#deco-bx5-implement-schema-validator)
- [⚫ deco-nbv Add YAML line number tracking](#deco-nbv-add-yaml-line-number-tracking)
- [⚫ deco-7au Setup Cobra root command](#deco-7au-setup-cobra-root-command)
- [⚫ deco-lwx Define Constraint domain type](#deco-lwx-define-constraint-domain-type)
- [⚫ deco-owf Implement HistoryRepository](#deco-owf-implement-historyrepository)
- [⚫ deco-pqf Implement reverse reference indexing](#deco-pqf-implement-reverse-reference-indexing)
- [⚫ deco-982 Define Issue domain type](#deco-982-define-issue-domain-type)
- [⚫ deco-0o3 Implement GraphBuilder service](#deco-0o3-implement-graphbuilder-service)
- [⚫ deco-px9 Implement ConfigRepository](#deco-px9-implement-configrepository)
- [⚫ deco-zty Define Graph domain type](#deco-zty-define-graph-domain-type)
- [⚫ deco-5ry Implement file discovery for nodes](#deco-5ry-implement-file-discovery-for-nodes)
- [⚫ deco-bqg Define Node domain type](#deco-bqg-define-node-domain-type)
- [⚫ deco-nkw Implement YAML NodeRepository](#deco-nkw-implement-yaml-noderepository)
- [⚫ deco-rba Initialize Go module and project structure](#deco-rba-initialize-go-module-and-project-structure)
- [⚫ deco-4bhp Investigate streaming/incremental validation](#deco-4bhp-investigate-streaming-incremental-validation)
- [⚫ deco-1se6 Add API specification example project](#deco-1se6-add-api-specification-example-project)
- [⚫ deco-oo2q Docker dev environment setup](#deco-oo2q-docker-dev-environment-setup)
- [⚫ deco-uvas Hash truncation to 64 bits risks collision](#deco-uvas-hash-truncation-to-64-bits-risks-collision)
- [⚫ deco-7jb Expand test coverage for validator/patcher/CLI](#deco-7jb-expand-test-coverage-for-validator-patcher-cli)
- [⚫ deco-04f Validation performance on large projects](#deco-04f-validation-performance-on-large-projects)
- [⚫ deco-xdp Validation performance on large projects](#deco-xdp-validation-performance-on-large-projects)
- [⚫ deco-ai9 Deterministic load/save ordering](#deco-ai9-deterministic-load-save-ordering)
- [⚫ deco-3ix Docs alignment: README/SPEC match schema & CLI behavior](#deco-3ix-docs-alignment-readme-spec-match-schema-cli-behavior)
- [⚫ deco-23m Saved queries / query presets](#deco-23m-saved-queries-query-presets)
- [⚫ deco-rmv Detect and warn on concurrent edit conflicts](#deco-rmv-detect-and-warn-on-concurrent-edit-conflicts)
- [⚫ deco-acy validate output lacks file/line context despite YAML location tracker](#deco-acy-validate-output-lacks-file-line-context-despite-yaml-location-tracker)
- [⚫ deco-1pi Add deco stats command for project overview](#deco-1pi-add-deco-stats-command-for-project-overview)
- [⚫ deco-t2h Add deco diff command to show node history](#deco-t2h-add-deco-diff-command-to-show-node-history)
- [⚫ deco-0ql Add deco export command for multiple output formats](#deco-0ql-add-deco-export-command-for-multiple-output-formats)
- [⚫ deco-5uc Add deco graph command to visualize dependencies](#deco-5uc-add-deco-graph-command-to-visualize-dependencies)
- [⚫ deco-gp6 Add deco issues command to list all TBDs](#deco-gp6-add-deco-issues-command-to-list-all-tbds)
- [⚫ deco-t5h Add deco rm command to delete nodes](#deco-t5h-add-deco-rm-command-to-delete-nodes)
- [⚫ deco-513 Add deco create command for scaffolding nodes](#deco-513-add-deco-create-command-for-scaffolding-nodes)
- [⚫ deco-h2c Add contract validation to deco validate command](#deco-h2c-add-contract-validation-to-deco-validate-command)
- [⚫ deco-wxh Validate contract scenarios reference valid nodes](#deco-wxh-validate-contract-scenarios-reference-valid-nodes)
- [⚫ deco-4ci Implement contract scenario parser](#deco-4ci-implement-contract-scenario-parser)
- [⚫ deco-0sk Add contract syntax validator](#deco-0sk-add-contract-syntax-validator)

---

## Dependency Graph

```mermaid
graph TD
    classDef open fill:#50FA7B,stroke:#333,color:#000
    classDef inprogress fill:#8BE9FD,stroke:#333,color:#000
    classDef blocked fill:#FF5555,stroke:#333,color:#000
    classDef closed fill:#6272A4,stroke:#333,color:#fff

    deco-02tq["deco-02tq<br/>Add propose/review mode for AI patches"]
    class deco-02tq open
    deco-03s["deco-03s<br/>Define DecoError structure"]
    class deco-03s closed
    deco-04f["deco-04f<br/>Validation performance on large projects"]
    class deco-04f closed
    deco-0hm["deco-0hm<br/>Write tests for Issue domain type"]
    class deco-0hm closed
    deco-0ja["deco-0ja<br/>deco set cannot update non-string fields"]
    class deco-0ja closed
    deco-0jp["deco-0jp<br/>Write tests for deco set command"]
    class deco-0jp closed
    deco-0kr["deco-0kr<br/>Implement deco apply command"]
    class deco-0kr closed
    deco-0o3["deco-0o3<br/>Implement GraphBuilder service"]
    class deco-0o3 closed
    deco-0ql["deco-0ql<br/>Add deco export command for multiple ..."]
    class deco-0ql closed
    deco-0s1["deco-0s1<br/>Implement QueryEngine filter"]
    class deco-0s1 closed
    deco-0sk["deco-0sk<br/>Add contract syntax validator"]
    class deco-0sk closed
    deco-0t5["deco-0t5<br/>CLI mutations don't record content_ha..."]
    class deco-0t5 closed
    deco-122q["deco-122q<br/>Ref validation ignores emits_events a..."]
    class deco-122q closed
    deco-142["deco-142<br/>Implement node rename with ref updates"]
    class deco-142 closed
    deco-16e["deco-16e<br/>Foundation: Project setup and domain ..."]
    class deco-16e closed
    deco-16w["deco-16w<br/>Write tests for DecoError structure"]
    class deco-16w closed
    deco-1eg["deco-1eg<br/>Implement deco query command"]
    class deco-1eg closed
    deco-1er["deco-1er<br/>Write tests for AuditEntry domain type"]
    class deco-1er closed
    deco-1fm["deco-1fm<br/>Implement deco validate command"]
    class deco-1fm closed
    deco-1iuf["deco-1iuf<br/>Session handover: strict block valida..."]
    class deco-1iuf open
    deco-1jc["deco-1jc<br/>Write tests for YAML line number trac..."]
    class deco-1jc closed
    deco-1pi["deco-1pi<br/>Add deco stats command for project ov..."]
    class deco-1pi closed
    deco-1se6["deco-1se6<br/>Add API specification example project"]
    class deco-1se6 closed
    deco-1v8["deco-1v8<br/>Implement deco append command"]
    class deco-1v8 closed
    deco-23bx["deco-23bx<br/>Clean invalid field from snake example"]
    class deco-23bx open
    deco-23m["deco-23m<br/>Saved queries / query presets"]
    class deco-23m closed
    deco-2g91["deco-2g91<br/>Add reference/schema discovery commands"]
    class deco-2g91 open
    deco-2t4["deco-2t4<br/>Write tests for schema Validator"]
    class deco-2t4 closed
    deco-2y3["deco-2y3<br/>Implement Validator orchestrator"]
    class deco-2y3 closed
    deco-31k["deco-31k<br/>Implement deco unset command"]
    class deco-31k closed
    deco-336["deco-336<br/>Write tests for Cobra root command"]
    class deco-336 closed
    deco-344["deco-344<br/>Write tests for file discovery"]
    class deco-344 closed
    deco-36rv["deco-36rv<br/>CEL constraint engine missing spec'd ..."]
    class deco-36rv closed
    deco-3co["deco-3co<br/>Write tests for Patcher append operation"]
    class deco-3co closed
    deco-3dy["deco-3dy<br/>Set up CI/CD pipeline"]
    class deco-3dy closed
    deco-3ip["deco-3ip<br/>Implement suggestion engine"]
    class deco-3ip closed
    deco-3ix["deco-3ix<br/>Docs alignment: README/SPEC match sch..."]
    class deco-3ix closed
    deco-3jm["deco-3jm<br/>Write tests for QueryEngine filter"]
    class deco-3jm closed
    deco-3o2h["deco-3o2h<br/>apply/rewrite don't write content has..."]
    class deco-3o2h closed
    deco-3vb["deco-3vb<br/>Write tests for deco unset command"]
    class deco-3vb closed
    deco-4bhp["deco-4bhp<br/>Investigate streaming/incremental val..."]
    class deco-4bhp closed
    deco-4ci["deco-4ci<br/>Implement contract scenario parser"]
    class deco-4ci closed
    deco-4dj["deco-4dj<br/>Write tests for suggestion engine"]
    class deco-4dj closed
    deco-4kg["deco-4kg<br/>Session handover: graph command added"]
    class deco-4kg closed
    deco-4p9["deco-4p9<br/>GDD playbook documentation"]
    class deco-4p9 open
    deco-4q8["deco-4q8<br/>Implement deco show command"]
    class deco-4q8 closed
    deco-4y9["deco-4y9<br/>Session handover: Block validation co..."]
    class deco-4y9 closed
    deco-513["deco-513<br/>Add deco create command for scaffoldi..."]
    class deco-513 closed
    deco-53i["deco-53i<br/>Implement deco set command"]
    class deco-53i closed
    deco-53k["deco-53k<br/>Write tests for YAML NodeRepository"]
    class deco-53k closed
    deco-5fv["deco-5fv<br/>Implement error aggregation"]
    class deco-5fv closed
    deco-5gxq["deco-5gxq<br/>Session handover: sync command fixes"]
    class deco-5gxq closed
    deco-5ns["deco-5ns<br/>Define ConfigRepository interface"]
    class deco-5ns closed
    deco-5ry["deco-5ry<br/>Implement file discovery for nodes"]
    class deco-5ry closed
    deco-5uc["deco-5uc<br/>Add deco graph command to visualize d..."]
    class deco-5uc closed
    deco-5wv["deco-5wv<br/>Write tests for error aggregation"]
    class deco-5wv closed
    deco-5xa["deco-5xa<br/>Errors: Rust-like error system"]
    class deco-5xa closed
    deco-603["deco-603<br/>Add type-specific validation for bloc..."]
    class deco-603 closed
    deco-6iv["deco-6iv<br/>Deco GDD: Project-lead vision for end..."]
    class deco-6iv open
    deco-6j1["deco-6j1<br/>Implement deco history command"]
    class deco-6j1 closed
    deco-747["deco-747<br/>Write tests for deco mv command"]
    class deco-747 closed
    deco-750["deco-750<br/>Implement Patcher unset operation"]
    class deco-750 closed
    deco-77y["deco-77y<br/>Write tests for ConfigRepository"]
    class deco-77y closed
    deco-79k["deco-79k<br/>Enforce constraint scope (node kind/p..."]
    class deco-79k closed
    deco-7au["deco-7au<br/>Setup Cobra root command"]
    class deco-7au closed
    deco-7ba["deco-7ba<br/>Implement deco init command"]
    class deco-7ba closed
    deco-7cp["deco-7cp<br/>Create error code registry"]
    class deco-7cp closed
    deco-7de["deco-7de<br/>Session handover: Error system and st..."]
    class deco-7de closed
    deco-7jb["deco-7jb<br/>Expand test coverage for validator/pa..."]
    class deco-7jb closed
    deco-7kj["deco-7kj<br/>Write tests for deco history command"]
    class deco-7kj closed
    deco-7rp["deco-7rp<br/>Implement Patcher set operation"]
    class deco-7rp closed
    deco-7yo["deco-7yo<br/>CLI: Command implementations"]
    class deco-7yo closed
    deco-7yp["deco-7yp<br/>Implement Patcher apply operation"]
    class deco-7yp closed
    deco-7z6u["deco-7z6u<br/>YAML error location infrastructure ex..."]
    class deco-7z6u closed
    deco-86e["deco-86e<br/>Block data lost when parsing YAML (bl..."]
    class deco-86e closed
    deco-87y["deco-87y<br/>Custom block types with validation hooks"]
    class deco-87y closed
    deco-8ca["deco-8ca<br/>Define AuditEntry domain type"]
    class deco-8ca closed
    deco-8ov["deco-8ov<br/>validate does not detect duplicate no..."]
    class deco-8ov closed
    deco-8rt["deco-8rt<br/>Session handover: CLI mutation comman..."]
    class deco-8rt closed
    deco-8sd["deco-8sd<br/>Write tests for error code registry"]
    class deco-8sd closed
    deco-8vg["deco-8vg<br/>Session handover: Storage and error s..."]
    class deco-8vg closed
    deco-8xo["deco-8xo<br/>Content hash excludes Glossary, Contr..."]
    class deco-8xo closed
    deco-934["deco-934<br/>Session handover: CLI foundation comp..."]
    class deco-934 closed
    deco-982["deco-982<br/>Define Issue domain type"]
    class deco-982 closed
    deco-9ls["deco-9ls<br/>Implement constraint Validator"]
    class deco-9ls closed
    deco-9rm["deco-9rm<br/>Compile GDD to LaTeX output"]
    class deco-9rm open
    deco-a3o["deco-a3o<br/>validate ignores unknown top-level ke..."]
    class deco-a3o closed
    deco-a5l["deco-a5l<br/>AI patch/rewrite safety: validate gat..."]
    class deco-a5l closed
    deco-a5n["deco-a5n<br/>Explain error / suggest fix output"]
    class deco-a5n open
    deco-acy["deco-acy<br/>validate output lacks file/line conte..."]
    class deco-acy closed
    deco-ai9["deco-ai9<br/>Deterministic load/save ordering"]
    class deco-ai9 closed
    deco-aoo["deco-aoo<br/>Write tests for Patcher set operation"]
    class deco-aoo closed
    deco-apc["deco-apc<br/>Configurable schema rules (org-level ..."]
    class deco-apc closed
    deco-awz["deco-awz<br/>Session handover: Review workflow imp..."]
    class deco-awz closed
    deco-az1["deco-az1<br/>Write tests for Validator orchestrator"]
    class deco-az1 closed
    deco-bhb["deco-bhb<br/>Write tests for node rename"]
    class deco-bhb closed
    deco-bqg["deco-bqg<br/>Define Node domain type"]
    class deco-bqg closed
    deco-bx5["deco-bx5<br/>Implement schema Validator"]
    class deco-bx5 closed
    deco-c4y["deco-c4y<br/>Session handover: File location in va..."]
    class deco-c4y closed
    deco-c5c["deco-c5c<br/>Add shell completion generation"]
    class deco-c5c closed
    deco-c5k["deco-c5k<br/>Session handover: Foundation complete..."]
    class deco-c5k closed
    deco-chwx["deco-chwx<br/>Status validation only checks presenc..."]
    class deco-chwx closed
    deco-cxk["deco-cxk<br/>Decide and implement block schema (in..."]
    class deco-cxk closed
    deco-cxl["deco-cxl<br/>sync O(nodes × history) performance -..."]
    class deco-cxl closed
    deco-dag0["deco-dag0<br/>Session handover: performance optimiz..."]
    class deco-dag0 closed
    deco-duf["deco-duf<br/>Session handover: mv command and CLI ..."]
    class deco-duf closed
    deco-dyi["deco-dyi<br/>Session handover: 3 CLI commands added"]
    class deco-dyi closed
    deco-e4g["deco-e4g<br/>Expand issues/TBD system: filters, se..."]
    class deco-e4g closed
    deco-e9t["deco-e9t<br/>Implement error formatter"]
    class deco-e9t closed
    deco-eft["deco-eft<br/>Write tests for constraint Validator"]
    class deco-eft closed
    deco-elv["deco-elv<br/>Storage: Repository implementations"]
    class deco-elv closed
    deco-epm4["deco-epm4<br/>README block examples use key/value b..."]
    class deco-epm4 closed
    deco-fgs5["deco-fgs5<br/>refactor: extend mv to update all ref..."]
    class deco-fgs5 closed
    deco-flq["deco-flq<br/>Session handover: CLI commands (valid..."]
    class deco-flq closed
    deco-fre["deco-fre<br/>Session handover: stats command added"]
    class deco-fre closed
    deco-g50["deco-g50<br/>Session handover: contract parser add..."]
    class deco-g50 closed
    deco-gbk["deco-gbk<br/>Session handover: audit history added"]
    class deco-gbk closed
    deco-gp6["deco-gp6<br/>Add deco issues command to list all TBDs"]
    class deco-gp6 closed
    deco-gpli["deco-gpli<br/>Allow explicit actor identity in audi..."]
    class deco-gpli open
    deco-gpu["deco-gpu<br/>Define HistoryRepository interface"]
    class deco-gpu closed
    deco-gw6["deco-gw6<br/>Session handover: Node rename service..."]
    class deco-gw6 closed
    deco-h2c["deco-h2c<br/>Add contract validation to deco valid..."]
    class deco-h2c closed
    deco-h91g["deco-h91g<br/>sync swallows errors - can exit clean..."]
    class deco-h91g closed
    deco-h9q["deco-h9q<br/>Write tests for Patcher apply operation"]
    class deco-h9q closed
    deco-hfu4["deco-hfu4<br/>Add optimistic concurrency controls f..."]
    class deco-hfu4 open
    deco-hond["deco-hond<br/>set resets status to draft but append..."]
    class deco-hond closed
    deco-hsx["deco-hsx<br/>Require content for approved/publishe..."]
    class deco-hsx closed
    deco-hxq["deco-hxq<br/>Write tests for HistoryRepository"]
    class deco-hxq closed
    deco-im0["deco-im0<br/>Write tests for reverse reference ind..."]
    class deco-im0 closed
    deco-inq["deco-inq<br/>Add error documentation generator"]
    class deco-inq closed
    deco-jfye["deco-jfye<br/>Validate node IDs to prevent path tra..."]
    class deco-jfye open
    deco-jna["deco-jna<br/>Implement reference Validator"]
    class deco-jna closed
    deco-jxy["deco-jxy<br/>Write tests for GraphBuilder service"]
    class deco-jxy closed
    deco-jyn["deco-jyn<br/>deco set panics on nested paths (poin..."]
    class deco-jyn closed
    deco-k3j["deco-k3j<br/>Implement Patcher append operation"]
    class deco-k3j closed
    deco-kz2i["deco-kz2i<br/>Add LLM context export command"]
    class deco-kz2i open
    deco-ll2["deco-ll2<br/>Write tests for Constraint domain type"]
    class deco-ll2 closed
    deco-lwl["deco-lwl<br/>Session handover: contract validation..."]
    class deco-lwl closed
    deco-lwx["deco-lwx<br/>Define Constraint domain type"]
    class deco-lwx closed
    deco-m9l["deco-m9l<br/>Implement YAML context extraction"]
    class deco-m9l closed
    deco-mbz["deco-mbz<br/>Add deco sync command for detecting a..."]
    class deco-mbz closed
    deco-mvn["deco-mvn<br/>Write tests for YAML context extraction"]
    class deco-mvn closed
    deco-nbv["deco-nbv<br/>Add YAML line number tracking"]
    class deco-nbv closed
    deco-nd7w["deco-nd7w<br/>Config paths (nodes_path, history_pat..."]
    class deco-nd7w closed
    deco-nkw["deco-nkw<br/>Implement YAML NodeRepository"]
    class deco-nkw closed
    deco-nop["deco-nop<br/>Write tests for QueryEngine search"]
    class deco-nop closed
    deco-oam["deco-oam<br/>Implement QueryEngine search"]
    class deco-oam closed
    deco-oh0q["deco-oh0q<br/>--dry-run always exits 0 even when ch..."]
    class deco-oh0q closed
    deco-oo2q["deco-oo2q<br/>Docker dev environment setup"]
    class deco-oo2q closed
    deco-ora["deco-ora<br/>Write tests for Node domain type"]
    class deco-ora closed
    deco-owf["deco-owf<br/>Implement HistoryRepository"]
    class deco-owf closed
    deco-oxu2["deco-oxu2<br/>Add tests for Resolve*Path with relat..."]
    class deco-oxu2 open
    deco-p54["deco-p54<br/>Define NodeRepository interface"]
    class deco-p54 closed
    deco-p75["deco-p75<br/>Backwards-compatible schema migrations"]
    class deco-p75 closed
    deco-pc3o["deco-pc3o<br/>Content hash uses non-deterministic m..."]
    class deco-pc3o closed
    deco-pqf["deco-pqf<br/>Implement reverse reference indexing"]
    class deco-pqf closed
    deco-px9["deco-px9<br/>Implement ConfigRepository"]
    class deco-px9 closed
    deco-q41["deco-q41<br/>Write tests for Patcher unset operation"]
    class deco-q41 closed
    deco-q5pk["deco-q5pk<br/>Validate unknown fields in RefLink an..."]
    class deco-q5pk closed
    deco-qll["deco-qll<br/>Implement deco list command"]
    class deco-qll closed
    deco-qw2b["deco-qw2b<br/>Validate examples/snake with deco val..."]
    class deco-qw2b closed
    deco-rb08["deco-rb08<br/>sync: detect manual file renames and ..."]
    class deco-rb08 closed
    deco-rba["deco-rba<br/>Initialize Go module and project stru..."]
    class deco-rba closed
    deco-rib["deco-rib<br/>Define strict top-level schema with e..."]
    class deco-rib closed
    deco-rmv["deco-rmv<br/>Detect and warn on concurrent edit co..."]
    class deco-rmv closed
    deco-rrk["deco-rrk<br/>Define Ref domain type"]
    class deco-rrk closed
    deco-rs0["deco-rs0<br/>Review workflow: approvals, status tr..."]
    class deco-rs0 closed
    deco-ry9["deco-ry9<br/>Session handover: Review workflow des..."]
    class deco-ry9 closed
    deco-s5c["deco-s5c<br/>Write tests for deco show command"]
    class deco-s5c closed
    deco-sbk["deco-sbk<br/>Write tests for deco validate command"]
    class deco-sbk closed
    deco-shb["deco-shb<br/>Write tests for Graph domain type"]
    class deco-shb closed
    deco-snlr["deco-snlr<br/>Add strict block field validation (co..."]
    class deco-snlr closed
    deco-sz1["deco-sz1<br/>Write tests for deco init command"]
    class deco-sz1 closed
    deco-t2h["deco-t2h<br/>Add deco diff command to show node hi..."]
    class deco-t2h closed
    deco-t5h["deco-t5h<br/>Add deco rm command to delete nodes"]
    class deco-t5h closed
    deco-t6q["deco-t6q<br/>Services: Business logic layer"]
    class deco-t6q closed
    deco-tw3["deco-tw3<br/>Session handover: Service layer progr..."]
    class deco-tw3 closed
    deco-tzg["deco-tzg<br/>Write tests for deco apply command"]
    class deco-tzg closed
    deco-u0jb["deco-u0jb<br/>Session handover: custom block types ..."]
    class deco-u0jb closed
    deco-uvas["deco-uvas<br/>Hash truncation to 64 bits risks coll..."]
    class deco-uvas closed
    deco-v1vi["deco-v1vi<br/>Add JSON output for validation and mu..."]
    class deco-v1vi open
    deco-v4o["deco-v4o<br/>Session handover: contract ref valida..."]
    class deco-v4o closed
    deco-vq3["deco-vq3<br/>Session handover: diff command added"]
    class deco-vq3 closed
    deco-vtn5["deco-vtn5<br/>Resolve*Path should be absolute or co..."]
    class deco-vtn5 open
    deco-vz6["deco-vz6<br/>Write tests for deco query command"]
    class deco-vz6 closed
    deco-vzo["deco-vzo<br/>Implement deco mv command"]
    class deco-vzo closed
    deco-wlu["deco-wlu<br/>Write tests for deco list command"]
    class deco-wlu closed
    deco-wxe["deco-wxe<br/>Session handover: Service layer compl..."]
    class deco-wxe closed
    deco-wxh["deco-wxh<br/>Validate contract scenarios reference..."]
    class deco-wxh closed
    deco-xdp["deco-xdp<br/>Validation performance on large projects"]
    class deco-xdp closed
    deco-xjv["deco-xjv<br/>Write tests for Ref domain type"]
    class deco-xjv closed
    deco-xqq["deco-xqq<br/>Session handover: Constraint scope en..."]
    class deco-xqq closed
    deco-xw8["deco-xw8<br/>Write tests for reference Validator"]
    class deco-xw8 closed
    deco-ykl["deco-ykl<br/>Session handover: contract validator ..."]
    class deco-ykl closed
    deco-yle3["deco-yle3<br/>Update CLI help text to reflect confi..."]
    class deco-yle3 open
    deco-yu5["deco-yu5<br/>Write tests for error formatter"]
    class deco-yu5 closed
    deco-zix["deco-zix<br/>Write tests for deco append command"]
    class deco-zix closed
    deco-zn8["deco-zn8<br/>Add audit history to all node-modifyi..."]
    class deco-zn8 closed
    deco-zty["deco-zty<br/>Define Graph domain type"]
    class deco-zty closed

    deco-03s ==> deco-16e
    deco-03s ==> deco-16w
    deco-0jp ==> deco-5xa
    deco-0jp ==> deco-t6q
    deco-0kr ==> deco-5xa
    deco-0kr ==> deco-t6q
    deco-0kr ==> deco-tzg
    deco-0o3 ==> deco-elv
    deco-0o3 ==> deco-jxy
    deco-0s1 ==> deco-3jm
    deco-0s1 ==> deco-elv
    deco-0sk ==> deco-4ci
    deco-142 ==> deco-bhb
    deco-142 ==> deco-elv
    deco-16e ==> deco-5ns
    deco-16e ==> deco-8ca
    deco-16e ==> deco-982
    deco-16e ==> deco-bqg
    deco-16e ==> deco-gpu
    deco-16e ==> deco-lwx
    deco-16e ==> deco-p54
    deco-16e ==> deco-rba
    deco-16e ==> deco-rrk
    deco-16e ==> deco-zty
    deco-16w ==> deco-16e
    deco-1eg ==> deco-5xa
    deco-1eg ==> deco-t6q
    deco-1eg ==> deco-vz6
    deco-1fm ==> deco-5xa
    deco-1fm ==> deco-sbk
    deco-1fm ==> deco-t6q
    deco-1jc ==> deco-16e
    deco-1v8 ==> deco-5xa
    deco-1v8 ==> deco-t6q
    deco-1v8 ==> deco-zix
    deco-2t4 ==> deco-elv
    deco-2y3 ==> deco-az1
    deco-2y3 ==> deco-elv
    deco-31k ==> deco-3vb
    deco-31k ==> deco-5xa
    deco-31k ==> deco-t6q
    deco-336 ==> deco-5xa
    deco-336 ==> deco-t6q
    deco-344 ==> deco-16e
    deco-3co ==> deco-elv
    deco-3ip ==> deco-16e
    deco-3ip ==> deco-4dj
    deco-3jm ==> deco-elv
    deco-3vb ==> deco-5xa
    deco-3vb ==> deco-t6q
    deco-4dj ==> deco-16e
    deco-4q8 ==> deco-5xa
    deco-4q8 ==> deco-s5c
    deco-4q8 ==> deco-t6q
    deco-53i ==> deco-0jp
    deco-53i ==> deco-5xa
    deco-53i ==> deco-t6q
    deco-53k ==> deco-16e
    deco-5fv ==> deco-16e
    deco-5fv ==> deco-5wv
    deco-5ry ==> deco-16e
    deco-5ry ==> deco-344
    deco-5wv ==> deco-16e
    deco-5xa ==> deco-03s
    deco-5xa ==> deco-16e
    deco-5xa ==> deco-3ip
    deco-5xa ==> deco-5fv
    deco-5xa ==> deco-7cp
    deco-5xa ==> deco-e9t
    deco-5xa ==> deco-inq
    deco-5xa ==> deco-m9l
    deco-603 ==> deco-cxk
    deco-6iv ==> deco-04f
    deco-6iv ==> deco-23m
    deco-6iv ==> deco-3ix
    deco-6iv ==> deco-4p9
    deco-6iv ==> deco-603
    deco-6iv ==> deco-79k
    deco-6iv ==> deco-7jb
    deco-6iv ==> deco-87y
    deco-6iv ==> deco-a5l
    deco-6iv ==> deco-a5n
    deco-6iv ==> deco-acy
    deco-6iv ==> deco-ai9
    deco-6iv ==> deco-apc
    deco-6iv ==> deco-e4g
    deco-6iv ==> deco-hsx
    deco-6iv ==> deco-p75
    deco-6iv ==> deco-rmv
    deco-6iv ==> deco-rs0
    deco-6iv ==> deco-xdp
    deco-6j1 ==> deco-5xa
    deco-6j1 ==> deco-7kj
    deco-6j1 ==> deco-t6q
    deco-747 ==> deco-5xa
    deco-747 ==> deco-t6q
    deco-750 ==> deco-elv
    deco-750 ==> deco-q41
    deco-77y ==> deco-16e
    deco-7au ==> deco-336
    deco-7au ==> deco-5xa
    deco-7au ==> deco-t6q
    deco-7ba ==> deco-5xa
    deco-7ba ==> deco-sz1
    deco-7ba ==> deco-t6q
    deco-7cp ==> deco-16e
    deco-7cp ==> deco-8sd
    deco-7kj ==> deco-5xa
    deco-7kj ==> deco-t6q
    deco-7rp ==> deco-aoo
    deco-7rp ==> deco-elv
    deco-7yo ==> deco-0kr
    deco-7yo ==> deco-16e
    deco-7yo ==> deco-1eg
    deco-7yo ==> deco-1fm
    deco-7yo ==> deco-1v8
    deco-7yo ==> deco-31k
    deco-7yo ==> deco-4q8
    deco-7yo ==> deco-53i
    deco-7yo ==> deco-5xa
    deco-7yo ==> deco-6j1
    deco-7yo ==> deco-7au
    deco-7yo ==> deco-7ba
    deco-7yo ==> deco-c5c
    deco-7yo ==> deco-qll
    deco-7yo ==> deco-t6q
    deco-7yo ==> deco-vzo
    deco-7yp ==> deco-elv
    deco-7yp ==> deco-h9q
    deco-86e ==> deco-cxk
    deco-8ca ==> deco-1er
    deco-8sd ==> deco-16e
    deco-982 ==> deco-0hm
    deco-9ls ==> deco-eft
    deco-9ls ==> deco-elv
    deco-a3o ==> deco-rib
    deco-aoo ==> deco-elv
    deco-az1 ==> deco-elv
    deco-bhb ==> deco-elv
    deco-bqg ==> deco-ora
    deco-bx5 ==> deco-2t4
    deco-bx5 ==> deco-elv
    deco-c5c ==> deco-5xa
    deco-c5c ==> deco-t6q
    deco-e9t ==> deco-16e
    deco-e9t ==> deco-yu5
    deco-eft ==> deco-elv
    deco-elv ==> deco-16e
    deco-elv ==> deco-5ry
    deco-elv ==> deco-nbv
    deco-elv ==> deco-nkw
    deco-elv ==> deco-owf
    deco-elv ==> deco-px9
    deco-h2c ==> deco-wxh
    deco-h9q ==> deco-elv
    deco-hxq ==> deco-16e
    deco-im0 ==> deco-elv
    deco-inq ==> deco-16e
    deco-jna ==> deco-elv
    deco-jna ==> deco-xw8
    deco-jxy ==> deco-elv
    deco-k3j ==> deco-3co
    deco-k3j ==> deco-elv
    deco-lwx ==> deco-ll2
    deco-m9l ==> deco-16e
    deco-m9l ==> deco-mvn
    deco-mvn ==> deco-16e
    deco-nbv ==> deco-16e
    deco-nbv ==> deco-1jc
    deco-nkw ==> deco-16e
    deco-nkw ==> deco-53k
    deco-nop ==> deco-elv
    deco-oam ==> deco-elv
    deco-oam ==> deco-nop
    deco-owf ==> deco-16e
    deco-owf ==> deco-hxq
    deco-pqf ==> deco-elv
    deco-pqf ==> deco-im0
    deco-px9 ==> deco-16e
    deco-px9 ==> deco-77y
    deco-q41 ==> deco-elv
    deco-qll ==> deco-5xa
    deco-qll ==> deco-t6q
    deco-qll ==> deco-wlu
    deco-rb08 ==> deco-fgs5
    deco-rrk ==> deco-xjv
    deco-s5c ==> deco-5xa
    deco-s5c ==> deco-t6q
    deco-sbk ==> deco-5xa
    deco-sbk ==> deco-t6q
    deco-sz1 ==> deco-5xa
    deco-sz1 ==> deco-t6q
    deco-t6q ==> deco-0o3
    deco-t6q ==> deco-0s1
    deco-t6q ==> deco-142
    deco-t6q ==> deco-16e
    deco-t6q ==> deco-2y3
    deco-t6q ==> deco-750
    deco-t6q ==> deco-7rp
    deco-t6q ==> deco-7yp
    deco-t6q ==> deco-9ls
    deco-t6q ==> deco-bx5
    deco-t6q ==> deco-elv
    deco-t6q ==> deco-jna
    deco-t6q ==> deco-k3j
    deco-t6q ==> deco-oam
    deco-t6q ==> deco-pqf
    deco-tzg ==> deco-5xa
    deco-tzg ==> deco-t6q
    deco-vz6 ==> deco-5xa
    deco-vz6 ==> deco-t6q
    deco-vzo ==> deco-5xa
    deco-vzo ==> deco-747
    deco-vzo ==> deco-t6q
    deco-wlu ==> deco-5xa
    deco-wlu ==> deco-t6q
    deco-wxh ==> deco-0sk
    deco-xdp -.-> deco-04f
    deco-xw8 ==> deco-elv
    deco-yu5 ==> deco-16e
    deco-zix ==> deco-5xa
    deco-zix ==> deco-t6q
    deco-zty ==> deco-shb
```

---

<a id="deco-1iuf-session-handover-strict-block-validation"></a>

## 📋 deco-1iuf Session handover: strict block validation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:55 |
| **Updated** | 2026-02-03 20:55 |

### Notes

Accomplished: implemented strict block field validation with E049, table column key allowlist, custom optional_fields support, updated schema hash and docs, removed invalid example field. Tests: go test ./internal/services/validator -run Block (with GOCACHE set) passed. State: working tree has code/doc changes plus beads updates; untracked prd-template.md present (pre-existing). Next steps: run full test suite if desired; finish landing the plane (bd sync, git commit, push). Decisions: used E049 for unknown block/column fields; built-in allowlists include id and optional mechanic inputs. Blockers: none.

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-1iuf -s in_progress

# Add a comment
bd comment deco-1iuf 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-1iuf -p 1

# View full details
bd show deco-1iuf
```

</details>

---

<a id="deco-jfye-validate-node-ids-to-prevent-path-traversal"></a>

## 🐛 deco-jfye Validate node IDs to prevent path traversal

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔥 Critical (P0) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:24 |
| **Updated** | 2026-02-03 20:24 |

### Description

**Summary**
Validate and sanitize node IDs to prevent path traversal and invalid filesystem writes.

**Problem / Opportunity**
Node IDs map directly to file paths. Without validation, a malicious or accidental ID like `../foo` could escape the nodes directory, which is unsafe for web usage.

**Goals**
- Reject IDs containing `..`, absolute paths, or invalid characters
- Ensure IDs only map within the nodes directory

**Non-goals**
- Renaming existing nodes
- Changing on-disk layout rules

**User Stories**
- As a backend operator, I want to ensure all node writes stay inside `.deco/nodes` to prevent security issues.

**Scope**
- In scope: validation in create/load/mutate paths, tests
- Out of scope: migration of existing bad IDs (handle explicitly if found)

**Constraints**
- Must be cross-platform (Windows + Unix)

**Proposed Approach**
1. Add ID validator (regex + path clean checks)
2. Enforce validation in CLI commands and repository save
3. Add explicit error code/message for invalid IDs

**TDD Plan (Required)**
- Tests to write first (red): reject `../x`, absolute paths, empty segments
- Minimal implementation (green): validation helper
- Refactor: reuse validator in all mutation commands

**Acceptance Criteria**
- [ ] Invalid IDs are rejected with clear error
- [ ] Valid IDs continue to work
- [ ] No writes outside `.deco/nodes`

**Test Plan**
- `go test ./internal/cli -run CreateInvalidID`
- `go test ./internal/storage/node -run PathForNode`

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli/create.go
- internal/storage/node/yaml_repository.go
- internal/domain (validator helper)

**Risks**
- Existing projects may contain invalid IDs

**Open Questions**
- Should invalid IDs be auto-normalized or strictly rejected?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-jfye -s in_progress

# Add a comment
bd comment deco-jfye 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-jfye -p 1

# View full details
bd show deco-jfye
```

</details>

---

<a id="deco-hfu4-add-optimistic-concurrency-controls-for-writes"></a>

## ✨ deco-hfu4 Add optimistic concurrency controls for writes

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ⚡ High (P1) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:24 |
| **Updated** | 2026-02-03 20:24 |

### Description

**Summary**
Add optimistic concurrency controls and file locking for write operations to prevent conflicting edits in multi-user workflows.

**Problem / Opportunity**
Web/UI use will introduce parallel edits. Current writes can silently clobber each other because there is no version/hash precondition or lock.

**Goals**
- Require `version` or `content_hash` preconditions for write operations
- Add optional file locking to prevent concurrent writes

**Non-goals**
- Full CRDT/merge engine
- Background sync daemon

**User Stories**
- As a UI user, I want to be warned when my draft is stale before overwriting another edit.

**Scope**
- In scope: precondition checks in CLI/services, history entries, tests
- Out of scope: auto-merge

**Constraints**
- Must preserve existing CLI behavior unless preconditions are supplied

**Proposed Approach**
1. Expose `--if-version` or `--if-hash` flags on write commands
2. Validate against current stored node before applying
3. Optional file lock (best-effort) during mutation

**TDD Plan (Required)**
- Tests to write first (red): reject stale version/hash
- Minimal implementation (green): precondition checks
- Refactor: shared mutation guard utility

**Acceptance Criteria**
- [ ] Writes fail when `--if-version` does not match current
- [ ] Writes fail when `--if-hash` does not match current content hash
- [ ] Error clearly indicates conflict

**Test Plan**
- `go test ./internal/cli -run Concurrency`

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli (set/append/unset/apply/rewrite/sync)
- internal/storage/node

**Risks**
- Locking may behave differently across OSes

**Open Questions**
- Prefer hash or version as the primary precondition?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-hfu4 -s in_progress

# Add a comment
bd comment deco-hfu4 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-hfu4 -p 1

# View full details
bd show deco-hfu4
```

</details>

---

<a id="deco-kz2i-add-llm-context-export-command"></a>

## ✨ deco-kz2i Add LLM context export command

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ⚡ High (P1) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:23 |
| **Updated** | 2026-02-03 20:23 |

### Description

**Summary**
Add a `deco export` command that emits LLM-ready context for a node, including referenced nodes, constraints, and schema hints in Markdown/JSON.

**Problem / Opportunity**
LLM workflows require stitched context; currently users must manually `deco show` and copy/paste, which is slow and error-prone.

**Goals**
- Provide a single command to export context for a node and its references
- Support formats suitable for LLM prompts (Markdown, JSON)

**Non-goals**
- Building a UI
- Implementing an LLM client

**User Stories**
- As a designer using LLMs, I want a single command to export full context, so I can generate accurate patches.

**Scope**
- In scope: CLI command, graph traversal, output formats, tests
- Out of scope: network calls, prompt templates

**Constraints**
- Must not change node data
- Must be deterministic and stable output order

**Proposed Approach**
1. Implement graph expansion (direct refs; optionally depth flag)
2. Render nodes + schema info into Markdown/JSON
3. Add CLI flags and tests

**TDD Plan (Required)**
- Tests to write first (red): export output includes target node, refs, constraints, stable ordering
- Minimal implementation (green): read nodes, traverse refs, emit output
- Refactor: extract renderer and traversal helpers

**Acceptance Criteria**
- [ ] `deco export <id>` outputs target node and direct refs
- [ ] `--format json|markdown` works
- [ ] Output order is deterministic for diffs

**Test Plan**
- `go test ./internal/cli -run Export`
- Edge cases: missing ref, cyclic refs, depth limit

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli (new command)
- internal/services/graph or new exporter
- internal/domain (output structs)

**Risks**
- Large graphs may be slow
- Output size could be huge

**Open Questions**
- Should default include only direct refs or full transitive closure?
- Include schema rules in output by default?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-kz2i -s in_progress

# Add a comment
bd comment deco-kz2i 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-kz2i -p 1

# View full details
bd show deco-kz2i
```

</details>

---

<a id="deco-02tq-add-propose-review-mode-for-ai-patches"></a>

## ✨ deco-02tq Add propose/review mode for AI patches

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:25 |
| **Updated** | 2026-02-03 20:25 |

### Description

**Summary**
Add a “suggest/propose” mode for AI-generated patches so humans can review before applying.

**Problem / Opportunity**
Current workflow is all-or-nothing. A web UI needs a safe review step for AI suggestions before writing to disk.

**Goals**
- Allow generating a diff/preview without applying
- Support approval workflow that applies the suggestion later

**Non-goals**
- Full multi-user review system
- Integration with external review tools

**User Stories**
- As a lead, I want AI to propose changes that I can approve before they are written.

**Scope**
- In scope: CLI support for proposal artifacts, diff output
- Out of scope: web UI implementation

**Constraints**
- Must not mutate nodes when in propose mode

**Proposed Approach**
1. Add `deco propose` or `deco apply --propose` to write a proposal artifact
2. Add `deco propose list/show/apply` for approval
3. Store proposals in `.deco/proposals/`

**TDD Plan (Required)**
- Tests to write first (red): proposal file created, no node changes
- Minimal implementation (green): write proposal + apply command
- Refactor: shared diff rendering

**Acceptance Criteria**
- [ ] Proposal mode never mutates nodes
- [ ] Proposal artifact can be applied later
- [ ] Proposal includes before/after snapshot

**Test Plan**
- `go test ./internal/cli -run Propose`

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli/apply.go
- internal/storage (proposal persistence)

**Risks**
- Proposal format design may change

**Open Questions**
- Should proposals be stored as patches or full snapshots?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-02tq -s in_progress

# Add a comment
bd comment deco-02tq 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-02tq -p 1

# View full details
bd show deco-02tq
```

</details>

---

<a id="deco-gpli-allow-explicit-actor-identity-in-audit-history"></a>

## ✨ deco-gpli Allow explicit actor identity in audit history

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:25 |
| **Updated** | 2026-02-03 20:25 |

### Description

**Summary**
Allow explicit actor identity to be passed into history entries instead of relying on OS username.

**Problem / Opportunity**
The audit trail currently uses `os/user` which is not suitable for web or multi-user contexts. We need explicit actors.

**Goals**
- Add `--actor` flag or env var that flows into history entries
- Preserve current behavior when not set

**Non-goals**
- Authentication/authorization
- Changing history format

**User Stories**
- As a web UI, I want to record the actual user who made a change.

**Scope**
- In scope: CLI flag plumbing, audit entry updates, tests
- Out of scope: auth system

**Constraints**
- Backward compatible history format

**Proposed Approach**
1. Add actor flag/env (e.g., `DECO_ACTOR`)
2. Pass actor into audit logging helpers
3. Add tests verifying actor override

**TDD Plan (Required)**
- Tests to write first (red): actor overrides OS username
- Minimal implementation (green): flag + plumbing
- Refactor: centralize actor retrieval

**Acceptance Criteria**
- [ ] Actor field uses provided value when set
- [ ] Defaults remain unchanged otherwise

**Test Plan**
- `go test ./internal/cli -run Actor`

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli/audit.go
- internal/cli/* mutation commands

**Risks**
- Inconsistent actor values if not standardized

**Open Questions**
- Flag name vs env name?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-gpli -s in_progress

# Add a comment
bd comment deco-gpli 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-gpli -p 1

# View full details
bd show deco-gpli
```

</details>

---

<a id="deco-v1vi-add-json-output-for-validation-and-mutation-errors"></a>

## ✨ deco-v1vi Add JSON output for validation and mutation errors

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:24 |
| **Updated** | 2026-02-03 20:24 |

### Description

**Summary**
Add `--json` output for `deco validate` and mutation commands so UIs/LLMs can consume structured error details.

**Problem / Opportunity**
Current validation errors are human-readable text only, which makes automated fixing or UI display hard.

**Goals**
- Provide structured JSON output for validation errors
- Include code, summary, detail, location, suggestion, related

**Non-goals**
- Changing existing text format
- Adding a REST API

**User Stories**
- As a UI integrator, I want JSON error payloads so I can display and resolve issues programmatically.

**Scope**
- In scope: `deco validate --json`, `deco apply --json`, `deco rewrite --json`
- Out of scope: streaming validation

**Constraints**
- Must preserve existing CLI output by default
- JSON schema should be stable

**Proposed Approach**
1. Add JSON renderer for `domain.DecoError`
2. Wire `--json` flag into validate/apply/rewrite
3. Add tests for JSON structure

**TDD Plan (Required)**
- Tests to write first (red): JSON output contains expected fields
- Minimal implementation (green): serializer + flag
- Refactor: reuse formatter utilities

**Acceptance Criteria**
- [ ] `deco validate --json` emits machine-readable error list
- [ ] JSON includes `code`, `summary`, `detail`, `location`, `suggestion`
- [ ] Exit codes unchanged

**Test Plan**
- `go test ./internal/cli -run ValidateJSON`
- Edge cases: no errors, multiple errors

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli/validate.go
- internal/domain/error_formatter.go or new JSON formatter

**Risks**
- Breaking tooling if output changes by default

**Open Questions**
- Should JSON include source snippets/lines?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-v1vi -s in_progress

# Add a comment
bd comment deco-v1vi 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-v1vi -p 1

# View full details
bd show deco-v1vi
```

</details>

---

<a id="deco-2g91-add-reference-schema-discovery-commands"></a>

## ✨ deco-2g91 Add reference/schema discovery commands

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:23 |
| **Updated** | 2026-02-03 20:23 |

### Description

**Summary**
Add CLI commands to list referenceable node IDs and schema rules (kinds, block types, required fields) for discoverability.

**Problem / Opportunity**
LLMs and users need to know “what can I reference?” and “what fields are valid?” Without this, edits hallucinate IDs or schema fields.

**Goals**
- Provide `deco refs` to list node IDs (filterable by kind)
- Provide `deco schema` to show required fields and custom block types

**Non-goals**
- Building UI components
- Auto-generating prompts

**User Stories**
- As a designer, I want a quick list of valid references so I can author without guessing.

**Scope**
- In scope: CLI output (text/JSON), filters, tests
- Out of scope: persistent indexes beyond current load

**Constraints**
- Must respect config schema rules and custom block types
- Deterministic output order

**Proposed Approach**
1. Implement refs listing via node repository + filters
2. Implement schema output from config + built-in block types
3. Add JSON flag for UI/LLM usage

**TDD Plan (Required)**
- Tests to write first (red): filters by kind, JSON output shape
- Minimal implementation (green): list nodes/config
- Refactor: shared renderer utilities

**Acceptance Criteria**
- [ ] `deco refs --kind item` lists only item IDs
- [ ] `deco schema` prints required fields and block types
- [ ] `--json` outputs machine-readable data

**Test Plan**
- `go test ./internal/cli -run Refs`
- `go test ./internal/cli -run Schema`

**Dependencies**
- None

**Files / Areas Touched**
- internal/cli (new commands)
- internal/storage/config

**Risks**
- Large projects may produce large outputs

**Open Questions**
- Should `deco schema` include built-in block definitions with examples?

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-2g91 -s in_progress

# Add a comment
bd comment deco-2g91 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-2g91 -p 1

# View full details
bd show deco-2g91
```

</details>

---

<a id="deco-6iv-deco-gdd-project-lead-vision-for-end-state-product"></a>

## 🚀 deco-6iv Deco GDD: Project-lead vision for end-state product

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | 🔹 Medium (P2) |
| **Status** | 🟢 open |
| **Created** | 2026-02-01 17:39 |
| **Updated** | 2026-02-02 21:20 |

### Description

Vision
As a project lead, I want Deco to be the authoritative, auditable, low‑friction system for creating and maintaining game design documents (GDDs). It should be strict enough to prevent ambiguity and drift, yet ergonomic for daily authoring, collaboration, and AI‑assisted workflows.

Primary Outcomes
- Single source of truth for design, with schema‑level guarantees and traceable changes.
- Fast authoring and refactoring without fear of silent data loss.
- Validation that catches structural/design errors early and points directly to the fix.
- Clear workflows for drafts → review → approved/published.

Core Authoring Features (Missing / Incomplete)
- Canonical block schema (inline fields vs data map) with zero data loss on load/save.
- Strict top‑level schema validation with explicit extension namespace (`custom` or `x_*`).
- Type‑specific validation for blocks (rule/table/param/mechanic/list, etc.).
- Status‑aware validation (draft relaxed, approved/published strict).
- Robust duplicate ID detection + helpful remediation.
- Rich error context (file:line:col + snippet) in validate output.
- Safe patching for nested paths (no panics), including pointer fields and map/slice updates.
- Typed `set/append/unset` with JSON/YAML parsing for non‑string fields.

Project Management & Collaboration
- Issue/TBD system: list, filter, severity, and per‑node issue tracking.
- Review workflow: approvals, status transitions, and changelog notes.
- Audit trail with who/what/when, and readable history diff per node.
- Conflict detection for concurrent edits.

Navigation & Discovery
- Strong query language (kind/status/tags/refs/text) with saved queries.
- Graph visualization of dependencies and reverse references.
- Stats/health report (coverage, missing content, orphan nodes).

Schema & Extensibility
- Configurable schema rules (org‑level constraints).
- Custom block types with validation hooks.
- CEL constraints scoped by node kinds and patterns.

AI Workflow Support
- Patch/rewrite modes are safe and transactional.
- Validation gate before write; explicit diffs for review.
- “Explain error” and “suggest fix” outputs for AI/humans.

Docs & Examples
- README/SPEC aligned with actual schema and CLI behaviors.
- Example projects that validate and demonstrate best practices.
- “How to structure your GDD” playbook.

Quality & Non‑Functional
- Deterministic load/save (no reordering surprises).
- Validation performance on large projects.
- Backwards‑compatible migrations for schema updates.
- Robust tests for validator, patcher, and CLI.

Acceptance Criteria (for epic completion)
- All blocking issues resolved (schema mismatch, validation gaps, CLI crashes).
- Strict validation with clear error locations by default.
- Authoring flows validated end‑to‑end (create/edit/validate/history).
- Example projects round‑trip without data loss.
- Docs reflect behavior and are consistent with tests.

### Notes

CLARIFICATIONS (from review):

1. **Acceptance Criteria Specificity**: 'All blocking issues resolved' means: all issues in DEPENDS ON list with status OPEN or IN_PROGRESS must be CLOSED before epic can close. No implicit issues - only named dependencies count.

2. **Critical Path**: deco-p75 (migrations) and deco-a5l (AI safety) are the blocking open issues. Sequence: p75 first (schema foundation), then a5l (depends on stable schema), then remaining P3 items can parallelize.

### Dependencies

- ⛔ **blocks**: `deco-79k`
- ⛔ **blocks**: `deco-87y`
- ⛔ **blocks**: `deco-a5l`
- ⛔ **blocks**: `deco-apc`
- ⛔ **blocks**: `deco-e4g`
- ⛔ **blocks**: `deco-hsx`
- ⛔ **blocks**: `deco-p75`
- ⛔ **blocks**: `deco-rs0`
- ⛔ **blocks**: `deco-04f`
- ⛔ **blocks**: `deco-23m`
- ⛔ **blocks**: `deco-3ix`
- ⛔ **blocks**: `deco-4p9`
- ⛔ **blocks**: `deco-7jb`
- ⛔ **blocks**: `deco-a5n`
- ⛔ **blocks**: `deco-acy`
- ⛔ **blocks**: `deco-ai9`
- ⛔ **blocks**: `deco-rmv`
- ⛔ **blocks**: `deco-xdp`
- ⛔ **blocks**: `deco-603`

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-6iv -s in_progress

# Add a comment
bd comment deco-6iv 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-6iv -p 1

# View full details
bd show deco-6iv
```

</details>

---

<a id="deco-23bx-clean-invalid-field-from-snake-example"></a>

## 🐛 deco-23bx Clean invalid field from snake example

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-03 20:25 |
| **Updated** | 2026-02-03 20:25 |

### Description

**Summary**
Remove the junk field `asdasdasd: 200` from the snake example node.

**Problem / Opportunity**
The example includes an obvious garbage field which undermines credibility and can confuse validation/LLM users.

**Goals**
- Clean the example so it reflects valid data

**Non-goals**
- Changing example structure or content beyond removal

**User Stories**
- As a new user, I want examples to be clean and realistic so I can trust the tool.

**Scope**
- In scope: remove the single field
- Out of scope: broader example refactor

**Constraints**
- None

**Proposed Approach**
1. Delete the invalid field from the example node
2. Ensure tests/examples still load

**TDD Plan (Required)**
- Tests to write first (red): none (doc/example change)
- Minimal implementation (green): remove line
- Refactor: n/a

**Acceptance Criteria**
- [ ] Example file no longer contains `asdasdasd`

**Test Plan**
- `deco validate` on examples (optional)

**Dependencies**
- None

**Files / Areas Touched**
- examples/snake/.deco/nodes/systems/core.yaml

**Risks**
- None

**Open Questions**
- None

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-23bx -s in_progress

# Add a comment
bd comment deco-23bx 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-23bx -p 1

# View full details
bd show deco-23bx
```

</details>

---

<a id="deco-oxu2-add-tests-for-resolve-path-with-relative-rootdir"></a>

## 📋 deco-oxu2 Add tests for Resolve*Path with relative rootDir

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-02 21:48 |
| **Updated** | 2026-02-02 21:48 |

### Description

Cover ResolveNodesPath/ResolveHistoryPath behavior for absolute config paths and relative rootDir in config tests to lock in expected path resolution.

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-oxu2 -s in_progress

# Add a comment
bd comment deco-oxu2 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-oxu2 -p 1

# View full details
bd show deco-oxu2
```

</details>

---

<a id="deco-yle3-update-cli-help-text-to-reflect-configurable-nodes-path"></a>

## 📋 deco-yle3 Update CLI help text to reflect configurable nodes_path

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-02 21:48 |
| **Updated** | 2026-02-02 21:48 |

### Description

Help text in create command still says nodes live in .deco/nodes even though nodes_path is configurable. Update CLI docs/messages to avoid misleading users.

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-yle3 -s in_progress

# Add a comment
bd comment deco-yle3 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-yle3 -p 1

# View full details
bd show deco-yle3
```

</details>

---

<a id="deco-vtn5-resolve-path-should-be-absolute-or-comments-updated"></a>

## 📋 deco-vtn5 Resolve*Path should be absolute or comments updated

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-02 21:48 |
| **Updated** | 2026-02-02 21:48 |

### Description

ResolveNodesPath/ResolveHistoryPath claim to return absolute paths but only join rootDir; if rootDir is relative, result stays relative. Either call filepath.Abs or adjust comments to avoid misleading callers.

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-vtn5 -s in_progress

# Add a comment
bd comment deco-vtn5 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-vtn5 -p 1

# View full details
bd show deco-vtn5
```

</details>

---

<a id="deco-4p9-gdd-playbook-documentation"></a>

## 📋 deco-4p9 GDD playbook documentation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 22:28 |

### Description

Goal: Provide guidance for teams adopting deco.

## Scope (Reduced)
Write docs/playbook.md covering:
1. Getting Started - project setup, first node walkthrough
2. Node Organization - kind taxonomy, ID naming, directory structure
3. Content Modeling - tables vs parameters, block type selection
4. References - uses vs related, avoiding circular refs
5. Team Workflow - draft -> review -> approved flow
6. Common Pitfalls - over-nesting, duplicate info, stale refs

## What's Deferred
- New example projects (RPG, city-builder, visual-novel) - add incrementally when features they'd showcase are shipped
- Feature showcase matrix - update existing examples as features land

## Acceptance Criteria
- [ ] docs/playbook.md exists with sections above
- [ ] Playbook references existing examples (snake, space-invaders)
- [ ] No new example projects required

## Non-Goals
- Full game specs
- New example projects (deferred)
- Tutorial videos

### Notes

CLARIFICATIONS (from review):

1. **Feature Showcase Scope**: 'All major features demonstrated across examples' means features that are SHIPPED at time of example creation. Unchecked items (constraints, custom block types, schema rules) are aspirational - examples updated when those features land. No dependency needed; list is a roadmap, not a requirement.

2. **Non-Goals Alignment**: 'Examples should be excerpts' means each example is a representative slice of a hypothetical full GDD, not a complete design. Feature showcase applies within that scope - demonstrate features using excerpt-sized content, not full specs.

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-4p9 -s in_progress

# Add a comment
bd comment deco-4p9 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-4p9 -p 1

# View full details
bd show deco-4p9
```

</details>

---

<a id="deco-a5n-explain-error-suggest-fix-output"></a>

## ✨ deco-a5n Explain error / suggest fix output

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | 🟢 open |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 21:20 |

### Description

Goal: Make validation failures actionable for humans and AI by providing context, explanations, and fix suggestions.

## User Stories
1. New users: When validation fails, understand WHY and HOW to fix without reading docs.
2. AI assistants: Receive structured error output that enables automatic fixes.
3. Documentation: Generate error reference docs from code.

## Current Infrastructure (Already Built)

Error System:
- DecoError struct with: Code, Summary, Detail, Location, Context, Suggestion
- ErrorCodeRegistry with categories: schema (E001-E019), refs (E020-E039), validation (E040-E059), io (E060-E079), graph (E080-E099)
- ErrorFormatter for consistent output
- ErrorDocsGenerator for markdown docs

Suggestion System:
- Suggester with Levenshtein distance
- Typo detection with prefix/suffix bonus
- FormatSuggestion() for 'Did you mean X?' messages

## What's Missing

### 1. Integration with Validator
Current deco validate doesn't populate Suggestion field. Need to wire up Suggester to validation errors.

### 2. CLI Flag for Verbose Explanation
deco validate              # Current: shows error
deco validate --explain    # New: shows error + explanation + fix

### 3. Structured Output for AI
deco validate --format=json   # Machine-readable for LLM consumption

## Implementation Phases

Phase 1: Wire Up Suggestions
- Add Suggester to validator orchestrator
- Unknown field -> 'Did you mean X?'
- Missing required field -> 'Add field X with value of type Y'
- Invalid reference -> 'Node X not found. Available: A, B, C'
- Invalid status -> 'Valid statuses: draft, approved, published'

Phase 2: --explain Flag
- Verbose output with explanations
- Show file location with line number
- Include fix suggestion
- Link to 'deco help errors CODE'

Phase 3: Common Mistake Patterns
- staus -> status
- refrence -> reference
- sumary -> summary
- Missing version: 1 (must be >0)
- status: done -> status: approved

Phase 4: --format=json
- Machine-readable output
- Include fix actions for automation

## Acceptance Criteria
- [ ] Validation errors include typo suggestions ('Did you mean?')
- [ ] --explain flag shows detailed explanation with context
- [ ] --format=json outputs structured errors
- [ ] deco help errors shows error code reference
- [ ] Common typos have specific suggestions
- [ ] Missing required fields list the field name and expected type

## Non-Goals
- Auto-fix without user confirmation (separate feature)
- IDE integration (CLI only)

### Notes

CLARIFICATIONS (from review):

1. **JSON Output Schema**: --format=json outputs:
   ```json
   {
     "valid": false,
     "errors": [
       {
         "code": "E001",
         "summary": "...",
         "detail": "...",
         "location": {"file": "...", "line": N, "column": N},
         "suggestion": "Did you mean 'status'?"
       }
     ]
   }
   ```
   Schema will be documented in docs/error-format.md and versioned if structure changes.

2. **Help Integration Plan**: `deco help errors CODE` requires:
   - New subcommand in cmd/deco/help.go
   - Reads from ErrorCodeRegistry (already exists in errors/registry.go)
   - Renders markdown description for given code
   - Lists related codes in same category

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-a5n -s in_progress

# Add a comment
bd comment deco-a5n 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-a5n -p 1

# View full details
bd show deco-a5n
```

</details>

---

<a id="deco-9rm-compile-gdd-to-latex-output"></a>

## ✨ deco-9rm Compile GDD to LaTeX output

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 💤 Backlog (P4) |
| **Status** | 🟢 open |
| **Created** | 2026-01-31 16:29 |
| **Updated** | 2026-02-02 21:20 |

### Description

Goal: Generate professional, print-ready GDD documents from deco projects.

## User Stories
1. Publishers/stakeholders: Receive polished PDF of game design for review
2. Physical reference: Print copy for design meetings
3. Documentation: Archive versioned design documents

## Command Interface

deco compile -o output.tex           # Generate LaTeX
deco compile -o output.tex --format pdf  # Future: direct PDF

Flags:
- -o, --output (required): Output file path
- --format: tex (default) or pdf (future)
- --include-contracts: Include Gherkin scenarios (default: false)
- --include-issues: Include TBDs (default: true)

## Document Structure

Hierarchy mapping:
- Top-level path segment -> chapter (systems, events, items)
- Second-level -> section
- Third-level -> subsection
- Deeper -> subsubsection or bold headings

Ordering: Alphabetical within each level

## Node Rendering

Header:
- Status badge: [DRAFT], [APPROVED], etc.
- Tags as inline labels
- Summary in italics

Block types -> LaTeX:
- table: tabular with headers, booktabs
- rule: block quote, indented
- param: key: value formatting
- mechanic: tcolorbox with name, description, conditions
- list: itemize environment

References:
- Uses: hyperlinked node IDs with context
- Related: same treatment
- Events/vocabulary if present

Issues section:
- Warning box (tcolorbox)
- TBD items with severity indicator
- Makes incomplete nodes obvious

## Content Scope

Included:
- Meta (title, status, tags)
- Summary
- Content blocks
- References
- Issues (optional)

Excluded by default:
- Contracts (implementation detail, optional flag to include)
- Constraints (internal validation)
- llm_context (internal)
- Reviewers (workflow metadata)

## Implementation

Code organization:
internal/
  cli/compile.go           # Command, flags
  services/compile/
    latex.go               # Generation logic
    latex_test.go          # Unit tests
    templates/             # LaTeX templates

LaTeX packages:
- hyperref: Cross-references
- booktabs: Professional tables
- xcolor: Status colors
- geometry: Margins
- enumitem: Lists
- tcolorbox: Boxes for mechanics/issues

No external Go dependencies - pure string building with proper escaping.

## Acceptance Criteria
- [ ] deco compile -o file.tex generates valid LaTeX
- [ ] Document has title, TOC, chapters by kind
- [ ] Tables render with booktabs
- [ ] Rules render as block quotes
- [ ] Parameters render as key-value
- [ ] Mechanics render in boxes
- [ ] Node IDs are hyperlinked
- [ ] Status badges visible
- [ ] Issues section shows TBDs
- [ ] Special characters escaped (%, $, &, #, etc.)
- [ ] Unit tests for each block type

## Future Extensions (Out of Scope)
- Direct PDF compilation (requires pdflatex)
- Custom templates
- Multiple output formats (HTML, Markdown)
- Diagrams from graph command

## Status
DEFERRED: Awaiting CEO approval before work begins

### Design

# LaTeX Compile Feature Design

**Issue:** deco-9rm
**Date:** 2026-01-31
**Status:** Approved

## Overview

Add a `deco compile` command that generates a print-ready LaTeX GDD document from the design graph. Produces professional, typeset output suitable for stakeholders, publishers, or physical reference.

## Command Interface

```
deco compile -o <output.tex>
```

**Flags:**
- `-o, --output` (required): Output file path for the generated LaTeX

**Behavior:**
- Loads all nodes from `.deco/nodes/`
- Validates project is initialized (`.deco/config.yaml` exists)
- Generates a complete LaTeX document
- Writes to the specified output file

**Error cases:**
- No `.deco/` directory → "not a deco project, run deco init"
- No nodes found → "no nodes to compile"
- Can't write output file → "failed to write output: <reason>"

## Document Structure

```latex
\documentclass[11pt,a4paper]{report}
% Preamble: packages for tables, colors, hyperlinks

\title{<project name from config>}
\author{Generated by Deco}
\date{\today}

\begin{document}
\maketitle
\tableofcontents

\chapter{Systems}
  \section{Settlement}
    \subsection{Colonists}
      % Node content here

\chapter{Events}
  % ...
\end{document}
```

**Hierarchy mapping:**
- Top-level path segment → `\chapter{}` (systems, events, etc.)
- Second-level → `\section{}`
- Third-level → `\subsection{}`
- Deeper levels → `\subsubsection{}` or bold headings

**Ordering:** Alphabetical within each level.

**Title casing:** Path segments capitalized. Node titles override when available.

## Node Rendering

Each node renders these elements:

**Header area:**
- Status badge: `[DRAFT]`, `[APPROVED]`, etc.
- Tags as inline labels
- Summary/one-liner in italics

**Content blocks:**

| Block Type | LaTeX Rendering |
|------------|-----------------|
| `table` | `tabular` with headers and borders |
| `rule` | Indented block quote |
| `param` | Key-value: **key**: value |
| `mechanic` | Named box with description, conditions, outputs |
| `list` | Bulleted `itemize` |

**References section:**
- "Uses:" with hyperlinked node IDs
- "Related:" same treatment
- Events and vocabulary if present

**Issues section:**
- Warning box at end of node
- Shows TBD/question items with severity
- Makes incomplete nodes visually obvious

**Cross-references:** Node IDs become `\hyperref` links.

## Content Scope

**Included:**
- Meta (title, status, tags)
- Summary
- Content blocks (tables, rules, params, mechanics, lists)
- References
- Issues

**Excluded:**
- Contracts (implementation detail)
- Constraints (internal validation)
- llm_context (internal)

## Implementation

**Code organization:**

```
internal/
  cli/
    compile.go          # Command definition, flags
  services/
    compile/
      latex.go          # LaTeX generation logic
      latex_test.go     # Unit tests
```

**Key components:**

1. **LatexCompiler** - Main struct
   - `Compile(nodes []domain.Node, config domain.Config) (string, error)`

2. **Block renderers** - Per-type functions
   - `renderTable()`, `renderRule()`, `renderParam()`, etc.

**LaTeX packages:**
- `hyperref` - Cross-references
- `booktabs` - Tables
- `xcolor` - Colors
- `geometry` - Margins
- `enumitem` - Lists
- `tcolorbox` - Boxes for mechanics/issues

**No external Go dependencies.** Pure string building with proper escaping of LaTeX special characters.

## Testing

Unit tests with example nodes verifying:
- Correct document structure
- Each block type renders properly
- Special characters escaped
- Cross-references generated

No PDF compilation in tests.

### Notes

CLARIFICATIONS (from review):

1. **Status Alignment**: Issue remains OPEN (not blocked). DEFERRED in description is a workflow note indicating CEO approval needed before claiming. To block on approval, would need a separate approval issue as blocker. Current state is correct: open and available for discussion/approval.

2. **Ordering Rules Clarification**: 
   - Structural ordering: alphabetical by path segment (systems/ before events/ because s > e? No - 'e' < 's', so events/ first)
   - Actually: alphabetical means a-z, so events/ < items/ < systems/
   - Node titles: used for display text in section headers, NOT for ordering
   - Order is ALWAYS by path, titles only affect rendered header text

Corrected: 'Ordering: Alphabetical by path segment within each level. Node titles used for display headers only, not ordering.'

<details>
<summary>📋 Commands</summary>

```bash
# Start working on this issue
bd update deco-9rm -s in_progress

# Add a comment
bd comment deco-9rm 'Your comment here'

# Change priority (0=Critical, 1=High, 2=Medium, 3=Low)
bd update deco-9rm -p 1

# View full details
bd show deco-9rm
```

</details>

---

<a id="deco-u0jb-session-handover-custom-block-types-implemented"></a>

## 📋 deco-u0jb Session handover: custom block types implemented

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 09:01 |
| **Updated** | 2026-02-02 13:33 |
| **Closed** | 2026-02-02 13:33 |

### Notes

## Accomplished

Closed 2 issues:
- **deco-dag0**: Previous session handover (absorbed)
- **deco-87y**: Custom block types with validation hooks

## Implementation details

Added custom block types feature that allows projects to define their own block types in config:

```yaml
custom_block_types:
  powerup:
    required_fields:
      - name
      - effect
      - duration
```

Key changes:
- `internal/storage/config/repository.go`: Added `BlockTypeConfig` and `CustomBlockTypes` to Config
- `internal/services/validator/block_validator.go`: Extended to support custom types
- `internal/services/validator/validator.go`: Added `NewOrchestratorWithFullConfig`
- `internal/cli/validate.go`: Uses new full config constructor
- Added comprehensive tests for custom block types
- Added example in `examples/snake/.deco/nodes/items/powerups.yaml`
- Documented in SPEC.md

## Current state

- All code committed locally (2 commits)
- **Push needed** - git push failed due to missing credentials
- Tests not run locally (Go not available) - will be verified by CI

## Commits pending push

1. `feat(validator): custom block types with validation hooks` - main feature
2. `chore(beads): sync issue tracker` - beads sync

## Recommended next steps

1. Push commits: `git push` (needs credentials)
2. Continue with P2 features:
   - deco-apc (Configurable schema rules) - similar pattern, can reuse approach
   - deco-a5l (AI patch safety)
   - deco-p75 (Schema migrations)

---

<a id="deco-dag0-session-handover-performance-optimizations"></a>

## 📋 deco-dag0 Session handover: performance optimizations

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 22:13 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 09:00 |

### Notes

## Accomplished

Closed 4 issues:
- **deco-5gxq**: Previous session handover (absorbed)
- **deco-cxl**: Sync performance O(nodes×history) → O(history+nodes)
- **deco-3ix**: Docs alignment (README/SPEC match actual CLI)
- **deco-04f**: Validation performance (CEL caching, 109x faster)

## Current state
- All tests passing
- Working tree clean
- Pushed to master

## Recommended next steps
- P2 features are ready: deco-apc (configurable schema), deco-87y (custom blocks), deco-a5l (AI safety)
- deco-7jb (test coverage) is a good smaller task
- deco-4p9 (example projects) is user-facing documentation

## Architecture notes
- Sync now uses QueryLatestHashes() for single-pass history reading
- Validator now caches CEL environment and compiled programs
- Benchmark suite added at internal/services/validator/validator_benchmark_test.go

---

<a id="deco-5gxq-session-handover-sync-command-fixes"></a>

## 📋 deco-5gxq Session handover: sync command fixes

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 22:03 |
| **Updated** | 2026-02-01 22:04 |
| **Closed** | 2026-02-01 22:04 |

### Notes

## Accomplished

Closed 4 sync-related issues:
- **deco-0t5**: CLI mutations now record content_hash in history
- **deco-h91g**: sync now properly reports errors (exit 2 instead of 0)
- **deco-oh0q**: --dry-run returns exit 1 when changes would occur
- **deco-8xo**: Content hash now includes Kind, Glossary, Contracts, LLMContext, Constraints, Custom

## Current state
- All tests passing
- Working tree clean
- Pushed to master

## Recommended next steps
- Check `bd ready` for remaining P2/P3 issues
- deco-cxl (sync performance) is a good optimization target
- deco-apc, deco-87y are larger features

## Architecture notes
- Created internal/cli/audit.go with shared ComputeContentHash and GetCurrentUser helpers
- All mutation commands now use these shared helpers

---

<a id="deco-awz-session-handover-review-workflow-implementation-in-progress"></a>

## 📋 deco-awz Session handover: Review workflow implementation in progress

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 20:14 |
| **Updated** | 2026-02-01 20:23 |
| **Closed** | 2026-02-01 20:23 |

### Notes

## What was accomplished
- Completed 6 of 13 tasks from review workflow plan
- All domain-level and validator changes done
- CLI validate command now uses config's required_approvals

## Commits this session
- a050b62 feat(domain): add Reviewer struct and Reviewers field to Node
- caa154e feat(config): add required_approvals setting with default of 1
- 720cbbf feat(audit): add submit, approve, reject operations for review workflow
- 1f8f014 feat(validator): add ApprovalValidator for review workflow
- d7551ff feat(validator): integrate ApprovalValidator into Orchestrator
- 498c071 feat(cli): use config required_approvals in validate command

## Current state
- All tests passing
- Working tree clean
- Pushed to remote

## Next steps
Continue with deco-rs0 using docs/plans/2026-02-01-review-workflow.md:
- Task 7: Create review.go CLI with submit subcommand
- Task 8: Add approve subcommand
- Task 9: Add reject subcommand
- Task 10: Add status subcommand
- Task 11: Register review command in main.go
- Task 12: Auto-reset status on edit
- Task 13: Final integration test

Use the executing-plans skill to continue.

---

<a id="deco-ry9-session-handover-review-workflow-design-complete"></a>

## 📋 deco-ry9 Session handover: Review workflow design complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:53 |
| **Updated** | 2026-02-01 19:56 |
| **Closed** | 2026-02-01 19:56 |

### Notes

## Accomplished
- Claimed deco-rs0 (Review workflow feature)
- Completed brainstorming session with CEO
- Documented full design in deco-rs0 issue

## Design Summary (deco-rs0)
- **States**: draft → review → approved
- **On edit**: Auto-reset to draft, bump version
- **Approvals needed**: Configurable per-project in config.yaml (default: 1)
- **Approval data**: Reviewer, timestamp, optional note, version number
- **Storage**: Both in node YAML (reviewers field) AND history for audit trail
- **CLI**: `deco review {submit,approve,reject,status}`

## Implementation Plan
1. Add Reviewer struct and reviewers field to Node (domain/node.go)
2. Add required_approvals to config schema
3. Add approve/reject/submit operations to audit.go
4. Create ApprovalValidator for status transition guards
5. Create internal/cli/review.go with subcommands
6. Auto-reset status to draft on any edit (patcher integration)

## Project State
- All tests passing (90+ tests)
- Working tree clean
- deco-rs0 is IN_PROGRESS with full design documented

## Next Steps
1. Read deco-rs0 design section (`bd show deco-rs0`)
2. Use superpowers:writing-plans to create implementation plan
3. Implement in TDD style per project conventions

---

<a id="deco-xqq-session-handover-constraint-scope-enforcement"></a>

## 📋 deco-xqq Session handover: Constraint scope enforcement

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:37 |
| **Updated** | 2026-02-01 19:42 |
| **Closed** | 2026-02-01 19:42 |

### Notes

## Accomplished
- Closed deco-79k: Enforce constraint scope in validator
  - Added matchesScope() to skip constraints that don't match node kind/pattern
  - Scope patterns: "all", exact kind match, glob patterns
  
- Closed deco-e4g: Expand issues/TBD system
  - Added --kind filter (node kind)
  - Added --tag filter (node tag)
  - Added --all flag (include resolved)
  - Added --json flag (JSON output)
  - Added --quiet flag (counts only)
  - Added --summary flag (per-node rollup)
  - Added 20+ new tests

## Project State
- All tests passing (90+ tests)
- Working tree has uncommitted changes

## Recommended Next Steps
- deco-rs0 (P2): Review workflow with approvals
- deco-apc (P2): Configurable schema rules

---

<a id="deco-c4y-session-handover-file-location-in-validation-errors"></a>

## 📋 deco-c4y Session handover: File location in validation errors

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:28 |
| **Updated** | 2026-02-01 19:34 |
| **Closed** | 2026-02-01 19:34 |

### Notes

## Accomplished
- Closed deco-acy: validate output lacks file/line context despite YAML location tracker
  - Added SourceFile field to domain.Node (yaml:"-" to skip serialization)
  - Updated YAMLRepository to set SourceFile when loading nodes
  - Updated all 8 validators to include Location in DecoError:
    - SchemaValidator, ContentValidator, ReferenceValidator
    - ConstraintValidator, DuplicateIDValidator, UnknownFieldValidator
    - ContractValidator, BlockValidator
  - Fixed collector deduplication to include summary for file-only locations
  
## Project State
- All tests passing (88+ tests)
- Working tree clean
- Pushed to master (ab86e12)

## Output Example
Validation errors now show file path:
```
[E008] Missing required field: Version at path/to/file.yaml: Node Version is required
[E020] Reference not found: target at path/to/file.yaml: Referenced node 'target' does not exist
```

## Recommended Next Steps
- deco-79k (P2): Enforce constraint scope in validator
- deco-e4g (P2): Expand issues/TBD system
- Future: Add line/column context using LocationTracker for even more precise error locations

---

<a id="deco-4y9-session-handover-block-validation-complete"></a>

## 📋 deco-4y9 Session handover: Block validation complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:18 |
| **Updated** | 2026-02-01 19:19 |
| **Closed** | 2026-02-01 19:19 |

### Notes

## Accomplished
- Closed deco-603: Add type-specific validation for blocks
  - Created BlockValidator with type-specific validation for rule/table/param/mechanic/list blocks
  - Integrated into orchestrator validation pipeline
  - 21 new tests for block validation, 1 orchestrator integration test
  - Updated error codes E047-E050 for block validation errors

## Project State
- All tests passing (88 validator tests total)
- Working tree clean
- Pushed to master

## Recommended Next Steps
Check `bd ready` for available work. Good candidates:
- deco-acy (P3): validate output lacks file/line context
- deco-79k (P2): Enforce constraint scope in validator
- deco-e4g (P2): Expand issues/TBD system

---

<a id="deco-lwl-session-handover-contract-validation-cli-complete"></a>

## 📋 deco-lwl Session handover: contract validation CLI complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:05 |
| **Updated** | 2026-02-01 19:06 |
| **Closed** | 2026-02-01 19:06 |

### Notes

## Accomplished

- Closed deco-0sk: Add contract syntax validator
- Closed deco-wxh: Validate contract scenarios reference valid nodes
- Closed deco-h2c: Add contract validation to deco validate command

Full contract validation pipeline now complete:
- Syntax validation (E100, E101, E103, E104)
- Node reference validation (E102) with typo suggestions
- CLI integration with deco validate command
- Updated help text to document contract validation
- Added 2 CLI tests for contract validation errors

## Current State

- All tests passing (57 validator tests, full suite green)
- Working tree has uncommitted changes

## Recommended Next

- deco-hsx: Require content for approved/published nodes (P2)
- deco-603: Add type-specific validation for blocks (P2)
- deco-e4g: Expand issues/TBD system (P2)

---

<a id="deco-v4o-session-handover-contract-ref-validation-added-deco-wxh"></a>

## 📋 deco-v4o Session handover: contract ref validation added (deco-wxh)

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:03 |
| **Updated** | 2026-02-01 19:04 |
| **Closed** | 2026-02-01 19:04 |

### Notes

## Accomplished

- Closed deco-0sk: Add contract syntax validator
- Closed deco-wxh: Validate contract scenarios reference valid nodes
- Extended ContractValidator to validate @node.id references (E102)
- Added typo suggestions for invalid references
- Added 6 tests for node reference validation
- Total 57 validator tests, all passing

## Current State

- All tests passing
- Working tree has uncommitted changes

## Recommended Next

From bd ready:
- deco-h2c: Add contract validation to deco validate command (unblocked by deco-wxh)
- deco-hsx: Require content for approved/published nodes (P2)
- deco-603: Add type-specific validation for blocks (P2)

---

<a id="deco-ykl-session-handover-contract-validator-added-deco-0sk"></a>

## 📋 deco-ykl Session handover: contract validator added (deco-0sk)

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 19:00 |
| **Updated** | 2026-02-01 19:03 |
| **Closed** | 2026-02-01 19:03 |

### Notes

## Accomplished

- Closed deco-0sk: Add contract syntax validator
- Created internal/services/validator/contract.go
- Validates: scenario names (E100), duplicate names (E103), empty steps (E101), no steps (E104)
- Integrated ContractValidator into Orchestrator
- Added 11 tests for contract validation
- All 51 validator tests pass, full suite green

## Current State

- All tests passing
- Working tree has uncommitted changes

## Recommended Next

From bd ready:
- deco-wxh: Validate contract scenarios reference valid nodes (unblocked by deco-0sk)
- deco-hsx: Require content for approved/published nodes (P2)
- deco-603: Add type-specific validation for blocks (P2)

---

<a id="deco-g50-session-handover-contract-parser-added-deco-4ci"></a>

## 📋 deco-g50 Session handover: contract parser added (deco-4ci)

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:54 |
| **Updated** | 2026-02-01 19:00 |
| **Closed** | 2026-02-01 19:00 |

### Notes

## Accomplished

- Closed deco-4ci: Implement contract scenario parser
- Created internal/domain/contract.go with Step, Scenario types
- Added ParseContract/ParseContracts functions
- Node reference extraction via @node.id syntax
- Added contract error codes E100-E119
- Full test coverage in contract_test.go
- Fixed epic deco-6iv dependencies (19 P2/P3 tasks now unblocked)

## Current State

- All tests passing
- Working tree clean

## Recommended Next

From bd ready (10 issues now available):
- deco-0sk: Add contract syntax validator (unblocked by deco-4ci)
- deco-603: Add type-specific validation for blocks (P2)
- deco-79k: Enforce constraint scope in validator (P2)
- deco-hsx: Require content for approved/published nodes (P2)
- And 6 more P2 features

---

<a id="deco-gbk-session-handover-audit-history-added"></a>

## 📋 deco-gbk Session handover: audit history added

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:47 |
| **Updated** | 2026-02-01 18:50 |
| **Closed** | 2026-02-01 18:50 |

### Notes

## Accomplished

- Closed deco-zn8: Add audit history to all node-modifying commands
- Commands now log to .deco/history.jsonl: create, set, append, unset, apply
- Each command records Before/After containing only changed fields
- deco diff shows meaningful output for all operations

## Current State

- All tests passing
- Working tree clean (after commit)

## Recommended Next

From bd ready:
- deco-6iv: Epic for end-state product vision (P2)
- deco-4ci: Implement contract scenario parser (P4)
- deco-9rm: Compile GDD to LaTeX output (P4)

---

<a id="deco-fre-session-handover-stats-command-added"></a>

## 📋 deco-fre Session handover: stats command added

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:37 |
| **Updated** | 2026-02-01 18:42 |
| **Closed** | 2026-02-01 18:42 |

### Notes

## Accomplished

- Closed deco-1pi: Add deco stats command for project overview
- Shows: node counts by kind/status, open issues by severity, dangling refs, constraint violations
- Full test coverage (8 test cases)

## Current State

- All tests passing
- Working tree clean

## Recommended Next

From bd ready:
- deco-6iv: Epic for end-state product vision
- deco-zn8: Add audit history to all node-modifying commands
- deco-4ci: Implement contract scenario parser

---

<a id="deco-vq3-session-handover-diff-command-added"></a>

## 📋 deco-vq3 Session handover: diff command added

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:32 |
| **Updated** | 2026-02-01 18:35 |
| **Closed** | 2026-02-01 18:35 |

### Notes

## Accomplished

- Closed deco-t2h: Add deco diff command to show node history
- Shows before/after for each change with +/- indicators
- Supports --since (RFC3339, date, relative: 2h, 1d, 1w)
- Supports --last N to limit output
- Full test coverage

## Current State

- All tests passing
- Working tree clean (except this handover issue)

## Recommended Next

From `bd ready`:
- deco-1pi: Add deco stats command for project overview
- deco-6iv: Epic for end-state product vision
- deco-4ci: Implement contract scenario parser

---

<a id="deco-4kg-session-handover-graph-command-added"></a>

## 📋 deco-4kg Session handover: graph command added

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:26 |
| **Updated** | 2026-02-01 18:28 |
| **Closed** | 2026-02-01 18:28 |

### Notes

## Accomplished

- Closed deco-5uc: Add deco graph command to visualize dependencies
- Outputs DOT (default) or Mermaid format
- Edges from refs.uses (solid) and refs.related (dashed)

## Current State

- All tests passing
- Working tree clean, pushed to origin/master

## Recommended Next

From `bd ready`:
- deco-t2h: Add deco diff command to show node history
- deco-1pi: Add deco stats command for project overview
- deco-6iv: Epic for end-state product vision

---

<a id="deco-dyi-session-handover-3-cli-commands-added"></a>

## 📋 deco-dyi Session handover: 3 CLI commands added

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:22 |
| **Updated** | 2026-02-01 18:23 |
| **Closed** | 2026-02-01 18:23 |

### Notes

## Accomplished

Closed 3 issues:
- deco-513: Add deco create command for scaffolding nodes
- deco-t5h: Add deco rm command to delete nodes  
- deco-gp6: Add deco issues command to list all TBDs

## Current State

- All tests passing
- Working tree clean (pushed to origin/master)
- New commands: `deco create`, `deco rm`, `deco issues`

## Recommended Next

From `bd ready`:
- deco-5uc: Add deco graph command to visualize dependencies
- deco-t2h: Add deco diff command to show node history
- deco-1pi: Add deco stats command for project overview

The graph command would be a good next step - it complements the existing ref system.

---

<a id="deco-duf-session-handover-mv-command-and-cli-epic-complete"></a>

## 📋 deco-duf Session handover: mv command and CLI epic complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:14 |
| **Updated** | 2026-01-31 16:15 |
| **Closed** | 2026-01-31 16:15 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed
1. **deco-747**: Wrote 17 comprehensive tests for mv command
   - Structure, rename, ref updates, history, errors, integration
2. **deco-vzo**: Implemented deco mv command
   - Renames node ID and file
   - Updates Uses and Related references across all nodes
   - Increments version on affected nodes
   - Records move operation in history
3. **deco-c5c**: Shell completion - Cobra provides automatically
4. **deco-7yo**: Closed CLI epic (all commands implemented)

### Project State
- Tests: All passing
- Working tree: Clean
- Pushed to origin/master

### Remaining Work (P4 backlog)
All 4 open issues are contract-related:
- deco-4ci: Implement contract scenario parser (ready)
- deco-0sk: Add contract syntax validator (blocked)
- deco-wxh: Validate contract scenarios reference valid nodes (blocked)
- deco-h2c: Add contract validation to deco validate command (blocked)

### Recommended Next Steps
1. Close this handover when received
2. Optional: Start contract feature chain with deco-4ci
3. Consider if contracts need more spec definition first

---

<a id="deco-gw6-session-handover-node-rename-service-complete"></a>

## 📋 deco-gw6 Session handover: Node rename service complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:07 |
| **Updated** | 2026-01-31 16:08 |
| **Closed** | 2026-01-31 16:08 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed: Node Rename Service (2 issues)
1. **deco-bhb**: Write tests for node rename - 26 comprehensive tests
2. **deco-142**: Implement node rename with ref updates
   - Renames node ID
   - Updates all Uses and Related references across the graph
   - Increments version on nodes whose references were updated
   - Deep copy to avoid modifying original nodes
   - Full validation (empty IDs, non-existent, collisions)

### In Progress
- **deco-747**: Write tests for deco mv command (claimed but not started)
  - Next session should continue this or unclaim if needed

### Project State
- Tests: All passing (26 new rename tests + existing tests)
- Working tree: Clean
- New files: internal/services/refactor/rename.go and rename_test.go

### Recommended Next Steps
1. Continue deco-747: Write mv command tests
2. deco-vzo: Implement deco mv command (uses Renamer service)
3. deco-c5c: Shell completion generation

---

<a id="deco-8rt-session-handover-cli-mutation-commands-complete"></a>

## 📋 deco-8rt Session handover: CLI mutation commands complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:01 |
| **Updated** | 2026-01-31 16:07 |
| **Closed** | 2026-01-31 16:07 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed: CLI Mutation Commands (12 issues closed)

**Query Command (2 issues):**
1. **deco-vz6**: Write tests for deco query command - 24 tests
2. **deco-1eg**: Implement deco query command
   - Text search in title/summary (case-insensitive)
   - Combined with kind/status/tag filters
   - Same table output as list

**Set Command (2 issues):**
3. **deco-0jp**: Write tests for deco set command - 20 tests
4. **deco-53i**: Implement deco set command
   - Set field values using Patcher.Set()
   - Supports simple fields and array indices
   - Auto-increments version

**Append Command (2 issues):**
5. **deco-zix**: Write tests for deco append command - 19 tests
6. **deco-1v8**: Implement deco append command
   - Append values to array fields
   - Auto-increments version

**Unset Command (2 issues):**
7. **deco-3vb**: Write tests for deco unset command - 22 tests
8. **deco-31k**: Implement deco unset command
   - Remove field values or array elements
   - Protects required fields (id, kind, version, status, title)
   - Auto-increments version

**Apply Command (2 issues):**
9. **deco-tzg**: Write tests for deco apply command - 24 tests
10. **deco-0kr**: Implement deco apply command
    - Batch operations from JSON file
    - Transactional (rollback on error)
    - --dry-run flag for validation

**History Command (2 issues):**
11. **deco-7kj**: Write tests for deco history command - 17 tests
12. **deco-6j1**: Implement deco history command
    - View audit log entries
    - Filter by --node and --limit

### Project State:
- **Tests:** All passing (146+ new CLI tests)
- **Working tree:** Clean, ready to push
- **Commands now available:**
  - init, validate, list, show, query
  - set, append, unset, apply
  - history

### Next Steps:
1. **deco-bhb/deco-142**: Node rename service (needed for mv command)
2. **deco-747/deco-vzo**: Mv command (depends on rename service)
3. **deco-c5c**: Shell completion generation

---

<a id="deco-flq-session-handover-cli-commands-validate-list-show"></a>

## 📋 deco-flq Session handover: CLI commands (validate, list, show)

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 15:46 |
| **Updated** | 2026-01-31 15:49 |
| **Closed** | 2026-01-31 15:49 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed: CLI Commands - Phase 1 (6 issues closed)

**Validate Command (2 issues):**
1. **deco-sbk**: Write tests for deco validate command
   - 24 comprehensive test cases
   - Tests schema, reference, constraint validation
   - Tests --quiet flag, error output, exit codes
   
2. **deco-1fm**: Implement deco validate command
   - Uses Validator Orchestrator from services
   - Exit code 0 on valid, 1 on errors
   - Rust-like error formatting
   - All 24 tests passing

**List Command (2 issues):**
3. **deco-wlu**: Write tests for deco list command
   - 23 comprehensive test cases
   - Tests --kind, --status, --tag filters
   - Tests combined filters, table output
   
4. **deco-qll**: Implement deco list command
   - Uses QueryEngine.Filter() from services
   - Dynamic table formatting
   - Multiple filter support with AND logic
   - All 23 tests passing

**Show Command (2 issues):**
5. **deco-s5c**: Write tests for deco show command
   - 19 comprehensive test cases
   - Tests node display, reverse refs, --json output
   
6. **deco-4q8**: Implement deco show command
   - Uses GraphBuilder.BuildReverseIndex() for reverse refs
   - Human-readable and JSON output modes
   - Displays all node fields comprehensively
   - All 19 tests passing

### Project State:
- **Tests:** All passing (66 new CLI tests + existing)
  - CLI total: 24 init + 5 root + 24 validate + 23 list + 19 show = 95 tests
  - All other test suites: passing (cached)
- **Coverage:** 83/105 issues closed (79%)
- **Working tree:** Clean, pushed to remote (commit 48bd30c)
- **Commands:** init, validate, list, show - all working end-to-end

### Architecture Status:
```
✅ Domain Layer: Complete
✅ Error System: Complete
✅ Storage Layer: Complete
✅ Service Layer: Complete
✅ CLI Layer: Core commands complete
   ✅ Root command: Help, version, global flags
   ✅ Init command: Project initialization
   ✅ Validate command: Schema + ref + constraint validation
   ✅ List command: Filter by kind/status/tag, table output
   ✅ Show command: Node details + reverse refs + JSON
   ⏳ Query command: Advanced filtering (next)
   ⏳ Set/Append/Unset/Apply: Mutation commands
   ⏳ Mv command: Node rename
   ⏳ History command: Audit log
```

## Next Steps: Continue CLI Commands

### Recommended Priority Order:

**1. Query Command (deco-vz6/deco-1eg)** - NEXT PRIORITY
   - Advanced filtering and search
   - Use QueryEngine.Filter() + Search()
   - Combine filters with text search
   - Issue IDs: deco-vz6 (tests), deco-1eg (impl)

**2. Set Command (deco-0jp/deco-53i)**
   - Use Patcher.Set() for field updates
   - Validate before and after
   - Issue IDs: deco-0jp (tests), deco-53i (impl)

**3. Append Command (deco-zix/deco-1v8)**
   - Use Patcher.Append() for arrays
   - Validate after append
   - Issue IDs: deco-zix (tests), deco-1v8 (impl)

**4. Unset Command (deco-3vb/deco-31k)**
   - Use Patcher.Unset() to remove fields
   - Validate after unset
   - Issue IDs: deco-3vb (tests), deco-31k (impl)

**5. Apply Command (deco-tzg/deco-0kr)**
   - Use Patcher.Apply() for batch operations
   - Support JSON patch format
   - Issue IDs: deco-tzg (tests), deco-0kr (impl)

**6. History Command (deco-7kj/deco-6j1)**
   - Use HistoryRepository for audit log
   - Issue IDs: deco-7kj (tests), deco-6j1 (impl)

**7. Mv Command (deco-747/deco-vzo + deco-bhb/deco-142)**
   - Requires node rename service (deco-142)
   - Update all references automatically
   - Issue IDs: deco-747 (tests), deco-vzo (impl), deco-bhb/deco-142 (rename service)

### Ready-to-Work Issues:
```bash
bd ready
```

Currently 10 issues ready to work (no blockers).

## Implementation Notes:

### CLI Pattern Established:
1. Write tests first (TDD) - verify they fail
2. Implement command using services (never direct storage)
3. Test structure: command structure, happy path, error cases, flags, integration
4. All commands follow this structure:
   - Load config to verify project
   - Load nodes as needed
   - Use service layer for business logic
   - Format output appropriately
   - Return proper exit codes

### Services Available:
- **QueryEngine**: Filter (kind/status/tags), Search (text)
- **Patcher**: Set, Append, Unset, Apply (batch with rollback)
- **GraphBuilder**: Dependencies, cycles, topological sort, reverse index
- **Validators**: Schema, Reference, Constraint, Orchestrator
- **ConfigRepository**: Load/save project config
- **NodeRepository**: Load/save/delete nodes
- **HistoryRepository**: Audit log queries

### Test Infrastructure:
- Use t.TempDir() for isolated test directories
- Create YAML files directly (no repositories in tests)
- Pass directory argument to commands
- Test all flags (long and short versions)
- Test error cases and edge cases
- Test integration with root command

### Manual Testing Pattern:
```bash
go test ./...                    # All tests
go build ./cmd/deco              # Build binary
./deco init /tmp/test-project    # Create test project
./deco validate /tmp/test-project
./deco list /tmp/test-project
./deco show <id> /tmp/test-project
```

## Technical Context:

### Recent Commits:
- bb64751: Implement validate command (24 tests)
- 0f7a39a: Implement list command (23 tests)
- 48bd30c: Implement show command (19 tests)

### Session Workflow:
1. Check `bd ready` for available work
2. Claim issue with `bd update <id> --status=in_progress`
3. Write tests first (TDD)
4. Verify tests fail
5. Implement to make tests pass
6. Add command to cmd/deco/main.go
7. Manual end-to-end testing
8. Close issue with `bd close <id>`
9. Stage, commit, sync, push

### Commit Message Format:
```
Implement <command> command with <key features>

<detailed implementation notes>
- Test file: X test cases covering...
- Implementation: Uses <services>...
- Integration: Wire up in main.go

All X tests passing. Manual end-to-end testing confirmed.

Closes: <issue-ids>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### Session End Checklist:
```
[x] 1. Close finished issues (bd close <id1> <id2> ...)
[x] 2. Create handover issue (this issue)
[x] 3. Sync beads (bd sync)
[x] 4. Stage and commit code (git add && git commit)
[x] 5. Final sync (bd sync)
[x] 6. Push to remote (git push)
```

## Query Command Details (Next Task):

The query command will be similar to list but with text search capability:
- Combines filtering (kind/status/tag) with text search
- Uses QueryEngine.Filter() + QueryEngine.Search()
- Search is case-insensitive, searches title and summary
- Output in same table format as list
- All filters are AND logic

Expected flags:
- --kind, --status, --tag (same as list)
- --search or -s for text search term
- Possibly --json for JSON output

Implementation should be straightforward given list and show examples.

---

<a id="deco-934-session-handover-cli-foundation-complete"></a>

## 📋 deco-934 Session handover: CLI foundation complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 15:31 |
| **Updated** | 2026-01-31 15:33 |
| **Closed** | 2026-01-31 15:33 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed: CLI Foundation (5 issues closed)

**Cobra Root Command (2 issues):**
1. **deco-336**: Write tests for Cobra root command
   - Tests for version flag, help output, global flags, subcommand registration
   - 13 comprehensive test cases, all passing

2. **deco-7au**: Implement Cobra root command
   - cmd/deco/main.go: CLI entry point
   - internal/cli/root.go: Root command implementation
   - Global flags: --config (-c), --verbose, --quiet (-q)
   - Version flag: --version (-v)
   - Automatic help generation

**Init Command (2 issues):**
3. **deco-sz1**: Write tests for deco init command
   - Tests for directory creation, config generation, force flag
   - Tests for existing project detection
   - 19 comprehensive test cases, all passing

4. **deco-7ba**: Implement deco init command
   - internal/cli/init.go: Init command implementation
   - Creates .deco/ directory structure
   - Generates config.yaml with sensible defaults
   - Creates nodes/ subdirectory
   - --force (-f) flag to reinitialize
   - Accepts optional directory argument

**Services Epic:**
5. **deco-t6q**: Closed Services epic (core complete)
   - All core services implemented and tested
   - Node rename (deco-142) deferred to Phase 4
   - Service layer ready for CLI implementation

### Project State:
- **Tests:** All passing (163 total tests)
  - CLI: 5 root tests + 19 init tests = 24 tests
  - Domain: 30 tests
  - Errors: 38 tests
  - Services: 102 tests
  - Storage: 21 tests

- **Coverage:** 76/104 issues closed (73%)
- **Working tree:** Clean, pushed to remote (commit fa7ab4e)
- **Dependencies:** Added github.com/spf13/cobra v1.10.2

### Architecture Status:
```
✅ Domain Layer: Complete
✅ Error System: Complete
✅ Storage Layer: Complete
✅ Service Layer: Complete (core)
⏳ CLI Layer: Foundation complete, commands in progress
   ✅ Root command: Help, version, global flags
   ✅ Init command: Project initialization
   ⏳ Validate command: Next priority
   ⏳ List/Show commands: Read operations
   ⏳ Set/Append/Unset/Apply: Mutation commands
   ⏳ Query/History: Advanced features
   ⏳ Mv command: Advanced (requires deco-142)
```

## Next Steps: CLI Commands - Phase 1

The CLI foundation is complete. Continue with core read commands:

### Recommended Priority Order:

**1. Validate Command (deco-sbk/deco-1fm)** - HIGHEST PRIORITY
   - Use Validator Orchestrator from services
   - Critical foundation for all mutation commands
   - Tests schema, references, and constraints
   - Issue IDs: deco-sbk (tests), deco-1fm (impl)

**2. List Command (deco-wlu/deco-qll)**
   - Use QueryEngine.Filter() for kind/status/tags
   - List nodes with filtering
   - Issue IDs: deco-wlu (tests), deco-qll (impl)

**3. Show Command (deco-s5c/deco-4q8)**
   - Display node details
   - Use GraphBuilder for reverse references
   - Issue IDs: deco-s5c (tests), deco-4q8 (impl)

**4. Query Command (deco-vz6/deco-1eg)**
   - Advanced filtering and search
   - Use QueryEngine.Filter() + Search()
   - Issue IDs: deco-vz6 (tests), deco-1eg (impl)

### Phase 2: Mutation Commands (After Validate)

5. **Set** (deco-0jp/deco-53i): Use Patcher.Set()
6. **Append** (deco-zix/deco-1v8): Use Patcher.Append()
7. **Unset** (deco-3vb/deco-31k): Use Patcher.Unset()
8. **Apply** (deco-tzg/deco-0kr): Use Patcher.Apply() for batch ops

### Phase 3: Advanced Commands

9. **History** (deco-7kj/deco-6j1): Use HistoryRepository
10. **Mv** (deco-747/deco-vzo + deco-bhb/deco-142): Requires node rename service

### Ready-to-Work Issues:
```bash
bd ready
```

Currently 14 issues ready to work (no blockers).

## Implementation Notes:

### CLI Pattern Established:
1. Write tests first (TDD)
2. Test structure, flags, integration
3. Implement command using services
4. All commands use services, never direct storage access
5. Clean separation: CLI → Services → Storage

### Services Available:
- **ConfigRepository**: Load/save project config
- **NodeRepository**: Load/save/delete nodes
- **HistoryRepository**: Audit log queries
- **GraphBuilder**: Dependencies, cycles, topological sort, reverse index
- **QueryEngine**: Filter (kind/status/tags), Search (text)
- **Patcher**: Set, Append, Unset, Apply (batch with rollback)
- **Validators**: Schema, Reference, Constraint, Orchestrator

### Test Infrastructure:
- Use t.TempDir() for isolated test directories
- Pass directory argument to avoid os.Chdir() issues
- Test all flags (long and short versions)
- Test error cases and edge cases
- Test integration with root command

### Next Session Checklist:
1. Start with validate command (deco-sbk)
2. Follow TDD: write tests, verify failure, implement
3. Use Validator Orchestrator from services
4. Continue with list and show commands
5. Build toward mutation commands

## Technical Context:

### Commit Message Format:
```
<summary line>

<detailed changes>

Closes: <issue-ids>
Related: <related-issues>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### Running Commands:
```bash
go test ./...              # All tests
go build ./cmd/deco        # Build binary
./deco --help              # Test CLI
bd ready                   # Find work
bd show <id>               # View issue
bd update <id> --status=in_progress  # Claim work
bd close <id>              # Complete work
```

### Validation Command Priority:
The validate command is CRITICAL because:
- All mutation commands should validate before AND after changes
- Tests the full error system (codes, formatting, suggestions)
- Exercises all three validators (schema, reference, constraint)
- Provides confidence in data integrity
- Required before implementing set/append/unset/apply commands

---

<a id="deco-wxe-session-handover-service-layer-complete-67-done"></a>

## 📋 deco-wxe Session handover: Service layer complete - 67% done

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 15:20 |
| **Updated** | 2026-01-31 15:30 |
| **Closed** | 2026-01-31 15:30 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed: Service Layer (18 issues closed)

**Patcher Service Complete (2 issues):**
1. **deco-h9q/deco-7yp**: Patcher Apply operation
   - Batch patch operations with transactional rollback
   - Supports set/append/unset in single atomic operation
   - Snapshot/restore using gob encoding for rollback
   - 12 comprehensive tests for apply, all passing

**QueryEngine Service Complete (4 issues):**
2. **deco-3jm/deco-0s1**: QueryEngine Filter
   - Filter by kind, status, tags (AND logic)
   - Combined criteria support
   - 10 comprehensive tests, all passing

3. **deco-nop/deco-oam**: QueryEngine Search  
   - Case-insensitive partial text search
   - Searches title and summary fields
   - 8 comprehensive tests, all passing

**Validator Services Complete (12 issues):**
4. **deco-2t4/deco-bx5**: SchemaValidator
   - Validates required fields (ID, Kind, Version, Status, Title)
   - Unique error summaries for deduplication
   - 8 comprehensive tests, all passing

5. **deco-xw8/deco-jna**: ReferenceValidator
   - Validates Uses and Related references resolve
   - Generates suggestions for typos using Levenshtein distance
   - 7 comprehensive tests, all passing

6. **deco-eft/deco-9ls**: ConstraintValidator
   - CEL (Common Expression Language) evaluation
   - Supports node field access in constraints
   - 5 comprehensive tests, all passing
   - Added google/cel-go v0.27.0 dependency

7. **deco-az1/deco-2y3**: Validator Orchestrator
   - Coordinates all three validators
   - Aggregates errors using Collector
   - Runs validators in order: schema → references → constraints
   - 5 comprehensive tests, all passing

### Project State:
- **Tests:** All passing (139 total tests)
  - Domain: 30 tests
  - Errors: 21 tests + 17 collector tests = 38 tests
  - Services:
    - Graph: 26 tests
    - Patcher: 30 tests
    - Query: 18 tests
    - Validator: 28 tests
  - Storage: 7 + 3 + 11 = 21 tests

- **Coverage:** 69/103 issues closed (67%)
- **Working tree:** Clean, pushed to remote (commit be200c6)

### Architecture Complete:
```
✅ Domain Layer: Complete (Node, Graph, Issue, Constraint, Ref, AuditEntry, DecoError)
✅ Error System: Complete (structure, codes, formatter, docs, YAML, suggestions, aggregation)
✅ Storage Layer: Complete (NodeRepository, ConfigRepository, HistoryRepository)
✅ Service Layer: COMPLETE
   ✅ GraphBuilder: Build, dependencies, cycle detection, topological sort, reverse index
   ✅ Patcher: Set, append, unset, apply (batch with rollback)
   ✅ QueryEngine: Filter (kind/status/tags), Search (text)
   ✅ Validators: Schema, Reference, Constraint (CEL), Orchestrator
⏳ CLI Layer: Not started (all dependencies complete!)
```

## Next Steps: CLI Implementation

The service layer is now **100% complete**. All business logic is implemented and tested. The next phase is CLI command implementation using Cobra.

### Recommended Priority Order:

**Phase 1: Core Read Commands (Foundation)**
1. **deco init** (deco-sz1/deco-7ba)
   - Initialize project structure
   - Create .deco/ directory
   - Unblocks: deco validate

2. **deco validate** (deco-sbk/deco-1fm)
   - Use Validator Orchestrator
   - Report schema + reference + constraint errors
   - Unblocks: all mutation commands (validation first)

3. **deco list** (deco-wlu/deco-qll)
   - Use QueryEngine for filtering
   - List nodes by kind/status/tags
   - Unblocks: workflow commands

4. **deco show <id>** (deco-s5c/deco-4q8)
   - Display node details
   - Use GraphBuilder for reverse references
   - Unblocks: none (standalone)

**Phase 2: Query Commands**
5. **deco query** (deco-vz6/deco-1eg)
   - Use QueryEngine filter + search
   - Advanced filtering and search
   - Unblocks: none

**Phase 3: Mutation Commands (Use Patcher)**
6. **deco set <id> <path> <value>** (deco-0jp/deco-53i)
   - Use Patcher.Set()
   - Validate before/after
   - Unblocks: none

7. **deco append <id> <path> <value>** (deco-zix/deco-1v8)
   - Use Patcher.Append()
   - Unblocks: none

8. **deco unset <id> <path>** (deco-3vb/deco-31k)
   - Use Patcher.Unset()
   - Unblocks: none

9. **deco apply <patch-file>** (deco-tzg/deco-0kr)
   - Use Patcher.Apply() for batch operations
   - Perfect for AI integration
   - Unblocks: none

**Phase 4: Advanced Commands**
10. **deco mv <old> <new>** (deco-747/deco-vzo + deco-bhb/deco-142)
    - Rename node with reference updates
    - Use ReferenceValidator to find refs
    - Use Patcher for updates
    - Complex but service layer ready

11. **deco history [<id>]** (deco-7kj/deco-6j1)
    - Use HistoryRepository
    - Display audit log
    - Unblocks: none

### Implementation Notes:

**Cobra Setup:**
- Root command (deco-336/deco-7au) - basic CLI structure
- Add shell completion (deco-c5c) - nice to have

**Testing Strategy:**
- Each command has test + implementation issue
- Test with real YAML files
- Validate end-to-end workflows
- Integration tests for command chains

**Service Integration:**
- All services are ready and tested
- NodeRepository for loading/saving
- ConfigRepository for project settings
- Validators for pre/post mutation checks
- QueryEngine for filtering/searching
- Patcher for all mutations
- GraphBuilder for dependency analysis

### Dependencies Ready:
All CLI commands can now be implemented as thin wrappers around services:
- ✅ Storage layer handles persistence
- ✅ QueryEngine handles filtering/searching
- ✅ Patcher handles mutations
- ✅ Validators handle validation
- ✅ GraphBuilder handles dependencies
- ✅ Error system handles reporting

### CLI Epic:
The CLI epic (deco-7yo) tracks overall CLI progress. Currently 0/21 CLI commands implemented. Start with init + validate to establish foundation.

## Technical Context:

### Key Design Decisions:
1. **Transactional Patcher**: Uses gob encoding for deep copy/restore
2. **CEL Integration**: google/cel-go for constraint expressions
3. **Error Deduplication**: Code + Summary for errors without location
4. **Validator Order**: Schema → References → Constraints (fail fast)

### Code Quality:
- Strict TDD: All tests written before implementation
- Comprehensive edge cases: nil, empty, invalid, boundaries
- All 139 tests passing
- Clear separation of concerns

---

<a id="deco-tw3-session-handover-service-layer-progress-54-complete"></a>

## 📋 deco-tw3 Session handover: Service layer progress - 54% complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 15:05 |
| **Updated** | 2026-01-31 15:19 |
| **Closed** | 2026-01-31 15:19 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed Issues (14 total, bringing total to 55/102 - 53.9%):

**Error System Complete:**
1. **deco-5wv/deco-5fv**: Error aggregation (tests + implementation)
   - Collector for aggregating multiple DecoErrors
   - Sorts by file, line, column; errors without location come last
   - Deduplicates by code + location (file:line:column)
   - Max error limit with truncation tracking
   - Methods: Add, AddBatch, Errors, HasErrors, Count, Truncated, Reset
   - 17 comprehensive tests, all passing

2. **deco-5xa**: Error system epic CLOSED
   - All error system work complete
   - Core error structure, codes, formatter, docs
   - YAML integration (location tracking, context extraction)
   - Suggestion engine (Levenshtein distance)
   - Error aggregation

**Graph Services Complete:**
3. **deco-jxy/deco-0o3**: GraphBuilder service (tests + implementation)
   - Build Graph from Node slices with duplicate detection
   - BuildDependencyMap: extract Refs.Uses relationships
   - DetectCycle: find circular dependencies with path
   - TopologicalSort: dependency-ordered node processing
   - 16 comprehensive tests, all passing

4. **deco-im0/deco-pqf**: Reverse reference indexing (tests + implementation)
   - BuildReverseIndex: "who references this node" lookup
   - Includes both Uses and Related references
   - Deduplicates multiple ref types to same target
   - Only tracks refs to existing nodes (no dangling refs)
   - 10 comprehensive tests, all passing

**Patcher Service Complete:**
5. **deco-aoo/deco-7rp**: Patcher set operation (tests + implementation)
6. **deco-3co/deco-k3j**: Patcher append operation (tests + implementation)
7. **deco-q41/deco-750**: Patcher unset operation (tests + implementation)
   - Unified Patcher service in internal/services/patcher/
   - Set: modify field values with path notation and type conversion
   - Append: add values to arrays (handles nil, empty, existing arrays)
   - Unset: remove fields or array elements (protects required fields)
   - Path support: dot notation (summary), array indexing (tags[0])
   - Reflection-based with proper type handling
   - 18 comprehensive tests, all passing

### Project State:
- **Tests:** All passing (93 total tests across all packages)
- **Coverage:** 55/102 issues closed (53.9%)
- **Working tree:** Clean, pushed to remote
- **New packages:**
  - internal/errors/collector.go
  - internal/services/graph/ (builder.go + tests)
  - internal/services/patcher/ (patcher.go + tests)

### Architecture Progress:
```
✅ Domain Layer: Complete (Node, Graph, Issue, Constraint, Ref, AuditEntry, DecoError)
✅ Error System: Complete (structure, codes, formatter, docs, YAML, suggestions, aggregation)
✅ Storage Layer: Complete (NodeRepository, ConfigRepository, HistoryRepository)
⏳ Service Layer: In Progress (3 of ~7 services complete)
   ✅ GraphBuilder: Build, dependencies, cycle detection, topological sort, reverse index
   ✅ Patcher: Set, append, unset operations
   ⏳ Remaining services needed:
      - QueryEngine (filter, search)
      - Validator orchestrator (schema, reference, constraint validation)
      - (Apply operation for Patcher - batch operations)
⏳ CLI Layer: Not started (depends on services)
```

## Next Steps (Recommended Priority):

### Immediate: Complete Remaining Service Layer

1. **Patcher Apply Operation** (deco-h9q/deco-7yp)
   - Apply batch patch operations from JSON/YAML
   - For AI integration (patch mode)
   - Should use existing Set/Append/Unset operations

2. **QueryEngine** (deco-3jm/deco-nop + deco-oam/deco-1eg)
   - Filter nodes by criteria (kind, status, tags, etc.)
   - Search by text (title, summary, content)
   - Tag-based queries
   - Integration with Graph for efficient lookups

3. **Validator Services** (deco-2t4/deco-az1 + implementations)
   - Schema validator: check required fields, types
   - Reference validator: verify all refs resolve
   - Constraint validator: CEL expression evaluation
   - Validator orchestrator: coordinate all validators, collect errors

### Then: CLI Layer
After all services are complete, the CLI commands can be implemented:

4. **Core Commands** (in dependency order):
   - deco init (initialize project)
   - deco validate (schema + refs + constraints)
   - deco list (list nodes)
   - deco show <id> (show node + reverse refs)
   - deco query <filter> (search/filter)
   
5. **Mutation Commands** (use Patcher):
   - deco set <id> <path> <value>
   - deco append <id> <path> <value>
   - deco unset <id> <path>
   - deco apply <patch-file>

6. **Advanced Commands**:
   - deco mv <old-id> <new-id> (with ref updates)
   - deco history [<id>] (audit log)

## Technical Notes:

### Patcher Design:
- Reflection-based field access via capitalizeFirst() helper
- Path parsing: splits by ".", extracts array indices from "field[N]"
- Array index notation: "tags[0]" -> field "Tags", index 0
- Map access: "glossary.term1" for nested map keys
- Type conversion: automatic when types are convertible
- Required fields: id, kind, version, status, title (cannot unset)

### GraphBuilder Design:
- Dependency map: nodeID -> []string (list of dependency IDs)
- Only Refs.Uses creates dependencies (Related does not)
- Cycle detection: DFS with visiting/visited state tracking
- Topological sort: Kahn's algorithm with in-degree counting
- Reverse index: targetID -> []string (nodes that reference target)

### Error Collector Design:
- Deduplication key: code + file + line + column
- Sort order: file (alpha), line (numeric), column (numeric)
- Nil locations come last in sort order
- Count tracks unique errors, not total attempts to add

### Testing Approach:
- Strict TDD: all tests written before implementation
- Comprehensive edge cases: empty, nil, invalid, boundaries
- Clear test names: TestComponent_Behavior pattern
- All 93 tests passing

## Ready to Work:
Run `bd ready` to see available tasks. Recommended next session:
1. Complete Patcher apply operation
2. Implement QueryEngine (filter + search)
3. Implement Validators (schema, reference, constraint, orchestrator)

This will complete the service layer and unblock all CLI commands.

---

<a id="deco-8vg-session-handover-storage-and-error-system-enhancements-complete"></a>

## 📋 deco-8vg Session handover: Storage and error system enhancements complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:53 |
| **Updated** | 2026-01-31 14:54 |
| **Closed** | 2026-01-31 14:54 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed Issues (10 total, bringing total to 41/101 - 40.6%):

**Storage Layer Complete:**
1. **deco-hxq/deco-owf**: HistoryRepository (tests + implementation)
   - Append-only JSONL format at .deco/history.jsonl
   - Thread-safe concurrent writes with mutex
   - Filtering by NodeID, Operation, User, time range, limit
   - 15 comprehensive tests, all passing

**Error System Enhancements:**
2. **deco-1jc/deco-nbv**: YAML line number tracking (tests + implementation)
   - LocationTracker using yaml.v3 Node with line/column info
   - Path-based location lookup (dot notation + array indexing)
   - Support for nested structures, multiline values
   - GetLocation() and GetValueLocation() methods
   - 13 comprehensive tests, all passing

3. **deco-mvn/deco-m9l**: YAML context extraction (tests + implementation)
   - ExtractContext() for surrounding lines
   - HighlightColumn() for visual error pointers
   - Edge case handling (start/end of file, CRLF)
   - Indentation preservation for accurate display
   - 15 comprehensive tests, all passing

4. **deco-4dj/deco-3ip**: Suggestion engine (tests + implementation)
   - Levenshtein distance calculation for typo detection
   - Smart suggestions with 'did you mean?' formatting
   - Threshold filtering, length similarity scoring
   - Prefix bonus for better matching
   - 17 comprehensive tests, all passing

**Epic Closed:**
5. **deco-elv**: Storage layer epic (all repository implementations complete)

### Project State:
- **Tests:** All passing (60 total tests across domain, errors, storage)
- **Coverage:** 41/101 issues closed (40.6%)
- **Working tree:** Clean, ready to commit
- **New packages:**
  - internal/errors/ (suggestions.go)
  - internal/errors/yaml/ (location.go, context.go)
  - internal/storage/history/ (jsonl_repository.go)

### Architecture Progress:
```
✅ Domain Layer: Complete (Node, Graph, Issue, Constraint, Ref, AuditEntry, DecoError)
✅ Error System: Enhanced with YAML integration
   ✅ Core: structure, registry, formatter, docs
   ✅ YAML: line tracking, context extraction
   ✅ Suggestions: Levenshtein-based typo suggestions
   ⏳ Remaining: error aggregation
✅ Storage Layer: Complete (all repositories implemented)
   ✅ NodeRepository: YAML + file discovery
   ✅ ConfigRepository: YAML
   ✅ HistoryRepository: JSONL append-only
⏳ Service Layer: Not started (ready to begin)
⏳ CLI Layer: Not started
```

## Next Steps (Recommended Priority):

### Immediate: Complete Error System
1. **deco-5wv/deco-5fv**: Error aggregation (tests + impl)
   - Collect and group related errors
   - Prevent error cascades
   - Integration with DecoError

### Then: Service Layer
After error system is complete, services are unblocked:

2. **deco-jxy/deco-0o3**: GraphBuilder service (tests + impl)
   - Build dependency graph from nodes
   - Detect cycles
   - Topological sort

3. **deco-3jm/deco-oam**: QueryEngine (tests + impl)
   - Filter nodes by criteria
   - Search by text
   - Tag-based queries

4. **deco-h9q/deco-7yp**: Patcher service (tests + impl)
   - Apply operations to nodes
   - Set, append, unset, move operations
   - Atomic updates

5. **deco-az1/deco-2y3**: Validator orchestrator (tests + impl)
   - Coordinate schema, reference, constraint validation
   - Collect and format errors
   - Integration with suggestion engine

### Then: CLI Layer
After services are ready, build the user-facing commands.

## Technical Notes:

### HistoryRepository Design:
- JSONL format (one JSON object per line)
- Thread-safe with sync.Mutex
- File location: .deco/history.jsonl
- Chronological ordering on query
- Efficient append-only operations

### YAML Error Integration:
- yaml.v3.Node preserves line/column info during parsing
- LocationTracker maps paths to source locations
- Context extraction for rich error display
- Column highlighting with ^ pointers

### Suggestion Engine:
- Levenshtein distance for typo detection
- Default threshold of 2 edits
- Prefix matching bonus reduces distance
- Length similarity as tiebreaker
- Top 3 suggestions returned

### Testing Approach:
- Strict TDD: all tests written before implementation
- Comprehensive edge case coverage
- Tests use t.TempDir() for isolation
- All 60 tests passing across packages

## Ready to Work:
Run `bd ready` to see available tasks. Next session should focus on completing error aggregation, then move to the service layer which is now unblocked.

---

<a id="deco-7de-session-handover-error-system-and-storage-layer-complete"></a>

## 📋 deco-7de Session handover: Error system and storage layer complete

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:37 |
| **Updated** | 2026-01-31 14:41 |
| **Closed** | 2026-01-31 14:41 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed Issues (12 total, bringing total to 31/100):

**Error System (Complete Core):**
1. **deco-03s**: DecoError structure with Rust-like pattern
   - Implemented Code, Summary, Detail, Location, Context, Suggestion, Related fields
   - Error() method for string representation

2. **deco-8sd/deco-7cp**: Error code registry
   - Defined E001-E099 error codes organized by category
   - Categories: schema, refs, validation, io, graph (E001-E019 per category)
   - Registry with Lookup, AllCodes, ByCategory, Categories methods

3. **deco-yu5/deco-e9t**: Error formatter
   - Rust-like formatting with ANSI colors
   - Source context display with line numbers and column pointers
   - FormatWithSource for showing code snippets
   - FormatMultiple for error aggregation

4. **deco-inq**: Error documentation generator
   - Markdown generation for error code documentation
   - GenerateMarkdown, GenerateDetailedMarkdown, GenerateCodeList methods
   - Filters out reserved error codes

**Storage Layer (NodeRepository + ConfigRepository):**
5. **deco-53k/deco-nkw**: YAML NodeRepository
   - Full CRUD: LoadAll, Load, Save, Delete, Exists
   - Nested directory support (e.g., systems/food.yaml)
   - Error handling for missing files, invalid YAML

6. **deco-344/deco-5ry**: File discovery for nodes
   - DiscoverAll for walking directory tree
   - DiscoverByPattern for filtered discovery
   - PathToID and IDToPath conversion utilities

7. **deco-77y/deco-px9**: ConfigRepository
   - Load/Save project configuration
   - Support for custom fields
   - Default values handling

### Project State:
- **Tests:** All passing (domain + storage layers)
- **Coverage:** 31/100 issues closed (31%)
- **Working tree:** Clean, all changes committed and pushed
- **Dependencies:** Added gopkg.in/yaml.v3 for YAML parsing

### Architecture Progress:
```
✅ Domain Layer: Complete (Node, Graph, Issue, Constraint, Ref, AuditEntry, DecoError)
✅ Error System: Core complete (structure, registry, formatter, docs)
   ⏳ Remaining: suggestion engine, error aggregation, YAML context extraction
✅ Storage Interfaces: Complete (NodeRepository, ConfigRepository, HistoryRepository)
✅ NodeRepository: Complete (YAML implementation + file discovery)
✅ ConfigRepository: Complete (YAML implementation)
⏳ HistoryRepository: Tests + implementation needed
⏳ Service Layer: Not started (depends on storage completion)
⏳ CLI Layer: Not started (depends on service + errors)
```

## Next Steps (Recommended Priority):

### Immediate: Complete Storage Layer
Continue with remaining repository implementations:

1. **deco-hxq/deco-owf**: HistoryRepository (tests + impl)
   - Append-only audit log
   - Query with filtering
   - JSONL format

2. **deco-1jc/deco-nbv**: YAML line number tracking (tests + impl)
   - Preserve source locations for error reporting
   - Integration with DecoError Location field

3. **deco-mvn/deco-m9l**: YAML context extraction (tests + impl)
   - Extract code context around error locations
   - Integration with error formatter

### Then: Complete Error System
Finish remaining error system components:

4. **deco-4dj/deco-3ip**: Suggestion engine (tests + impl)
   - Generate helpful suggestions based on error context
   - Integration with DecoError Suggestion field

5. **deco-5wv/deco-5fv**: Error aggregation (tests + impl)
   - Collect and group related errors
   - Prevent error cascades

### Then: Service Layer
After storage layer is complete, start on services:

6. **deco-jxy**: GraphBuilder service
7. **deco-0s1/deco-oam**: QueryEngine
8. **deco-7yp/deco-h9q**: Patcher service
9. **deco-az1/deco-2t4/deco-xw8/deco-eft**: Validator orchestrator

## Technical Notes:

### Error System Design:
- Following Rust compiler error format
- DecoError is the core error type used throughout
- ErrorFormatter provides rich, colorized output
- Error codes E001-E099 organized by category with room for expansion

### Storage Layer Design:
- Repository pattern for clean abstraction
- YAML as primary persistence format
- File-based organization: .deco/nodes/{id}.yaml, .deco/config.yaml, .deco/history.jsonl
- Support for nested directories matching node ID structure

### Testing Approach:
- Strict TDD: tests written before implementation
- All tests use t.TempDir() for isolation
- Comprehensive coverage of edge cases (missing files, invalid YAML, nested paths)

## Ready to Work:
Run `bd ready` to see available tasks. Next session should focus on completing the storage layer before moving to services.

---

<a id="deco-c5k-session-handover-foundation-complete-error-system-next"></a>

## 📋 deco-c5k Session handover: Foundation complete, error system next

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔥 Critical (P0) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:24 |
| **Updated** | 2026-01-31 14:36 |
| **Closed** | 2026-01-31 14:36 |

### Notes

## Session Accomplishments (2026-01-31)

### Completed (18 issues closed):
1. **deco-rba**: Initialized Go module and project structure
   - Created folder structure: cmd/deco/, internal/{domain,storage,service,cli}/
   - Enhanced .gitignore

2. **Domain Types (TDD)** - All tests written first, then implementations:
   - **deco-ora/deco-bqg**: Node domain type with validation
   - **deco-shb/deco-zty**: Graph with CRUD operations
   - **deco-0hm/deco-982**: Issue tracking (TBDs/questions)
   - **deco-ll2/deco-lwx**: Constraint (CEL validation rules)
   - **deco-xjv/deco-rrk**: Ref/RefLink (inter-node references)
   - **deco-1er/deco-8ca**: AuditEntry (audit log)

3. **Repository Interfaces**:
   - **deco-p54**: NodeRepository (LoadAll, Load, Save, Delete, Exists)
   - **deco-5ns**: ConfigRepository (Load, Save) + Config type
   - **deco-gpu**: HistoryRepository (Append, Query) + Filter type

4. **Error System (started)**:
   - **deco-16w**: DecoError tests written (Rust-like error structure)

5. **deco-16e**: Closed Foundation epic (all dependencies complete)

### Project State:
- All tests passing (domain layer complete)
- Clean architecture: domain types → repository interfaces → (next: implementations)
- 18/99 issues closed, 13 ready to work
- Working tree clean, all pushed to remote

## Next Steps (Recommended Priority):

### Immediate: Complete Error System
The error system tests are written (deco-16w ✓), next is implementation:
1. **deco-03s**: Define DecoError structure (BLOCKED by deco-16w, now ready)
   - Implement DecoError, Location, Related types
   - Implement Error() method
   - Make tests pass

Then continue with error system tasks (all now unblocked):
2. **deco-8sd/deco-7cp**: Error code registry (tests + impl)
3. **deco-yu5/deco-e9t**: Error formatter (tests + impl)
4. **deco-inq**: Error documentation generator

### Then: Storage Implementations
After error system, implement repository interfaces:
1. **deco-53k/deco-nkw**: YAML NodeRepository (tests + impl)
2. **deco-344/deco-5ry**: File discovery for nodes (tests + impl)
3. **deco-77y/deco-px9**: ConfigRepository (tests + impl)
4. **deco-hxq/deco-owf**: HistoryRepository (tests + impl)
5. **deco-1jc/deco-nbv**: YAML line number tracking (tests + impl)
6. **deco-mvn/deco-m9l**: YAML context extraction (tests + impl)

### Architecture Notes:
- Following strict TDD: tests first, then implementation
- Error system follows Rust-like pattern (Code, Summary, Detail, Location, Context, Suggestion, Related)
- Storage layer uses repository pattern for testability
- All domain types have validation methods

### Working Relationship Reminder:
- CEO: Strategic decisions, unblocking
- CTO (Claude): Task selection, execution, git management
- Staff: Subagents for coding tasks

---

<a id="deco-q5pk-validate-unknown-fields-in-reflink-and-other-nested-structures"></a>

## 🐛 deco-q5pk Validate unknown fields in RefLink and other nested structures

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-04 18:26 |
| **Updated** | 2026-02-04 18:30 |
| **Closed** | 2026-02-04 18:30 |

### Notes

Typos in nested structures like refs.uses[].context (e.g., 'conext') are silently ignored by Go's YAML unmarshaler. Strict field validation (E049) exists for blocks and top-level fields, but not for RefLink or other nested types. Reproduction: add 'conext: ...' instead of 'context: ...' in a refs.uses entry - deco validate passes. Fix: extend strict field checking to RefLink, Reviewer, and other nested structs, either via yaml.Decoder.KnownFields(true) or by walking raw YAML and comparing keys.

---

<a id="deco-snlr-add-strict-block-field-validation-configurable"></a>

## ✨ deco-snlr Add strict block field validation (configurable)

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-03 20:24 |
| **Updated** | 2026-02-03 20:55 |
| **Closed** | 2026-02-03 20:55 |

### Description

**Summary**
Enforce strict block field validation by default so unknown block keys are rejected (LLM safety + typo catching).

**Problem / Opportunity**
Blocks currently accept arbitrary fields (e.g., `asdasdasd: 200`), so typos and LLM hallucinations persist silently. This weakens the reliability of structured docs.

**Goals**
- Default-on strict validation that rejects unknown block fields for built-in and custom block types
- Default-on validation for table column objects (reject unknown column keys)
- Provide helpful suggestions for near-miss field names

**Non-goals**
- Auto-fixing or migrating existing data
- Changing block serialization format

**User Stories**
- As a designer using LLMs, I want unknown block fields rejected so invalid data doesn’t slip in.

**Scope**
- In scope: default strict behavior, allowlists for built-in blocks, optional fields for custom blocks, strict table column validation, validator errors, tests, docs
- Out of scope: auto-migration

**Constraints**
- Deterministic, stable error output
- Backward compatibility is not required (no external users yet)

**Proposed Approach**
1. Make strict block field validation the default behavior.
2. Define allowlists for built-in block types (always allow `id`):
   - `rule`: `text`
   - `table`: `columns`, `rows`
   - `param`: `name`, `datatype`, `min`, `max`, `default`, `unit`, `description`
   - `mechanic`: `name`, `description`, `conditions`, `outputs` (and optionally `inputs` if desired)
   - `list`: `items`
3. Add strict validation for table column objects (allow only `key`, `type`, `enum`, `display`).
4. Extend custom block config with `optional_fields` (allowlist extension) so custom blocks allow `required_fields + optional_fields + id`.
5. Unknown block field emits a validation error (use `E049` for unknown block field or introduce a dedicated code) and includes:
   - node/section/block location
   - suggestion via `errors.Suggester`
6. Add tests in `internal/services/validator/block_validator_test.go`.
7. Update `README.md` and `SPEC.md`.

**TDD Plan (Required)**
- Tests to write first (red): unknown block field fails; unknown table column key fails; allowed fields pass; custom optional fields pass
- Minimal implementation (green): allowlists + validator check + optional_fields support + column validation
- Refactor: centralize allowlists and suggestion logic

**Acceptance Criteria**
- [ ] Unknown block fields are rejected by default
- [ ] Table column objects reject unknown keys by default
- [ ] Custom block types accept `required_fields` and `optional_fields`
- [ ] Errors include location and suggestions

**Test Plan**
- `go test ./internal/services/validator -run Block`
- Cases: unknown key in param block, unknown key in table column, valid fields in mechanic, custom block with optional field

**Dependencies**
- If examples fail due to unknown fields, update or remove them (see `deco-23bx`)

**Files / Areas Touched**
- internal/services/validator/block_validator.go
- internal/storage/config (add `optional_fields`)
- internal/services/validator/validator.go (wire custom config)
- README.md, SPEC.md

**Risks**
- Allowlist too strict could block legitimate data

**Open Questions**
- Should we allow a per-project opt-out for strict columns?

---

<a id="deco-pc3o-content-hash-uses-non-deterministic-map-ordering-causing-hash-churn"></a>

## 🐛 deco-pc3o Content hash uses non-deterministic map ordering causing hash churn

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 16:06 |
| **Closed** | 2026-02-02 16:06 |

### Description

Content hash calculation uses YAML marshaling of maps (Glossary, Custom) without deterministic ordering. Go map iteration order is random, causing hash churn even when data is unchanged. This directly affects audit integrity and sync behavior.

Files:
- internal/cli/audit.go

---

<a id="deco-3o2h-apply-rewrite-don-t-write-content-hash-breaking-sync-detection"></a>

## 🐛 deco-3o2h apply/rewrite don't write content hash, breaking sync detection

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 16:10 |
| **Closed** | 2026-02-02 16:10 |

### Description

apply and rewrite commands don't write a content hash to the audit log, but sync relies on history hashes to detect manual edits. After a rewrite/apply, sync can falsely treat the node as a manual edit. This is a workflow correctness bug.

Files:
- internal/cli/apply.go
- internal/cli/rewrite.go
- internal/cli/sync.go

---

<a id="deco-36rv-cel-constraint-engine-missing-spec-d-capabilities-allnodes-custom-fields"></a>

## 🐛 deco-36rv CEL constraint engine missing spec'd capabilities (allNodes, custom fields)

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:40 |
| **Closed** | 2026-02-02 21:40 |

### Description

CEL constraints only see id/kind/version/status/title/tags. The allNodes variable is never populated, so cross-node constraints are impossible. Access to custom, content, or refs fields is not available. The SPEC example constraint (checking node references exist) cannot be implemented.

Files:
- internal/services/validator/validator.go

---

<a id="deco-nd7w-config-paths-nodes-path-history-path-are-hardcoded-not-configurable"></a>

## 🐛 deco-nd7w Config paths (nodes_path, history_path) are hardcoded, not configurable

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:36 |
| **Closed** | 2026-02-02 21:36 |

### Description

The spec says nodes_path and history_path are configurable in .deco/config.yaml, but the repositories are hardcoded to .deco/nodes and .deco/history.jsonl. This breaks the 'config is authoritative' premise.

Files:
- internal/storage/node/yaml_repository.go
- internal/storage/history/jsonl_repository.go
- internal/storage/config/repository.go

---

<a id="deco-h91g-sync-swallows-errors-can-exit-clean-despite-failures"></a>

## 🐛 deco-h91g sync swallows errors - can exit clean despite failures

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-01 21:58 |
| **Closed** | 2026-02-01 21:58 |

### Description

If node fails to load or sync/baseline fails, code prints warning and continues. Exit code can still be 'clean', hiding errors from CI.

Acceptance:
- Any failure to sync/baseline results in non-zero exit
- Track failures and signal at end
- Tests cover error accumulation

---

<a id="deco-0t5-cli-mutations-don-t-record-content-hash-in-history"></a>

## 🐛 deco-0t5 CLI mutations don't record content_hash in history

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-01 21:55 |
| **Closed** | 2026-02-01 21:55 |

### Description

Normal CLI commands (set, append, create, mv, etc.) don't write ContentHash to history. After any CLI edit, sync compares against stale hash and false-positives, bumping version + resetting review status unnecessarily.

Acceptance:
- All node-mutating commands record content_hash in history
- After CLI edit, deco sync immediately returns clean (exit 0)
- Shared helper for 'log entry with hash' to avoid drift

---

<a id="deco-jyn-deco-set-panics-on-nested-paths-pointer-deref-in-patcher"></a>

## 🐛 deco-jyn deco set panics on nested paths (pointer deref in patcher)

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:29 |
| **Updated** | 2026-02-01 17:44 |
| **Closed** | 2026-02-01 17:44 |

### Description

Repro:
1) Copy examples\snake to a scratch dir (or use any project)
2) Run: deco set systems/core content.sections[0].blocks[0].text "Test" <dir>

Actual: panic: "reflect: call of reflect.Value.FieldByName on ptr Value" with stack in internal/services/patcher/patcher.go.
Expected: command should succeed or return a structured error, never panic.
Likely cause: patcher.setValue recurses into pointer fields (Content *Content) without handling Ptr.
Refs: internal/services/patcher/patcher.go, internal/cli/set.go.

---

<a id="deco-86e-block-data-lost-when-parsing-yaml-blocks-data-null"></a>

## 🐛 deco-86e Block data lost when parsing YAML (blocks[].data null)

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:29 |
| **Updated** | 2026-02-01 17:49 |
| **Closed** | 2026-02-01 17:49 |

### Description

Repro:
1) cd examples\snake
2) deco show systems/core --json
3) content.sections[*].blocks[*].data is null even though YAML defines text/columns/rows.

Expected: block fields are preserved (either inline fields mapped into data, or schema uses explicit data: ... and docs/examples follow).
Actual: unknown block keys are dropped by yaml.Unmarshal, so content is lost and re-saving would strip it.
Refs: internal/domain/node.go (Block has Data only); README example and examples/snake use inline block fields.

### Dependencies

- ⛔ **blocks**: `deco-cxk`

---

<a id="deco-5xa-errors-rust-like-error-system"></a>

## 🚀 deco-5xa Errors: Rust-like error system

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:57 |

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-03s`
- ⛔ **blocks**: `deco-7cp`
- ⛔ **blocks**: `deco-e9t`
- ⛔ **blocks**: `deco-m9l`
- ⛔ **blocks**: `deco-3ip`
- ⛔ **blocks**: `deco-5fv`
- ⛔ **blocks**: `deco-inq`

---

<a id="deco-7yo-cli-command-implementations"></a>

## 🚀 deco-7yo CLI: Command implementations

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:12 |

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-7au`
- ⛔ **blocks**: `deco-7ba`
- ⛔ **blocks**: `deco-1fm`
- ⛔ **blocks**: `deco-qll`
- ⛔ **blocks**: `deco-4q8`
- ⛔ **blocks**: `deco-1eg`
- ⛔ **blocks**: `deco-53i`
- ⛔ **blocks**: `deco-1v8`
- ⛔ **blocks**: `deco-31k`
- ⛔ **blocks**: `deco-vzo`
- ⛔ **blocks**: `deco-6j1`
- ⛔ **blocks**: `deco-0kr`
- ⛔ **blocks**: `deco-c5c`

---

<a id="deco-t6q-services-business-logic-layer"></a>

## 🚀 deco-t6q Services: Business logic layer

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:29 |

### Notes

Core services complete: GraphBuilder, Patcher (set/append/unset/apply), QueryEngine (filter/search), Validators (schema/reference/constraint/orchestrator). Node rename (deco-142) deferred to Phase 4 as advanced CLI feature. Service layer is ready for CLI implementation.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-0o3`
- ⛔ **blocks**: `deco-pqf`
- ⛔ **blocks**: `deco-bx5`
- ⛔ **blocks**: `deco-jna`
- ⛔ **blocks**: `deco-9ls`
- ⛔ **blocks**: `deco-2y3`
- ⛔ **blocks**: `deco-7rp`
- ⛔ **blocks**: `deco-k3j`
- ⛔ **blocks**: `deco-750`
- ⛔ **blocks**: `deco-7yp`
- ⛔ **blocks**: `deco-0s1`
- ⛔ **blocks**: `deco-oam`
- ⛔ **blocks**: `deco-142`

---

<a id="deco-elv-storage-repository-implementations"></a>

## 🚀 deco-elv Storage: Repository implementations

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:50 |

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-nkw`
- ⛔ **blocks**: `deco-5ry`
- ⛔ **blocks**: `deco-px9`
- ⛔ **blocks**: `deco-owf`
- ⛔ **blocks**: `deco-nbv`

---

<a id="deco-16e-foundation-project-setup-and-domain-types"></a>

## 🚀 deco-16e Foundation: Project setup and domain types

| Property | Value |
|----------|-------|
| **Type** | 🚀 epic |
| **Priority** | ⚡ High (P1) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:21 |

### Dependencies

- ⛔ **blocks**: `deco-rba`
- ⛔ **blocks**: `deco-bqg`
- ⛔ **blocks**: `deco-zty`
- ⛔ **blocks**: `deco-982`
- ⛔ **blocks**: `deco-lwx`
- ⛔ **blocks**: `deco-rrk`
- ⛔ **blocks**: `deco-8ca`
- ⛔ **blocks**: `deco-p54`
- ⛔ **blocks**: `deco-5ns`
- ⛔ **blocks**: `deco-gpu`

---

<a id="deco-qw2b-validate-examples-snake-with-deco-validate"></a>

## 📋 deco-qw2b Validate examples/snake with deco validate

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-03 20:59 |
| **Updated** | 2026-02-03 21:08 |
| **Closed** | 2026-02-03 21:08 |

### Notes

Build CLI and run deco validate in examples/snake; confirm strict block validation works with intentional break/fix.

---

<a id="deco-hond-set-resets-status-to-draft-but-append-unset-don-t-inconsistent-workflow"></a>

## 🐛 deco-hond set resets status to draft but append/unset don't - inconsistent workflow

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:44 |
| **Closed** | 2026-02-02 21:44 |

### Description

The set command resets approved/review status to draft, but append and unset don't apply the same logic. Same underlying edit operation, different workflow behavior. This inconsistency can cause unexpected state.

Files:
- internal/cli/set.go
- internal/cli/append.go
- internal/cli/unset.go

---

<a id="deco-122q-ref-validation-ignores-emits-events-and-vocabulary-fields"></a>

## 🐛 deco-122q Ref validation ignores emits_events and vocabulary fields

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:42 |
| **Closed** | 2026-02-02 21:42 |

### Description

The spec includes emits_events and vocabulary in refs, but only uses/related are checked for existence. This means invalid references in these fields won't be caught.

Files:
- internal/services/validator/validator.go

---

<a id="deco-epm4-readme-block-examples-use-key-value-but-validator-requires-name-datatype"></a>

## 🐛 deco-epm4 README block examples use key/value but validator requires name/datatype

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:51 |
| **Closed** | 2026-02-02 21:51 |

### Description

README's param block example uses key/value fields, but the block validator requires name/datatype. This spec drift will confuse users and cause unexpected validation failures.

Files:
- README.md
- internal/services/validator/block_validator.go

---

<a id="deco-7z6u-yaml-error-location-infrastructure-exists-but-unused-in-validate-output"></a>

## 🐛 deco-7z6u YAML error location infrastructure exists but unused in validate output

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 21:49 |
| **Closed** | 2026-02-02 21:49 |

### Description

There's a YAML location tracker in internal/errors/yaml/location.go, but validate output never includes line/column information. The README shows line-based errors which creates a UX mismatch.

Files:
- internal/errors/yaml/location.go
- internal/cli/validate.go

---

<a id="deco-chwx-status-validation-only-checks-presence-not-allowed-values"></a>

## 🐛 deco-chwx Status validation only checks presence, not allowed values

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 15:49 |
| **Updated** | 2026-02-02 16:11 |
| **Closed** | 2026-02-02 16:11 |

### Description

The spec and README define allowed statuses (draft, review, approved, deprecated, archived), but validation only checks that status field is present, not that it contains a valid value. This undermines workflow guarantees.

Files:
- internal/services/validator/validator.go

---

<a id="deco-fgs5-refactor-extend-mv-to-update-all-reference-types"></a>

## ✨ deco-fgs5 refactor: extend mv to update all reference types

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 22:15 |
| **Updated** | 2026-02-02 21:54 |
| **Closed** | 2026-02-02 21:54 |

### Description

Currently `deco mv` only updates refs.uses[].target and refs.related[].target.

It does NOT update:
- refs.emits_events (event node IDs)
- refs.vocabulary (glossary node IDs)  
- IDs mentioned in summary, content blocks, custom, or other free-text fields

The README claims "Rename nodes and all references update automatically" which overstates what's actually implemented.

Acceptance:
- mv updates refs.emits_events and refs.vocabulary
- Consider: optional text search/replace for ID mentions in free-text fields (with user confirmation)
- Update README to accurately reflect capabilities

---

<a id="deco-rb08-sync-detect-manual-file-renames-and-update-references"></a>

## ✨ deco-rb08 sync: detect manual file renames and update references

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 22:14 |
| **Updated** | 2026-02-02 22:08 |
| **Closed** | 2026-02-02 22:08 |

### Description

When users manually rename a node file (bypassing `deco mv`), sync should detect this and offer to update references.

Currently sync only detects content changes via hash comparison. It doesn't detect:
- File renames (node-a.yaml → node-b.yaml)
- ID changes within a file

Possible approach:
- Track file paths in history alongside node IDs
- On sync, detect when a file exists but its ID doesn't match history
- Offer to run refactor/reference update for the rename

This would make the README claim "Refactorable" more robust - users wouldn't be locked into using `deco mv`.

### Notes

Context: deco mv currently only updates refs.uses and refs.related, not all reference types. See deco-fgs5 for extending refactor coverage first.

This issue is specifically about detecting manual file renames during sync and triggering the refactor service (whatever it supports at that time).

### Dependencies

- ⛔ **blocks**: `deco-fgs5`

---

<a id="deco-oh0q-dry-run-always-exits-0-even-when-changes-would-occur"></a>

## 🐛 deco-oh0q --dry-run always exits 0 even when changes would occur

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-01 22:00 |
| **Closed** | 2026-02-01 22:00 |

### Description

Dry-run unconditionally returns syncExitClean, making it useless for CI drift checks.

Acceptance:
- Dry-run returns exit 1 when changes would be made
- Tests cover dry-run exit codes

---

<a id="deco-cxl-sync-o-nodes-history-performance-needs-single-pass-indexing"></a>

## 📋 deco-cxl sync O(nodes × history) performance - needs single-pass indexing

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-01 22:06 |
| **Closed** | 2026-02-01 22:06 |

### Description

sync queries history per node via JSONL scan. Runtime is O(nodes × history_entries). Large projects will be slow.

Acceptance:
- sync reads history once per run
- Add 'latest hash by node' scan/cache
- Complexity reduced to O(history + nodes)

---

<a id="deco-8xo-content-hash-excludes-glossary-contracts-llmcontext-constraints-custom-kind"></a>

## 🐛 deco-8xo Content hash excludes Glossary, Contracts, LLMContext, Constraints, Custom, Kind

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-01 22:02 |
| **Closed** | 2026-02-01 22:02 |

### Description

contentFields struct only covers title, summary, tags, refs, issues, content. Manual edits to Glossary, Contracts, LLMContext, Constraints, Custom, Kind go undetected by sync.

Acceptance:
- Decide and document precise definition of 'content'
- Either include these fields in hash OR document exclusions
- Tests cover newly included fields

### Notes

WIP: Added Kind, Glossary, Contracts, LLMContext, Constraints, Custom to contentFields. Tests pass. Ready to close - just needs final commit.

---

<a id="deco-mbz-add-deco-sync-command-for-detecting-and-fixing-unversioned-changes"></a>

## ✨ deco-mbz Add deco sync command for detecting and fixing unversioned changes

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 20:32 |
| **Updated** | 2026-02-01 20:56 |
| **Closed** | 2026-02-01 20:56 |

### Description

## Summary
Add a `deco sync` command that detects manually-edited nodes and fixes their metadata before committing. Used as a pre-commit hook.

## Problem
When LLMs or humans edit YAML files directly, they bypass the CLI's version bumping and review reset logic. This leaves nodes with stale approvals.

## Solution

### Core Flow
1. Find node files modified since HEAD (`git diff --name-only HEAD`)
2. For each modified .yaml in .deco/nodes/:
   - Parse current file (working tree)
   - Parse HEAD version (`git show HEAD:<path>`)
   - Compare semantically (ignore formatting)
   - If content differs: bump version, reset to draft, clear reviewers, log to history
3. Output synced nodes with version changes
4. Exit 1 if files modified (so commit aborts, user re-commits)

### Content Fields (trigger sync)
title, summary, content, tags, refs, issues

### Metadata Fields (ignored)
version, status, reviewers, kind, id

### CLI Interface
```
deco sync [directory]
  --dry-run    Show what would change
  --quiet      Suppress output
```

### Exit Codes
- 0: No changes needed
- 1: Files modified, re-commit needed  
- 2: Error (not git repo, invalid nodes)

### Output
- `Synced: sword-001 (v2→v3), shield-002 (v1→v2)`
- Dry-run: `Would sync: sword-001 (v2→v3)`
- Nothing to sync: silent, exit 0

### History Entry
```json
{"operation": "sync", "node_id": "...", "before": {"version": 2, "status": "approved"}, "after": {"version": 3, "status": "draft"}}
```

## Acceptance Criteria
- [ ] Detects nodes modified since HEAD (semantic comparison)
- [ ] Bumps version for modified nodes
- [ ] Resets approved/review status to draft
- [ ] Clears reviewers on modified nodes
- [ ] Logs sync operation to history
- [ ] --dry-run shows changes without applying
- [ ] Exit codes: 0 (clean), 1 (synced), 2 (error)
- [ ] Works in pre-commit hook workflow

### Design

## Design Decisions (from brainstorming)

### Detection
- Compare working tree against Git HEAD
- Semantic comparison only (ignore formatting)
- Ignore new and deleted files - only sync modified nodes

### Content Fields (trigger sync when changed)
- title, summary, content, tags, refs, issues

### Metadata Fields (ignored in comparison)
- version, status, reviewers, kind, id

### Actions on Modified Node
1. Bump version (always engine-driven)
2. Reset status to "draft" (if was approved/review)
3. Clear reviewers array
4. Log "sync" operation to history

### CLI Interface
```
deco sync [directory]
  --dry-run    Show what would change
  --quiet      Suppress output
```

### Exit Codes
- 0: No changes needed
- 1: Files modified, re-commit needed
- 2: Error (not git repo, invalid nodes)

### Output
- Normal: `Synced: sword-001 (v2→v3), shield-002 (v1→v2)`
- Dry-run: `Would sync: sword-001 (v2→v3)`
- Nothing: (silent, exit 0)

### Pre-commit Hook
Auto-fix and exit 1 so user re-commits with fixes included.

---

<a id="deco-zn8-add-audit-history-to-all-node-modifying-commands"></a>

## 📋 deco-zn8 Add audit history to all node-modifying commands

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:34 |
| **Updated** | 2026-02-01 18:46 |
| **Closed** | 2026-02-01 18:46 |

### Description

Commands that modify nodes should record before/after state in the audit log so deco diff can show changes.

### Notes

## Commands needing history recording

Already have it:
- rm (delete)
- mv (move)

Need to add:
- create (record after state)
- set (record before/after for changed field)
- append (record before/after for array)
- unset (record before state of removed field)
- apply (record before/after for each change)

## Implementation pattern

From rm.go:
```go
entry := domain.AuditEntry{
    Timestamp: time.Now(),
    NodeID:    node.ID,
    Operation: "set",
    User:      username,
    Before:    map[string]interface{}{...},
    After:     map[string]interface{}{...},
}
historyRepo.Append(entry)
```

## Acceptance criteria

- All modifying commands write to .deco/history.jsonl
- Before/After contain only the changed fields (not full node)
- deco diff shows meaningful output for all operations

---

<a id="deco-p75-backwards-compatible-schema-migrations"></a>

## ✨ deco-p75 Backwards-compatible schema migrations

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Assignee** | @claude |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 15:47 |

### Description

Goal: Evolve deco's schema without breaking existing GDD projects.

## User Stories
1. **Deco upgrades**: When users install a new deco version with schema changes, their existing GDD projects should continue to work or be cleanly migrated.
2. **Team evolution**: As teams refine their GDD structure, they can update the schema and migrate all nodes in one operation.

## Problem
Currently, if we change:
- Required fields in domain.Node
- Block type structure
- Config format
- Validation rules

...existing projects break with no migration path.

## Scope
Migration support for:
- [x] Adding new required fields (with defaults or prompts)
- [x] Renaming fields (mapping old -> new)
- [x] Restructuring node hierarchies (e.g., moving nested fields)
- [x] Changing validation rules (may require data transformation)

## Technical Context
- Config already has 'version: 1' field (internal/storage/config/repository.go:30)
- Content hashing uses SHA-256 truncated to 64 bits (internal/cli/audit.go:30-58)
- Nodes have Version field for edit tracking, separate from schema version
- Audit log tracks changes with ContentHash for integrity

## Proposed Implementation

### Schema Version Tracking
- Add 'schema_version' to config (distinct from 'version' which is format version)
- Store as hash of schema definition, not sequential number
- Compute hash from: required fields, block types, validation rules

### Migration Detection
- 'deco validate' checks current schema hash vs project's stored hash
- If mismatch: warn and recommend 'deco migrate'
- Exit code distinguishes: 0=valid, 1=invalid, 2=needs-migration

### Migration Execution
- 'deco migrate' command (whole-project, not per-file)
- Loads all nodes, applies transformation functions, writes back
- Creates backup before migration
- Logs migration as audit entries (new 'migrate' operation type)

### Migration Definitions
- migrations/ package with versioned transformations
- Each migration has: source_hash, target_hash, transform functions
- Chain migrations if jumping multiple versions

## CLI Interface
```
deco migrate              # Migrate project to current schema
deco migrate --dry-run    # Show what would change
deco migrate --backup     # Create backup before migrating (default: true)
```

## Acceptance Criteria
- [ ] Schema version tracked in config (hash-based)
- [ ] 'deco validate' detects schema mismatch, recommends migration
- [ ] 'deco migrate' command transforms all nodes in project
- [ ] Migration creates backup by default
- [ ] Audit log records migration operations
- [ ] At least one test migration (e.g., adding a new required field)

## Non-Goals
- Rollback/downgrade support (out of scope per discussion)
- Per-file migration (always whole-project)
- Automatic migration without user action (explicit command required)
- GUI or interactive migration wizard

### Notes

CLARIFICATIONS (from review):

1. **Schema Hash Canonicalization**: Hash computed from deterministic JSON encoding of schema definition (sorted keys, no whitespace). Fields ordered: required_fields[], block_types[], validation_rules[]. Each sorted alphabetically.

2. **Migration Mapping Scheme**: migrations/ package contains files named `<source_hash>_to_<target_hash>.go`. Registry maps source→target hashes. When chaining, resolver finds shortest path from current to target.

3. **Exit Code Compatibility**: Exit code 2 (needs-migration) is new for deco migrate ecosystem. Other commands (validate, query) unaffected. Matrix: 0=success, 1=validation_error, 2=migration_needed (migrate only).

4. **Backup Behavior**: --backup (default true) creates `.deco/backup-<timestamp>/` before migration. --no-backup=true disables. Existing backup paths are preserved (timestamped, never overwritten).

---

<a id="deco-a5l-ai-patch-rewrite-safety-validate-gate-transactional-apply-explicit-diff"></a>

## ✨ deco-a5l AI patch/rewrite safety: validate gate + transactional apply + explicit diff

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 14:37 |

### Description

Goal: AI changes are safe, validated, and reviewable before being written.

## User Stories
1. **AI-assisted editing**: When an LLM generates patches/rewrites, humans can review the exact changes before they're applied.
2. **Validation safety**: Invalid changes are rejected before corrupting the GDD, even if the AI is confident.
3. **Audit trail**: All AI changes are logged with clear before/after states for debugging.

## Problem
Currently:
- `deco apply` writes changes without validating the resulting node
- `--dry-run` only says 'N operations would be applied' - doesn't show actual diff
- No dedicated 'full rewrite' command for when AI rewrites entire nodes
- If AI produces invalid YAML/schema, corruption happens before validation

## Current State (apply.go)
Already implemented:
- ✅ Transactional rollback on operation failure (line 144-170 in tests)
- ✅ `--dry-run` flag (but limited output)
- ✅ Version increment after apply
- ✅ History logging with before/after values

Missing:
- ❌ Schema validation before write
- ❌ Diff output showing actual field changes
- ❌ Full rewrite support (replace entire node content)

## Proposed Implementation

### 1. Validation Gate
Add validation step between apply and save:
```go
// In runApply(), after p.Apply() succeeds:
orchestrator := validator.NewOrchestrator(...)
if err := orchestrator.ValidateNode(n); err != nil {
    return fmt.Errorf("patch would create invalid node: %w", err)
}
// Only then save
```

### 2. Explicit Diff Output
Enhance `--dry-run` to show field-level diff:
```
$ deco apply sword-001 patch.json --dry-run

Proposed changes to sword-001:
  title: "Iron Sword" → "Enchanted Blade"
  tags: [weapon, combat] → [weapon, combat, magic]
  + summary: "A blade imbued with fire magic"

Validation: ✓ Valid
Run without --dry-run to apply.
```

### 3. Rewrite Command
New `deco rewrite` for full node replacement:
```bash
deco rewrite <node-id> <new-content.yaml>  # Replace entire node
deco rewrite <node-id> <new-content.yaml> --dry-run  # Show full diff
```

Rewrite vs Apply:
- Apply: Surgical patches (set/append/unset operations)
- Rewrite: Full replacement (when AI rewrites entire node)

### 4. Diff Library
Use go-diff or similar for readable diffs:
- Unified diff format for large changes
- Inline format for small changes
- Color output when terminal supports it

## CLI Interface
```
deco apply <id> <patch.json>           # Apply patch with validation
deco apply <id> <patch.json> --dry-run # Show diff + validation result
deco rewrite <id> <file.yaml>          # Replace node with validation
deco rewrite <id> <file.yaml> --dry-run # Show full diff
deco diff <id> <file.yaml>             # Just show diff without applying
```

## Acceptance Criteria
- [ ] `deco apply` validates resulting node before write; aborts if invalid
- [ ] `--dry-run` shows field-level diff of proposed changes
- [ ] `--dry-run` shows validation result (valid/invalid with errors)
- [ ] New `deco rewrite` command for full node replacement
- [ ] Rewrite validates before write
- [ ] Both commands log to audit history with full before/after

## Edge Cases
- Patch creates valid intermediate state but invalid final state (multi-op)
- Rewrite changes node ID (should be rejected or handled specially)
- Concurrent edit detection (out of scope - see deco-rmv)

## Non-Goals
- Interactive diff approval (CLI only, not TUI)
- Merge/conflict resolution (separate feature)
- AI prompt engineering (just the safety layer)

### Notes

CLARIFICATIONS (from review):

1. **Transactional Rollback**: Production implementation exists in apply.go:144-170 (tested in apply_test.go). If any operation fails mid-batch, prior mutations are not persisted. Tests verify this behavior.

2. **Scope Alignment**: `deco diff` removed from CLI Interface section - it's convenience sugar, not core. Acceptance criteria unchanged. May add in future PR if needed.

3. **Rewrite Format Specification**: 
   - Input file MUST be valid YAML
   - MUST include schema_version field matching project's current schema
   - MUST include all required node fields (id, title, status, kind, version)
   - File is fully validated before any write occurs

---

<a id="deco-79k-enforce-constraint-scope-node-kind-pattern-in-validator"></a>

## ✨ deco-79k Enforce constraint scope (node kind/pattern) in validator

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 19:37 |
| **Closed** | 2026-02-01 19:37 |

### Description

Goal: respect Constraint.scope for CEL constraints.

Acceptance:
- Constraints apply only to matching node kinds/patterns.
- Non-matching nodes are skipped.
- Tests cover scope matching rules.

---

<a id="deco-87y-custom-block-types-with-validation-hooks"></a>

## ✨ deco-87y Custom block types with validation hooks

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 09:00 |

### Description

Goal: extend content blocks beyond core types.

Acceptance:
- Config defines custom block types + required fields.
- Validator enforces per custom type.
- Example demonstrates a custom block.

---

<a id="deco-apc-configurable-schema-rules-org-level-constraints"></a>

## ✨ deco-apc Configurable schema rules (org-level constraints)

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 13:53 |

### Description

Goal: allow org/project to define additional schema rules.

Acceptance:
- Project config can declare required fields per kind.
- Validation enforces these rules.
- Documented in SPEC.

---

<a id="deco-rs0-review-workflow-approvals-status-transitions-changelog-notes"></a>

## ✨ deco-rs0 Review workflow: approvals, status transitions, changelog notes

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 20:23 |
| **Closed** | 2026-02-01 20:23 |

### Description

Goal: formalize draft -> approved/published workflow.

Needs:
- Review approvals (at least 1 required) and reviewer metadata.
- Guarded status transitions with rules.
- Changelog notes stored with history.

Acceptance:
- CLI supports approve/reject and status transitions.
- Validation enforces required approvals for approved/published.
- History shows review notes.

### Design

## Design Decisions

### Use Case
- Team collaboration with multiple approvers
- Audit trail tracking who blessed what and when
- Solo-friendly (works with single approver)
- Versioned: each version cycles through states independently

### Status Flow
`draft` → `review` → `approved`

- **draft**: Work in progress
- **review**: Ready for review (submitted for approval)
- **approved**: Has required approvals for this version

### Edit Behavior
- Any edit bumps version and **auto-resets status to draft**
- Requires fresh approval cycle
- History preserves what was approved before the edit

### Approval Requirements
- Configurable per-project in `.deco/config.yaml`
- Example: `required_approvals: 2`
- Default: 1 (solo-friendly)

### Approval Data Structure
Each approval includes:
- Reviewer name/email
- Timestamp
- Optional note/comment
- Version number approved

### Storage
- **Node YAML**: `reviewers` field shows current approvals for this version
- **History**: Full audit trail of all approvals across versions

### CLI Commands
```
deco review submit <node-id>     # draft → review
deco review approve <node-id>    # add approval (with optional --notes)
deco review reject <node-id>     # review → draft (requires --notes)
deco review status <node-id>     # show approval state
```

### Implementation Notes
- Add `reviewers` field to Node struct
- Add `required_approvals` to config.yaml schema
- Add ApprovalValidator to check approval count for status transitions
- Add "approve", "reject", "submit" operations to audit.go
- Create internal/cli/review.go with subcommands

### Notes

## Implementation Progress (2026-02-01 - Session 2)

**Completed Tasks (6/13):**
1. ✅ Add Reviewer struct and Reviewers field to Node (a050b62)
2. ✅ Add required_approvals to Config with default of 1 (caa154e)
3. ✅ Add submit, approve, reject operations to Audit (720cbbf)
4. ✅ Create ApprovalValidator (1f8f014)
5. ✅ Integrate ApprovalValidator into Orchestrator (d7551ff)
6. ✅ Update CLI validate command to use config (498c071)

**Remaining Tasks (7/13):**
7. Create review.go CLI with submit subcommand
8. Add approve subcommand to review
9. Add reject subcommand to review
10. Add status subcommand to review
11. Register review command in main.go
12. Auto-reset status on edit
13. Final integration test

**Plan file:** docs/plans/2026-02-01-review-workflow.md

All tests passing. Core domain + validator done. CLI commands remain.

---

<a id="deco-e4g-expand-issues-tbd-system-filters-severity-per-node-tracking"></a>

## ✨ deco-e4g Expand issues/TBD system: filters, severity, per-node tracking

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 19:41 |
| **Closed** | 2026-02-01 19:41 |

### Description

Goal: make issues/TBDs actionable at scale.

Needs:
- Filters by severity/status/kind/tag/node.
- Per-node issue listing and rollups.
- Clear severity semantics (low/med/high/critical).

Acceptance:
- `deco issues` supports filters and per-node queries.
- Output supports quiet/json formats.
- Docs updated.

---

<a id="deco-603-add-type-specific-validation-for-blocks-rule-table-param-etc"></a>

## ✨ deco-603 Add type-specific validation for blocks (rule/table/param/etc.)

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:32 |
| **Updated** | 2026-02-01 19:18 |
| **Closed** | 2026-02-01 19:18 |

### Description

Goal: Prevent structurally invalid blocks from passing validation.

Examples:
- rule blocks require text
- table blocks require columns and rows
- param blocks require name/datatype/min/max as applicable

Acceptance:
- Validator enforces required fields per block type.
- Errors are clear and point to the offending block.
- Tests cover each block type.

### Dependencies

- ⛔ **blocks**: `deco-cxk`

---

<a id="deco-hsx-require-content-for-approved-published-nodes-allow-drafts-without-content"></a>

## ✨ deco-hsx Require content for approved/published nodes (allow drafts without content)

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:32 |
| **Updated** | 2026-02-01 19:09 |
| **Closed** | 2026-02-01 19:09 |

### Description

As a GDD author, I want strictness to increase with status.

Acceptance:
- status=draft: content optional.
- status=approved or published: content (and at least one section) required.
- Validation error should be clear and actionable.

---

<a id="deco-cxk-decide-and-implement-block-schema-inline-fields-vs-data-map"></a>

## 📋 deco-cxk Decide and implement block schema (inline fields vs data map)

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:32 |
| **Updated** | 2026-02-01 17:49 |
| **Closed** | 2026-02-01 17:49 |

### Description

Need a clear canonical block representation to prevent data loss.

Decision to make:
- Keep inline block fields (text/columns/rows/etc.) as canonical and map them into Block.Data, OR
- Require explicit `data:` for all blocks and update docs/examples accordingly.

Acceptance:
- Chosen schema is implemented consistently in load/save.
- README/examples updated to match the schema.
- No silent dropping of block fields.

### Notes

Decision: Keep inline block fields as canonical. Implement custom YAML unmarshaler for Block to capture all fields (except 'type') into Data map.

---

<a id="deco-rib-define-strict-top-level-schema-with-explicit-extension-mechanism"></a>

## ✨ deco-rib Define strict top-level schema with explicit extension mechanism

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:32 |
| **Updated** | 2026-02-01 18:06 |
| **Closed** | 2026-02-01 18:06 |

### Description

Goal: Reject unknown top-level keys during validation while still allowing intentional extensions.

Proposal:
- Validation should fail on unknown root keys by default.
- Provide an explicit escape hatch: e.g. `custom:` map or `x_*` namespace, or a `--allow-unknown` flag.
- Document the extension mechanism in README/SPEC.

Acceptance:
- Unknown top-level key fails validation in default mode.
- Keys under the extension namespace pass.
- Docs updated to explain how to extend safely.

---

<a id="deco-0ja-deco-set-cannot-update-non-string-fields"></a>

## 🐛 deco-0ja deco set cannot update non-string fields

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:30 |
| **Updated** | 2026-02-01 17:59 |
| **Closed** | 2026-02-01 17:59 |

### Description

Repro:
1) Run: deco set systems/core version 2 <dir>

Actual: Error "cannot convert string to int".
Expected: CLI should parse numeric/boolean/JSON values based on target field type (or allow a --json flag) so common fields like version can be updated.
Refs: internal/cli/set.go, internal/services/patcher/patcher.go.

---

<a id="deco-a3o-validate-ignores-unknown-top-level-keys-and-structural-typos"></a>

## 🐛 deco-a3o validate ignores unknown top-level keys and structural typos

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:30 |
| **Updated** | 2026-02-01 18:07 |
| **Closed** | 2026-02-01 18:07 |

### Description

Repro:
1) In a node YAML, change `content:` to `contentx:` or add a random top-level key.
2) Run: deco validate <dir>

Actual: validation passes.
Expected: validation should fail (or at least warn) for unknown top-level keys / missing known sections, to prevent silent corruption.
Example: examples/snake/.deco/nodes/systems/core.yaml has `asdasdasd: 200` under a block and validate passes.
Refs: internal/services/validator/validator.go.

### Dependencies

- ⛔ **blocks**: `deco-rib`

---

<a id="deco-8ov-validate-does-not-detect-duplicate-node-ids"></a>

## 🐛 deco-8ov validate does not detect duplicate node IDs

| Property | Value |
|----------|-------|
| **Type** | 🐛 bug |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:30 |
| **Updated** | 2026-02-01 17:46 |
| **Closed** | 2026-02-01 17:46 |

### Description

Repro:
1) Copy any node file to a new filename but keep the same id field (e.g., systems/core.yaml -> systems/core-dup.yaml with id: systems/core).
2) Run: deco validate <dir>

Actual: validates successfully.
Expected: validation should fail because node IDs must be unique; duplicates cause ambiguous references and data loss.
Refs: internal/services/validator/validator.go.

---

<a id="deco-3dy-set-up-ci-cd-pipeline"></a>

## 📋 deco-3dy Set up CI/CD pipeline

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-01-31 16:53 |
| **Closed** | 2026-01-31 16:53 |

### Description

Set up GitHub Actions CI/CD pipeline for deco. Should include: run tests on PR/push, build binaries for multiple platforms (Linux, macOS, Windows), create releases with binaries on tag, run deco validate on example projects.

---

<a id="deco-tzg-write-tests-for-deco-apply-command"></a>

## 📋 deco-tzg Write tests for deco apply command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:57 |

### Description

TDD: Write failing tests for apply command. Test: applies patch file, validates after apply, --dry-run flag, error handling.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-7kj-write-tests-for-deco-history-command"></a>

## 📋 deco-7kj Write tests for deco history command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:59 |

### Description

TDD: Write failing tests for history command. Test: shows all history, --node filter, --limit flag, output format.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-747-write-tests-for-deco-mv-command"></a>

## 📋 deco-747 Write tests for deco mv command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:10 |

### Description

TDD: Write failing tests for mv command. Test: renames node, updates all refs, moves file, creates history entry.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-3vb-write-tests-for-deco-unset-command"></a>

## 📋 deco-3vb Write tests for deco unset command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:56 |

### Description

TDD: Write failing tests for unset command. Test: removes field, validates after unset, error on required field.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-zix-write-tests-for-deco-append-command"></a>

## 📋 deco-zix Write tests for deco append command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:54 |

### Description

TDD: Write failing tests for append command. Test: appends to array, validates after append, error on non-array.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-bhb-write-tests-for-node-rename"></a>

## 📋 deco-bhb Write tests for node rename

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:05 |

### Description

TDD: Write failing tests for rename operation. Test: ID change, all refs updated, file moved, history entry created.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-0jp-write-tests-for-deco-set-command"></a>

## 📋 deco-0jp Write tests for deco set command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:53 |

### Description

TDD: Write failing tests for set command. Test: sets field value, validates after set, creates history entry, error on invalid path.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-nop-write-tests-for-queryengine-search"></a>

## 📋 deco-nop Write tests for QueryEngine search

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:10 |

### Description

TDD: Write failing tests for text search. Test: search titles, search content, case insensitivity, partial matches.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-vz6-write-tests-for-deco-query-command"></a>

## 📋 deco-vz6 Write tests for deco query command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:51 |

### Description

TDD: Write failing tests for query command. Test: filter expressions, text search, combined filters, output formats.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-5wv-write-tests-for-error-aggregation"></a>

## 📋 deco-5wv Write tests for error aggregation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:57 |

### Description

TDD: Write failing tests for error collection. Test: multiple errors collected, sorted by file/line, deduplication, max errors limit.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-3jm-write-tests-for-queryengine-filter"></a>

## 📋 deco-3jm Write tests for QueryEngine filter

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:10 |

### Description

TDD: Write failing tests for node filtering. Test: filter by kind, status, tags, combined filters, empty results.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-4dj-write-tests-for-suggestion-engine"></a>

## 📋 deco-4dj Write tests for suggestion engine

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:50 |

### Description

TDD: Write failing tests for did-you-mean suggestions. Test: Levenshtein distance, threshold for suggestions, multiple candidates, no suggestion when too different.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-s5c-write-tests-for-deco-show-command"></a>

## 📋 deco-s5c Write tests for deco show command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:43 |

### Description

TDD: Write failing tests for show command. Test: displays node details, shows reverse refs, handles missing node, --json output.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-h9q-write-tests-for-patcher-apply-operation"></a>

## 📋 deco-h9q Write tests for Patcher apply operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:08 |

### Description

TDD: Write failing tests for patch file application. Test: JSON patch format, multiple operations, rollback on error.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-mvn-write-tests-for-yaml-context-extraction"></a>

## 📋 deco-mvn Write tests for YAML context extraction

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:26 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:48 |

### Description

TDD: Write failing tests for context extraction. Test: surrounding lines, start of file, end of file, column highlighting.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-wlu-write-tests-for-deco-list-command"></a>

## 📋 deco-wlu Write tests for deco list command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:40 |

### Description

TDD: Write failing tests for list command. Test: lists all nodes, --kind filter, --status filter, --tag filter, table output format.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-q41-write-tests-for-patcher-unset-operation"></a>

## 📋 deco-q41 Write tests for Patcher unset operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

TDD: Write failing tests for unset operation. Test: remove field, remove nested field, remove from array, missing field handling.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-yu5-write-tests-for-error-formatter"></a>

## 📋 deco-yu5 Write tests for error formatter

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:29 |

### Description

TDD: Write failing tests for Rust-like error formatting. Test: line numbers, context lines, color output, suggestion display, related locations.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-3co-write-tests-for-patcher-append-operation"></a>

## 📋 deco-3co Write tests for Patcher append operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

TDD: Write failing tests for append operation. Test: append to array, append to empty array, append to non-array error.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-sbk-write-tests-for-deco-validate-command"></a>

## 📋 deco-sbk Write tests for deco validate command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:37 |

### Description

TDD: Write failing tests for validate command. Test: exit code 0 on valid, exit code 1 on errors, error output format, --quiet flag.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-8sd-write-tests-for-error-code-registry"></a>

## 📋 deco-8sd Write tests for error code registry

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:27 |

### Description

TDD: Write failing tests for error codes. Test: code uniqueness, code lookup, category ranges (schema E001-E019, refs E020-E039, etc).

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-aoo-write-tests-for-patcher-set-operation"></a>

## 📋 deco-aoo Write tests for Patcher set operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

TDD: Write failing tests for set operation. Test: set existing field, set nested field, set new field, invalid path handling.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-sz1-write-tests-for-deco-init-command"></a>

## 📋 deco-sz1 Write tests for deco init command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:29 |

### Description

TDD: Write failing tests for init command. Test: creates .deco/, creates config.yaml, creates nodes/, detects existing project, --force flag.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-16w-write-tests-for-decoerror-structure"></a>

## 📋 deco-16w Write tests for DecoError structure

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:22 |

### Description

TDD: Write failing tests for DecoError struct. Test: all fields populated, Location type, Related type, Error() method output.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-az1-write-tests-for-validator-orchestrator"></a>

## 📋 deco-az1 Write tests for Validator orchestrator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:15 |

### Description

TDD: Write failing tests for combined validation. Test: all validators run, errors aggregated, proper ordering.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-336-write-tests-for-cobra-root-command"></a>

## 📋 deco-336 Write tests for Cobra root command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:29 |

### Description

TDD: Write failing tests for CLI root setup. Test: version flag, help output, global flags, subcommand registration.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-eft-write-tests-for-constraint-validator"></a>

## 📋 deco-eft Write tests for constraint Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:15 |

### Description

TDD: Write failing tests for constraint evaluation. Test: expression evaluation, cross-node constraints, violation reporting with both locations.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-xw8-write-tests-for-reference-validator"></a>

## 📋 deco-xw8 Write tests for reference Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:15 |

### Description

TDD: Write failing tests for ref validation. Test: valid refs resolve, broken refs detected, suggestion generation for typos.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-2t4-write-tests-for-schema-validator"></a>

## 📋 deco-2t4 Write tests for schema Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:15 |

### Description

TDD: Write failing tests for schema validation. Test: required fields (id, kind, version, status, title), missing fields, invalid types.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-im0-write-tests-for-reverse-reference-indexing"></a>

## 📋 deco-im0 Write tests for reverse reference indexing

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:01 |

### Description

TDD: Write failing tests for reverse ref computation. Test: used-by index, multiple refs to same node, circular refs, orphan nodes.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-jxy-write-tests-for-graphbuilder-service"></a>

## 📋 deco-jxy Write tests for GraphBuilder service

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:00 |

### Description

TDD: Write failing tests for graph building. Test: building from node slice, handling empty input, populating node map.

### Dependencies

- ⛔ **blocks**: `deco-elv`

---

<a id="deco-1jc-write-tests-for-yaml-line-number-tracking"></a>

## 📋 deco-1jc Write tests for YAML line number tracking

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:45 |

### Description

TDD: Write failing tests for line/column tracking during YAML parsing. Test: accurate line numbers, column positions, nested structures, multiline values.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-hxq-write-tests-for-historyrepository"></a>

## 📋 deco-hxq Write tests for HistoryRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:43 |

### Description

TDD: Write failing tests for append-only history repository. Test: Append(entry), Query(filter), file creation, append-only guarantee, concurrent writes.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-77y-write-tests-for-configrepository"></a>

## 📋 deco-77y Write tests for ConfigRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:36 |

### Description

TDD: Write failing tests for config repository. Test: Load, Save, missing config file handling, invalid config, default values.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-1er-write-tests-for-auditentry-domain-type"></a>

## 📋 deco-1er Write tests for AuditEntry domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:18 |

### Description

TDD: Write failing tests for AuditEntry struct. Test: timestamp handling, operation types, before/after state capture.

---

<a id="deco-344-write-tests-for-file-discovery"></a>

## 📋 deco-344 Write tests for file discovery

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:34 |

### Description

TDD: Write failing tests for node file discovery. Test: walking directory tree, .yaml file detection, path-to-ID mapping, ignoring non-yaml files, empty directories.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-xjv-write-tests-for-ref-domain-type"></a>

## 📋 deco-xjv Write tests for Ref domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:17 |

### Description

TDD: Write failing tests for Ref and RefLink types. Test: ref types (uses, related, emits_events, vocabulary), resolution state.

---

<a id="deco-53k-write-tests-for-yaml-noderepository"></a>

## 📋 deco-53k Write tests for YAML NodeRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:32 |

### Description

TDD: Write failing tests for YAML file-based node repository. Test: LoadAll, Load(id), Save(node), Delete(id), Exists(id), handling missing files, invalid YAML, nested directories.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-ll2-write-tests-for-constraint-domain-type"></a>

## 📋 deco-ll2 Write tests for Constraint domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:16 |

### Description

TDD: Write failing tests for Constraint struct. Test: expression storage, message, scope definition.

---

<a id="deco-0hm-write-tests-for-issue-domain-type"></a>

## 📋 deco-0hm Write tests for Issue domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:16 |

### Description

TDD: Write failing tests for Issue struct. Test: creation, severity levels, location tracking, resolved state.

---

<a id="deco-shb-write-tests-for-graph-domain-type"></a>

## 📋 deco-shb Write tests for Graph domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:15 |

### Description

TDD: Write failing tests for Graph struct. Test: node storage, lookup by ID, iteration, empty graph handling, duplicate ID handling.

---

<a id="deco-ora-write-tests-for-node-domain-type"></a>

## 📋 deco-ora Write tests for Node domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:25 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:14 |

### Description

TDD: Write failing tests for Node struct before implementation. Test: creation, field access, validation of required fields (ID, Kind, Version, Status, Title), serialization/deserialization.

---

<a id="deco-c5c-add-shell-completion-generation"></a>

## 📋 deco-c5c Add shell completion generation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:16 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:12 |

### Description

Add Cobra completion command for bash, zsh, fish, powershell.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`

---

<a id="deco-142-implement-node-rename-with-ref-updates"></a>

## 📋 deco-142 Implement node rename with ref updates

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:06 |

### Description

Create internal/service/refactor/rename.go. Rename node ID and update all references across the graph.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-bhb`

---

<a id="deco-0kr-implement-deco-apply-command"></a>

## 📋 deco-0kr Implement deco apply command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:58 |

### Description

Create internal/cli/apply.go. Apply patch file for AI integration.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-tzg`

---

<a id="deco-6j1-implement-deco-history-command"></a>

## 📋 deco-6j1 Implement deco history command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:00 |

### Description

Create internal/cli/history.go. Show audit log, optionally filtered by node ID.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-7kj`

---

<a id="deco-oam-implement-queryengine-search"></a>

## 📋 deco-oam Implement QueryEngine search

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:11 |

### Description

Add text search across node titles, content, and issues.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-nop`

---

<a id="deco-vzo-implement-deco-mv-command"></a>

## 📋 deco-vzo Implement deco mv command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 16:11 |

### Description

Create internal/cli/mv.go. Rename node with automatic ref updates across all nodes.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-747`

---

<a id="deco-0s1-implement-queryengine-filter"></a>

## 📋 deco-0s1 Implement QueryEngine filter

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:11 |

### Description

Create internal/service/query/engine.go. Filter nodes by kind, status, tags, or custom predicates.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-3jm`

---

<a id="deco-31k-implement-deco-unset-command"></a>

## 📋 deco-31k Implement deco unset command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:56 |

### Description

Create internal/cli/unset.go. Remove field: deco unset <id> <path>.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-3vb`

---

<a id="deco-7yp-implement-patcher-apply-operation"></a>

## 📋 deco-7yp Implement Patcher apply operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:09 |

### Description

Create internal/service/patcher/apply.go. Apply a patch file (JSON patch format) to nodes.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-h9q`

---

<a id="deco-inq-add-error-documentation-generator"></a>

## 📋 deco-inq Add error documentation generator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:31 |

### Description

Create cmd/deco/docs.go or script to generate error code documentation from registry. Output markdown for docs site.

### Dependencies

- ⛔ **blocks**: `deco-16e`

---

<a id="deco-1v8-implement-deco-append-command"></a>

## 📋 deco-1v8 Implement deco append command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:55 |

### Description

Create internal/cli/append.go. Append to array field: deco append <id> <path> <value>.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-zix`

---

<a id="deco-750-implement-patcher-unset-operation"></a>

## 📋 deco-750 Implement Patcher unset operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

Create internal/service/patcher/unset.go. Remove field at path from node.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-q41`

---

<a id="deco-5fv-implement-error-aggregation"></a>

## 📋 deco-5fv Implement error aggregation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:57 |

### Description

Create internal/errors/collector.go. Collect multiple errors during validation, sort by file and line, deduplicate related errors.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-5wv`

---

<a id="deco-53i-implement-deco-set-command"></a>

## 📋 deco-53i Implement deco set command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:53 |

### Description

Create internal/cli/set.go. Patch a node field: deco set <id> <path> <value>. Validate after patch.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-0jp`

---

<a id="deco-k3j-implement-patcher-append-operation"></a>

## 📋 deco-k3j Implement Patcher append operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

Create internal/service/patcher/append.go. Append value to array field at path.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-3co`

---

<a id="deco-3ip-implement-suggestion-engine"></a>

## 📋 deco-3ip Implement suggestion engine

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:52 |

### Description

Create internal/errors/suggest.go. Generate 'did you mean?' suggestions using Levenshtein distance for typos in refs, kinds, status values.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-4dj`

---

<a id="deco-gpu-define-historyrepository-interface"></a>

## 📋 deco-gpu Define HistoryRepository interface

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:20 |

### Description

Create internal/storage/history/repository.go with interface: Append(entry), Query(filter) for append-only audit log.

---

<a id="deco-m9l-implement-yaml-context-extraction"></a>

## 📋 deco-m9l Implement YAML context extraction

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:49 |

### Description

Create internal/errors/context.go. Extract surrounding lines from YAML file for error display. Handle edge cases (start/end of file).

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-mvn`

---

<a id="deco-1eg-implement-deco-query-command"></a>

## 📋 deco-1eg Implement deco query command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:51 |

### Description

Create internal/cli/query.go. Search nodes with filter expressions, output matching nodes.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-vz6`

---

<a id="deco-7rp-implement-patcher-set-operation"></a>

## 📋 deco-7rp Implement Patcher set operation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:04 |

### Description

Create internal/service/patcher/set.go. Set a field value at a path in a node. Validate path exists.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-aoo`

---

<a id="deco-5ns-define-configrepository-interface"></a>

## 📋 deco-5ns Define ConfigRepository interface

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:20 |

### Description

Create internal/storage/config/repository.go with interface: Load, Save for project config.

---

<a id="deco-4q8-implement-deco-show-command"></a>

## 📋 deco-4q8 Implement deco show command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:45 |

### Description

Create internal/cli/show.go. Show single node details including reverse refs (what uses this node).

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-s5c`

---

<a id="deco-e9t-implement-error-formatter"></a>

## 📋 deco-e9t Implement error formatter

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:30 |

### Description

Create internal/errors/formatter.go. Format DecoError in Rust-like style with colors, line numbers, context lines, and suggestions.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-yu5`

---

<a id="deco-2y3-implement-validator-orchestrator"></a>

## 📋 deco-2y3 Implement Validator orchestrator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:18 |

### Description

Create internal/service/validator/validator.go. Combine schema, refs, constraint validation. Return collected DecoErrors.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-az1`

---

<a id="deco-p54-define-noderepository-interface"></a>

## 📋 deco-p54 Define NodeRepository interface

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:20 |

### Description

Create internal/storage/node/repository.go with interface: LoadAll, Load(id), Save(node), Delete(id), Exists(id).

---

<a id="deco-qll-implement-deco-list-command"></a>

## 📋 deco-qll Implement deco list command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:41 |

### Description

Create internal/cli/list.go. List all nodes with filters: --kind, --status, --tag. Table output format.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-wlu`

---

<a id="deco-7cp-create-error-code-registry"></a>

## 📋 deco-7cp Create error code registry

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:28 |

### Description

Create internal/errors/codes.go. Define error codes E001-E099 with constants and descriptions. Categories: schema (E001-E019), refs (E020-E039), constraints (E040-E059), parse (E060-E079).

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-8sd`

---

<a id="deco-9ls-implement-constraint-validator"></a>

## 📋 deco-9ls Implement constraint Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:18 |

### Description

Create internal/service/validator/constraints.go. Evaluate constraint expressions across nodes, report violations with both locations.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-eft`

---

<a id="deco-8ca-define-auditentry-domain-type"></a>

## 📋 deco-8ca Define AuditEntry domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/history.go for audit log entries. Fields: Timestamp, Operation, NodeID, User, Before, After.

### Dependencies

- ⛔ **blocks**: `deco-1er`

---

<a id="deco-1fm-implement-deco-validate-command"></a>

## 📋 deco-1fm Implement deco validate command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:38 |

### Description

Create internal/cli/validate.go. Load all nodes, build graph, run validation, output errors in Rust-like format.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-sbk`

---

<a id="deco-03s-define-decoerror-structure"></a>

## 📋 deco-03s Define DecoError structure

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:27 |

### Description

Create internal/domain/error.go with DecoError struct: Code, Summary, Detail, Location, Context, Suggestion, Related. Include Location and Related types.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-16w`

---

<a id="deco-jna-implement-reference-validator"></a>

## 📋 deco-jna Implement reference Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:18 |

### Description

Create internal/service/validator/refs.go. Check all refs resolve to existing nodes, collect broken refs with suggestions.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-xw8`

---

<a id="deco-rrk-define-ref-domain-type"></a>

## 📋 deco-rrk Define Ref domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:15 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/ref.go for node references. Types: uses, related, emits_events, vocabulary. Include RefLink for resolved refs.

### Dependencies

- ⛔ **blocks**: `deco-xjv`

---

<a id="deco-7ba-implement-deco-init-command"></a>

## 📋 deco-7ba Implement deco init command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:29 |

### Description

Create internal/cli/init.go. Initialize .deco/ directory with config.yaml and nodes/ folder. Check for existing project.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-sz1`

---

<a id="deco-bx5-implement-schema-validator"></a>

## 📋 deco-bx5 Implement schema Validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:18 |

### Description

Create internal/service/validator/schema.go. Validate nodes against required fields: meta.id, meta.kind, meta.version, meta.status, meta.title.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-2t4`

---

<a id="deco-nbv-add-yaml-line-number-tracking"></a>

## 📋 deco-nbv Add YAML line number tracking

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:47 |

### Description

Wrap yaml.v3 decoder to track line/column numbers for each parsed field. Essential for error reporting with locations.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-1jc`

---

<a id="deco-7au-setup-cobra-root-command"></a>

## 📋 deco-7au Setup Cobra root command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:29 |

### Description

Create cmd/deco/main.go and internal/cli/root.go. Initialize Cobra with root command, version flag, global flags.

### Dependencies

- ⛔ **blocks**: `deco-t6q`
- ⛔ **blocks**: `deco-5xa`
- ⛔ **blocks**: `deco-336`

---

<a id="deco-lwx-define-constraint-domain-type"></a>

## 📋 deco-lwx Define Constraint domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/constraint.go for validation constraints. Fields: Expr, Message, Scope.

### Dependencies

- ⛔ **blocks**: `deco-ll2`

---

<a id="deco-owf-implement-historyrepository"></a>

## 📋 deco-owf Implement HistoryRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:44 |

### Description

Create internal/storage/history/jsonl_repository.go implementing HistoryRepository. Append-only JSONL file at .deco/history.jsonl.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-hxq`

---

<a id="deco-pqf-implement-reverse-reference-indexing"></a>

## 📋 deco-pqf Implement reverse reference indexing

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:01 |

### Description

In GraphBuilder, compute 'used-by' reverse refs so we can show what references a node.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-im0`

---

<a id="deco-982-define-issue-domain-type"></a>

## 📋 deco-982 Define Issue domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/issue.go for tracking TBDs/questions in nodes. Fields: ID, Severity, Message, Location, Resolved.

### Dependencies

- ⛔ **blocks**: `deco-0hm`

---

<a id="deco-0o3-implement-graphbuilder-service"></a>

## 📋 deco-0o3 Implement GraphBuilder service

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 15:00 |

### Description

Create internal/service/graph/builder.go. Build Graph from slice of Nodes, populate reverse reference index.

### Dependencies

- ⛔ **blocks**: `deco-elv`
- ⛔ **blocks**: `deco-jxy`

---

<a id="deco-px9-implement-configrepository"></a>

## 📋 deco-px9 Implement ConfigRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:36 |

### Description

Create internal/storage/config/yaml_config.go implementing ConfigRepository. Load/save .deco/config.yaml with project settings.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-77y`

---

<a id="deco-zty-define-graph-domain-type"></a>

## 📋 deco-zty Define Graph domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/graph.go with Graph struct holding nodes map and reverse reference index. Methods for node lookup and traversal.

### Dependencies

- ⛔ **blocks**: `deco-shb`

---

<a id="deco-5ry-implement-file-discovery-for-nodes"></a>

## 📋 deco-5ry Implement file discovery for nodes

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:35 |

### Description

Walk .deco/nodes/ directory tree, discover all .yaml files, map file paths to node IDs based on relative path.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-344`

---

<a id="deco-bqg-define-node-domain-type"></a>

## 📋 deco-bqg Define Node domain type

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:19 |

### Description

Create internal/domain/node.go with Node struct: ID, Kind, Version, Status, Title, Tags, Refs, Content, Issues. Include Section and Block types.

### Dependencies

- ⛔ **blocks**: `deco-ora`

---

<a id="deco-nkw-implement-yaml-noderepository"></a>

## 📋 deco-nkw Implement YAML NodeRepository

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:33 |

### Description

Create internal/storage/node/yaml_repository.go implementing NodeRepository. Parse YAML files from .deco/nodes/, handle nested directories, preserve line numbers for error reporting.

### Dependencies

- ⛔ **blocks**: `deco-16e`
- ⛔ **blocks**: `deco-53k`

---

<a id="deco-rba-initialize-go-module-and-project-structure"></a>

## 📋 deco-rba Initialize Go module and project structure

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 🔹 Medium (P2) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-30 15:14 |
| **Updated** | 2026-02-01 17:09 |
| **Closed** | 2026-01-31 14:13 |

### Description

Run go mod init, create folder structure: cmd/deco/, internal/domain/, internal/storage/, internal/service/, internal/cli/. Add .gitignore for Go.

---

<a id="deco-4bhp-investigate-streaming-incremental-validation"></a>

## 📋 deco-4bhp Investigate streaming/incremental validation

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-03 20:25 |
| **Updated** | 2026-02-03 20:26 |
| **Closed** | 2026-02-03 20:26 |

### Description

**Summary**
Explore incremental/streaming validation for long AI-generated outputs to fail fast on obvious errors.

**Problem / Opportunity**
Current validation happens after full output. For large AI generations, early feedback could save time and cost.

**Goals**
- Prototype a streaming validator that can reject invalid patches early

**Non-goals**
- Replacing full validation
- Implementing an LLM client

**User Stories**
- As an AI-integrator, I want early validation so I can stop invalid generations sooner.

**Scope**
- In scope: design/prototype, feasibility assessment
- Out of scope: production-ready streaming parser

**Constraints**
- Must not increase complexity of core validation for normal CLI usage

**Proposed Approach**
1. Identify which validations can be done incrementally
2. Prototype a streaming patch validator
3. Document tradeoffs

**TDD Plan (Required)**
- Tests to write first (red): n/a (research spike)
- Minimal implementation (green): prototype or doc
- Refactor: n/a

**Acceptance Criteria**
- [ ] Documented feasibility and recommended next steps

**Test Plan**
- None (research)

**Dependencies**
- None

**Files / Areas Touched**
- docs/ or internal/services/validator

**Risks**
- Streaming YAML/JSON parsing complexity

**Open Questions**
- Which validations are safe to run incrementally?

---

<a id="deco-1se6-add-api-specification-example-project"></a>

## 📋 deco-1se6 Add API specification example project

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 22:33 |
| **Updated** | 2026-02-03 21:34 |
| **Closed** | 2026-02-03 21:34 |

### Description

## Context / Problem

Current examples (snake, space-invaders) are game-focused, which doesn't showcase the broader use cases described in the updated README/SPEC. Need a non-game example to demonstrate Deco's value for:
- Technical teams documenting APIs
- The AI-assisted workflow angle
- Reference validation across endpoints/schemas

**Current behavior:** Only game design examples exist in examples/
**Why it matters:** Broader audience needs to see themselves in the tool. API specs are familiar, relatable, and demonstrate the validation/reference value clearly.
**Relevant files:** examples/ directory

## Goal / Outcome

A complete, validating API specification example in examples/api-spec/ that demonstrates:
- Endpoint definitions with method/path/response
- Schema nodes referenced by endpoints
- Authentication system with cross-references
- Error handling patterns
- Rate limiting documentation
- At least one constraint showing cross-node validation

## Scope

**In:**
- Create examples/api-spec/ directory structure
- 8-12 nodes covering a realistic REST API
- README explaining the example
- Must pass deco validate

**Out:**
- OpenAPI/Swagger generation (separate feature)
- Actual runnable API code
- Client SDK documentation

## TDD Plan

**Tests to write first:**
- Validation test: `deco validate` in examples/api-spec/ returns 0 exit code
- Reference test: All refs resolve (no broken links)
- Completeness test: No open issues with severity >= high

**Expected failures:** N/A for documentation task - validation is the test

## Acceptance Criteria

- [ ] examples/api-spec/ directory exists with structured nodes
- [ ] `deco validate` passes with zero errors
- [ ] At least 8 nodes covering: auth, users, endpoints, schemas, errors
- [ ] Cross-references between endpoints and schemas
- [ ] At least one CEL constraint demonstrating cross-node validation
- [ ] README.md in example explaining structure and purpose
- [ ] Example demonstrates all core node features (refs, issues, contracts)

## Proposed Approach

Directory structure:
```
examples/api-spec/
├── .deco/
│   ├── config.yaml
│   └── nodes/
│       ├── systems/
│       │   ├── auth.yaml          # JWT authentication
│       │   ├── rate-limiting.yaml # Rate limit rules
│       │   └── errors.yaml        # Error response patterns
│       ├── endpoints/
│       │   ├── users.yaml         # /users CRUD
│       │   ├── auth.yaml          # /auth/login, /auth/refresh
│       │   └── products.yaml      # /products CRUD
│       ├── schemas/
│       │   ├── user.yaml          # User object schema
│       │   ├── product.yaml       # Product object schema
│       │   └── error.yaml         # Error response schema
│       └── glossaries/
│           └── api-terms.yaml     # JWT, Bearer, etc.
└── README.md
```

Each endpoint node references its request/response schemas. Auth system is referenced by all endpoints. Demonstrates the reference validation value.

## Test / Verification Plan

**Commands:**
```bash
cd examples/api-spec
deco validate
deco stats
deco graph --format mermaid
```

**Expected outcome:**
- validate: exit 0, no errors
- stats: shows 8+ nodes, 0 open high-severity issues
- graph: shows interconnected dependency structure

## Risks / Edge Cases

- Scope creep into too many endpoints (keep it focused)
- Over-engineering the example (should be learnable, not comprehensive)
- Ensure example stays valid as deco evolves (add to CI?)

## Notes / Links

- Updated README.md lists API Specs as primary use case
- Existing game examples in examples/snake/ and examples/space-invaders/ for reference
- Should complement, not replace, existing examples

---

<a id="deco-oo2q-docker-dev-environment-setup"></a>

## 📋 deco-oo2q Docker dev environment setup

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-02 13:34 |
| **Updated** | 2026-02-02 13:35 |
| **Closed** | 2026-02-02 13:35 |

### Notes

Containerized Claude Code environment for consistent development

---

<a id="deco-uvas-hash-truncation-to-64-bits-risks-collision"></a>

## 📋 deco-uvas Hash truncation to 64 bits risks collision

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 21:43 |
| **Updated** | 2026-02-02 22:28 |
| **Closed** | 2026-02-02 22:28 |

### Description

Goal: Ensure content hash truncation doesn't cause missed change detection.

## Current Implementation
File: internal/cli/audit.go:38-61

Uses SHA-256, truncated to first 8 bytes (64 bits / 16 hex chars):
hash := sha256.Sum256(data)
return hex.EncodeToString(hash[:8])

## Collision Risk Analysis

### Birthday Problem
With 64-bit hash, collision probability exceeds:
- 50% at ~5 billion documents
- 0.1% at ~6 million documents
- 0.001% at ~600K documents

### Deco Context
- Typical project: 10-1000 nodes
- Typical history: 10-100 entries per node
- Large project: 10K nodes * 100 entries = 1M total hashes

Collision probability for 1M hashes: ~0.003% (1 in 30,000)
This is LOW but non-zero over project lifetime.

### Collision Impact
If two different contents produce same hash:
- deco sync would miss a manual edit
- Concurrent edit detection could false-positive

NOT a security concern (not using hash for auth), but could cause data integrity issues.

## Options

### Option 1: Keep 64 bits (Status Quo)
- Pro: Compact output, readable hashes
- Pro: Sufficient for most projects
- Con: Non-zero collision risk at scale
- Decision: Document the risk, accept for typical use

### Option 2: Increase to 128 bits (16 bytes / 32 hex chars)
- Pro: Collision probability drops to negligible (1 in 10^19 for 1M docs)
- Con: Longer hash strings in output/logs
- Con: Breaking change for existing history

### Option 3: Full SHA-256 (32 bytes / 64 hex chars)
- Pro: No truncation, maximum security
- Con: Very long strings
- Con: Overkill for change detection

### Option 4: Use different hash (e.g., xxHash for speed)
- Pro: Faster for large content
- Con: Less collision resistance
- Not recommended for integrity checking

## Recommendation
**Option 2: 128 bits** - Best balance of safety and readability.

Migration path:
1. Add version flag to hash format
2. New hashes use 128 bits
3. Comparison handles both formats during transition
4. Eventually deprecate 64-bit format

## Acceptance Criteria
- [ ] Document collision probability analysis (this issue)
- [ ] Make decision: keep 64 bits OR upgrade to 128 bits
- [ ] If upgrading: implement backward-compatible hash comparison
- [ ] Update audit.go comments with rationale
- [ ] Add to docs if keeping 64-bit (explain trade-off)

## References
- Birthday problem: https://en.wikipedia.org/wiki/Birthday_problem
- Git uses 160-bit SHA-1 (but for security, not just change detection)
- UUID collision risk similar discussion

## Non-Goals
- Cryptographic security (hash is for change detection, not auth)
- Zero collision guarantee (impractical without infinite storage)

### Notes

CLARIFICATIONS (from review):

1. **Sync Logic Linkage**: `deco sync` in audit.go:89-112 compares computed hash of current file content against last recorded hash in history. If different and history.last_operation != 'manual_edit', it's detected as manual edit. Collision would cause two different contents to match same hash → missed detection.

2. **Documentation Criterion**: 'Document collision probability analysis' means: add comment block in audit.go above computeContentHash() explaining 64-bit birthday problem math and design decision rationale. Not a separate doc file.

---

<a id="deco-7jb-expand-test-coverage-for-validator-patcher-cli"></a>

## 📋 deco-7jb Expand test coverage for validator/patcher/CLI

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:14 |
| **Updated** | 2026-02-02 21:20 |
| **Closed** | 2026-02-02 14:46 |

### Description

Goal: Increase confidence in core workflows through comprehensive test coverage.

## Current Coverage Status
```
Package                         Coverage
─────────────────────────────────────────
services/graph                  100.0%  ✅
services/query                  100.0%  ✅
domain                          95.4%   ✅
errors                          94.8%   ✅
storage/history                 86.4%
cli                             85.8%
storage/config                  85.2%
services/validator              84.5%
errors/yaml                     84.6%
storage/node                    84.0%
services/refactor               81.9%
services/patcher                65.7%   ⚠️  Lowest
cmd/deco                        0.0%    (main entry, expected)
```

## Priority Areas (by coverage gap)

### 1. services/patcher (65.7%) - HIGH PRIORITY
Critical for AI workflows. Needs tests for:
- [ ] Nested path operations (e.g., `content.sections[0].blocks[1].data.value`)
- [ ] Edge cases: empty arrays, nil maps, type mismatches
- [ ] Rollback behavior when operations fail mid-batch
- [ ] Boundary conditions for array indices

### 2. services/refactor (81.9%)
Ref updates when nodes are moved/renamed:
- [ ] Circular reference handling
- [ ] Partial match scenarios
- [ ] Large graph performance

### 3. services/validator (84.5%)
Schema and constraint validation:
- [ ] Custom block type validation edge cases
- [ ] CEL constraint evaluation errors
- [ ] Contract validation with malformed given/when/then
- [ ] Unknown field detection in nested structures

### 4. storage/node (84.0%)
Node persistence:
- [ ] Concurrent read/write scenarios
- [ ] Malformed YAML recovery
- [ ] File system errors (permissions, disk full)

### 5. errors/yaml (84.6%)
Error context extraction:
- [ ] Multi-line YAML errors
- [ ] Deeply nested error locations
- [ ] Unicode/special character handling

## Test Categories to Add

### Edge Case Tests
- Empty projects (0 nodes)
- Single node with self-reference
- Maximum depth nesting in content blocks
- Very large nodes (>1MB YAML)
- Unicode in all string fields

### Error Path Tests
- Disk full during write
- Permission denied
- Corrupted YAML files
- Invalid UTF-8 sequences
- Circular dependencies

### Integration Tests
- Full CLI workflow: init → create → validate → set → apply → sync
- Multi-node reference integrity through operations
- History consistency after failures

## Acceptance Criteria
- [ ] patcher coverage ≥ 85%
- [ ] All packages ≥ 80% coverage
- [ ] Edge cases documented in test names
- [ ] No test flakiness (run 10x without failure)
- [ ] CI reports coverage (if available)

## Non-Goals
- 100% coverage (diminishing returns)
- Mocking external systems (file system is acceptable)
- Performance benchmarks (separate concern)

## Testing Guidelines
- Use table-driven tests for similar cases
- Test names should describe the scenario: `TestPatcher_ApplySet_NestedPath_CreatesIntermediateNodes`
- Include both happy path and error path in each test file
- Prefer real file system operations over mocks for CLI tests

### Notes

CLARIFICATIONS (from review):

1. Flakiness Test Protocol: 'No test flakiness (run 10x without failure)' means:
   Command: for i in 1 2 3 4 5 6 7 8 9 10; do go test ./... -count=1 || exit 1; done
   Run from repo root, in CI or clean local env. -count=1 disables test caching. All 10 runs must pass with exit 0.

---

<a id="deco-04f-validation-performance-on-large-projects"></a>

## 📋 deco-04f Validation performance on large projects

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 22:12 |
| **Closed** | 2026-02-01 22:12 |

### Description

Goal: validation scales to large graphs.

Acceptance:
- Benchmark suite added.
- Identify/resolve hot spots.

---

<a id="deco-xdp-validation-performance-on-large-projects"></a>

## 📋 deco-xdp Validation performance on large projects

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 20:48 |
| **Closed** | 2026-02-01 20:48 |

### Description

Goal: validation scales to large graphs.

Acceptance:
- Benchmark suite added.
- Identify/resolve hot spots.

### Dependencies

- 🔗 **related**: `deco-04f`

---

<a id="deco-ai9-deterministic-load-save-ordering"></a>

## 📋 deco-ai9 Deterministic load/save ordering

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 21:41 |
| **Closed** | 2026-02-01 21:41 |

### Description

Goal: avoid noisy diffs and surprising reorderings.

Acceptance:
- Save order stable across runs.
- Map ordering deterministic (or normalized) where possible.

---

<a id="deco-3ix-docs-alignment-readme-spec-match-schema-cli-behavior"></a>

## 📋 deco-3ix Docs alignment: README/SPEC match schema & CLI behavior

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-01 22:10 |
| **Closed** | 2026-02-01 22:10 |

### Description

Goal: docs reflect reality and are consistent with tests.

Acceptance:
- README/SPEC examples validate.
- CLI flags/behavior documented accurately.

---

<a id="deco-23m-saved-queries-query-presets"></a>

## ✨ deco-23m Saved queries / query presets

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 22:28 |
| **Closed** | 2026-02-02 22:28 |

### Description

Goal: Let teams store and reuse common queries for consistent project health monitoring.

## User Stories
1. Project leads: Run the same health check queries each sprint (orphan nodes, open issues, draft items).
2. New team members: Discover useful queries without learning the syntax.
3. CI/CD: Run named validation queries in pipelines.

## Problem
Current workflow for 'find all items without approval':
deco query --kind item --status draft
Must remember/retype flags every time. Can't share queries across team.

## Current Query Capabilities
Flags available:
- --kind / -k: Filter by node type
- --status / -s: Filter by status
- --tag / -t: Filter by tag
- Search term: Text search in title/summary

All filters combine with AND logic.

## Proposed Implementation

### Config Schema Addition
In .deco/config.yaml:

saved_queries:
  drafts:
    description: All nodes still in draft
    filters:
      status: draft
  orphans:
    description: Nodes with no incoming references
    builtin: orphan_nodes
  stale-approved:
    description: Approved nodes not updated in 30 days
    filters:
      status: approved
    search: ''  # No text search
  combat-items:
    description: Items tagged with combat
    filters:
      kind: item
      tag: combat

### Builtin Queries
Special queries that require logic beyond simple filters:
- orphan_nodes: Nodes with no refs.uses pointing to them
- unreachable: Nodes not reachable from root/entry nodes
- cyclic_refs: Nodes involved in reference cycles
- open_issues: Nodes with unresolved issues array items
- missing_refs: Nodes with broken references

### CLI Interface

Run saved query:
deco query --saved drafts
deco query -S drafts          # Short form

List saved queries:
deco query --list-saved
# Output:
# Saved queries:
#   drafts       - All nodes still in draft
#   orphans      - Nodes with no incoming references
#   combat-items - Items tagged with combat

Add new saved query (from current flags):
deco query --kind item --status draft --save my-query
# Saves current filters as 'my-query'

Delete saved query:
deco query --delete-saved my-query

Describe a saved query:
deco query --describe orphans
# Output:
# Query: orphans
# Type: builtin
# Description: Nodes with no incoming references

### Query Composition
Allow saved queries as base + additional filters:
deco query --saved drafts --kind item
# Runs 'drafts' query filtered further by kind=item

## Acceptance Criteria
- [ ] saved_queries section in config.yaml
- [ ] --saved / -S flag executes named query
- [ ] --list-saved shows available queries
- [ ] --save <name> creates query from current flags
- [ ] --delete-saved <name> removes query
- [ ] At least 3 builtin queries (orphans, open_issues, missing_refs)
- [ ] Saved queries composable with additional flags

## Example Saved Queries to Ship By Default
1. orphans - Nodes nobody references
2. open-issues - Nodes with issues[].resolved = false
3. needs-review - status=draft with version > 1 (edited but not reviewed)
4. stale - approved nodes, could check last modified if we add that

## Non-Goals
- Complex query language (keep it simple, use flags)
- Query history/recent queries
- Cross-project query sharing

### Notes

CLARIFICATIONS (from review):

1. **Composition Precedence**: Saved query provides base filters. CLI flags override/extend:
   - Same key (e.g., both have kind): CLI wins
   - Different keys: merged (AND logic)
   - search: CLI appends to saved (space-separated terms)
   Example: saved has kind=item, CLI adds --status=draft → kind=item AND status=draft

2. **Empty Search Clarification**: `search: ''` (empty string) means 'no text search filter' - equivalent to omitting the search field entirely. NOT 'match all' (which is the default behavior when no search specified anyway).

---

<a id="deco-rmv-detect-and-warn-on-concurrent-edit-conflicts"></a>

## ✨ deco-rmv Detect and warn on concurrent edit conflicts

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 18:13 |
| **Updated** | 2026-02-02 22:25 |
| **Closed** | 2026-02-02 22:25 |

### Description

Goal: Prevent silent data loss when multiple editors modify the same node concurrently.

## User Stories
1. Team collaboration: When Alice and Bob both edit sword-001, the second save should warn instead of silently overwriting.
2. AI safety: When LLM patches a node that was manually edited since it was read, abort and report the conflict.
3. Debugging: When conflicts occur, understand what changed and when.

## Problem
Current flow (no protection):
1. Alice loads sword-001 (version 3, hash ABC123)
2. Bob loads sword-001 (version 3, hash ABC123)
3. Alice saves changes -> version 4, hash DEF456
4. Bob saves changes -> version 5, hash GHI789 (Alice's changes silently lost!)

## Current State

### What Exists
- ContentHash computed for every mutation (internal/cli/audit.go:38)
- History tracks hash per operation
- sync command detects manual edits by comparing hashes
- Node.Version field (but only for display, not concurrency)

### What's Missing
- No hash comparison before write
- No CLI flag to supply expected hash
- No conflict detection during apply/set/append/unset

## Proposed Implementation

### Optimistic Concurrency Control
Add --expect-hash flag to mutation commands:
deco set sword-001 title 'New Name' --expect-hash ABC123

If current node hash != expected hash -> abort with conflict error.

### Implementation Details

1. Add to mutation commands (set, append, unset, apply):
   - --expect-hash flag (optional)
   - Before applying changes, compute current hash
   - If flag provided and hashes don't match, abort

2. Error output:
   Conflict detected on sword-001
   
   Expected hash:  ABC123
   Current hash:   XYZ789
   
   The node was modified since you last read it.
   
   Options:
   1. Reload the node and reapply your changes
   2. Use --force to overwrite (loses concurrent changes)
   3. Use 'deco diff sword-001' to see current state

3. Exit codes:
   - 0: Success
   - 1: Validation/operation error
   - 3: Conflict detected (new)

### For AI Workflows
LLM flow becomes:
1. Read node, capture hash from response
2. Generate patch
3. Apply with --expect-hash
4. If conflict, re-read and regenerate

### show Command Enhancement
Include hash in output for easy capture:
deco show sword-001
# ...
Content hash: ABC123 (use with --expect-hash)

## CLI Interface

Mutation commands gain:
--expect-hash <hash>   # Abort if current hash differs
--force                # Overwrite even if conflict (existing, maybe rename)

New output:
deco show sword-001 --format=json  # Include content_hash field

## Acceptance Criteria
- [ ] --expect-hash flag on set, append, unset, apply commands
- [ ] Conflict error with clear message when hashes don't match
- [ ] Exit code 3 for conflicts (distinct from validation errors)
- [ ] deco show includes content hash
- [ ] --force flag to override conflict check
- [ ] Tests for conflict scenarios

## Edge Cases
- First write to new node (no expected hash needed)
- Node deleted between read and write
- Multiple rapid edits (hash chain)
- --expect-hash with invalid hash format

## Non-Goals
- Merge conflict resolution (show conflict, don't auto-merge)
- Three-way merge (too complex for YAML)
- File locking (Git handles this at commit time)
- Real-time collaboration (deco is CLI, not collaborative editor)

### Notes

Won't implement - git already handles merge conflicts at commit/push time. This was over-engineering.

---

<a id="deco-acy-validate-output-lacks-file-line-context-despite-yaml-location-tracker"></a>

## 📋 deco-acy validate output lacks file/line context despite YAML location tracker

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-02-01 17:30 |
| **Updated** | 2026-02-01 19:27 |
| **Closed** | 2026-02-01 19:27 |

### Description

Observed: deco validate prints only error summary text (no file/line), even though README examples show file/line context and a YAML location tracker exists.
Expected: validation errors include file:line (and column when available) to match docs and speed fixing.
Refs: internal/cli/validate.go (printing), internal/errors/yaml/location.go (tracker).

---

<a id="deco-1pi-add-deco-stats-command-for-project-overview"></a>

## ✨ deco-1pi Add deco stats command for project overview

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:37 |
| **Closed** | 2026-02-01 18:37 |

### Description

Add a 'deco stats' command showing project health: node count by kind, node count by status, open issues count by severity, reference health (dangling refs), constraint violations summary.

---

<a id="deco-t2h-add-deco-diff-command-to-show-node-history"></a>

## ✨ deco-t2h Add deco diff command to show node history

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:31 |
| **Closed** | 2026-02-01 18:31 |

### Description

Add a 'deco diff <id>' command that shows changes to a node over time from the history log. Support --since=<timestamp> and --last=N to limit output. Show before/after for each change.

---

<a id="deco-0ql-add-deco-export-command-for-multiple-output-formats"></a>

## ✨ deco-0ql Add deco export command for multiple output formats

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-01-31 16:32 |
| **Closed** | 2026-01-31 16:32 |

### Description

Add a 'deco export' command supporting multiple output formats: --format=html, --format=markdown, --format=pdf. Complements the LaTeX feature. Should produce navigable documentation with cross-references.

---

<a id="deco-5uc-add-deco-graph-command-to-visualize-dependencies"></a>

## ✨ deco-5uc Add deco graph command to visualize dependencies

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:26 |
| **Closed** | 2026-02-01 18:26 |

### Description

Add a 'deco graph' command that outputs the node dependency graph. Support DOT/Graphviz format for visualization. Could also support --format=mermaid for Markdown embedding.

---

<a id="deco-gp6-add-deco-issues-command-to-list-all-tbds"></a>

## ✨ deco-gp6 Add deco issues command to list all TBDs

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:22 |
| **Closed** | 2026-02-01 18:22 |

### Description

Add a 'deco issues' command that lists all open issues/TBDs across the entire design graph. Support filtering by --severity and --node. Show location and context for each issue.

---

<a id="deco-t5h-add-deco-rm-command-to-delete-nodes"></a>

## ✨ deco-t5h Add deco rm command to delete nodes

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:20 |
| **Closed** | 2026-02-01 18:20 |

### Description

Add a 'deco rm <id>' command to delete a node. Should warn if other nodes reference it (show reverse refs) and require --force to delete anyway. Log deletion in history.

---

<a id="deco-513-add-deco-create-command-for-scaffolding-nodes"></a>

## ✨ deco-513 Add deco create command for scaffolding nodes

| Property | Value |
|----------|-------|
| **Type** | ✨ feature |
| **Priority** | ☕ Low (P3) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 16:32 |
| **Updated** | 2026-02-01 18:18 |
| **Closed** | 2026-02-01 18:18 |

### Description

Add a 'deco create <id>' command that scaffolds a new node with required fields. Could support --kind flag and optionally --template for predefined node templates.

---

<a id="deco-h2c-add-contract-validation-to-deco-validate-command"></a>

## 📋 deco-h2c Add contract validation to deco validate command

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 💤 Backlog (P4) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:01 |
| **Updated** | 2026-02-01 19:04 |
| **Closed** | 2026-02-01 19:04 |

### Description

Integrate contract validator into deco validate orchestrator. Report errors for invalid contract syntax or references.

### Dependencies

- ⛔ **blocks**: `deco-wxh`

---

<a id="deco-wxh-validate-contract-scenarios-reference-valid-nodes"></a>

## 📋 deco-wxh Validate contract scenarios reference valid nodes

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 💤 Backlog (P4) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:01 |
| **Updated** | 2026-02-01 19:03 |
| **Closed** | 2026-02-01 19:03 |

### Description

Ensure contract scenarios only reference nodes/fields that exist in the design graph. Validate expect_event references to actual event definitions.

### Dependencies

- ⛔ **blocks**: `deco-0sk`

---

<a id="deco-4ci-implement-contract-scenario-parser"></a>

## 📋 deco-4ci Implement contract scenario parser

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 💤 Backlog (P4) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:01 |
| **Updated** | 2026-02-01 18:54 |
| **Closed** | 2026-02-01 18:54 |

### Description

Create internal/domain/contract.go with Contract, Scenario, Step types. Parse contract.scenarios from YAML into structured domain types.

---

<a id="deco-0sk-add-contract-syntax-validator"></a>

## 📋 deco-0sk Add contract syntax validator

| Property | Value |
|----------|-------|
| **Type** | 📋 task |
| **Priority** | 💤 Backlog (P4) |
| **Status** | ⚫ closed |
| **Created** | 2026-01-31 14:01 |
| **Updated** | 2026-02-01 19:00 |
| **Closed** | 2026-02-01 19:00 |

### Description

Create internal/service/validator/contract.go. Validate contract scenarios syntax: given/when/then structure, field types, scenario IDs are unique within node.

### Dependencies

- ⛔ **blocks**: `deco-4ci`

---

