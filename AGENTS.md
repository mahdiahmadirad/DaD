# Repository Working Agreement

This file governs all work in the DaD repository. It applies equally to human
contributors and AI coding agents. More specific instructions may be added in
subdirectories later, but they must not silently contradict this file.

## Source of direction

Before changing the repository:

1. Read `PROJECT-VISION.md`.
2. Read the active `docs/tasks/TASK-*.md` file.
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
- `docs/tasks/TASK-*.md` defines bounded changes, their constraints, acceptance
  conditions, and completion evidence.
- `docs/adr/` contains accepted architectural decisions and their rationale.
- `docs/specs/` contains approved, implementable behavior and interfaces.
- Code implements approved behavior.
- Tests and other checks provide repeatable evidence about the implementation.

Do not use one document type as a substitute for another. Create a new artifact
only when an approved task requires it.

## Task lifecycle

Each implementation phase must have one active task under `docs/tasks/`. A task
should state:

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

## CLI implementation

The DaD CLI is a Go module governed by `ADR-0001`, `SPEC-0001`, and
`SPEC-0002`. Its public entry point is `cmd/dad`; reusable behavior belongs in
`internal/`. Canonical Markdown and prompts are embedded from their
authoritative repository files through `resources.go`.

For CLI changes, run:

```text
gofmt -w <changed Go files>
go vet ./...
go test ./...
```

Add platform-specific implementation only when the required behavior cannot be
expressed portably. Keep public behavior equivalent on Windows, Linux, and
macOS. Packaging, publishing, CI, integrations, and new commands require their
own explicitly approved tasks.

## Release implementation

`ADR-0002` and `SPEC-0003` govern versioning, changelog maintenance, release
artifacts, and GitHub publication. Release configuration is authoritative in
`.goreleaser.yaml`; `CHANGELOG.md` is the authoritative user-facing release
history.

For release-automation changes, also run:

```text
goreleaser check
goreleaser release --snapshot --clean
go run ./internal/releasetool verify 0.0.0-snapshot dist
```

A snapshot is local verification, not evidence of publication. Creating or
pushing a release tag and publishing a GitHub Release require the separately
approved release task and explicit maintainer authorization.
