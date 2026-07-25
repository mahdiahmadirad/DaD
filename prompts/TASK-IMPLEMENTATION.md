# Task Implementation Prompt

## Purpose

Use this prompt to implement one approved DaD TASK. Replace `TASK-NNNN` with
the canonical identifier of the task to implement.

## Prompt

Implement `TASK-NNNN` only. Do not perform work outside that task.

Before changing the repository:

1. Read the repository working agreement and project vision.
2. Read `TASK-NNNN` completely.
3. Confirm that its status is `Ready` or `In Progress`. If it is `Proposed`,
   `Blocked`, `Complete`, or `Cancelled`, do not begin implementation; report
   the lifecycle condition that prevents work.
4. Read every accepted ADR and approved SPEC referenced by the TASK, following
   material references needed to understand their constraints.
5. Inspect the affected implementation, tests, local instructions, and
   version-control status.
6. Confirm that the outcome, scope, constraints, and acceptance conditions are
   sufficient and mutually consistent.

If governing sources conflict, required behavior is ambiguous, a
consequential architectural decision is missing, or completion requires
broader scope, stop and report the exact authority gap. Do not invent policy
inside the implementation.

When the task is implementable:

1. Mark it `In Progress` if it is `Ready`.
2. Make the smallest coherent change that satisfies its outcome.
3. Preserve unrelated contributor work and existing project patterns.
4. Do not add anticipated features, speculative abstractions, unrelated
   cleanup, or future-task preparation.
5. Add or update proportionate verification for changed behavior.
6. Reconcile documentation made inaccurate by the authorized change, keeping
   each concern in its authoritative document type.
7. Evaluate every acceptance condition and run the relevant available checks.
8. Record only checks actually performed, their results, and any limitations
   in the TASK's completion evidence.
9. Mark the TASK `Complete` only when every acceptance condition is satisfied.

Review the complete diff for scope, consistency, accidental changes, and
unresolved placeholders. Do not commit, push, publish, or begin follow-up work
unless separately authorized.

At completion, report:

- the outcome delivered;
- files changed;
- verification performed and results;
- remaining limitations or blockers; and
- the final TASK status.

Stop at the `TASK-NNNN` boundary.
