# DaD Workflow

DaD governs a change from proposed outcome to verified completion. The same
workflow applies to human contributors and AI coding agents; projects may map
review and approval steps onto their existing collaboration process.

Document statuses and allowed transitions are defined in
[DOCUMENTATION.md](DOCUMENTATION.md).

## 1. Frame the task

Create a TASK in `Proposed` status using the
[TASK template](templates/TASK-TEMPLATE.md). Define:

- one concrete outcome;
- work that is in and out of scope;
- known constraints and dependencies;
- governing ADRs and SPECs;
- verifiable acceptance conditions.

Resolve ambiguous boundaries before marking the task `Ready`. The task should
be small enough that a reviewer can determine whether its result belongs
within the stated scope.

## 2. Assemble the working context

Before implementation, read:

1. the project vision and repository working agreement;
2. the active TASK;
3. every accepted ADR and approved SPEC referenced by the TASK;
4. directly affected code, tests, and local instructions.

Follow material references needed to understand a governing constraint. Do not
load unrelated documents merely because they exist. The goal is complete
relevant context, not maximum context.

Confirm that referenced documents exist, have a governing lifecycle status,
and do not contradict one another. If they conflict, stop and request a
resolution.

## 3. Resolve authority gaps

The working context must be sufficient to implement without inventing durable
project policy.

- If the work requires a consequential architectural choice, prepare an ADR
  and obtain acceptance before relying on it.
- If required behavior or an interface is ambiguous, prepare or revise a SPEC
  and obtain approval before implementing that behavior.
- If the task outcome, scope, or acceptance conditions need material change,
  update the TASK and obtain review before continuing.

Not every uncertainty requires a document. Local, reversible implementation
details can remain implementation choices when they do not alter architecture,
approved behavior, or task boundaries.

## 4. Implement within the boundary

Mark the TASK `In Progress` when work begins. Make the smallest coherent change
that satisfies the outcome and follows existing project patterns.

During implementation:

- do not add anticipated features or unrelated cleanup;
- preserve work outside the task;
- keep implementation consistent with accepted ADRs and approved SPECs;
- add or update proportionate verification with behavior changes; and
- record material discoveries that affect the governing context.

An implementation shortcut does not override a document. If code cannot
reasonably conform, return to step 3.

## 5. Verify the result

Evaluate every acceptance condition with repeatable evidence where practical.
Verification can include automated tests, static checks, builds, inspections,
or documented manual checks, depending on the change.

Verification must answer:

- Does the result deliver the TASK outcome?
- Is all changed work in scope?
- Does behavior conform to applicable SPECs?
- Does the implementation respect applicable ADRs?
- Are documentation, code, and tests materially consistent?
- Are failures, skipped checks, and limitations reported accurately?

Failed verification returns the work to implementation or, when it exposes an
authority gap, to context resolution.

## 6. Reconcile the repository

Before completion, update any governing document made inaccurate by the
change. Keep authoritative content in its proper document type and use
references rather than duplicate requirements or rationale.

Review the diff as a whole. Remove accidental scope, stale claims, temporary
notes, and generated output that is not intended to become authoritative
source.

## 7. Complete the task

Record concise completion evidence in the TASK, including:

- the resulting behavior or artifacts;
- checks actually performed and their results;
- relevant limitations or follow-up boundaries.

Mark the TASK `Complete` only after every acceptance condition is satisfied.
If the outcome is abandoned, mark it `Cancelled` and record why. If work cannot
proceed, mark it `Blocked` and identify the condition that must change.

Completed and cancelled tasks remain historical records. Further work receives
a new TASK.

## Handling discoveries

Use the destination that matches the discovery:

| Discovery | Action |
| --- | --- |
| Existing decision is missing or must change | Propose a new ADR |
| Approved behavior is ambiguous or must change | Revise or supersede the SPEC through review |
| Current outcome or scope is wrong | Revise the active TASK or create a follow-up TASK |
| Implementation defect within current scope | Fix and verify within the active TASK |
| Useful work outside current scope | Do not implement it; propose separate work |
| Governing documents conflict | Stop until the conflict is resolved |

The workflow deliberately makes stopping a valid engineering action. Pausing
at an authority or scope boundary prevents hidden decisions from becoming
maintenance debt.

## Review responsibilities

The contributor demonstrates that the change follows the governing context.
The reviewer checks both implementation quality and alignment:

- the TASK was ready and remained bounded;
- applicable ADRs and SPECs were followed;
- new policy was not hidden in code;
- acceptance evidence supports completion; and
- the repository is understandable without the implementation conversation.

Approval authority belongs to the project's human governance process.
Automation and agents may provide evidence or recommendations but do not
replace that responsibility.
