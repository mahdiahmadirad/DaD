# Repository Working Agreement

This file governs all work in the DaD repository. It applies equally to human
contributors and AI coding agents. More specific instructions may be added in
subdirectories later, but they must not silently contradict this file.

## Source of direction

Before changing the repository:

1. Read `PROJECT-VISION.md`.
2. Read the active `TASK-*.md` file.
3. Read any accepted decisions and approved specifications referenced by that
   task.
4. Inspect the implementation and tests affected by the change.

Do not treat chat history, generated plans, or tool-specific instructions as
durable project authority. If important context exists only outside the
repository, capture it in the appropriate project artifact before relying on
it.

If sources disagree, stop and surface the conflict. Do not choose whichever
source makes implementation easiest.

## Implementation rules

- Work only within the active task's scope and acceptance conditions.
- Do not implement anticipated phases, speculative abstractions, or unrelated
  cleanup.
- Prefer the smallest coherent change that satisfies the task.
- Preserve agent neutrality in framework semantics and repository conventions.
- Follow existing project patterns unless the task explicitly changes them.
- Keep authoritative information in one place and link to it elsewhere.
- Update governing documentation when a change makes it inaccurate.
- Add or update proportionate verification for behavior changes.
- Keep generated artifacts reproducible and clearly distinguish them from
  authoritative source files.
- Preserve user and contributor work that is unrelated to the active task.

## Document roles

- `PROJECT-VISION.md` defines durable purpose, audience, principles, scope, and
  non-goals.
- `TASK-*.md` defines a bounded change, its constraints, and its acceptance
  conditions.
- `docs/adr/` will contain accepted architectural decisions and their
  rationale.
- `docs/specs/` will contain approved, implementable behavior and interfaces.
- Code implements approved behavior.
- Tests and other checks provide repeatable evidence about the implementation.

Do not use one document type as a substitute for another. Create a new artifact
only when an approved task requires it.

## Task lifecycle

Each implementation phase must have one active task. A task should state:

- the problem or outcome;
- in-scope and out-of-scope work;
- relevant constraints and dependencies;
- acceptance conditions; and
- the completion evidence or outcome.

During implementation, record material discoveries that change the task's
meaning. If satisfying the task requires a new architectural decision, a
specification, or broader scope, pause and obtain approval rather than
inventing policy inside the implementation.

At completion:

1. Verify every acceptance condition.
2. Record the resulting evidence or concise outcome in the task.
3. Confirm that repository guidance still matches reality.
4. Stop at the phase boundary and wait for approval before beginning later
   work.

## Quality standard

A completed change should be understandable by a maintainer who was not part
of the conversation that produced it. It should be appropriately scoped,
internally consistent, verifiable, and free of dependencies on a particular AI
agent unless an approved integration explicitly requires one.

When no automated check exists, report the manual verification performed and
its limitations. Never claim to have run a check that was not run.

## Phase 1 boundary

The repository is currently bootstrapped, not implemented. Until a later task
is explicitly approved, do not create ADRs, specifications, framework code,
tools, CI workflows, a CLI, or examples.
