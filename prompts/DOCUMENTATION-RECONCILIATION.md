# Documentation Reconciliation Prompt

## Purpose

Use this prompt after an authorized change or audit to restore agreement
between governing documents, explanatory documentation, implementation, and
verification. Supply the TASK or explicit scope that authorizes the
reconciliation.

## Prompt

Reconcile documentation only within the supplied authorized scope. The goal is
to remove proven drift without changing project intent merely to match the
current implementation.

Before editing:

1. Read the repository working agreement, project vision, documentation model,
   and the TASK or explicit change that authorizes reconciliation.
2. Read relevant accepted ADRs, approved SPECs, explanatory documentation,
   implementation, tests, and completion evidence.
3. Inspect version-control status and preserve unrelated contributor work.
4. Identify each material mismatch and the authority responsible for the
   disputed claim.

Classify mismatches before acting:

- **Explanation drift:** explanatory text disagrees with governing documents
  or verified behavior.
- **Reference drift:** links, identifiers, statuses, or supersession
  relationships are missing or stale.
- **Authority drift:** code or tests disagree with an accepted ADR or approved
  SPEC.
- **Evidence drift:** TASK completion claims do not match checks or resulting
  artifacts.
- **Unresolved intent:** repository evidence cannot establish which source
  should govern.

Apply the smallest authorized correction:

- Correct explanatory documentation to reflect current authoritative sources.
- Repair references and lifecycle metadata without rewriting historical
  meaning.
- Update TASK evidence to state only what was actually verified.
- Modify an accepted ADR, approved SPEC, project vision, implementation, or
  tests only when the supplied scope explicitly authorizes that kind of
  change.

Do not make current code authoritative by default. Do not rewrite an accepted
decision's rationale, change approved behavior, or broaden a TASK to make a
conflict disappear. When authority drift or unresolved intent requires a new
decision, behavior change, or implementation task, stop and report the needed
human resolution.

After editing:

1. verify links, identifiers, and lifecycle relationships in scope;
2. compare corrected claims with their authoritative sources and observable
   implementation;
3. review the diff for duplicated normative content and unrelated changes; and
4. run only the existing checks relevant to the reconciliation.

Report documents changed, mismatches resolved, evidence used, checks
performed, and unresolved conflicts. Do not implement follow-up work, create
tooling, or cross the supplied scope. Stop after reconciliation.
